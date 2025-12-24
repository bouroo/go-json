package json_test

import (
	"strings"
	"testing"

	"github.com/goccy/go-json"
)

func TestArrayNilPointerUnmarshal(t *testing.T) {
	t.Run("nil pointer with array JSON", func(t *testing.T) {
		var arr *[3]int
		err := json.Unmarshal([]byte("[1,2,3]"), arr)
		if err == nil {
			t.Fatal("expected InvalidUnmarshalError, got nil")
		}
	})

	t.Run("nil pointer with empty array JSON", func(t *testing.T) {
		var arr *[3]int
		err := json.Unmarshal([]byte("[]"), arr)
		if err == nil {
			t.Fatal("expected InvalidUnmarshalError, got nil")
		}
	})

	t.Run("nil pointer with null JSON", func(t *testing.T) {
		var arr *[3]int
		err := json.Unmarshal([]byte("null"), arr)
		if err == nil {
			t.Fatal("expected InvalidUnmarshalError, got nil")
		}
	})

	t.Run("valid pointer with array JSON", func(t *testing.T) {
		arr := [3]int{99, 99, 99}
		err := json.Unmarshal([]byte("[1,2,3]"), &arr)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if arr[0] != 1 || arr[1] != 2 || arr[2] != 3 {
			t.Fatalf("expected [1 2 3], got %v", arr)
		}
	})

	t.Run("valid pointer with empty array JSON", func(t *testing.T) {
		arr := [3]int{1, 2, 3}
		err := json.Unmarshal([]byte("[]"), &arr)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if arr[0] != 0 || arr[1] != 0 || arr[2] != 0 {
			t.Fatalf("expected [0 0 0], got %v", arr)
		}
	})

	t.Run("valid pointer with null JSON", func(t *testing.T) {
		arr := [3]int{1, 2, 3}
		err := json.Unmarshal([]byte("null"), &arr)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if arr[0] != 0 || arr[1] != 0 || arr[2] != 0 {
			t.Fatalf("expected zero values after null, got %v", arr)
		}
	})

	t.Run("Decoder with nil pointer", func(t *testing.T) {
		var arr *[3]int
		decoder := json.NewDecoder(strings.NewReader("[1,2,3]"))
		err := decoder.Decode(arr)
		if err == nil {
			t.Fatal("expected InvalidUnmarshalError, got nil")
		}
	})

	t.Run("Decoder with valid pointer", func(t *testing.T) {
		arr := [3]int{99, 99, 99}
		decoder := json.NewDecoder(strings.NewReader("[1,2,3]"))
		err := decoder.Decode(&arr)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if arr[0] != 1 || arr[1] != 2 || arr[2] != 3 {
			t.Fatalf("expected [1 2 3], got %v", arr)
		}
	})

	t.Run("nested struct with nil array pointer", func(t *testing.T) {
		type S struct {
			Arr [3]int
		}
		var s *S
		err := json.Unmarshal([]byte(`{"Arr":[1,2,3]}`), s)
		if err == nil {
			t.Fatal("expected InvalidUnmarshalError, got nil")
		}
	})
}
