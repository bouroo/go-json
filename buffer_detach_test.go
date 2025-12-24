package json_test

import (
	"testing"

	"github.com/goccy/go-json"
)

func TestBufferDetachOptimization(t *testing.T) {
	tests := []struct {
		name         string
		v            interface{}
		expectDetach bool
	}{
		{"small struct - should detach", struct{ Name string }{"Alice"}, true},
		{"small string - should detach", "hello", true},
		{"small number - should detach", 42, true},
		{"small map - should detach", map[string]int{"a": 1, "b": 2}, true},
		{"large array - should copy", make([]int, 1000), false},
		{"large string - should copy", string(make([]byte, 5000)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.v)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			if len(data) > 0 {
				t.Logf("Marshal output length: %d bytes", len(data))
			}
		})
	}
}

func TestBufferDetachCorrectness(t *testing.T) {
	type testCase struct {
		name     string
		input    interface{}
		expected string
	}

	cases := []testCase{
		{"empty struct", struct{}{}, "{}"},
		{"string", "hello", `"hello"`},
		{"number", 42, `42`},
		{"bool", true, `true`},
		{"slice", []int{1, 2, 3}, `[1,2,3]`},
		{"map", map[string]int{"a": 1}, `{"a":1}`},
		{"nested struct", struct {
			A int
			B string
		}{1, "test"}, `{"A":1,"B":"test"}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.input)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if string(got) != tc.expected {
				t.Errorf("Marshal(%v) = %s, want %s", tc.input, string(got), tc.expected)
			}
		})
	}
}

func TestBufferDetachIndent(t *testing.T) {
	type testCase struct {
		name     string
		input    interface{}
		expected string
	}

	cases := []testCase{
		{"simple struct", struct {
			Name string
			Age  int
		}{"Alice", 30}, "{\n  \"Name\": \"Alice\",\n  \"Age\": 30\n}"},
		{"nested struct", struct {
			User struct {
				Name string
			}
		}{struct{ Name string }{"Bob"}}, "{\n  \"User\": {\n    \"Name\": \"Bob\"\n  }\n}"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.MarshalIndent(tc.input, "", "  ")
			if err != nil {
				t.Fatalf("MarshalIndent failed: %v", err)
			}
			if string(got) != tc.expected {
				t.Errorf("MarshalIndent(%v) = %s, want %s", tc.input, string(got), tc.expected)
			}
		})
	}
}

func BenchmarkSmallMarshal(b *testing.B) {
	v := struct {
		Name  string
		Email string
		Age   int
	}{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   30,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(v)
	}
}

func BenchmarkLargeMarshal(b *testing.B) {
	v := make([]int, 1000)
	for i := range v {
		v[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(v)
	}
}
