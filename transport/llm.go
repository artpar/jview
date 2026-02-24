package transport

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"jview/jlog"
	"jview/protocol"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	anyllm "github.com/mozilla-ai/any-llm-go"
)

// LLMConfig configures the LLM transport.
type LLMConfig struct {
	Provider anyllm.Provider
	Model    string
	Prompt   string
	Mode     string // "tools" (default) or "raw"
	RecordTo string // path to write JSONL recording (empty = no recording)
}

type actionPayload struct {
	SurfaceID string
	Event     *protocol.EventDef
	Data      map[string]interface{}
}

// LLMTransport connects to an LLM provider and streams A2UI messages
// from the LLM's responses. User actions trigger new conversation turns.
type LLMTransport struct {
	config   LLMConfig
	messages chan *protocol.Message
	errors   chan error
	actions  chan actionPayload
	done     chan struct{}
	stopOnce sync.Once
	cancel   context.CancelFunc

	recorder *os.File // lazily opened JSONL recorder
	// OnInitialTurnDone is called after the first doTurn completes.
	// Use this to finalize cache files.
	OnInitialTurnDone func()
}

func NewLLMTransport(cfg LLMConfig) *LLMTransport {
	if cfg.Mode == "" {
		cfg.Mode = "tools"
	}
	return &LLMTransport{
		config:   cfg,
		messages: make(chan *protocol.Message, 64),
		errors:   make(chan error, 8),
		actions:  make(chan actionPayload, 16),
		done:     make(chan struct{}),
	}
}

func (t *LLMTransport) Messages() <-chan *protocol.Message {
	return t.messages
}

func (t *LLMTransport) Errors() <-chan error {
	return t.errors
}

func (t *LLMTransport) Start() {
	go t.run()
}

func (t *LLMTransport) Stop() {
	t.stopOnce.Do(func() {
		close(t.done)
		if t.cancel != nil {
			t.cancel()
		}
	})
}

func (t *LLMTransport) SendAction(surfaceID string, event *protocol.EventDef, data map[string]interface{}) {
	select {
	case t.actions <- actionPayload{SurfaceID: surfaceID, Event: event, Data: data}:
	case <-t.done:
	}
}

// recordLine writes a JSONL line to the recorder file. No-op if no RecordTo path.
func (t *LLMTransport) recordLine(line []byte) {
	if t.config.RecordTo == "" {
		return
	}
	if t.recorder == nil {
		f, err := os.Create(t.config.RecordTo)
		if err != nil {
			jlog.Errorf("transport", "", "failed to open recorder: %v", err)
			return
		}
		t.recorder = f
	}
	t.recorder.Write(line)
	t.recorder.Write([]byte("\n"))
}

// CloseRecorder closes the recording file if open.
func (t *LLMTransport) CloseRecorder() {
	if t.recorder != nil {
		t.recorder.Close()
		t.recorder = nil
	}
}

func (t *LLMTransport) run() {
	defer close(t.messages)
	defer close(t.errors)

	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel
	defer cancel()

	jlog.Infof("transport", "", "starting conversation with provider, model=%s mode=%s", t.config.Model, t.config.Mode)

	history := []anyllm.Message{
		{Role: anyllm.RoleSystem, Content: systemPrompt(t.config.Prompt)},
		{Role: anyllm.RoleUser, Content: t.config.Prompt},
	}

	firstTurn := true
	for {
		select {
		case <-t.done:
			return
		default:
		}

		history = t.doTurn(ctx, history)
		if history == nil {
			return
		}

		if firstTurn {
			firstTurn = false
			if t.OnInitialTurnDone != nil {
				t.OnInitialTurnDone()
			}
		}

		// Wait for a user action to trigger the next turn
		select {
		case ap := <-t.actions:
			userMsg := t.formatAction(ap)
			history = append(history, anyllm.Message{
				Role:    anyllm.RoleUser,
				Content: userMsg,
			})
		case <-t.done:
			return
		}
	}
}

// doTurn executes one LLM turn. Returns the updated history, or nil to stop.
func (t *LLMTransport) doTurn(ctx context.Context, history []anyllm.Message) []anyllm.Message {
	if t.config.Mode == "raw" {
		return t.doTurnRaw(ctx, history)
	}
	return t.doTurnTools(ctx, history)
}

// doTurnTools uses tool calling mode. Handles the tool call loop —
// the LLM may make multiple tool calls before finishing.
func (t *LLMTransport) doTurnTools(ctx context.Context, history []anyllm.Message) []anyllm.Message {
	for {
		select {
		case <-t.done:
			return nil
		default:
		}

		jlog.Infof("transport", "", "sending completion request (%d messages in history)", len(history))
		maxTok := 16384
		params := anyllm.CompletionParams{
			Model:     t.config.Model,
			Messages:  history,
			Tools:     a2uiTools(),
			MaxTokens: &maxTok,
		}

		var resp *anyllm.ChatCompletion
		var err error
		for attempt := 0; attempt < 3; attempt++ {
			resp, err = t.config.Provider.Completion(ctx, params)
			if err == nil || !isTransient(err) {
				break
			}
			delay := time.Duration(1<<uint(attempt)) * time.Second
			jlog.Warnf("transport", "", "transient error (attempt %d/3), retrying in %v: %v", attempt+1, delay, err)
			select {
			case <-time.After(delay):
			case <-t.done:
				return nil
			}
		}
		if err != nil {
			jlog.Errorf("transport", "", "completion error: %v", err)
			select {
			case t.errors <- fmt.Errorf("llm completion: %w", err):
			case <-t.done:
			}
			return nil
		}

		if len(resp.Choices) == 0 {
			jlog.Warn("transport", "", "empty response (no choices)")
			return history
		}

		choice := resp.Choices[0]
		jlog.Infof("transport", "", "got response, finish_reason=%s, tool_calls=%d", choice.FinishReason, len(choice.Message.ToolCalls))

		// Append assistant message to history
		history = append(history, choice.Message)

		// Process tool calls
		if len(choice.Message.ToolCalls) > 0 {
			for _, tc := range choice.Message.ToolCalls {
				jlog.Infof("transport", "", "processing tool call: %s", tc.Function.Name)

				// Handle utility tools that return data to the LLM (not protocol messages)
				if tc.Function.Name == "a2ui_inspectLibrary" {
					result := handleInspectLibrary(tc)
					history = append(history, anyllm.Message{
						Role:       anyllm.RoleTool,
						Content:    result,
						ToolCallID: tc.ID,
					})
					continue
				}
				if tc.Function.Name == "a2ui_getLogs" {
					result := handleGetLogs(tc)
					history = append(history, anyllm.Message{
						Role:       anyllm.RoleTool,
						Content:    result,
						ToolCallID: tc.ID,
					})
					continue
				}

				msg, rawBytes, err := toolCallToMessage(tc)
				if err != nil {
					jlog.Warnf("transport", "", "tool call parse error: %v", err)
					// Send error as tool result so the LLM knows
					history = append(history, anyllm.Message{
						Role:       anyllm.RoleTool,
						Content:    fmt.Sprintf("error: %v", err),
						ToolCallID: tc.ID,
					})
					continue
				}

				t.recordLine(rawBytes)

				select {
				case t.messages <- msg:
				case <-t.done:
					return nil
				}

				// Send success as tool result
				history = append(history, anyllm.Message{
					Role:       anyllm.RoleTool,
					Content:    "ok",
					ToolCallID: tc.ID,
				})
			}

			// If finish reason is tool_calls or length, loop to let the LLM continue
			if choice.FinishReason == anyllm.FinishReasonToolCalls || choice.FinishReason == anyllm.FinishReasonLength {
				continue
			}
		}

		// LLM is done for this turn
		return history
	}
}

// doTurnRaw uses raw text mode — the LLM outputs JSONL directly in its response.
func (t *LLMTransport) doTurnRaw(ctx context.Context, history []anyllm.Message) []anyllm.Message {
	chunks, errs := t.config.Provider.CompletionStream(ctx, anyllm.CompletionParams{
		Model:    t.config.Model,
		Messages: history,
	})

	var fullContent strings.Builder
	var lineBuf strings.Builder

	for chunk := range chunks {
		select {
		case <-t.done:
			return nil
		default:
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		text := chunk.Choices[0].Delta.Content
		if text == "" {
			continue
		}

		fullContent.WriteString(text)
		lineBuf.WriteString(text)

		// Process complete lines
		for {
			content := lineBuf.String()
			idx := strings.Index(content, "\n")
			if idx < 0 {
				break
			}
			line := strings.TrimSpace(content[:idx])
			lineBuf.Reset()
			lineBuf.WriteString(content[idx+1:])

			if line == "" {
				continue
			}

			parser := protocol.NewParser(strings.NewReader(line))
			msg, err := parser.Next()
			if err != nil {
				// Non-fatal — LLM may output non-JSONL text
				jlog.Debugf("transport", "", "raw parse skip: %v", err)
				continue
			}
			t.recordLine([]byte(line))
			select {
			case t.messages <- msg:
			case <-t.done:
				return nil
			}
		}
	}

	// Check for stream error
	if err := <-errs; err != nil {
		select {
		case t.errors <- fmt.Errorf("llm stream: %w", err):
		case <-t.done:
		}
		return nil
	}

	// Append assistant response to history
	history = append(history, anyllm.Message{
		Role:    anyllm.RoleAssistant,
		Content: fullContent.String(),
	})

	return history
}

// formatAction formats a user event into a message string for the LLM.
func (t *LLMTransport) formatAction(ap actionPayload) string {
	parts := []string{
		fmt.Sprintf("User action on surface %q:", ap.SurfaceID),
		fmt.Sprintf("  event: %s", ap.Event.Name),
	}
	if len(ap.Data) > 0 {
		data, _ := json.MarshalIndent(ap.Data, "  ", "  ")
		parts = append(parts, fmt.Sprintf("  data:\n  %s", string(data)))
	}
	return strings.Join(parts, "\n")
}

// isTransient returns true for errors that are likely to succeed on retry:
// rate limits, server errors, timeouts, and connection failures.
func isTransient(err error) bool {
	if err == nil {
		return false
	}
	// anyllm typed errors
	var rateLimitErr *anyllm.RateLimitError
	var providerErr *anyllm.ProviderError
	if stderrors.As(err, &rateLimitErr) {
		return true
	}
	if stderrors.As(err, &providerErr) {
		return true
	}
	// Network errors: connection refused, reset, timeout
	var netErr *net.OpError
	if stderrors.As(err, &netErr) {
		return true
	}
	if stderrors.Is(err, context.DeadlineExceeded) {
		return true
	}
	// String heuristics for wrapped errors
	msg := err.Error()
	for _, substr := range []string{"429", "500", "502", "503", "504", "timeout", "connection refused", "connection reset"} {
		if strings.Contains(msg, substr) {
			return true
		}
	}
	return false
}
