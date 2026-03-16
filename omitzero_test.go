package json

import (
	"bytes"
	"testing"
)

// TestOmitZero tests the omitzero struct tag functionality
func TestOmitZero(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		// Basic primitive fields - omitted when zero
		{
			name: "int zero",
			value: struct {
				X int `json:"x,omitzero"`
			}{X: 0},
			expected: "{}",
		},
		{
			name: "int non-zero",
			value: struct {
				X int `json:"x,omitzero"`
			}{X: 42},
			expected: `{"x":42}`,
		},
		{
			name: "string zero",
			value: struct {
				X string `json:"x,omitzero"`
			}{X: ""},
			expected: "{}",
		},
		{
			name: "string non-zero",
			value: struct {
				X string `json:"x,omitzero"`
			}{X: "hello"},
			expected: `{"x":"hello"}`,
		},
		{
			name: "bool zero",
			value: struct {
				X bool `json:"x,omitzero"`
			}{X: false},
			expected: "{}",
		},
		{
			name: "bool non-zero",
			value: struct {
				X bool `json:"x,omitzero"`
			}{X: true},
			expected: `{"x":true}`,
		},
		{
			name: "float64 zero",
			value: struct {
				X float64 `json:"x,omitzero"`
			}{X: 0.0},
			expected: "{}",
		},
		{
			name: "float64 non-zero",
			value: struct {
				X float64 `json:"x,omitzero"`
			}{X: 3.14},
			expected: `{"x":3.14}`,
		},

		// Pointer fields - nil is zero, but non-nil pointing to zero is NOT zero
		{
			name: "pointer nil",
			value: struct {
				X *int `json:"x,omitzero"`
				Y int  `json:"y,omitzero"`
			}{X: nil, Y: 0},
			expected: "{}",
		},
		{
			name: "pointer to zero",
			value: struct {
				X *int `json:"x,omitzero"`
				Y int  `json:"y,omitzero"`
			}{X: intPtr(0), Y: 0},
			expected: `{"x":0}`,
		},
		{
			name: "pointer to non-zero",
			value: struct {
				X *int `json:"x,omitzero"`
				Y int  `json:"y,omitzero"`
			}{X: intPtr(42), Y: 0},
			expected: `{"x":42}`,
		},

		// Slice fields - nil is zero, but empty slice is NOT zero (key difference from omitempty)
		{
			name: "slice nil",
			value: struct {
				X []int `json:"x,omitzero"`
				Y int  `json:"y,omitzero"`
			}{X: nil, Y: 0},
			expected: "{}",
		},
		{
			name: "slice empty",
			value: struct {
				X []int `json:"x,omitzero"`
				Y int  `json:"y,omitzero"`
			}{X: []int{}, Y: 0},
			expected: `{"x":[]}`,
		},
		{
			name: "slice with values",
			value: struct {
				X []int `json:"x,omitzero"`
				Y int  `json:"y,omitzero"`
			}{X: []int{1, 2, 3}, Y: 0},
			expected: `{"x":[1,2,3]}`,
		},

		// Comparison: omitempty vs omitzero for empty slice
		{
			name: "omitempty with empty slice",
			value: struct {
				X []int `json:"x,omitempty"`
				Y int  `json:"y,omitempty"`
			}{X: []int{}, Y: 0},
			expected: "{}",
		},
		{
			name: "omitzero with empty slice",
			value: struct {
				X []int `json:"x,omitzero"`
				Y int  `json:"y,omitzero"`
			}{X: []int{}, Y: 0},
			expected: `{"x":[]}`,
		},

		// Mixed tags
		{
			name: "mixed omitempty and omitzero",
			value: struct {
				A int       `json:"a,omitempty"`
				B int       `json:"b,omitzero"`
				C []int     `json:"c,omitempty"`
				D []int     `json:"d,omitzero"`
			}{
				A: 0,
				B: 0,
				C: []int{},
				D: []int{},
			},
			expected: `{"d":[]}`,
		},

		// Multiple fields with omitzero
		{
			name: "multiple fields omitzero",
			value: struct {
				X int    `json:"x,omitzero"`
				Y string `json:"y,omitzero"`
				Z bool   `json:"z,omitzero"`
			}{
				X: 0,
				Y: "",
				Z: false,
			},
			expected: "{}",
		},
		{
			name: "multiple fields omitzero mixed",
			value: struct {
				X int    `json:"x,omitzero"`
				Y string `json:"y,omitzero"`
				Z bool   `json:"z,omitzero"`
			}{
				X: 42,
				Y: "hello",
				Z: true,
			},
			expected: `{"x":42,"y":"hello","z":true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := Marshal(tt.value)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}
			if !bytes.Equal(b, []byte(tt.expected)) {
				t.Errorf("expected %s, got %s", tt.expected, string(b))
			}
		})
	}
}

// TestOmitZeroUnmarshal tests that unmarshal still works correctly with omitzero fields
func TestOmitZeroUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:  "unmarshal with omitzero field",
			input: `{"x":42}`,
			expected: struct {
				X int `json:"x,omitzero"`
			}{X: 42},
		},
		{
			name:  "unmarshal omitted field",
			input: `{}`,
			expected: struct {
				X int `json:"x,omitzero"`
			}{X: 0},
		},
		{
			name:  "unmarshal empty slice",
			input: `{"x":[]}`,
			expected: struct {
				X []int `json:"x,omitzero"`
			}{X: []int{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new instance of the expected type
			var result interface{}
			switch exp := tt.expected.(type) {
			case struct {
				X int `json:"x,omitzero"`
			}:
				var r struct {
					X int `json:"x,omitzero"`
				}
				if err := Unmarshal([]byte(tt.input), &r); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if r.X != exp.X {
					t.Errorf("expected X=%d, got X=%d", exp.X, r.X)
				}
				return
			case struct {
				X []int `json:"x,omitzero"`
			}:
				var r struct {
					X []int `json:"x,omitzero"`
				}
				if err := Unmarshal([]byte(tt.input), &r); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if len(r.X) != len(exp.X) {
					t.Errorf("expected X length=%d, got X length=%d", len(exp.X), len(r.X))
				}
				return
			}
			_ = result
		})
	}
}

// Benchmark omitzero
func BenchmarkOmitZero(b *testing.B) {
	type testStruct struct {
		X int       `json:"x,omitzero"`
		Y string    `json:"y,omitzero"`
		Z []int     `json:"z,omitzero"`
		W bool      `json:"w,omitzero"`
	}

	value := testStruct{
		X: 0,
		Y: "",
		Z: nil,
		W: false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(value)
		if err != nil {
			b.Fatalf("failed to marshal: %v", err)
		}
	}
}

// Helper function to create an int pointer
func intPtr(v int) *int {
	return &v
}
