package engine

import (
	"encoding/json"
	"testing"
)

func TestValidateRequired(t *testing.T) {
	v := NewValidator()
	rules := json.RawMessage(`[{"type":"required"}]`)

	errors := v.Validate("", rules)
	if len(errors) != 1 {
		t.Fatalf("errors = %d, want 1", len(errors))
	}
	if errors[0] != "This field is required" {
		t.Errorf("error = %q", errors[0])
	}

	errors = v.Validate("hello", rules)
	if len(errors) != 0 {
		t.Errorf("errors = %d, want 0", len(errors))
	}
}

func TestValidateRequiredCustomMessage(t *testing.T) {
	v := NewValidator()
	rules := json.RawMessage(`[{"type":"required","message":"Please fill this in"}]`)

	errors := v.Validate("", rules)
	if len(errors) != 1 || errors[0] != "Please fill this in" {
		t.Errorf("error = %v", errors)
	}
}

func TestValidateMinLength(t *testing.T) {
	v := NewValidator()
	rules := json.RawMessage(`[{"type":"minLength","value":3}]`)

	errors := v.Validate("ab", rules)
	if len(errors) != 1 {
		t.Fatalf("errors = %d, want 1", len(errors))
	}

	errors = v.Validate("abc", rules)
	if len(errors) != 0 {
		t.Errorf("errors = %d, want 0", len(errors))
	}
}

func TestValidateMaxLength(t *testing.T) {
	v := NewValidator()
	rules := json.RawMessage(`[{"type":"maxLength","value":5}]`)

	errors := v.Validate("123456", rules)
	if len(errors) != 1 {
		t.Fatalf("errors = %d, want 1", len(errors))
	}

	errors = v.Validate("12345", rules)
	if len(errors) != 0 {
		t.Errorf("errors = %d, want 0", len(errors))
	}
}

func TestValidatePattern(t *testing.T) {
	v := NewValidator()
	rules := json.RawMessage(`[{"type":"pattern","value":"^[0-9]+$"}]`)

	errors := v.Validate("abc", rules)
	if len(errors) != 1 {
		t.Fatalf("errors = %d, want 1", len(errors))
	}

	errors = v.Validate("123", rules)
	if len(errors) != 0 {
		t.Errorf("errors = %d, want 0", len(errors))
	}
}

func TestValidateEmail(t *testing.T) {
	v := NewValidator()
	rules := json.RawMessage(`[{"type":"email"}]`)

	errors := v.Validate("notanemail", rules)
	if len(errors) != 1 {
		t.Fatalf("errors = %d, want 1", len(errors))
	}

	errors = v.Validate("user@example.com", rules)
	if len(errors) != 0 {
		t.Errorf("errors = %d, want 0", len(errors))
	}
}

func TestValidateMultipleRules(t *testing.T) {
	v := NewValidator()
	rules := json.RawMessage(`[{"type":"required"},{"type":"minLength","value":3}]`)

	errors := v.Validate("", rules)
	if len(errors) != 2 {
		t.Fatalf("errors = %d, want 2", len(errors))
	}
}

func TestValidateClearErrors(t *testing.T) {
	v := NewValidator()
	rules := json.RawMessage(`[{"type":"required"}]`)

	errors := v.Validate("", rules)
	if len(errors) != 1 {
		t.Fatalf("errors = %d, want 1", len(errors))
	}

	// Now valid
	errors = v.Validate("hello", rules)
	if len(errors) != 0 {
		t.Errorf("errors = %d, want 0 after valid input", len(errors))
	}
}

func TestValidateNoRules(t *testing.T) {
	v := NewValidator()

	errors := v.Validate("anything", nil)
	if errors != nil {
		t.Errorf("errors = %v, want nil", errors)
	}
}
