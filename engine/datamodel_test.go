package engine

import "testing"

func TestDataModelSetGet(t *testing.T) {
	dm := NewDataModel()

	dm.Set("/name", "Alice")
	val, ok := dm.Get("/name")
	if !ok || val != "Alice" {
		t.Errorf("Get(/name) = %v, %v; want Alice, true", val, ok)
	}
}

func TestDataModelNestedPath(t *testing.T) {
	dm := NewDataModel()

	dm.Set("/user/name", "Bob")
	dm.Set("/user/age", 30.0)

	val, ok := dm.Get("/user/name")
	if !ok || val != "Bob" {
		t.Errorf("Get(/user/name) = %v, %v; want Bob, true", val, ok)
	}

	val, ok = dm.Get("/user/age")
	if !ok || val != 30.0 {
		t.Errorf("Get(/user/age) = %v, %v; want 30, true", val, ok)
	}
}

func TestDataModelDelete(t *testing.T) {
	dm := NewDataModel()

	dm.Set("/x", "val")
	dm.Delete("/x")

	_, ok := dm.Get("/x")
	if ok {
		t.Error("expected /x to be deleted")
	}
}

func TestDataModelReturnsPaths(t *testing.T) {
	dm := NewDataModel()

	changed, err := dm.Set("/a/b/c", "deep")
	if err != nil {
		t.Fatal(err)
	}
	// Should include at least the target path
	found := false
	for _, p := range changed {
		if p == "/a/b/c" {
			found = true
		}
	}
	if !found {
		t.Errorf("changed paths %v don't include /a/b/c", changed)
	}
}

func TestDataModelGetRoot(t *testing.T) {
	dm := NewDataModel()
	dm.Set("/x", 1.0)

	root, ok := dm.Get("")
	if !ok {
		t.Fatal("expected root")
	}
	m, ok := root.(map[string]interface{})
	if !ok {
		t.Fatalf("root is %T, want map", root)
	}
	if m["x"] != 1.0 {
		t.Errorf("root[x] = %v, want 1", m["x"])
	}
}

func TestDataModelMissing(t *testing.T) {
	dm := NewDataModel()
	_, ok := dm.Get("/nonexistent")
	if ok {
		t.Error("expected missing path to return false")
	}
}

func TestDataModelDeleteArrayElement(t *testing.T) {
	dm := NewDataModel()
	dm.Set("/items", []interface{}{"a", "b", "c"})

	dm.Delete("/items/1") // delete "b"

	val, ok := dm.Get("/items")
	if !ok {
		t.Fatal("expected /items to exist")
	}
	arr := val.([]interface{})
	if len(arr) != 2 {
		t.Fatalf("len = %d, want 2", len(arr))
	}
	if arr[0] != "a" || arr[1] != "c" {
		t.Errorf("items = %v, want [a c]", arr)
	}
}

func TestDataModelDeleteArrayFirst(t *testing.T) {
	dm := NewDataModel()
	dm.Set("/items", []interface{}{"x", "y", "z"})

	dm.Delete("/items/0")

	val, _ := dm.Get("/items")
	arr := val.([]interface{})
	if len(arr) != 2 {
		t.Fatalf("len = %d, want 2", len(arr))
	}
	if arr[0] != "y" || arr[1] != "z" {
		t.Errorf("items = %v, want [y z]", arr)
	}
}

func TestDataModelDeleteArrayLast(t *testing.T) {
	dm := NewDataModel()
	dm.Set("/items", []interface{}{1.0, 2.0, 3.0})

	dm.Delete("/items/2")

	val, _ := dm.Get("/items")
	arr := val.([]interface{})
	if len(arr) != 2 {
		t.Fatalf("len = %d, want 2", len(arr))
	}
	if arr[0] != 1.0 || arr[1] != 2.0 {
		t.Errorf("items = %v, want [1 2]", arr)
	}
}

func TestDataModelDeleteArrayOutOfBounds(t *testing.T) {
	dm := NewDataModel()
	dm.Set("/items", []interface{}{"a", "b"})

	changed, err := dm.Delete("/items/5")
	if err != nil {
		t.Fatal(err)
	}
	// Out of bounds: nothing deleted, no changed paths
	if changed != nil {
		t.Errorf("expected nil changed, got %v", changed)
	}

	val, _ := dm.Get("/items")
	arr := val.([]interface{})
	if len(arr) != 2 {
		t.Errorf("len = %d, want 2 (unchanged)", len(arr))
	}
}

func TestDataModelDeleteNestedInArray(t *testing.T) {
	dm := NewDataModel()
	dm.Set("/users", []interface{}{
		map[string]interface{}{"name": "Alice", "age": 30.0},
		map[string]interface{}{"name": "Bob", "age": 25.0},
	})

	dm.Delete("/users/0/age")

	val, ok := dm.Get("/users/0")
	if !ok {
		t.Fatal("expected /users/0 to exist")
	}
	obj := val.(map[string]interface{})
	if _, exists := obj["age"]; exists {
		t.Error("expected age to be deleted")
	}
	if obj["name"] != "Alice" {
		t.Errorf("name = %v, want Alice", obj["name"])
	}
}

func TestDataModelDeepMissingPath(t *testing.T) {
	dm := NewDataModel()
	val, ok := dm.Get("/a/b/c/d/e")
	if ok {
		t.Errorf("expected false, got true with val=%v", val)
	}
}

func TestDataModelSetNegativeArrayIndex(t *testing.T) {
	dm := NewDataModel()
	dm.Set("/arr", []interface{}{"a", "b"})

	_, err := dm.Set("/arr/-1", "x")
	if err == nil {
		t.Error("expected error for negative array index")
	}
}

func TestDataModelSetStringIndexOnArray(t *testing.T) {
	dm := NewDataModel()
	dm.Set("/arr", []interface{}{"a", "b"})

	_, err := dm.Set("/arr/notanumber", "x")
	if err == nil {
		t.Error("expected error for non-numeric array index")
	}
}

func TestDataModelSetNonMapRoot(t *testing.T) {
	dm := NewDataModel()
	// Replace root with a string
	dm.Set("", "just a string")

	// Now try to set a nested path — should create intermediate maps
	_, err := dm.Set("/key", "value")
	// This may or may not error depending on implementation, but must not panic
	if err != nil {
		// Acceptable: root was a string, can't index into it
		return
	}
	val, ok := dm.Get("/key")
	if !ok || val != "value" {
		t.Errorf("Get(/key) = %v, %v", val, ok)
	}
}
