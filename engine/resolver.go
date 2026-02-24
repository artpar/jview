package engine

import (
	"fmt"
	"jview/protocol"
	"jview/renderer"
)

// Resolver evaluates dynamic values against a data model and tracks bindings.
type Resolver struct {
	dm      *DataModel
	tracker *BindingTracker
}

func NewResolver(dm *DataModel, tracker *BindingTracker) *Resolver {
	return &Resolver{dm: dm, tracker: tracker}
}

// Resolve takes a protocol Component and returns a RenderNode with all
// dynamic values resolved to concrete values.
func (r *Resolver) Resolve(comp *protocol.Component) *renderer.RenderNode {
	node := &renderer.RenderNode{
		ComponentID: comp.ComponentID,
		Type:        comp.Type,
	}

	// Resolve children
	if comp.Children != nil && comp.Children.Static != nil {
		node.ChildIDs = comp.Children.Static
	}

	// Resolve props
	p := &node.Props
	cp := &comp.Props

	switch comp.Type {
	case protocol.CompText:
		p.Content = r.resolveString(comp.ComponentID, cp.Content)
		p.Variant = cp.Variant
		if p.Variant == "" {
			p.Variant = "body"
		}

	case protocol.CompRow:
		p.Justify = cp.Justify
		p.Align = cp.Align
		p.Gap = cp.Gap
		if p.Gap == 0 {
			p.Gap = 8
		}
		p.Padding = cp.Padding

	case protocol.CompColumn:
		p.Justify = cp.Justify
		p.Align = cp.Align
		p.Gap = cp.Gap
		if p.Gap == 0 {
			p.Gap = 8
		}
		p.Padding = cp.Padding

	case protocol.CompCard:
		p.Title = r.resolveString(comp.ComponentID, cp.Title)
		p.Subtitle = r.resolveString(comp.ComponentID, cp.Subtitle)
		p.Collapsible = r.resolveBool(comp.ComponentID, cp.Collapsible)
		p.Collapsed = r.resolveBool(comp.ComponentID, cp.Collapsed)
		p.Padding = cp.Padding

	case protocol.CompButton:
		p.Label = r.resolveString(comp.ComponentID, cp.Label)
		p.Style = cp.Style
		if p.Style == "" {
			p.Style = "secondary"
		}
		p.Disabled = r.resolveBool(comp.ComponentID, cp.Disabled)

	case protocol.CompTextField:
		p.Placeholder = r.resolveString(comp.ComponentID, cp.Placeholder)
		p.Value = r.resolveString(comp.ComponentID, cp.Value)
		p.InputType = cp.InputType
		if p.InputType == "" {
			p.InputType = "shortText"
		}
		p.ReadOnly = r.resolveBool(comp.ComponentID, cp.ReadOnly)
		p.DataBinding = cp.DataBinding

	case protocol.CompCheckBox:
		p.Label = r.resolveString(comp.ComponentID, cp.Label)
		p.Checked = r.resolveBool(comp.ComponentID, cp.Checked)
		p.DataBinding = cp.DataBinding
	}

	return node
}

func (r *Resolver) resolveString(componentID string, dv *protocol.DynamicString) string {
	if dv == nil {
		return ""
	}
	if dv.IsPath {
		r.tracker.Register(dv.Path, componentID)
		val, ok := r.dm.Get(dv.Path)
		if !ok {
			return ""
		}
		return fmt.Sprintf("%v", val)
	}
	return dv.Literal
}

func (r *Resolver) resolveNumber(componentID string, dv *protocol.DynamicNumber) float64 {
	if dv == nil {
		return 0
	}
	if dv.IsPath {
		r.tracker.Register(dv.Path, componentID)
		val, ok := r.dm.Get(dv.Path)
		if !ok {
			return 0
		}
		if f, ok := val.(float64); ok {
			return f
		}
		return 0
	}
	return dv.Literal
}

func (r *Resolver) resolveBool(componentID string, dv *protocol.DynamicBoolean) bool {
	if dv == nil {
		return false
	}
	if dv.IsPath {
		r.tracker.Register(dv.Path, componentID)
		val, ok := r.dm.Get(dv.Path)
		if !ok {
			return false
		}
		if b, ok := val.(bool); ok {
			return b
		}
		return false
	}
	return dv.Literal
}
