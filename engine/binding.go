package engine

import "strings"

// BindingTracker maintains a reverse index: data model path → set of component IDs.
type BindingTracker struct {
	pathToComponents map[string]map[string]bool
}

func NewBindingTracker() *BindingTracker {
	return &BindingTracker{
		pathToComponents: make(map[string]map[string]bool),
	}
}

// Register binds a component to a data model path.
func (bt *BindingTracker) Register(path string, componentID string) {
	set, ok := bt.pathToComponents[path]
	if !ok {
		set = make(map[string]bool)
		bt.pathToComponents[path] = set
	}
	set[componentID] = true
}

// Unregister removes all bindings for a component.
func (bt *BindingTracker) Unregister(componentID string) {
	for path, set := range bt.pathToComponents {
		delete(set, componentID)
		if len(set) == 0 {
			delete(bt.pathToComponents, path)
		}
	}
}

// Affected returns all component IDs bound to any of the changed paths.
// A component is affected if its bound path is a prefix of or equal to a changed path,
// or if a changed path is a prefix of its bound path.
func (bt *BindingTracker) Affected(changedPaths []string) []string {
	seen := make(map[string]bool)
	for _, changed := range changedPaths {
		for boundPath, components := range bt.pathToComponents {
			if pathOverlaps(changed, boundPath) {
				for id := range components {
					seen[id] = true
				}
			}
		}
	}

	result := make([]string, 0, len(seen))
	for id := range seen {
		result = append(result, id)
	}
	return result
}

// pathOverlaps returns true if a and b are equal, or one is a prefix of the other
// at a path segment boundary.
func pathOverlaps(a, b string) bool {
	if a == b {
		return true
	}
	// Normalize: ensure leading /
	if !strings.HasPrefix(a, "/") {
		a = "/" + a
	}
	if !strings.HasPrefix(b, "/") {
		b = "/" + b
	}

	if strings.HasPrefix(a, b) {
		rest := a[len(b):]
		return rest == "" || rest[0] == '/'
	}
	if strings.HasPrefix(b, a) {
		rest := b[len(a):]
		return rest == "" || rest[0] == '/'
	}
	return false
}
