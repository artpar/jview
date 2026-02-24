package protocol

import "encoding/json"

// ProcessTransportConfig describes how a process obtains its messages.
type ProcessTransportConfig struct {
	Type     string          `json:"type"`               // "file", "llm", "interval"
	Path     string          `json:"path,omitempty"`     // file transport: JSONL path
	Provider string          `json:"provider,omitempty"` // llm transport: provider name
	Model    string          `json:"model,omitempty"`    // llm transport: model name
	Prompt   string          `json:"prompt,omitempty"`   // llm transport: system prompt
	Interval int             `json:"interval,omitempty"` // interval transport: milliseconds
	Message  json.RawMessage `json:"message,omitempty"`  // interval transport: message to send on each tick
}

// CreateProcess starts a named, managed goroutine with its own transport.
type CreateProcess struct {
	Type      MessageType            `json:"type"`
	ProcessID string                 `json:"processId"`
	Transport ProcessTransportConfig `json:"transport"`
}

// StopProcess terminates a running process.
type StopProcess struct {
	Type      MessageType `json:"type"`
	ProcessID string      `json:"processId"`
}

// SendToProcess routes a message to a running process's transport.
type SendToProcess struct {
	Type      MessageType     `json:"type"`
	ProcessID string          `json:"processId"`
	Message   json.RawMessage `json:"message"`
}
