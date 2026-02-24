package protocol

// EventAction is what happens when a user interacts with a component.
type EventAction struct {
	Action *Action `json:"action,omitempty"`
}

// Action describes a server-side or client-side action.
type Action struct {
	Type     string                 `json:"type"`               // "serverAction", "clientAction"
	Name     string                 `json:"name,omitempty"`     // action name
	Params   map[string]interface{} `json:"params,omitempty"`   // parameters
	DataRefs []string               `json:"dataRefs,omitempty"` // data model paths to include
}

// FunctionCall represents a built-in function evaluation.
type FunctionCall struct {
	Name string        `json:"name"`
	Args []interface{} `json:"args"`
}
