package types

import (
	"encoding/json"
	"testing"
)

func TestFlexibleInt_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected FlexibleInt
	}{
		{
			name:     "integer value",
			input:    `123`,
			expected: FlexibleInt{Value: 123, IsSet: true},
		},
		{
			name:     "string integer",
			input:    `"456"`,
			expected: FlexibleInt{Value: 456, IsSet: true},
		},
		{
			name:     "empty string",
			input:    `""`,
			expected: FlexibleInt{Value: 0, IsSet: false},
		},
		{
			name:     "zero integer",
			input:    `0`,
			expected: FlexibleInt{Value: 0, IsSet: true},
		},
		{
			name:     "negative integer",
			input:    `-123`,
			expected: FlexibleInt{Value: -123, IsSet: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f FlexibleInt
			err := json.Unmarshal([]byte(tt.input), &f)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if f.Value != tt.expected.Value {
				t.Errorf("Expected value %d, got %d", tt.expected.Value, f.Value)
			}

			if f.IsSet != tt.expected.IsSet {
				t.Errorf("Expected IsSet %v, got %v", tt.expected.IsSet, f.IsSet)
			}
		})
	}
}

func TestFlexibleInt_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    FlexibleInt
		expected string
	}{
		{
			name:     "set value",
			input:    FlexibleInt{Value: 123, IsSet: true},
			expected: `123`,
		},
		{
			name:     "unset value",
			input:    FlexibleInt{Value: 0, IsSet: false},
			expected: `""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result))
			}
		})
	}
}

func TestFlexibleInt_Methods(t *testing.T) {
	// Test set value
	f1 := FlexibleInt{Value: 123, IsSet: true}
	if f1.Int64() != 123 {
		t.Errorf("Expected Int64() to return 123, got %d", f1.Int64())
	}
	if f1.String() != "123" {
		t.Errorf("Expected String() to return '123', got '%s'", f1.String())
	}

	// Test unset value
	f2 := FlexibleInt{Value: 0, IsSet: false}
	if f2.Int64() != 0 {
		t.Errorf("Expected Int64() to return 0 for unset, got %d", f2.Int64())
	}
	if f2.String() != "" {
		t.Errorf("Expected String() to return empty string for unset, got '%s'", f2.String())
	}
}
