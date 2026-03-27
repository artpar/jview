package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractComponent(t *testing.T) {
	dir := t.TempDir()
	jsonlPath := filepath.Join(dir, "test.jsonl")

	// Write a simple JSONL file with createSurface + updateDataModel + updateComponents
	content := `{"type":"createSurface","surfaceId":"s1","title":"Test","width":400,"height":300}
{"type":"updateDataModel","surfaceId":"s1","ops":[{"op":"add","path":"/val","value":""}]}
{"type":"updateComponents","surfaceId":"s1","components":[{"componentId":"root","type":"Column","props":{"gap":4},"children":["btn1","btn2"]},{"componentId":"btn1","type":"Button","props":{"label":"A"}},{"componentId":"btn2","type":"Button","props":{"label":"B"}}]}
`
	if err := os.WriteFile(jsonlPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	dc, err := ExtractComponent(jsonlPath, "TestWidget")
	if err != nil {
		t.Fatalf("ExtractComponent failed: %v", err)
	}

	if dc.Name != "TestWidget" {
		t.Errorf("name = %q, want TestWidget", dc.Name)
	}
	if dc.Type != "defineComponent" {
		t.Errorf("type = %q, want defineComponent", dc.Type)
	}
	if len(dc.Components) != 3 {
		t.Fatalf("got %d components, want 3", len(dc.Components))
	}

	// Check that root was renamed to _root
	var rootComp map[string]interface{}
	if err := json.Unmarshal(dc.Components[0], &rootComp); err != nil {
		t.Fatal(err)
	}
	if rootComp["componentId"] != "_root" {
		t.Errorf("root componentId = %q, want _root", rootComp["componentId"])
	}

	// Check other components are preserved
	var btn1 map[string]interface{}
	if err := json.Unmarshal(dc.Components[1], &btn1); err != nil {
		t.Fatal(err)
	}
	if btn1["componentId"] != "btn1" {
		t.Errorf("btn1 componentId = %q, want btn1", btn1["componentId"])
	}
}

func TestExtractComponentMultipleBatches(t *testing.T) {
	dir := t.TempDir()
	jsonlPath := filepath.Join(dir, "test.jsonl")

	// Two separate updateComponents batches
	content := `{"type":"createSurface","surfaceId":"s1","title":"Test","width":400,"height":300}
{"type":"updateComponents","surfaceId":"s1","components":[{"componentId":"root","type":"Column","props":{},"children":["a"]}]}
{"type":"updateComponents","surfaceId":"s1","components":[{"componentId":"a","type":"Text","props":{"content":"hello"}}]}
`
	if err := os.WriteFile(jsonlPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	dc, err := ExtractComponent(jsonlPath, "Multi")
	if err != nil {
		t.Fatalf("ExtractComponent failed: %v", err)
	}

	if len(dc.Components) != 2 {
		t.Fatalf("got %d components, want 2", len(dc.Components))
	}

	// root should be _root
	var root map[string]interface{}
	json.Unmarshal(dc.Components[0], &root)
	if root["componentId"] != "_root" {
		t.Errorf("root componentId = %q, want _root", root["componentId"])
	}
}

func TestExtractComponentRealSampleApp(t *testing.T) {
	// Test against the real calculator sample app JSONL
	jsonlPath := filepath.Join("..", "sample_apps", "calculator", "prompt.jsonl")
	if _, err := os.Stat(jsonlPath); os.IsNotExist(err) {
		t.Skip("sample_apps/calculator/prompt.jsonl not found")
	}

	dc, err := ExtractComponent(jsonlPath, "Calculator")
	if err != nil {
		t.Fatalf("ExtractComponent failed: %v", err)
	}

	if dc.Name != "Calculator" {
		t.Errorf("name = %q, want Calculator", dc.Name)
	}
	if len(dc.Components) == 0 {
		t.Fatal("expected components, got 0")
	}

	// Verify root is _root
	var root map[string]interface{}
	if err := json.Unmarshal(dc.Components[0], &root); err != nil {
		t.Fatal(err)
	}
	if root["componentId"] != "_root" {
		t.Errorf("root componentId = %q, want _root", root["componentId"])
	}
	t.Logf("extracted %d components from calculator sample", len(dc.Components))
}

func TestExtractComponentNoComponents(t *testing.T) {
	dir := t.TempDir()
	jsonlPath := filepath.Join(dir, "empty.jsonl")

	content := `{"type":"createSurface","surfaceId":"s1","title":"Empty","width":400,"height":300}
`
	if err := os.WriteFile(jsonlPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ExtractComponent(jsonlPath, "Empty")
	if err == nil {
		t.Fatal("expected error for no components, got nil")
	}
}
