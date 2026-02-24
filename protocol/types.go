package protocol

import "encoding/json"

// MessageType identifies A2UI message kinds.
type MessageType string

const (
	MsgCreateSurface    MessageType = "createSurface"
	MsgDeleteSurface    MessageType = "deleteSurface"
	MsgUpdateComponents MessageType = "updateComponents"
	MsgUpdateDataModel  MessageType = "updateDataModel"
	MsgSetTheme         MessageType = "setTheme"
)

// Envelope wraps every A2UI JSONL line.
type Envelope struct {
	Type      MessageType     `json:"type"`
	SurfaceID string          `json:"surfaceId,omitempty"`
	Payload   json.RawMessage `json:"-"`
}

func (e *Envelope) UnmarshalJSON(data []byte) error {
	type alias Envelope
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*e = Envelope(a)
	e.Payload = data
	return nil
}

// CreateSurface tells the renderer to open a new surface (window).
type CreateSurface struct {
	Type      MessageType `json:"type"`
	SurfaceID string      `json:"surfaceId"`
	Title     string      `json:"title"`
	Width     int         `json:"width,omitempty"`
	Height    int         `json:"height,omitempty"`
}

// DeleteSurface removes a surface.
type DeleteSurface struct {
	Type      MessageType `json:"type"`
	SurfaceID string      `json:"surfaceId"`
}

// UpdateComponents sends a batch of component definitions.
type UpdateComponents struct {
	Type       MessageType `json:"type"`
	SurfaceID  string      `json:"surfaceId"`
	Components []Component `json:"components"`
}

// DataModelOp is a single data model operation.
type DataModelOp struct {
	Op    string      `json:"op"`    // "add", "replace", "remove"
	Path  string      `json:"path"`  // JSON Pointer
	Value interface{} `json:"value,omitempty"`
}

// UpdateDataModel applies operations to a surface's data model.
type UpdateDataModel struct {
	Type      MessageType   `json:"type"`
	SurfaceID string        `json:"surfaceId"`
	Ops       []DataModelOp `json:"ops"`
}

// SetTheme changes the visual theme.
type SetTheme struct {
	Type      MessageType `json:"type"`
	SurfaceID string      `json:"surfaceId"`
	Theme     string      `json:"theme"` // "light", "dark", "system"
}
