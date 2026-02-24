package transport

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func fixtureDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "testdata")
}

func TestFileTransportHappyPath(t *testing.T) {
	defer goroutineLeakCheck(t)()

	ft := NewFileTransport(filepath.Join(fixtureDir(), "hello.jsonl"))
	ft.Start()

	msgs := drainWithTimeout(t, ft.Messages(), 5*time.Second)
	if len(msgs) == 0 {
		t.Fatal("expected at least one message")
	}

	// Errors channel must also close
	errs := drainErrorsWithTimeout(t, ft.Errors(), time.Second)
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
}

func TestFileTransportMissingFile(t *testing.T) {
	defer goroutineLeakCheck(t)()

	ft := NewFileTransport("/nonexistent/path/to/file.jsonl")
	ft.Start()

	// Messages channel must close (no messages)
	msgs := drainWithTimeout(t, ft.Messages(), 5*time.Second)
	if len(msgs) != 0 {
		t.Errorf("expected 0 messages from missing file, got %d", len(msgs))
	}

	// Errors channel must close with one error
	errs := drainErrorsWithTimeout(t, ft.Errors(), time.Second)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestFileTransportStopBeforeComplete(t *testing.T) {
	defer goroutineLeakCheck(t)()

	ft := NewFileTransport(filepath.Join(fixtureDir(), "hello.jsonl"))
	ft.Start()

	// Stop immediately before reading all messages
	ft.Stop()

	// Both channels must close
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-ft.Messages():
			if !ok {
				goto msgsDone
			}
		case <-timer.C:
			t.Fatal("messages channel did not close after Stop()")
		}
	}
msgsDone:
	for {
		select {
		case _, ok := <-ft.Errors():
			if !ok {
				return
			}
		case <-timer.C:
			t.Fatal("errors channel did not close after Stop()")
		}
	}
}

func TestFileTransportMalformedJSON(t *testing.T) {
	defer goroutineLeakCheck(t)()

	// Create a temp file with bad JSON
	tmp, err := os.CreateTemp("", "jview-test-*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	tmp.WriteString("this is not json\n")
	tmp.Close()

	ft := NewFileTransport(tmp.Name())
	ft.Start()

	// Messages closes (possibly with 0 messages since first line fails)
	drainWithTimeout(t, ft.Messages(), 5*time.Second)

	// Errors channel must close with at least one error
	errs := drainErrorsWithTimeout(t, ft.Errors(), time.Second)
	if len(errs) == 0 {
		t.Fatal("expected parse error for malformed JSON")
	}
}

func TestFileTransportEmptyFile(t *testing.T) {
	defer goroutineLeakCheck(t)()

	tmp, err := os.CreateTemp("", "jview-test-*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	tmp.Close()

	ft := NewFileTransport(tmp.Name())
	ft.Start()

	msgs := drainWithTimeout(t, ft.Messages(), 5*time.Second)
	if len(msgs) != 0 {
		t.Errorf("expected 0 messages from empty file, got %d", len(msgs))
	}

	errs := drainErrorsWithTimeout(t, ft.Errors(), time.Second)
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
}
