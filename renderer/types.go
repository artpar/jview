package renderer

import "jview/protocol"

// ViewHandle is an opaque pointer to a native view.
type ViewHandle uintptr

// CallbackID identifies a registered callback in the platform layer.
type CallbackID uint64

// RenderNode holds the resolved properties for a single component,
// ready for the platform renderer to create or update a native view.
type RenderNode struct {
	ComponentID string
	Type        protocol.ComponentType
	Props       ResolvedProps
	ChildIDs    []string
	Callbacks   map[string]CallbackID // eventType → CallbackID
}

// ResolvedProps contains all resolved (concrete) property values.
// Dynamic values have been evaluated against the data model.
type ResolvedProps struct {
	// Text
	Content string
	Variant string

	// Layout
	Justify string
	Align   string
	Gap     int
	Padding int

	// Card
	Title       string
	Subtitle    string
	Collapsible bool
	Collapsed   bool

	// Button
	Label    string
	Style    string
	Disabled bool

	// TextField
	Placeholder string
	Value       string
	InputType   string
	ReadOnly    bool
	DataBinding string

	// CheckBox
	Checked bool

	// Slider
	Min         float64
	Max         float64
	Step        float64
	SliderValue float64

	// Image
	Src    string
	Alt    string
	Width  int
	Height int

	// Icon
	Name string
	Size int
}

// RenderOp is a single rendering operation to be dispatched to the main thread.
type RenderOp struct {
	Kind        RenderOpKind
	Node        *RenderNode
	Handle      ViewHandle
	ParentID    string
	CallbackID  CallbackID
}

type RenderOpKind int

const (
	OpCreateView RenderOpKind = iota
	OpUpdateView
	OpSetChildren
	OpRemoveView
	OpCreateWindow
)

// WindowSpec describes a new window to create.
type WindowSpec struct {
	SurfaceID string
	Title     string
	Width     int
	Height    int
}
