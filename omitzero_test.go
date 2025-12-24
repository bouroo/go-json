package json_test

import (
	stdjson "encoding/json"
	"testing"
	"time"

	gojson "github.com/goccy/go-json"
)

type OmitZeroTest struct {
	IntField    int       `json:"int_field,omitzero"`
	StringField string    `json:"string_field,omitzero"`
	BoolField   bool      `json:"bool_field,omitzero"`
	FloatField  float64   `json:"float_field,omitzero"`
	TimeField   time.Time `json:"time_field,omitzero"`
	NormalField string    `json:"normal_field"`
}

type OmitZeroWithValue struct {
	IntField    int       `json:"int_field,omitzero"`
	StringField string    `json:"string_field,omitzero"`
	TimeField   time.Time `json:"time_field,omitzero"`
	NormalField string    `json:"normal_field"`
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
				TimeField:   time.Time{},
				NormalField: "test",
			},
			expected: `{"normal_field":"test"}`,
		},
		{
			name: "non-zero values",
			value: OmitZeroWithValue{
				IntField:    42,
				StringField: "hello",
				TimeField:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				NormalField: "world",
			},
			expected: `{"int_field":42,"string_field":"hello","time_field":"2023-01-01T00:00:00Z","normal_field":"world"}`,
		},
		{
			name: "mixed zero and non-zero",
			value: OmitZeroWithValue{
				IntField:    0,
				StringField: "hello",
				TimeField:   time.Time{},
				NormalField: "world",
			},
			expected: `{"string_field":"hello","normal_field":"world"}`,
		},
		{
			name: "time field with non-zero value",
			value: OmitZeroWithValue{
				IntField:    0,
				StringField: "",
				TimeField:   time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC),
				NormalField: "world",
			},
			expected: `{"time_field":"2023-12-25T15:30:45Z","normal_field":"world"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test go-json
			b, err := gojson.Marshal(tt.value)
			if err != nil {
				t.Fatalf("go-json.Marshal failed: %v", err)
			}

			got := string(b)
			if got != tt.expected {
				t.Errorf("go-json.Marshal() = %s, want %s", got, tt.expected)
			}

			// Compare with stdlib json for basic structure
			stdlibB, _ := stdjson.Marshal(tt.value)
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

	b, err := gojson.Marshal(data)
	if err != nil {
		t.Fatalf("go-json.Marshal failed: %v", err)
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
		TimeField:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		NormalField: "world",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gojson.Marshal(data)
	}
}

func BenchmarkOmitZeroAllZero(b *testing.B) {
	data := OmitZeroTest{
		NormalField: "world",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gojson.Marshal(data)
	}
}
