package engine

import "jview/protocol"

// Tree maintains the component hierarchy for a surface.
// Components are stored in a flat map; the tree structure is derived from ParentID + Children.
type Tree struct {
	components map[string]*protocol.Component
	rootIDs    []string // top-level components (no parent)
}

func NewTree() *Tree {
	return &Tree{
		components: make(map[string]*protocol.Component),
	}
}

// Update adds or replaces components. Returns the IDs of components that changed.
func (t *Tree) Update(comps []protocol.Component) []string {
	var changed []string
	for i := range comps {
		comp := &comps[i]
		t.components[comp.ComponentID] = comp
		changed = append(changed, comp.ComponentID)
	}
	t.rebuildRoots()
	return changed
}

// Get returns a component by ID.
func (t *Tree) Get(id string) (*protocol.Component, bool) {
	c, ok := t.components[id]
	return c, ok
}

// Children returns the child component IDs for a given component.
func (t *Tree) Children(id string) []string {
	comp, ok := t.components[id]
	if !ok || comp.Children == nil {
		return nil
	}
	return comp.Children.Static
}

// RootIDs returns the top-level component IDs.
func (t *Tree) RootIDs() []string {
	return t.rootIDs
}

// All returns all component IDs.
func (t *Tree) All() []string {
	ids := make([]string, 0, len(t.components))
	for id := range t.components {
		ids = append(ids, id)
	}
	return ids
}

// rebuildRoots recalculates which components are root-level.
// A root component either has no parentId or its parentId references a non-existent component.
func (t *Tree) rebuildRoots() {
	t.rootIDs = nil
	// Build set of IDs that are children of something
	childSet := make(map[string]bool)
	for _, comp := range t.components {
		if comp.Children != nil && comp.Children.Static != nil {
			for _, childID := range comp.Children.Static {
				childSet[childID] = true
			}
		}
	}

	for id := range t.components {
		if !childSet[id] {
			t.rootIDs = append(t.rootIDs, id)
		}
	}
}
