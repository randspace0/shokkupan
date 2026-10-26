package dto

import (
	"strings"
	"testing"
)

func TestAppNameLength(t *testing.T) {
	v := NewRequestValidator()
	type req struct {
		Name string `validate:"appName"`
	}
	for name, ok := range map[string]bool{
		"ordering-system-frontend-production": true, // 35 chars
		strings.Repeat("a", 63):               true,
		strings.Repeat("a", 64):               false,
		"":                                    false,
	} {
		err := v.validator.Struct(req{Name: name})
		if (err == nil) != ok {
			t.Errorf("name %q (len %d): valid=%v, want %v", name, len(name), err == nil, ok)
		}
	}
}
