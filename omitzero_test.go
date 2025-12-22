package json

import (
	"encoding/json"
	"testing"
)

type OmitZeroTest struct {
	IntField    int    `json:"int_field,omitzero"`
	StringField string `json:"string_field,omitzero"`
	BoolField   bool   `json:"bool_field,omitzero"`
	FloatField  float64 `json:"float_field,omitzero"`
	NormalField string `json:"normal_field"`
}

type OmitZeroWithValue struct {
	IntField    int    `json:"int_field,omitzero"`
	StringField string `json:"string_field,omitzero"`
	NormalField string `json:"normal_field"`
}

func TestOmitZeroBasic(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name: "all zero values",
			value: OmitZeroTest{
				IntField:    0,
				StringField: "",
				BoolField:   false,
				FloatField:  0.0,
				NormalField: "test",
			},
			expected: `{"normal_field":"test"}`,
		},
		{
			name: "non-zero values",
			value: OmitZeroWithValue{
				IntField:    42,
				StringField: "hello",
				NormalField: "world",
			},
			expected: `{"int_field":42,"string_field":"hello","normal_field":"world"}`,
		},
		{
			name: "mixed zero and non-zero",
			value: OmitZeroWithValue{
				IntField:    0,
				StringField: "hello",
				NormalField: "world",
			},
			expected: `{"string_field":"hello","normal_field":"world"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test go-json
			b, err := Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			got := string(b)
			if got != tt.expected {
				t.Errorf("go-json.Marshal() = %s, want %s", got, tt.expected)
			}

			// Compare with stdlib json for basic structure
			stdlibB, _ := json.Marshal(tt.value)
			t.Logf("stdlib json.Marshal() = %s", string(stdlibB))
			t.Logf("go-json.Marshal()   = %s", got)
		})
	}
}

func TestOmitZeroEmptyStruct(t *testing.T) {
	// Test that a struct with all omitzero fields and zero values produces minimal output
	data := OmitZeroTest{
		NormalField: "value",
	}

	b, err := Marshal(data)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	expected := `{"normal_field":"value"}`
	got := string(b)

	if got != expected {
		t.Errorf("got %s, want %s", got, expected)
	}
}

func BenchmarkOmitZero(b *testing.B) {
	data := OmitZeroTest{
		IntField:    42,
		StringField: "hello",
		NormalField: "world",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Marshal(data)
	}
}

func BenchmarkOmitZeroAllZero(b *testing.B) {
	data := OmitZeroTest{
		NormalField: "world",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Marshal(data)
	}
}
