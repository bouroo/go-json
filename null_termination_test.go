package json_test

import (
	"testing"

	"github.com/goccy/go-json"
)

func TestNullTerminationWithCapacity(t *testing.T) {
	data := []byte(`{"key":"value","number":42}`)
	data = append(data, 0)
	data = data[:len(data)-1]

	var v struct {
		Key    string
		Number int
	}

	err := json.Unmarshal(data, &v)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if v.Key != "value" || v.Number != 42 {
		t.Errorf("Unmarshal result: Key=%q, Number=%d, want Key='value', Number=42", v.Key, v.Number)
	}
}

func TestNullTerminationWithoutCapacity(t *testing.T) {
	data := []byte(`{"key":"value","number":42}`)

	var v struct {
		Key    string
		Number int
	}

	err := json.Unmarshal(data, &v)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if v.Key != "value" || v.Number != 42 {
		t.Errorf("Unmarshal result: Key=%q, Number=%d, want Key='value', Number=42", v.Key, v.Number)
	}
}

func TestUnmarshalBasicTypes(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		v    interface{}
	}{
		{"string", []byte(`"hello"`), new(string)},
		{"number", []byte(`42`), new(int)},
		{"bool", []byte(`true`), new(bool)},
		{"null", []byte(`null`), new(*int)},
		{"object", []byte(`{"a":1}`), new(map[string]int)},
		{"array", []byte(`[1,2,3]`), new([]int)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := json.Unmarshal(tc.data, tc.v)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
		})
	}
}

func BenchmarkUnmarshalSmallWithCapacity(b *testing.B) {
	data := []byte(`{"key":"value","number":42}`)
	data = append(data, 0)
	data = data[:len(data)-1]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v struct {
			Key    string
			Number int
		}
		_ = json.Unmarshal(data, &v)
	}
}

func BenchmarkUnmarshalSmallWithoutCapacity(b *testing.B) {
	data := []byte(`{"key":"value","number":42}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v struct {
			Key    string
			Number int
		}
		_ = json.Unmarshal(data, &v)
	}
}
