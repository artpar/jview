package engine

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Validator runs validation checks on field values.
type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

// ValidationRule defines a single validation constraint.
type ValidationRule struct {
	Type    string `json:"type"`
	Value   interface{} `json:"value,omitempty"`
	Message string `json:"message,omitempty"`
}

// Validate checks a value against parsed validation rules.
// Returns a list of error messages (empty if valid).
func (v *Validator) Validate(value string, rulesJSON json.RawMessage) []string {
	if len(rulesJSON) == 0 {
		return nil
	}

	var rules []ValidationRule
	if err := json.Unmarshal(rulesJSON, &rules); err != nil {
		return nil
	}

	var errors []string
	for _, rule := range rules {
		if msg := v.checkRule(value, rule); msg != "" {
			errors = append(errors, msg)
		}
	}
	return errors
}

func (v *Validator) checkRule(value string, rule ValidationRule) string {
	switch rule.Type {
	case "required":
		if strings.TrimSpace(value) == "" {
			return ruleMessage(rule, "This field is required")
		}

	case "minLength":
		minLen := toInt(rule.Value)
		if len(value) < minLen {
			return ruleMessage(rule, fmt.Sprintf("Must be at least %d characters", minLen))
		}

	case "maxLength":
		maxLen := toInt(rule.Value)
		if len(value) > maxLen {
			return ruleMessage(rule, fmt.Sprintf("Must be at most %d characters", maxLen))
		}

	case "pattern":
		pattern, ok := rule.Value.(string)
		if !ok {
			return ""
		}
		matched, err := regexp.MatchString(pattern, value)
		if err != nil || !matched {
			return ruleMessage(rule, "Invalid format")
		}

	case "email":
		if !isValidEmail(value) {
			return ruleMessage(rule, "Invalid email address")
		}
	}
	return ""
}

func ruleMessage(rule ValidationRule, defaultMsg string) string {
	if rule.Message != "" {
		return rule.Message
	}
	return defaultMsg
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	case string:
		var n int
		fmt.Sscanf(val, "%d", &n)
		return n
	default:
		return 0
	}
}

func isValidEmail(s string) bool {
	if s == "" {
		return false
	}
	at := strings.LastIndex(s, "@")
	if at < 1 {
		return false
	}
	domain := s[at+1:]
	dot := strings.LastIndex(domain, ".")
	return dot > 0 && dot < len(domain)-1
}
