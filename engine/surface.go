package engine

import (
	"fmt"
	"jview/protocol"
	"jview/renderer"
	"log"
)

// Surface manages a single A2UI surface: its component tree, data model, and bindings.
type Surface struct {
	id       string
	tree     *Tree
	dm       *DataModel
	tracker  *BindingTracker
	resolver *Resolver
	rend     renderer.Renderer
	dispatch renderer.Dispatcher

	// activeCallbacks tracks registered callbacks: componentID → eventType → CallbackID
	activeCallbacks map[string]map[string]renderer.CallbackID

	// ActionHandler is called when a component triggers a server action.
	ActionHandler func(surfaceID string, action *protocol.Action, eventData map[string]interface{})
}

func NewSurface(id string, rend renderer.Renderer, dispatch renderer.Dispatcher) *Surface {
	dm := NewDataModel()
	tracker := NewBindingTracker()
	return &Surface{
		id:              id,
		tree:            NewTree(),
		dm:              dm,
		tracker:         tracker,
		resolver:        NewResolver(dm, tracker),
		rend:            rend,
		dispatch:        dispatch,
		activeCallbacks: make(map[string]map[string]renderer.CallbackID),
	}
}

// HandleUpdateComponents processes a batch of component definitions.
func (s *Surface) HandleUpdateComponents(msg protocol.UpdateComponents) {
	changed := s.tree.Update(msg.Components)
	s.renderComponents(changed)
}

// HandleUpdateDataModel applies data model operations and re-renders affected components.
func (s *Surface) HandleUpdateDataModel(msg protocol.UpdateDataModel) {
	var allChanged []string
	for _, op := range msg.Ops {
		var changed []string
		var err error
		switch op.Op {
		case "add", "replace":
			changed, err = s.dm.Set(op.Path, op.Value)
		case "remove":
			changed, err = s.dm.Delete(op.Path)
		default:
			log.Printf("surface %s: unknown data op %q", s.id, op.Op)
			continue
		}
		if err != nil {
			log.Printf("surface %s: data op error: %v", s.id, err)
			continue
		}
		allChanged = append(allChanged, changed...)
	}

	if len(allChanged) == 0 {
		return
	}

	affected := s.tracker.Affected(allChanged)
	if len(affected) > 0 {
		s.renderComponents(affected)
	}
}

// renderComponents resolves and dispatches render operations for the given component IDs.
func (s *Surface) renderComponents(componentIDs []string) {
	type renderWork struct {
		node *renderer.RenderNode
		comp *protocol.Component
	}

	// Build work items and sort: children before parents (leaves first).
	// We do a topological sort based on the children references.
	nodeMap := make(map[string]*renderWork, len(componentIDs))
	for _, id := range componentIDs {
		comp, ok := s.tree.Get(id)
		if !ok {
			continue
		}
		node := s.resolver.Resolve(comp)
		s.registerCallbacks(comp, node)
		nodeMap[id] = &renderWork{node: node, comp: comp}
	}

	// Topological order: leaves first, roots last
	var ordered []*renderWork
	visited := make(map[string]bool)
	var visit func(id string)
	visit = func(id string) {
		if visited[id] {
			return
		}
		visited[id] = true
		w, ok := nodeMap[id]
		if !ok {
			return
		}
		// Visit children first
		for _, childID := range w.node.ChildIDs {
			visit(childID)
		}
		ordered = append(ordered, w)
	}
	for _, id := range componentIDs {
		visit(id)
	}

	s.dispatch.RunOnMain(func() {
		// First pass: create/update all views (leaves first)
		for _, w := range ordered {
			handle := s.rend.GetHandle(s.id, w.node.ComponentID)
			if handle == 0 {
				s.rend.CreateView(s.id, w.node)
			} else {
				s.rend.UpdateView(s.id, handle, w.node)
			}
		}

		// Second pass: set children for containers (now all children exist)
		for _, w := range ordered {
			if len(w.node.ChildIDs) == 0 {
				continue
			}
			parentHandle := s.rend.GetHandle(s.id, w.node.ComponentID)
			if parentHandle == 0 {
				continue
			}
			childHandles := make([]renderer.ViewHandle, 0, len(w.node.ChildIDs))
			for _, childID := range w.node.ChildIDs {
				ch := s.rend.GetHandle(s.id, childID)
				if ch != 0 {
					childHandles = append(childHandles, ch)
				}
			}
			if len(childHandles) > 0 {
				s.rend.SetChildren(s.id, parentHandle, childHandles)
			}
		}

		// Set root view(s)
		roots := s.tree.RootIDs()
		if len(roots) == 1 {
			h := s.rend.GetHandle(s.id, roots[0])
			if h != 0 {
				s.rend.SetRootView(s.id, h)
			}
		} else if len(roots) > 1 {
			wrapperNode := &renderer.RenderNode{
				ComponentID: "__root_wrapper__",
				Type:        protocol.CompColumn,
				Props: renderer.ResolvedProps{
					Gap:     8,
					Padding: 16,
				},
			}
			wrapperHandle := s.rend.GetHandle(s.id, "__root_wrapper__")
			if wrapperHandle == 0 {
				wrapperHandle = s.rend.CreateView(s.id, wrapperNode)
			}
			var rootHandles []renderer.ViewHandle
			for _, rid := range roots {
				h := s.rend.GetHandle(s.id, rid)
				if h != 0 {
					rootHandles = append(rootHandles, h)
				}
			}
			s.rend.SetChildren(s.id, wrapperHandle, rootHandles)
			s.rend.SetRootView(s.id, wrapperHandle)
		}
	})
}

func (s *Surface) registerCallbacks(comp *protocol.Component, node *renderer.RenderNode) {
	node.Callbacks = make(map[string]renderer.CallbackID)

	// Unregister any existing callbacks for this component before re-registering
	if old, exists := s.activeCallbacks[comp.ComponentID]; exists {
		for _, cbID := range old {
			s.rend.UnregisterCallback(cbID)
		}
		delete(s.activeCallbacks, comp.ComponentID)
	}

	switch comp.Type {
	case protocol.CompButton:
		if comp.Props.OnClick != nil && comp.Props.OnClick.Action != nil {
			action := comp.Props.OnClick.Action
			cbID := s.rend.RegisterCallback(s.id, comp.ComponentID, "click", func(data string) {
				if s.ActionHandler != nil {
					s.ActionHandler(s.id, action, map[string]interface{}{"data": data})
				} else {
					fmt.Printf("action: %s %s (surface %s, component %s)\n",
						action.Type, action.Name, s.id, comp.ComponentID)
				}
			})
			node.Callbacks["click"] = cbID
			s.trackCallback(comp.ComponentID, "click", cbID)
		}

	case protocol.CompTextField:
		if comp.Props.DataBinding != "" {
			binding := comp.Props.DataBinding
			compID := comp.ComponentID
			cbID := s.rend.RegisterCallback(s.id, comp.ComponentID, "change", func(value string) {
				changed, err := s.dm.Set(binding, value)
				if err != nil {
					log.Printf("surface %s: binding set error: %v", s.id, err)
					return
				}
				affected := s.tracker.Affected(changed)
				var toRender []string
				for _, id := range affected {
					if id != compID {
						toRender = append(toRender, id)
					}
				}
				if len(toRender) > 0 {
					s.renderComponents(toRender)
				}
			})
			node.Callbacks["change"] = cbID
			s.trackCallback(comp.ComponentID, "change", cbID)
		}

	case protocol.CompCheckBox:
		if comp.Props.DataBinding != "" {
			binding := comp.Props.DataBinding
			compID := comp.ComponentID
			cbID := s.rend.RegisterCallback(s.id, comp.ComponentID, "toggle", func(value string) {
				boolVal := value == "true" || value == "1"
				changed, err := s.dm.Set(binding, boolVal)
				if err != nil {
					log.Printf("surface %s: binding set error: %v", s.id, err)
					return
				}
				affected := s.tracker.Affected(changed)
				var toRender []string
				for _, id := range affected {
					if id != compID {
						toRender = append(toRender, id)
					}
				}
				if len(toRender) > 0 {
					s.renderComponents(toRender)
				}
			})
			node.Callbacks["toggle"] = cbID
			s.trackCallback(comp.ComponentID, "toggle", cbID)
		}
	}
}

func (s *Surface) trackCallback(componentID, eventType string, cbID renderer.CallbackID) {
	if s.activeCallbacks[componentID] == nil {
		s.activeCallbacks[componentID] = make(map[string]renderer.CallbackID)
	}
	s.activeCallbacks[componentID][eventType] = cbID
}
