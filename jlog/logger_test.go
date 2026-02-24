package jlog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestLogger(maxEntries int) *Logger {
	return &Logger{
		entries:  make([]Entry, maxEntries),
		maxSize:  maxEntries,
		minLevel: LevelDebug,
	}
}

func TestLogAndQuery(t *testing.T) {
	l := newTestLogger(100)

	l.log(LevelInfo, "session", "main", "surface created")
	l.log(LevelWarn, "transport", "", "connection lost")
	l.log(LevelError, "darwin", "main", "view creation failed")
	l.log(LevelDebug, "resolver", "main", "resolving bindings")

	result := l.query(QueryOpts{MinLevel: LevelInfo})
	if result.Total != 3 {
		t.Fatalf("expected 3 entries at Info+, got %d", result.Total)
	}

	result = l.query(QueryOpts{MinLevel: LevelWarn})
	if result.Total != 2 {
		t.Fatalf("expected 2 entries at Warn+, got %d", result.Total)
	}

	result = l.query(QueryOpts{MinLevel: LevelError})
	if result.Total != 1 {
		t.Fatalf("expected 1 entry at Error+, got %d", result.Total)
	}
}

func TestQueryByComponent(t *testing.T) {
	l := newTestLogger(100)

	l.log(LevelInfo, "session", "", "msg1")
	l.log(LevelInfo, "transport", "", "msg2")
	l.log(LevelInfo, "session", "", "msg3")

	result := l.query(QueryOpts{Component: "session"})
	if result.Total != 2 {
		t.Fatalf("expected 2 session entries, got %d", result.Total)
	}
}

func TestQueryBySurface(t *testing.T) {
	l := newTestLogger(100)

	l.log(LevelInfo, "session", "main", "msg1")
	l.log(LevelInfo, "session", "modal", "msg2")
	l.log(LevelInfo, "session", "main", "msg3")

	result := l.query(QueryOpts{Surface: "main"})
	if result.Total != 2 {
		t.Fatalf("expected 2 main entries, got %d", result.Total)
	}
}

func TestQueryWithPattern(t *testing.T) {
	l := newTestLogger(100)

	l.log(LevelInfo, "session", "", "surface created")
	l.log(LevelInfo, "session", "", "surface deleted")
	l.log(LevelWarn, "transport", "", "connection lost")

	result := l.query(QueryOpts{Pattern: "surface"})
	if result.Total != 2 {
		t.Fatalf("expected 2 pattern matches, got %d", result.Total)
	}

	result = l.query(QueryOpts{Pattern: "surface.*deleted"})
	if result.Total != 1 {
		t.Fatalf("expected 1 regex match, got %d", result.Total)
	}
}

func TestQueryPagination(t *testing.T) {
	l := newTestLogger(100)

	for i := 0; i < 10; i++ {
		l.log(LevelInfo, "test", "", "msg")
	}

	result := l.query(QueryOpts{Limit: 3})
	if len(result.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result.Entries))
	}
	if result.Total != 10 {
		t.Fatalf("expected total 10, got %d", result.Total)
	}

	result = l.query(QueryOpts{Limit: 3, Offset: 8})
	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries at offset 8, got %d", len(result.Entries))
	}

	result = l.query(QueryOpts{Offset: 100})
	if len(result.Entries) != 0 {
		t.Fatalf("expected 0 entries at offset 100, got %d", len(result.Entries))
	}
}

func TestRingBufferOverflow(t *testing.T) {
	l := newTestLogger(5)

	// Write 8 entries — first 3 should be overwritten
	for i := 0; i < 8; i++ {
		l.log(LevelInfo, "test", "", strings.Repeat("x", i+1))
	}

	result := l.query(QueryOpts{})
	if result.Total != 5 {
		t.Fatalf("expected 5 entries after overflow, got %d", result.Total)
	}

	// Oldest entry should be the 4th write (message "xxxx")
	if result.Entries[0].Message != "xxxx" {
		t.Fatalf("expected oldest entry to be 'xxxx', got %q", result.Entries[0].Message)
	}
	// Newest should be the 8th write (message "xxxxxxxx")
	if result.Entries[4].Message != "xxxxxxxx" {
		t.Fatalf("expected newest entry to be 'xxxxxxxx', got %q", result.Entries[4].Message)
	}
}

func TestMinLevel(t *testing.T) {
	l := newTestLogger(100)
	l.minLevel = LevelWarn

	l.log(LevelDebug, "test", "", "debug msg")
	l.log(LevelInfo, "test", "", "info msg")
	l.log(LevelWarn, "test", "", "warn msg")
	l.log(LevelError, "test", "", "error msg")

	result := l.query(QueryOpts{})
	if result.Total != 2 {
		t.Fatalf("expected 2 entries (warn+error), got %d", result.Total)
	}
}

func TestFileOutput(t *testing.T) {
	dir := t.TempDir()
	l := &Logger{
		entries:  make([]Entry, 100),
		maxSize:  100,
		minLevel: LevelInfo,
	}

	path := filepath.Join(dir, "test.log")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create log file: %v", err)
	}
	l.file = f

	l.log(LevelInfo, "session", "main", "hello world")
	l.log(LevelWarn, "transport", "", "warning message")
	l.file.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "INFO [session/main] hello world") {
		t.Fatalf("log file missing expected content, got:\n%s", content)
	}
	if !strings.Contains(content, "WARN [transport] warning message") {
		t.Fatalf("log file missing warning, got:\n%s", content)
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"warn", LevelWarn},
		{"warning", LevelWarn},
		{"error", LevelError},
		{"ERROR", LevelError},
		{"unknown", LevelInfo},
	}
	for _, tt := range tests {
		got := ParseLevel(tt.input)
		if got != tt.expected {
			t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestSequenceNumbers(t *testing.T) {
	l := newTestLogger(100)

	l.log(LevelInfo, "test", "", "first")
	l.log(LevelInfo, "test", "", "second")
	l.log(LevelInfo, "test", "", "third")

	result := l.query(QueryOpts{})
	if result.Total != 3 {
		t.Fatalf("expected 3, got %d", result.Total)
	}

	for i := 1; i < len(result.Entries); i++ {
		if result.Entries[i].Seq <= result.Entries[i-1].Seq {
			t.Fatalf("sequence numbers not monotonically increasing: %d <= %d",
				result.Entries[i].Seq, result.Entries[i-1].Seq)
		}
	}
}
