package protocol

// EventAction is what happens when a user interacts with a component.
// It supports optional filtering (for keyboard/mouse events), data model writes,
// throttle/debounce for high-frequency events, and event consumption.
type EventAction struct {
	Action         *Action      `json:"action,omitempty"`
	Filter         *EventFilter `json:"filter,omitempty"`
	DataPath       string       `json:"dataPath,omitempty"`       // write to DataModel on fire
	DataValue      interface{}  `json:"dataValue,omitempty"`      // value to write (nil = native event data)
	Throttle       int          `json:"throttle,omitempty"`       // ms, max fire rate
	Debounce       int          `json:"debounce,omitempty"`       // ms, quiet period before fire
	PreventDefault bool         `json:"preventDefault,omitempty"` // consume event, don't propagate
}

// EventFilter restricts when an event handler fires based on event properties.
// Used primarily for keyboard (key + modifiers) and mouse (button) filtering.
type EventFilter struct {
	Key       string   `json:"key,omitempty"`       // "Enter", "Escape", "a", "ArrowDown", etc.
	Modifiers []string `json:"modifiers,omitempty"` // "cmd", "shift", "option", "ctrl"
	Button    int      `json:"button,omitempty"`    // mouse button: 0=left, 1=right, 2=middle
}

// Action describes an interaction outcome — either a server-bound event, a client-side function call,
// or a standard AppKit action routed through the responder chain.
type Action struct {
	Event          *EventDef       `json:"event,omitempty"`
	FunctionCall   *ActionFuncCall `json:"functionCall,omitempty"`
	StandardAction string          `json:"standardAction,omitempty"`
}

// EventDef is a server-bound event with optional data references.
type EventDef struct {
	Name      string                 `json:"name"`
	Context   map[string]interface{} `json:"context,omitempty"`
	DataRefs  []string               `json:"dataRefs,omitempty"`
	ProcessID string                 `json:"processId,omitempty"`
}

// ActionFuncCall is a client-side function call (e.g., updateDataModel).
type ActionFuncCall struct {
	Call string      `json:"call"`
	Args interface{} `json:"args,omitempty"`
}

// FunctionCall represents a built-in function evaluation (used in dynamic values).
type FunctionCall struct {
	Name string        `json:"name"`
	Args []interface{} `json:"args"`
}
