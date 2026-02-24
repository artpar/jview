package protocol

import "encoding/json"

// DynamicValue can be a literal or a data model path reference.
// Literal: "hello" or 42 or true
// Path reference: {"path": "/some/pointer"}
// Function call: {"functionCall": {...}}
type DynamicString struct {
	Literal      string
	Path         string
	FunctionCall *FunctionCall
	IsPath       bool
	IsFunc       bool
}

func (d *DynamicString) UnmarshalJSON(data []byte) error {
	// Try literal string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		d.Literal = s
		return nil
	}

	// Try object with "path" or "functionCall"
	var obj struct {
		Path         string        `json:"path"`
		FunctionCall *FunctionCall `json:"functionCall"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if obj.Path != "" {
		d.Path = obj.Path
		d.IsPath = true
	}
	if obj.FunctionCall != nil {
		d.FunctionCall = obj.FunctionCall
		d.IsFunc = true
	}
	return nil
}

type DynamicNumber struct {
	Literal float64
	Path    string
	IsPath  bool
}

func (d *DynamicNumber) UnmarshalJSON(data []byte) error {
	var n float64
	if err := json.Unmarshal(data, &n); err == nil {
		d.Literal = n
		return nil
	}

	var obj struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if obj.Path != "" {
		d.Path = obj.Path
		d.IsPath = true
	}
	return nil
}

type DynamicBoolean struct {
	Literal bool
	Path    string
	IsPath  bool
}

func (d *DynamicBoolean) UnmarshalJSON(data []byte) error {
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		d.Literal = b
		return nil
	}

	var obj struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if obj.Path != "" {
		d.Path = obj.Path
		d.IsPath = true
	}
	return nil
}

type DynamicStringList struct {
	Literal []string
	Path    string
	IsPath  bool
}

func (d *DynamicStringList) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		d.Literal = arr
		return nil
	}

	var obj struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if obj.Path != "" {
		d.Path = obj.Path
		d.IsPath = true
	}
	return nil
}
