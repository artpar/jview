package engine

import (
	"jview/protocol"
	"jview/renderer"
	"jview/transport"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// parseMessage parses a single JSONL line into a Message.
func parseMessage(t *testing.T, jsonl string) *protocol.Message {
	t.Helper()
	p := protocol.NewParser(strings.NewReader(jsonl))
	msg, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}
	return msg
}

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

// TestE2EFunctionCallsFixture validates function call evaluation.
func TestE2EFunctionCallsFixture(t *testing.T) {
	mock := renderer.NewMockRenderer()
	disp := &renderer.MockDispatcher{}
	sess := NewSession(mock, disp)

	ft := transport.NewFileTransport(filepath.Join(fixtureDir(), "function_calls.jsonl"))
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

	created := map[string]string{}
	for _, c := range mock.Created {
		created[c.Node.ComponentID] = c.Node.Props.Content
	}

	// greeting: concat("Hello, ", {path:"/name"}, "!") → "Hello, world!"
	if created["greeting"] != "Hello, world!" {
		t.Errorf("greeting = %q, want 'Hello, world!'", created["greeting"])
	}

	// upper: toUpperCase({path:"/name"}) → "WORLD"
	if created["upper"] != "WORLD" {
		t.Errorf("upper = %q, want 'WORLD'", created["upper"])
	}

	// computed: format("Name: {0}, Length: {1}", /name, length(/name)) → "Name: world, Length: 5"
	if created["computed"] != "Name: world, Length: 5" {
		t.Errorf("computed = %q, want 'Name: world, Length: 5'", created["computed"])
	}

	// Verify data binding: updating /name should re-render function call components
	before := len(mock.Updated)
	sess.HandleMessage(parseMessage(t, `{"type":"updateDataModel","surfaceId":"funcs","ops":[{"op":"replace","path":"/name","value":"test"}]}`))

	foundGreeting := false
	for _, u := range mock.Updated[before:] {
		if u.Node != nil && u.Node.ComponentID == "greeting" {
			foundGreeting = true
			if u.Node.Props.Content != "Hello, test!" {
				t.Errorf("greeting after update = %q, want 'Hello, test!'", u.Node.Props.Content)
			}
		}
	}
	if !foundGreeting {
		t.Error("function call component not re-rendered after data model update")
	}
}

// TestE2EListFixture validates template expansion and List component.
func TestE2EListFixture(t *testing.T) {
	mock := renderer.NewMockRenderer()
	disp := &renderer.MockDispatcher{}
	sess := NewSession(mock, disp)

	ft := transport.NewFileTransport(filepath.Join(fixtureDir(), "list.jsonl"))
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

	// Check expanded components exist
	created := map[string]string{}
	createdTitles := map[string]string{}
	for _, c := range mock.Created {
		created[c.Node.ComponentID] = c.Node.Props.Content
		createdTitles[c.Node.ComponentID] = c.Node.Props.Title
	}

	// 3 items in data model → 3 card+role pairs
	if createdTitles["itemCard_0"] != "Alice" {
		t.Errorf("itemCard_0 title = %q, want Alice", createdTitles["itemCard_0"])
	}
	if createdTitles["itemCard_1"] != "Bob" {
		t.Errorf("itemCard_1 title = %q, want Bob", createdTitles["itemCard_1"])
	}
	if createdTitles["itemCard_2"] != "Charlie" {
		t.Errorf("itemCard_2 title = %q, want Charlie", createdTitles["itemCard_2"])
	}
	if created["itemRole_0"] != "Engineer" {
		t.Errorf("itemRole_0 = %q, want Engineer", created["itemRole_0"])
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

// TestE2ECalculatorTests runs the calculator test fixture through the test runner.
func TestE2ECalculatorTests(t *testing.T) {
	mock := renderer.NewMockRenderer()
	disp := &renderer.MockDispatcher{}

	results, err := RunTestFile(filepath.Join(fixtureDir(), "calculator_test.jsonl"), mock, disp)
	if err != nil {
		t.Fatalf("RunTestFile: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("no tests found in calculator_test.jsonl")
	}

	passed := 0
	totalAssertions := 0
	for _, r := range results {
		totalAssertions += r.Assertions
		if r.Passed {
			passed++
		} else {
			t.Errorf("FAIL: %s: %s", r.Name, r.Error)
		}
	}

	t.Logf("%d/%d tests passed, %d assertions total", passed, len(results), totalAssertions)
}
