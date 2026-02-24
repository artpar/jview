package engine

import (
	"jview/renderer"
	"jview/transport"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func fixtureDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "testdata")
}

// TestE2EHelloFixture reads the actual hello.jsonl through file transport
// into the engine with a mock renderer and validates the output.
func TestE2EHelloFixture(t *testing.T) {
	mock := renderer.NewMockRenderer()
	disp := &renderer.MockDispatcher{}
	sess := NewSession(mock, disp)

	ft := transport.NewFileTransport(filepath.Join(fixtureDir(), "hello.jsonl"))
	ft.Start()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for msg := range ft.Messages() {
			sess.HandleMessage(msg)
		}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for fixture")
	}

	if len(mock.Windows) != 1 {
		t.Fatalf("windows = %d, want 1", len(mock.Windows))
	}
	if mock.Windows[0].Title != "Hello jview" {
		t.Errorf("title = %q", mock.Windows[0].Title)
	}

	created := map[string]string{}
	for _, c := range mock.Created {
		created[c.Node.ComponentID] = c.Node.Props.Content
	}

	if _, ok := created["card1"]; !ok {
		t.Error("card1 not created")
	}
	if created["heading"] != "Hello, jview!" {
		t.Errorf("heading content = %q", created["heading"])
	}
}

// TestE2EContactFormFixture validates the contact form with data binding.
func TestE2EContactFormFixture(t *testing.T) {
	mock := renderer.NewMockRenderer()
	disp := &renderer.MockDispatcher{}
	sess := NewSession(mock, disp)

	ft := transport.NewFileTransport(filepath.Join(fixtureDir(), "contact_form.jsonl"))
	ft.Start()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for msg := range ft.Messages() {
			sess.HandleMessage(msg)
		}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}

	// All expected components exist
	created := map[string]bool{}
	for _, c := range mock.Created {
		created[c.Node.ComponentID] = true
	}
	for _, id := range []string{"root", "title", "nameField", "emailField", "previewCard", "subscribeCheck", "submitBtn"} {
		if !created[id] {
			t.Errorf("missing component: %s", id)
		}
	}

	// Data binding works: type in name field, preview updates
	before := len(mock.Updated)
	mock.InvokeCallback("form", "nameField", "change", "Test User")

	foundUpdate := false
	for _, u := range mock.Updated[before:] {
		if u.Node != nil && u.Node.ComponentID == "previewName" && u.Node.Props.Content == "Test User" {
			foundUpdate = true
		}
	}
	if !foundUpdate {
		t.Error("name binding did not propagate to previewName")
	}
}

// TestE2ELayoutFixture validates nested row/column layout.
func TestE2ELayoutFixture(t *testing.T) {
	mock := renderer.NewMockRenderer()
	disp := &renderer.MockDispatcher{}
	sess := NewSession(mock, disp)

	ft := transport.NewFileTransport(filepath.Join(fixtureDir(), "layout.jsonl"))
	ft.Start()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for msg := range ft.Messages() {
			sess.HandleMessage(msg)
		}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}

	// Find the Row component and verify it has 2 children
	contentHandle := mock.GetHandle("layout", "content")
	if contentHandle == 0 {
		t.Fatal("content (Row) not created")
	}

	foundRow := false
	for _, cs := range mock.Children {
		if cs.ParentHandle == contentHandle {
			foundRow = true
			if len(cs.ChildHandles) != 2 {
				t.Errorf("Row children = %d, want 2", len(cs.ChildHandles))
			}
		}
	}
	if !foundRow {
		t.Error("no SetChildren for content Row")
	}
}
