package engine

import (
	"encoding/json"
	"testing"
)

func TestSubstituteParamsSimple(t *testing.T) {
	body := map[string]interface{}{
		"functionCall": map[string]interface{}{
			"name": "concat",
			"args": []interface{}{
				map[string]interface{}{"path": "/display"},
				map[string]interface{}{"param": "digit"},
			},
		},
	}
	args := map[string]interface{}{"digit": "7"}
	result := substituteParams(body, args)
	rm := result.(map[string]interface{})
	fc := rm["functionCall"].(map[string]interface{})
	fcArgs := fc["args"].([]interface{})
	if fcArgs[1] != "7" {
		t.Errorf("param not substituted: got %v", fcArgs[1])
	}
	// Path should be untouched
	pathObj := fcArgs[0].(map[string]interface{})
	if pathObj["path"] != "/display" {
		t.Errorf("path mutated: got %v", pathObj["path"])
	}
}

func TestSubstituteParamsNested(t *testing.T) {
	body := map[string]interface{}{
		"functionCall": map[string]interface{}{
			"name": "if",
			"args": []interface{}{
				true,
				map[string]interface{}{"param": "a"},
				map[string]interface{}{
					"functionCall": map[string]interface{}{
						"name": "concat",
						"args": []interface{}{
							"x",
							map[string]interface{}{"param": "b"},
						},
					},
				},
			},
		},
	}
	args := map[string]interface{}{"a": "hello", "b": "world"}
	result := substituteParams(body, args)
	rm := result.(map[string]interface{})
	fc := rm["functionCall"].(map[string]interface{})
	fcArgs := fc["args"].([]interface{})
	if fcArgs[1] != "hello" {
		t.Errorf("param a not substituted: got %v", fcArgs[1])
	}
	inner := fcArgs[2].(map[string]interface{})["functionCall"].(map[string]interface{})
	innerArgs := inner["args"].([]interface{})
	if innerArgs[1] != "world" {
		t.Errorf("param b not substituted: got %v", innerArgs[1])
	}
}

func TestSubstituteParamsMissing(t *testing.T) {
	body := map[string]interface{}{
		"param": "missing",
	}
	args := map[string]interface{}{"other": "val"}
	result := substituteParams(body, args)
	rm := result.(map[string]interface{})
	if rm["param"] != "missing" {
		t.Errorf("missing param should be untouched: got %v", rm)
	}
}

func TestRewriteComponentIDs(t *testing.T) {
	trees := []map[string]interface{}{
		{"componentId": "_root", "type": "Button", "children": []interface{}{"_label"}},
		{"componentId": "_label", "type": "Text", "parentId": "_root"},
	}
	idMap := map[string]string{
		"_root":  "btn7",
		"_label": "btn7__label",
	}
	rewriteComponentIDs(trees, idMap)

	if trees[0]["componentId"] != "btn7" {
		t.Errorf("root not rewritten: got %v", trees[0]["componentId"])
	}
	if trees[1]["componentId"] != "btn7__label" {
		t.Errorf("label not rewritten: got %v", trees[1]["componentId"])
	}
	if trees[1]["parentId"] != "btn7" {
		t.Errorf("parentId not rewritten: got %v", trees[1]["parentId"])
	}
	children := trees[0]["children"].([]interface{})
	if children[0] != "btn7__label" {
		t.Errorf("child ref not rewritten: got %v", children[0])
	}
}

func TestRewriteComponentIDsStaticFormat(t *testing.T) {
	trees := []map[string]interface{}{
		{
			"componentId": "_root",
			"children":    map[string]interface{}{"static": []interface{}{"_a", "_b"}},
		},
	}
	idMap := map[string]string{"_root": "x", "_a": "x__a", "_b": "x__b"}
	rewriteComponentIDs(trees, idMap)

	children := trees[0]["children"].(map[string]interface{})
	static := children["static"].([]interface{})
	if static[0] != "x__a" || static[1] != "x__b" {
		t.Errorf("static children not rewritten: got %v", static)
	}
}

func TestRewriteScopedPaths(t *testing.T) {
	tree := map[string]interface{}{
		"props": map[string]interface{}{
			"content":     map[string]interface{}{"path": "$/display"},
			"dataBinding": "$/name",
		},
		"children": map[string]interface{}{
			"forEach": "$/items",
		},
	}
	result := rewriteScopedPaths(tree, "/calc1").(map[string]interface{})
	props := result["props"].(map[string]interface{})
	content := props["content"].(map[string]interface{})
	if content["path"] != "/calc1/display" {
		t.Errorf("path not rewritten: got %v", content["path"])
	}
	if props["dataBinding"] != "/calc1/name" {
		t.Errorf("dataBinding not rewritten: got %v", props["dataBinding"])
	}
	children := result["children"].(map[string]interface{})
	if children["forEach"] != "/calc1/items" {
		t.Errorf("forEach not rewritten: got %v", children["forEach"])
	}
}

func TestRewriteScopedPathsNoScope(t *testing.T) {
	tree := map[string]interface{}{
		"path": "/normal/path",
	}
	result := rewriteScopedPaths(tree, "/scope").(map[string]interface{})
	if result["path"] != "/normal/path" {
		t.Errorf("non-scoped path should be untouched: got %v", result["path"])
	}
}

func TestDeepCopyJSON(t *testing.T) {
	original := map[string]interface{}{
		"a": []interface{}{"x", "y"},
		"b": map[string]interface{}{"c": "d"},
	}
	copied := deepCopyJSON(original).(map[string]interface{})

	// Modify copy
	copied["a"].([]interface{})[0] = "MODIFIED"
	copied["b"].(map[string]interface{})["c"] = "MODIFIED"

	// Original should be unchanged
	if original["a"].([]interface{})[0] != "x" {
		t.Error("deep copy failed: original modified")
	}
	if original["b"].(map[string]interface{})["c"] != "d" {
		t.Error("deep copy failed: original nested modified")
	}
}

func TestJsonToMaps(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage(`{"componentId":"_root","type":"Button"}`),
		json.RawMessage(`{"componentId":"_label","type":"Text"}`),
	}
	maps, err := jsonToMaps(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(maps) != 2 {
		t.Fatalf("expected 2 maps, got %d", len(maps))
	}
	if maps[0]["componentId"] != "_root" {
		t.Errorf("first map componentId = %v", maps[0]["componentId"])
	}
}
