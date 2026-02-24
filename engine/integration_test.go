package engine

import (
	"jview/protocol"
	"jview/renderer"
	"strings"
	"testing"
)

// feedMessages parses JSONL and feeds to session. Returns after all processed.
func feedMessages(t *testing.T, sess *Session, jsonl string) {
	t.Helper()
	p := protocol.NewParser(strings.NewReader(jsonl))
	for {
		msg, err := p.Next()
		if err != nil {
			break
		}
		sess.HandleMessage(msg)
	}
}

func TestHelloFixtureCreatesWindowAndViews(t *testing.T) {
	mock := renderer.NewMockRenderer()
	disp := &renderer.MockDispatcher{}
	sess := NewSession(mock, disp)

	feedMessages(t, sess, `{"type":"createSurface","surfaceId":"main","title":"Hello","width":600,"height":400}
{"type":"updateComponents","surfaceId":"main","components":[{"componentId":"card1","type":"Card","props":{"title":"Welcome"},"children":["heading","body"]},{"componentId":"heading","type":"Text","props":{"content":"Hello, jview!","variant":"h1"}},{"componentId":"body","type":"Text","props":{"content":"Body text.","variant":"body"}}]}`)

	// Window created
	if len(mock.Windows) != 1 {
		t.Fatalf("windows = %d, want 1", len(mock.Windows))
	}
	if mock.Windows[0].Title != "Hello" {
		t.Errorf("window title = %q, want Hello", mock.Windows[0].Title)
	}

	// All 3 components created
	if len(mock.Created) < 3 {
		t.Fatalf("created views = %d, want >= 3", len(mock.Created))
	}

	// Find the created components by ID
	created := map[string]*renderer.RenderNode{}
	for _, c := range mock.Created {
		created[c.Node.ComponentID] = c.Node
	}

	if _, ok := created["heading"]; !ok {
		t.Error("heading not created")
	}
	if _, ok := created["body"]; !ok {
		t.Error("body not created")
	}
	if _, ok := created["card1"]; !ok {
		t.Error("card1 not created")
	}

	// Heading has correct resolved props
	if h := created["heading"]; h != nil {
		if h.Props.Content != "Hello, jview!" {
			t.Errorf("heading content = %q, want 'Hello, jview!'", h.Props.Content)
		}
		if h.Props.Variant != "h1" {
			t.Errorf("heading variant = %q, want h1", h.Props.Variant)
		}
	}

	// Card has children set
	if len(mock.Children) == 0 {
		t.Fatal("no SetChildren calls")
	}
	cardHandle := mock.GetHandle("main", "card1")
	foundCardChildren := false
	for _, cs := range mock.Children {
		if cs.ParentHandle == cardHandle {
			foundCardChildren = true
			if len(cs.ChildHandles) != 2 {
				t.Errorf("card child count = %d, want 2", len(cs.ChildHandles))
			}
		}
	}
	if !foundCardChildren {
		t.Error("no SetChildren call for card1")
	}

	// Root set
	if len(mock.RootSets) == 0 {
		t.Fatal("no SetRootView calls")
	}
}

func TestDataBindingPropagates(t *testing.T) {
	mock := renderer.NewMockRenderer()
	disp := &renderer.MockDispatcher{}
	sess := NewSession(mock, disp)

	feedMessages(t, sess, `{"type":"createSurface","surfaceId":"form","title":"Form"}
{"type":"updateDataModel","surfaceId":"form","ops":[{"op":"add","path":"/name","value":""}]}
{"type":"updateComponents","surfaceId":"form","components":[{"componentId":"field","type":"TextField","props":{"value":{"path":"/name"},"dataBinding":"/name"}},{"componentId":"display","type":"Text","props":{"content":{"path":"/name"},"variant":"body"}}]}`)

	// Both components created
	created := map[string]bool{}
	for _, c := range mock.Created {
		created[c.Node.ComponentID] = true
	}
	if !created["field"] || !created["display"] {
		t.Fatalf("missing components: field=%v display=%v", created["field"], created["display"])
	}

	// Simulate typing in the text field
	initialUpdates := len(mock.Updated)
	mock.InvokeCallback("form", "field", "change", "Alice")

	// The display text should have been re-rendered with new value
	newUpdates := mock.Updated[initialUpdates:]
	foundDisplayUpdate := false
	for _, u := range newUpdates {
		node := u.Node
		if node != nil && node.ComponentID == "display" {
			foundDisplayUpdate = true
			if node.Props.Content != "Alice" {
				t.Errorf("display content after binding = %q, want Alice", node.Props.Content)
			}
		}
	}
	if !foundDisplayUpdate {
		t.Error("data binding did not trigger display update")
	}
}

func TestCheckBoxBindingPropagates(t *testing.T) {
	mock := renderer.NewMockRenderer()
	disp := &renderer.MockDispatcher{}
	sess := NewSession(mock, disp)

	feedMessages(t, sess, `{"type":"createSurface","surfaceId":"s1","title":"T"}
{"type":"updateDataModel","surfaceId":"s1","ops":[{"op":"add","path":"/agree","value":false}]}
{"type":"updateComponents","surfaceId":"s1","components":[{"componentId":"cb","type":"CheckBox","props":{"label":"Agree","checked":{"path":"/agree"},"dataBinding":"/agree"}},{"componentId":"status","type":"Text","props":{"content":{"path":"/agree"},"variant":"body"}}]}`)

	initialUpdates := len(mock.Updated)
	mock.InvokeCallback("s1", "cb", "toggle", "true")

	newUpdates := mock.Updated[initialUpdates:]
	foundStatusUpdate := false
	for _, u := range newUpdates {
		if u.Node != nil && u.Node.ComponentID == "status" {
			foundStatusUpdate = true
			if u.Node.Props.Content != "true" {
				t.Errorf("status content = %q, want 'true'", u.Node.Props.Content)
			}
		}
	}
	if !foundStatusUpdate {
		t.Error("checkbox binding did not trigger status update")
	}
}

func TestButtonCallbackFires(t *testing.T) {
	mock := renderer.NewMockRenderer()
	disp := &renderer.MockDispatcher{}
	sess := NewSession(mock, disp)

	var firedAction string
	// Hook into the surface's action handler after creation
	feedMessages(t, sess, `{"type":"createSurface","surfaceId":"s1","title":"T"}
{"type":"updateComponents","surfaceId":"s1","components":[{"componentId":"btn","type":"Button","props":{"label":"Go","style":"primary","onClick":{"action":{"type":"serverAction","name":"doThing"}}}}]}`)

	// Set action handler on the surface
	if surf, ok := sess.surfaces["s1"]; ok {
		surf.ActionHandler = func(sid string, action *protocol.Action, data map[string]interface{}) {
			firedAction = action.Name
		}
	}

	mock.InvokeCallback("s1", "btn", "click", "")
	if firedAction != "doThing" {
		t.Errorf("firedAction = %q, want doThing", firedAction)
	}
}

func TestMultipleDataModelOps(t *testing.T) {
	mock := renderer.NewMockRenderer()
	disp := &renderer.MockDispatcher{}
	sess := NewSession(mock, disp)

	feedMessages(t, sess, `{"type":"createSurface","surfaceId":"s1","title":"T"}
{"type":"updateDataModel","surfaceId":"s1","ops":[{"op":"add","path":"/a","value":"1"}]}
{"type":"updateComponents","surfaceId":"s1","components":[{"componentId":"t1","type":"Text","props":{"content":{"path":"/a"}}}]}`)

	// Find initial content
	var initialContent string
	for _, c := range mock.Created {
		if c.Node.ComponentID == "t1" {
			initialContent = c.Node.Props.Content
		}
	}
	if initialContent != "1" {
		t.Errorf("initial content = %q, want '1'", initialContent)
	}

	// Update data model
	beforeUpdates := len(mock.Updated)
	feedMessages(t, sess, `{"type":"updateDataModel","surfaceId":"s1","ops":[{"op":"replace","path":"/a","value":"2"}]}`)

	newUpdates := mock.Updated[beforeUpdates:]
	found := false
	for _, u := range newUpdates {
		if u.Node != nil && u.Node.ComponentID == "t1" {
			found = true
			if u.Node.Props.Content != "2" {
				t.Errorf("updated content = %q, want '2'", u.Node.Props.Content)
			}
		}
	}
	if !found {
		t.Error("data model update did not trigger re-render")
	}
}

func TestRenderOrderLeavesBeforeParents(t *testing.T) {
	mock := renderer.NewMockRenderer()
	disp := &renderer.MockDispatcher{}
	sess := NewSession(mock, disp)

	// Parent references children that come later in the array
	feedMessages(t, sess, `{"type":"createSurface","surfaceId":"s1","title":"T"}
{"type":"updateComponents","surfaceId":"s1","components":[{"componentId":"parent","type":"Column","children":["child1","child2"]},{"componentId":"child1","type":"Text","props":{"content":"A"}},{"componentId":"child2","type":"Text","props":{"content":"B"}}]}`)

	// All 3 must be created
	if len(mock.Created) < 3 {
		t.Fatalf("created = %d, want >= 3", len(mock.Created))
	}

	// Children must be created before parent
	order := map[string]int{}
	for i, c := range mock.Created {
		order[c.Node.ComponentID] = i
	}
	if order["child1"] > order["parent"] {
		t.Error("child1 created after parent — topological sort failed")
	}
	if order["child2"] > order["parent"] {
		t.Error("child2 created after parent — topological sort failed")
	}
}

func TestMultiRootWrapping(t *testing.T) {
	sess, mock := newTestSession()

	feedMessages(t, sess, `{"type":"createSurface","surfaceId":"s1","title":"T"}
{"type":"updateComponents","surfaceId":"s1","components":[{"componentId":"t1","type":"Text","props":{"content":"One"}},{"componentId":"t2","type":"Text","props":{"content":"Two"}}]}`)

	// __root_wrapper__ Column should be created
	wrapperCreated := false
	for _, c := range mock.Created {
		if c.Node.ComponentID == "__root_wrapper__" {
			wrapperCreated = true
			if c.Node.Type != protocol.CompColumn {
				t.Errorf("wrapper type = %q, want Column", c.Node.Type)
			}
		}
	}
	if !wrapperCreated {
		t.Fatal("__root_wrapper__ not created for multi-root surface")
	}

	// SetChildren called on wrapper with both root handles
	wrapperHandle := mock.GetHandle("s1", "__root_wrapper__")
	if wrapperHandle == 0 {
		t.Fatal("wrapper handle not found")
	}
	foundWrapperChildren := false
	for _, cs := range mock.Children {
		if cs.ParentHandle == wrapperHandle {
			foundWrapperChildren = true
			if len(cs.ChildHandles) != 2 {
				t.Errorf("wrapper child count = %d, want 2", len(cs.ChildHandles))
			}
		}
	}
	if !foundWrapperChildren {
		t.Error("no SetChildren call for wrapper")
	}

	// SetRootView called with wrapper handle
	foundRoot := false
	for _, rs := range mock.RootSets {
		if rs.Handle == wrapperHandle {
			foundRoot = true
		}
	}
	if !foundRoot {
		t.Error("SetRootView not called with wrapper handle")
	}
}

func TestCallbackReregistrationCleansUp(t *testing.T) {
	sess, mock := newTestSession()

	feedMessages(t, sess, `{"type":"createSurface","surfaceId":"s1","title":"T"}
{"type":"updateDataModel","surfaceId":"s1","ops":[{"op":"add","path":"/name","value":""}]}
{"type":"updateComponents","surfaceId":"s1","components":[{"componentId":"field","type":"TextField","props":{"value":{"path":"/name"},"dataBinding":"/name"}}]}`)

	// Get the first callback ID
	firstCBID := mock.GetCallbackID("s1", "field", "change")
	if firstCBID == 0 {
		t.Fatal("no callback registered for field change")
	}

	// Re-render the same TextField (simulating an update)
	feedMessages(t, sess, `{"type":"updateComponents","surfaceId":"s1","components":[{"componentId":"field","type":"TextField","props":{"value":{"path":"/name"},"dataBinding":"/name"}}]}`)

	// Old callback should be unregistered
	if mock.HasCallback(firstCBID) {
		t.Error("old callback still registered after re-render")
	}

	// New callback should exist
	newCBID := mock.GetCallbackID("s1", "field", "change")
	if newCBID == 0 {
		t.Fatal("no new callback registered")
	}
	if newCBID == firstCBID {
		t.Error("new callback ID same as old — expected fresh registration")
	}
}

func TestCallbackBindingPathChange(t *testing.T) {
	sess, mock := newTestSession()

	feedMessages(t, sess, `{"type":"createSurface","surfaceId":"s1","title":"T"}
{"type":"updateDataModel","surfaceId":"s1","ops":[{"op":"add","path":"/name","value":"Alice"},{"op":"add","path":"/nickname","value":""}]}
{"type":"updateComponents","surfaceId":"s1","components":[{"componentId":"field","type":"TextField","props":{"value":{"path":"/name"},"dataBinding":"/name"}},{"componentId":"display","type":"Text","props":{"content":{"path":"/nickname"}}}]}`)

	// Now change binding from /name to /nickname
	feedMessages(t, sess, `{"type":"updateComponents","surfaceId":"s1","components":[{"componentId":"field","type":"TextField","props":{"value":{"path":"/nickname"},"dataBinding":"/nickname"}}]}`)

	// Typing should write to /nickname
	beforeUpdates := len(mock.Updated)
	mock.InvokeCallback("s1", "field", "change", "Ally")

	// Check the data model was updated at /nickname
	if surf, ok := sess.surfaces["s1"]; ok {
		val, found := surf.dm.Get("/nickname")
		if !found || val != "Ally" {
			t.Errorf("/nickname = %v (%v), want Ally", val, found)
		}
	}

	// The display component bound to /nickname should be updated
	foundDisplay := false
	for _, u := range mock.Updated[beforeUpdates:] {
		if u.Node != nil && u.Node.ComponentID == "display" {
			foundDisplay = true
			if u.Node.Props.Content != "Ally" {
				t.Errorf("display content = %q, want Ally", u.Node.Props.Content)
			}
		}
	}
	if !foundDisplay {
		t.Error("display not updated after binding path change")
	}
}

func TestContactFormFixtureFull(t *testing.T) {
	mock := renderer.NewMockRenderer()
	disp := &renderer.MockDispatcher{}
	sess := NewSession(mock, disp)

	feedMessages(t, sess, `{"type":"createSurface","surfaceId":"form","title":"Contact Form","width":500,"height":500}
{"type":"updateDataModel","surfaceId":"form","ops":[{"op":"add","path":"/name","value":""},{"op":"add","path":"/email","value":""},{"op":"add","path":"/subscribe","value":false}]}
{"type":"updateComponents","surfaceId":"form","components":[{"componentId":"root","type":"Column","props":{"gap":16,"padding":20},"children":["title","nameField","emailField","previewName","submitBtn"]},{"componentId":"title","type":"Text","props":{"content":"Contact Us","variant":"h2"}},{"componentId":"nameField","type":"TextField","props":{"placeholder":"Enter your name","value":{"path":"/name"},"dataBinding":"/name"}},{"componentId":"emailField","type":"TextField","props":{"placeholder":"you@example.com","value":{"path":"/email"},"dataBinding":"/email"}},{"componentId":"previewName","type":"Text","props":{"content":{"path":"/name"},"variant":"body"}},{"componentId":"submitBtn","type":"Button","props":{"label":"Submit","style":"primary","onClick":{"action":{"type":"serverAction","name":"submitForm","dataRefs":["/name","/email","/subscribe"]}}}}]}`)

	// Window + all 6 components created
	if len(mock.Windows) != 1 {
		t.Fatalf("windows = %d", len(mock.Windows))
	}

	created := map[string]bool{}
	for _, c := range mock.Created {
		created[c.Node.ComponentID] = true
	}
	for _, id := range []string{"root", "title", "nameField", "emailField", "previewName", "submitBtn"} {
		if !created[id] {
			t.Errorf("component %q not created", id)
		}
	}

	// Type in name field → preview updates
	before := len(mock.Updated)
	mock.InvokeCallback("form", "nameField", "change", "Jane")

	foundPreview := false
	for _, u := range mock.Updated[before:] {
		if u.Node != nil && u.Node.ComponentID == "previewName" {
			foundPreview = true
			if u.Node.Props.Content != "Jane" {
				t.Errorf("preview = %q, want Jane", u.Node.Props.Content)
			}
		}
	}
	if !foundPreview {
		t.Error("typing in nameField did not update previewName")
	}
}

func TestSessionUnknownSurface(t *testing.T) {
	sess, mock := newTestSession()

	// Sending updateComponents to a non-existent surface should not panic
	feedMessages(t, sess, `{"type":"updateComponents","surfaceId":"nonexistent","components":[{"componentId":"t1","type":"Text","props":{"content":"hello"}}]}`)

	if len(mock.Created) != 0 {
		t.Errorf("expected no created views, got %d", len(mock.Created))
	}
}

func TestSessionDuplicateCreateSurface(t *testing.T) {
	sess, mock := newTestSession()

	feedMessages(t, sess, `{"type":"createSurface","surfaceId":"s1","title":"First"}
{"type":"createSurface","surfaceId":"s1","title":"Second"}`)

	// Should only create one window (duplicate ignored)
	if len(mock.Windows) != 1 {
		t.Errorf("windows = %d, want 1", len(mock.Windows))
	}
	if mock.Windows[0].Title != "First" {
		t.Errorf("title = %q, want First", mock.Windows[0].Title)
	}
}

func TestSessionDeleteSurface(t *testing.T) {
	sess, mock := newTestSession()

	feedMessages(t, sess, `{"type":"createSurface","surfaceId":"s1","title":"T"}
{"type":"deleteSurface","surfaceId":"s1"}`)

	// After deletion, updateComponents should be ignored
	feedMessages(t, sess, `{"type":"updateComponents","surfaceId":"s1","components":[{"componentId":"t1","type":"Text","props":{"content":"hello"}}]}`)

	if len(mock.Created) != 0 {
		t.Errorf("expected no created views after surface deletion, got %d", len(mock.Created))
	}
}

func TestSessionDeleteNonexistent(t *testing.T) {
	sess, _ := newTestSession()

	// Should not panic
	feedMessages(t, sess, `{"type":"deleteSurface","surfaceId":"nope"}`)
}

func TestSessionUnknownMessageType(t *testing.T) {
	sess, _ := newTestSession()

	// Should not panic — unknown type handled gracefully
	msg := &protocol.Message{
		Type:      "totallyBogus",
		SurfaceID: "s1",
	}
	sess.HandleMessage(msg)
}
