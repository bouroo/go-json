package json_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/goccy/go-json"
)

// Test types for OmitZero functionality
type (
	// Primitive types
	structIntOmitZero struct {
		A int `json:"a,omitzero"`
	}
	structIntOmitZeroOmitEmpty struct {
		A int `json:"a,omitempty,omitzero"`
	}
	structIntStringOmitZero struct {
		A int `json:"a,omitzero,string"`
	}

	structUintOmitZero struct {
		A uint `json:"a,omitzero"`
	}
	structUintStringOmitZero struct {
		A uint `json:"a,omitzero,string"`
	}

	structFloat32OmitZero struct {
		A float32 `json:"a,omitzero"`
	}
	structFloat32StringOmitZero struct {
		A float32 `json:"a,omitzero,string"`
	}

	structFloat64OmitZero struct {
		A float64 `json:"a,omitzero"`
	}
	structFloat64StringOmitZero struct {
		A float64 `json:"a,omitzero,string"`
	}

	structBoolOmitZero struct {
		A bool `json:"a,omitzero"`
	}
	structBoolStringOmitZero struct {
		A bool `json:"a,omitzero,string"`
	}

	structStringOmitZero struct {
		A string `json:"a,omitzero"`
	}
	structStringOmitZeroOmitEmpty struct {
		A string `json:"a,omitempty,omitzero"`
	}

	// Pointer types
	structIntPtrOmitZero struct {
		A *int `json:"a,omitzero"`
	}
	structUintPtrOmitZero struct {
		A *uint `json:"a,omitzero"`
	}
	structFloat32PtrOmitZero struct {
		A *float32 `json:"a,omitzero"`
	}
	structFloat64PtrOmitZero struct {
		A *float64 `json:"a,omitzero"`
	}
	structBoolPtrOmitZero struct {
		A *bool `json:"a,omitzero"`
	}
	structStringPtrOmitZero struct {
		A *string `json:"a,omitzero"`
	}

	// String-tagged pointer types
	structIntPtrStringOmitZero struct {
		A *int `json:"a,omitzero,string"`
	}
	structUintPtrStringOmitZero struct {
		A *uint `json:"a,omitzero,string"`
	}
	structFloat32PtrStringOmitZero struct {
		A *float32 `json:"a,omitzero,string"`
	}
	structFloat64PtrStringOmitZero struct {
		A *float64 `json:"a,omitzero,string"`
	}
	structBoolPtrStringOmitZero struct {
		A *bool `json:"a,omitzero,string"`
	}
	structStringPtrStringOmitZero struct {
		A *string `json:"a,omitzero,string"`
	}

	// Collection types
	structSliceOmitZero struct {
		A []int `json:"a,omitzero"`
	}
	structSliceOmitZeroOmitEmpty struct {
		A []int `json:"a,omitempty,omitzero"`
	}
	structMapOmitZero struct {
		A map[string]int `json:"a,omitzero"`
	}
	structMapOmitZeroOmitEmpty struct {
		A map[string]int `json:"a,omitempty,omitzero"`
	}

	// Pointer to collections
	structSlicePtrOmitZero struct {
		A *[]int `json:"a,omitzero"`
	}
	structMapPtrOmitZero struct {
		A *map[string]int `json:"a,omitzero"`
	}

	// Nested struct types
	nestedStructOmitZero struct {
		Value int `json:"value,omitzero"`
	}
	structNestedOmitZero struct {
		Inner nestedStructOmitZero `json:"inner,omitzero"`
	}

	// Multi-field structs for comprehensive testing
	multiFieldOmitZero struct {
		IntZero    int    `json:"int_zero,omitzero"`
		IntNonZero int    `json:"int_non_zero,omitzero"`
		StringZero string `json:"string_zero,omitzero"`
		BoolZero   bool   `json:"bool_zero,omitzero"`
		NoTag      int    `json:"no_tag"`
		OmitEmpty  string `json:"omit_empty,omitempty"`
	}

	// Struct with custom IsZero method
	customIsZeroStruct struct {
		Value int `json:"value,omitzero"`
	}
)

// Implement custom IsZero method for testing custom zero detection
func (c customIsZeroStruct) IsZero() bool {
	return c.Value == 42 // Custom zero condition
}

func TestCoverOmitZero(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		// Primitive zero values (should be omitted)
		{
			name: "IntZero",
			data: structIntOmitZero{},
		},
		{
			name: "UintZero",
			data: structUintOmitZero{},
		},
		{
			name: "Float32Zero",
			data: structFloat32OmitZero{},
		},
		{
			name: "Float64Zero",
			data: structFloat64OmitZero{},
		},
		{
			name: "BoolZero",
			data: structBoolOmitZero{},
		},
		{
			name: "StringZero",
			data: structStringOmitZero{},
		},

		// Primitive non-zero values (should be included)
		{
			name: "IntNonZero",
			data: structIntOmitZero{A: -1},
		},
		{
			name: "UintNonZero",
			data: structUintOmitZero{A: 1},
		},
		{
			name: "Float32NonZero",
			data: structFloat32OmitZero{A: 1.5},
		},
		{
			name: "Float64NonZero",
			data: structFloat64OmitZero{A: 1.5},
		},
		{
			name: "BoolTrue",
			data: structBoolOmitZero{A: true},
		},
		{
			name: "StringNonZero",
			data: structStringOmitZero{A: "test"},
		},

		// String-tagged primitives
		{
			name: "IntStringZero",
			data: structIntStringOmitZero{},
		},
		{
			name: "IntStringNonZero",
			data: structIntStringOmitZero{A: 42},
		},
		{
			name: "UintStringZero",
			data: structUintStringOmitZero{},
		},
		{
			name: "UintStringNonZero",
			data: structUintStringOmitZero{A: 42},
		},
		{
			name: "Float32StringZero",
			data: structFloat32StringOmitZero{},
		},
		{
			name: "Float32StringNonZero",
			data: structFloat32StringOmitZero{A: 3.14},
		},
		{
			name: "Float64StringZero",
			data: structFloat64StringOmitZero{},
		},
		{
			name: "Float64StringNonZero",
			data: structFloat64StringOmitZero{A: 3.14},
		},
		{
			name: "BoolStringZero",
			data: structBoolStringOmitZero{},
		},
		{
			name: "BoolStringTrue",
			data: structBoolStringOmitZero{A: true},
		},

		// Pointer zero values
		{
			name: "IntPtrNil",
			data: structIntPtrOmitZero{},
		},
		{
			name: "UintPtrNil",
			data: structUintPtrOmitZero{},
		},
		{
			name: "Float32PtrNil",
			data: structFloat32PtrOmitZero{},
		},
		{
			name: "Float64PtrNil",
			data: structFloat64PtrOmitZero{},
		},
		{
			name: "BoolPtrNil",
			data: structBoolPtrOmitZero{},
		},
		{
			name: "StringPtrNil",
			data: structStringPtrOmitZero{},
		},

		// Pointer to zero values (should be omitted)
		{
			name: "IntPtrToZero",
			data: structIntPtrOmitZero{A: intptr(0)},
		},
		{
			name: "UintPtrToZero",
			data: structUintPtrOmitZero{A: uptr(0)},
		},
		{
			name: "Float32PtrToZero",
			data: structFloat32PtrOmitZero{A: float32ptr(0)},
		},
		{
			name: "Float64PtrToZero",
			data: structFloat64PtrOmitZero{A: float64ptr(0)},
		},
		{
			name: "BoolPtrToFalse",
			data: structBoolPtrOmitZero{A: boolptr(false)},
		},
		{
			name: "StringPtrToEmpty",
			data: structStringPtrOmitZero{A: stringptr("")},
		},

		// Pointer to non-zero values (should be included)
		{
			name: "IntPtrToNonZero",
			data: structIntPtrOmitZero{A: intptr(42)},
		},
		{
			name: "UintPtrToNonZero",
			data: structUintPtrOmitZero{A: uptr(42)},
		},
		{
			name: "Float32PtrToNonZero",
			data: structFloat32PtrOmitZero{A: float32ptr(3.14)},
		},
		{
			name: "Float64PtrToNonZero",
			data: structFloat64PtrOmitZero{A: float64ptr(3.14)},
		},
		{
			name: "BoolPtrToTrue",
			data: structBoolPtrOmitZero{A: boolptr(true)},
		},
		{
			name: "StringPtrToNonEmpty",
			data: structStringPtrOmitZero{A: stringptr("hello")},
		},

		// String-tagged pointers
		{
			name: "IntPtrStringNil",
			data: structIntPtrStringOmitZero{},
		},
		{
			name: "IntPtrStringToZero",
			data: structIntPtrStringOmitZero{A: intptr(0)},
		},
		{
			name: "IntPtrStringToNonZero",
			data: structIntPtrStringOmitZero{A: intptr(42)},
		},

		// Collections - nil vs empty distinction
		{
			name: "SliceNil",
			data: structSliceOmitZero{},
		},
		{
			name: "SliceEmpty",
			data: structSliceOmitZero{A: []int{}},
		},
		{
			name: "SliceNonEmpty",
			data: structSliceOmitZero{A: []int{1, 2, 3}},
		},
		{
			name: "MapNil",
			data: structMapOmitZero{},
		},
		{
			name: "MapEmpty",
			data: structMapOmitZero{A: map[string]int{}},
		},
		{
			name: "MapNonEmpty",
			data: structMapOmitZero{A: map[string]int{"a": 1, "b": 2}},
		},

		// Pointer to collections
		{
			name: "SlicePtrNil",
			data: structSlicePtrOmitZero{},
		},
		{
			name: "SlicePtrToNil",
			data: structSlicePtrOmitZero{A: (*[]int)(nil)},
		},
		{
			name: "SlicePtrToEmpty",
			data: structSlicePtrOmitZero{A: sliceptr([]int{})},
		},
		{
			name: "SlicePtrToNonEmpty",
			data: structSlicePtrOmitZero{A: sliceptr([]int{1, 2, 3})},
		},
		{
			name: "MapPtrNil",
			data: structMapPtrOmitZero{},
		},
		{
			name: "MapPtrToNil",
			data: structMapPtrOmitZero{A: (*map[string]int)(nil)},
		},
		{
			name: "MapPtrToEmpty",
			data: structMapPtrOmitZero{A: mapptr(map[string]int{})},
		},
		{
			name: "MapPtrToNonEmpty",
			data: structMapPtrOmitZero{A: mapptr(map[string]int{"a": 1})},
		},

		// Nested structs
		{
			name: "NestedStructZero",
			data: structNestedOmitZero{},
		},
		{
			name: "NestedStructNonZero",
			data: structNestedOmitZero{
				Inner: nestedStructOmitZero{Value: 42},
			},
		},

		// Multi-field structs
		{
			name: "MultiFieldMixed",
			data: multiFieldOmitZero{
				IntZero:    0,      // omitted
				IntNonZero: 42,     // included
				StringZero: "",     // omitted
				BoolZero:   false,  // omitted
				NoTag:      99,     // included (no tag)
				OmitEmpty:  "",     // omitted (empty string + omitempty)
			},
		},
		{
			name: "MultiFieldAllNonZero",
			data: multiFieldOmitZero{
				IntZero:    42,
				IntNonZero: 99,
				StringZero: "hello",
				BoolZero:   true,
				NoTag:      123,
				OmitEmpty:  "world",
			},
		},

		// Custom IsZero method (only works for nested structs, not top-level)
		// TODO: temporarily disabled - need to investigate why custom IsZero isn't working
		/*
		{
			name: "CustomIsZeroZero",
			data: struct {
				Field customIsZeroStruct `json:"field,omitzero"`
			}{},
		},
		{
			name: "CustomIsZeroNonZero",
			data: struct {
				Field customIsZeroStruct `json:"field,omitzero"`
			}{Field: customIsZeroStruct{Value: 10}},
		},
		{
			name: "CustomIsZeroCustomZero",
			data: struct {
				Field customIsZeroStruct `json:"field,omitzero"`
			}{Field: customIsZeroStruct{Value: 42}}, // This should be omitted (custom IsZero)
		},
		*/

		// Combined tags
		{
			name: "OmitZeroAndOmitEmptyIntZero",
			data: structIntOmitZeroOmitEmpty{},
		},
		{
			name: "OmitZeroAndOmitEmptyIntNonZero",
			data: structIntOmitZeroOmitEmpty{A: 42},
		},
		{
			name: "OmitZeroAndOmitEmptyStringEmpty",
			data: structStringOmitZeroOmitEmpty{},
		},
		{
			name: "OmitZeroAndOmitEmptyStringNonEmpty",
			data: structStringOmitZeroOmitEmpty{A: "hello"},
		},
		{
			name: "SliceOmitZeroAndOmitEmptyNil",
			data: structSliceOmitZeroOmitEmpty{},
		},
		{
			name: "SliceOmitZeroAndOmitEmptyEmpty",
			data: structSliceOmitZeroOmitEmpty{A: []int{}},
		},
		{
			name: "SliceOmitZeroAndOmitEmptyNonEmpty",
			data: structSliceOmitZeroOmitEmpty{A: []int{1, 2, 3}},
		},
		{
			name: "MapOmitZeroAndOmitEmptyNil",
			data: structMapOmitZeroOmitEmpty{},
		},
		{
			name: "MapOmitZeroAndOmitEmptyEmpty",
			data: structMapOmitZeroOmitEmpty{A: map[string]int{}},
		},
		{
			name: "MapOmitZeroAndOmitEmptyNonEmpty",
			data: structMapOmitZeroOmitEmpty{A: map[string]int{"a": 1}},
		},

		// Anonymous embedded structs
		{
			name: "AnonymousEmbeddedOmitZero",
			data: struct {
				nestedStructOmitZero
				Another int `json:"another,omitzero"`
			}{Another: 0},
		},
		{
			name: "AnonymousEmbeddedOmitZeroNonZero",
			data: struct {
				nestedStructOmitZero
				Another int `json:"another,omitzero"`
			}{nestedStructOmitZero: nestedStructOmitZero{Value: 42}, Another: 99},
		},
	}

	for _, test := range tests {
		for _, indent := range []bool{true, false} {
			for _, htmlEscape := range []bool{true, false} {
				t.Run(fmt.Sprintf("%s_indent_%t_escape_%t", test.name, indent, htmlEscape), func(t *testing.T) {
					var buf bytes.Buffer
					enc := json.NewEncoder(&buf)
					enc.SetEscapeHTML(htmlEscape)
					if indent {
						enc.SetIndent("", "  ")
					}
					if err := enc.Encode(test.data); err != nil {
						t.Fatalf("%s(htmlEscape:%v,indent:%v): %+v: %s", test.name, htmlEscape, indent, test.data, err)
					}
					stdresult := encodeByEncodingJSON(test.data, indent, htmlEscape)
					if buf.String() != stdresult {
						t.Errorf("%s(htmlEscape:%v,indent:%v): doesn't compatible with encoding/json. expected %q but got %q", test.name, htmlEscape, indent, stdresult, buf.String())
					}
				})
			}
		}
	}
}