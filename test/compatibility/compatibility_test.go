package compatibility_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strings"
	"testing"
	"time"

	gojson "github.com/goccy/go-json"
)

// ============================================================================
// Type Definitions for Testing
// ============================================================================

// Primitive types
type Primitives struct {
	Int     int     `json:"int"`
	Int8    int8    `json:"int8"`
	Int16   int16   `json:"int16"`
	Int32   int32   `json:"int32"`
	Int64   int64   `json:"int64"`
	Uint    uint    `json:"uint"`
	Uint8   uint8   `json:"uint8"`
	Uint16  uint16  `json:"uint16"`
	Uint32  uint32  `json:"uint32"`
	Uint64  uint64  `json:"uint64"`
	Float32 float32 `json:"float32"`
	Float64 float64 `json:"float64"`
	Bool    bool    `json:"bool"`
	String  string  `json:"string"`
}

// Primitives with omitzero
type PrimitivesOmitZero struct {
	Int     int     `json:"int,omitzero"`
	Int8    int8    `json:"int8,omitzero"`
	Int16   int16   `json:"int16,omitzero"`
	Int32   int32   `json:"int32,omitzero"`
	Int64   int64   `json:"int64,omitzero"`
	Uint    uint    `json:"uint,omitzero"`
	Uint8   uint8   `json:"uint8,omitzero"`
	Uint16  uint16  `json:"uint16,omitzero"`
	Uint32  uint32  `json:"uint32,omitzero"`
	Uint64  uint64  `json:"uint64,omitzero"`
	Float32 float32 `json:"float32,omitzero"`
	Float64 float64 `json:"float64,omitzero"`
	Bool    bool    `json:"bool,omitzero"`
	String  string  `json:"string,omitzero"`
}

// Primitives with omitempty
type PrimitivesOmitEmpty struct {
	Int     int     `json:"int,omitempty"`
	Int8    int8    `json:"int8,omitempty"`
	Int16   int16   `json:"int16,omitempty"`
	Int32   int32   `json:"int32,omitempty"`
	Int64   int64   `json:"int64,omitempty"`
	Uint    uint    `json:"uint,omitempty"`
	Uint8   uint8   `json:"uint8,omitempty"`
	Uint16  uint16  `json:"uint16,omitempty"`
	Uint32  uint32  `json:"uint32,omitempty"`
	Uint64  uint64  `json:"uint64,omitempty"`
	Float32 float32 `json:"float32,omitempty"`
	Float64 float64 `json:"float64,omitempty"`
	Bool    bool    `json:"bool,omitempty"`
	String  string  `json:"string,omitempty"`
}

// Primitives with string tag
type PrimitivesString struct {
	Int     int     `json:"int,string"`
	Int8    int8    `json:"int8,string"`
	Int16   int16   `json:"int16,string"`
	Int32   int32   `json:"int32,string"`
	Int64   int64   `json:"int64,string"`
	Uint    uint    `json:"uint,string"`
	Uint8   uint8   `json:"uint8,string"`
	Uint16  uint16  `json:"uint16,string"`
	Uint32  uint32  `json:"uint32,string"`
	Uint64  uint64  `json:"uint64,string"`
	Float32 float32 `json:"float32,string"`
	Float64 float64 `json:"float64,string"`
	Bool    bool    `json:"bool,string"`
}

// Pointer types
type Pointers struct {
	Int     *int     `json:"int"`
	Int8    *int8    `json:"int8"`
	Int16   *int16   `json:"int16"`
	Int32   *int32   `json:"int32"`
	Int64   *int64   `json:"int64"`
	Uint    *uint    `json:"uint"`
	Uint8   *uint8   `json:"uint8"`
	Uint16  *uint16  `json:"uint16"`
	Uint32  *uint32  `json:"uint32"`
	Uint64  *uint64  `json:"uint64"`
	Float32 *float32 `json:"float32"`
	Float64 *float64 `json:"float64"`
	Bool    *bool    `json:"bool"`
	String  *string  `json:"string"`
}

// Pointer types with omitzero
type PointersOmitZero struct {
	Int     *int     `json:"int,omitzero"`
	Int8    *int8    `json:"int8,omitzero"`
	Int16   *int16   `json:"int16,omitzero"`
	Int32   *int32   `json:"int32,omitzero"`
	Int64   *int64   `json:"int64,omitzero"`
	Uint    *uint    `json:"uint,omitzero"`
	Uint8   *uint8   `json:"uint8,omitzero"`
	Uint16  *uint16  `json:"uint16,omitzero"`
	Uint32  *uint32  `json:"uint32,omitzero"`
	Uint64  *uint64  `json:"uint64,omitzero"`
	Float32 *float32 `json:"float32,omitzero"`
	Float64 *float64 `json:"float64,omitzero"`
	Bool    *bool    `json:"bool,omitzero"`
	String  *string  `json:"string,omitzero"`
}

// Multi-level pointers
type MultiLevelPointers struct {
	Int ***int `json:"int"`
}

// Collections
type Collections struct {
	Slice    []int               `json:"slice"`
	SlicePtr []*int              `json:"slice_ptr"`
	Map      map[string]int      `json:"map"`
	MapPtr   map[string]*int     `json:"map_ptr"`
	Array    [3]int              `json:"array"`
	ArrayPtr [2]*string          `json:"array_ptr"`
}

// Collections with omitzero (note: arrays don't support omitzero like slices in standard library)
type CollectionsOmitZero struct {
	Slice    []int           `json:"slice,omitzero"`
	SlicePtr []*int          `json:"slice_ptr,omitzero"`
	Map      map[string]int  `json:"map,omitzero"`
	MapPtr   map[string]*int `json:"map_ptr,omitzero"`
}

// Collections with omitempty
type CollectionsOmitEmpty struct {
	Slice    []int               `json:"slice,omitempty"`
	SlicePtr []*int              `json:"slice_ptr,omitempty"`
	Map      map[string]int      `json:"map,omitempty"`
	MapPtr   map[string]*int     `json:"map_ptr,omitempty"`
	Array    [3]int              `json:"array,omitempty"`
}

// Nested structs - Level 1
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

type Person struct {
	Name    string    `json:"name"`
	Age     int       `json:"age"`
	Email   string    `json:"email"`
	Address Address  `json:"address"`
}

// Nested structs - Level 2 (embedded)
type Contact struct {
	Phone string `json:"phone"`
	Fax   string `json:"fax"`
}

type Employee struct {
	Person
	Contact    Contact  `json:"contact"`
	Department string   `json:"department"`
	Salary     float64  `json:"salary"`
	HireDate   time.Time `json:"hire_date"`
}

// Nested structs - Level 3
type Company struct {
	Name     string     `json:"name"`
	Employees []Employee `json:"employees"`
	Location Address    `json:"location"`
}

// Deeply nested (5 levels)
type Level1 struct {
	Name string `json:"name"`
}

type Level2 struct {
	Level1 Level1 `json:"level1"`
	Depth  int    `json:"depth"`
}

type Level3 struct {
	Level2 Level2 `json:"level2"`
	Depth  int    `json:"depth"`
}

type Level4 struct {
	Level3 Level3 `json:"level3"`
	Depth  int    `json:"depth"`
}

type Level5 struct {
	Level4 Level4 `json:"level4"`
	Depth  int    `json:"depth"`
	Value  string `json:"value"`
}

// Embedded struct
type BaseInfo struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

type EmbeddedUser struct {
	BaseInfo `json:"base_info"`
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// Anonymous fields
type AnonymousInner struct {
	InnerValue string `json:"inner_value"`
}

type AnonymousOuter struct {
	AnonymousInner
	OuterValue string `json:"outer_value"`
}

type AnonymousTagged struct {
	Inner string `json:"inner"`
}

// Interface types
type InterfaceStruct struct {
	Value interface{} `json:"value"`
}

type InterfaceWithType struct {
	Number  interface{} `json:"number"`
	String  interface{} `json:"string"`
	Bool    interface{} `json:"bool"`
	Slice   interface{} `json:"slice"`
	Map     interface{} `json:"map"`
	Null    interface{} `json:"null"`
}

// Special types
type TimeHolder struct {
	Created time.Time `json:"created"`
	Updated *time.Time `json:"updated,omitempty"`
}

type RawMessageHolder struct {
	Data json.RawMessage `json:"data"`
}

type NumberHolder struct {
	Num json.Number `json:"num"`
}

// OmitZero with nested structs
type OmitZeroNested struct {
	Name    string  `json:"name,omitzero"`
	Age     int     `json:"age,omitzero"`
	Address Address `json:"address,omitzero"`
}

// OmitZero with custom IsZero
type CustomZero struct {
	Threshold float64 `json:"threshold"`
}

func (c CustomZero) IsZero() bool {
	return math.Abs(c.Threshold) < 0.001
}

type OmitZeroCustom struct {
	Name      string     `json:"name,omitzero"`
	Threshold CustomZero `json:"threshold,omitzero"`
}

// Custom Marshaler
type CustomMarshal struct {
	Value int
}

func (c CustomMarshal) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"custom":%d}`, c.Value)), nil
}

type WithCustomMarshal struct {
	Items []CustomMarshal `json:"items"`
}

// Custom TextMarshaler
type TextMarshalValue struct {
	Val string
}

func (t TextMarshalValue) MarshalText() ([]byte, error) {
	return []byte("text:" + t.Val), nil
}

type WithTextMarshal struct {
	Name  string           `json:"name"`
	Value TextMarshalValue `json:"value"`
}

// Ignored field
type WithIgnored struct {
	Public  string `json:"public"`
	Private string `json:"-"`
}

// Field named with special character
type SpecialFieldNames struct {
	Dash    string `json:"-,"`
	Empty   string `json:",omitempty"`
	WithTag string `json:"my_field_name"`
}

// Recursive structure
type Node struct {
	Value int     `json:"value"`
	Left  *Node   `json:"left,omitempty"`
	Right *Node   `json:"right,omitempty"`
}

// Mutually recursive structures
type Tree struct {
	Value int   `json:"value"`
	Left  *Tree `json:"left,omitempty"`
}

type Forest struct {
	Trees []Tree `json:"trees"`
}

// ============================================================================
// Test Helpers
// ============================================================================

func compareMarshal(t *testing.T, name string, v interface{}) {
	stdJSON, stdErr := json.Marshal(v)
	goJSON, goErr := gojson.Marshal(v)

	// Compare errors
	if (stdErr == nil) != (goErr == nil) {
		t.Errorf("%s: error mismatch - stdlib: %v, go-json: %v", name, stdErr, goErr)
		return
	}

	// Compare output
	if string(stdJSON) != string(goJSON) {
		t.Errorf("%s: output mismatch\nstdlib: %s\ngo-json: %s", name, string(stdJSON), string(goJSON))
	}
}

func compareMarshalIndent(t *testing.T, name string, v interface{}, prefix, indent string) {
	stdBuf := &bytes.Buffer{}
	stdEnc := json.NewEncoder(stdBuf)
	stdEnc.SetIndent(prefix, indent)
	stdErr := stdEnc.Encode(v)

	goBuf := &bytes.Buffer{}
	goEnc := gojson.NewEncoder(goBuf)
	goEnc.SetIndent(prefix, indent)
	goErr := goEnc.Encode(v)

	// Compare errors
	if (stdErr == nil) != (goErr == nil) {
		t.Errorf("%s: error mismatch - stdlib: %v, go-json: %v", name, stdErr, goErr)
		return
	}

	// Compare output
	if stdBuf.String() != goBuf.String() {
		t.Errorf("%s: indent output mismatch\nstdlib: %s\ngo-json: %s", name, stdBuf.String(), goBuf.String())
	}
}

func compareUnmarshal(t *testing.T, name string, jsonData string, v interface{}) {
	stdErr := json.Unmarshal([]byte(jsonData), v)
	goErr := gojson.Unmarshal([]byte(jsonData), v)

	// Compare errors
	if (stdErr == nil) != (goErr == nil) {
		t.Errorf("%s: unmarshal error mismatch - stdlib: %v, go-json: %v", name, stdErr, goErr)
		return
	}

	// For successful unmarshal, verify the values match
	if stdErr == nil {
		// Create copies to compare
		stdV := reflect.New(reflect.TypeOf(v).Elem()).Interface()
		goV := reflect.New(reflect.TypeOf(v).Elem()).Interface()

		json.Unmarshal([]byte(jsonData), stdV)
		gojson.Unmarshal([]byte(jsonData), goV)

		if !reflect.DeepEqual(stdV, goV) {
			t.Errorf("%s: unmarshal value mismatch\nstdlib: %+v\ngo-json: %+v", name, stdV, goV)
		}
	}
}

func compareRoundTrip(t *testing.T, name string, original interface{}) {
	// Marshal with go-json
	marshaled, err := gojson.Marshal(original)
	if err != nil {
		t.Errorf("%s: go-json marshal error: %v", name, err)
		return
	}

	// Get the underlying type and create a new instance
	v := reflect.ValueOf(original)
	isPtr := v.Kind() == reflect.Ptr

	// For pointer types, use the element type
	if isPtr {
		v = v.Elem()
	}

	resultPtr := reflect.New(v.Type()).Interface()
	err = gojson.Unmarshal(marshaled, resultPtr)
	if err != nil {
		t.Errorf("%s: go-json unmarshal error: %v", name, err)
		return
	}

	// If original was a pointer, compare with pointer result
	// Otherwise, dereference for comparison
	if isPtr {
		result := reflect.ValueOf(resultPtr).Elem().Interface()
		origVal := reflect.ValueOf(original).Elem().Interface()
		if !reflect.DeepEqual(origVal, result) {
			t.Errorf("%s: round-trip mismatch\noriginal: %+v\nresult: %+v", name, origVal, result)
		}
	} else {
		result := reflect.Indirect(reflect.ValueOf(resultPtr)).Interface()
		if !reflect.DeepEqual(original, result) {
			t.Errorf("%s: round-trip mismatch\noriginal: %+v\nresult: %+v", name, original, result)
		}
	}
}

// ============================================================================
// Primitive Type Tests
// ============================================================================

func TestCompatibility_Primitives(t *testing.T) {
	tests := []struct {
		name  string
		value Primitives
	}{
		{"zero values", Primitives{}},
		{"non-zero values", Primitives{
			Int: -1, Int8: -8, Int16: -16, Int32: -32, Int64: -64,
			Uint: 1, Uint8: 8, Uint16: 16, Uint32: 32, Uint64: 64,
			Float32: 3.14, Float64: 3.14159,
			Bool: true, String: "hello",
		}},
		{"max values", Primitives{
			Int: math.MaxInt, Int8: math.MaxInt8, Int16: math.MaxInt16,
			Int32: math.MaxInt32, Int64: math.MaxInt64,
			Uint: math.MaxUint, Uint8: math.MaxUint8, Uint16: math.MaxUint16,
			Uint32: math.MaxUint32, Uint64: math.MaxUint64,
			Float32: math.MaxFloat32, Float64: math.MaxFloat64,
		}},
		{"min values", Primitives{
			Int: math.MinInt, Int8: math.MinInt8, Int16: math.MinInt16,
			Int32: math.MinInt32, Int64: math.MinInt64,
		}},
		{"special floats", Primitives{
			Float32: float32(math.Inf(1)), Float64: math.Inf(1),
		}},
		{"negative zero", Primitives{
			Float32: float32(math.Copysign(0, -1)),
			Float64: math.Copysign(0, -1),
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, "Primitives", tt.value)
		})
	}
}

func TestCompatibility_PrimitivesOmitZero(t *testing.T) {
	tests := []struct {
		name  string
		value PrimitivesOmitZero
	}{
		{"all zero", PrimitivesOmitZero{}},
		{"int non-zero", PrimitivesOmitZero{Int: 42}},
		{"all non-zero", PrimitivesOmitZero{
			Int: 1, Int8: 2, Int16: 3, Int32: 4, Int64: 5,
			Uint: 6, Uint8: 7, Uint16: 8, Uint32: 9, Uint64: 10,
			Float32: 1.5, Float64: 2.5, Bool: true, String: "test",
		}},
		{"partial zero", PrimitivesOmitZero{
			Int: 0, String: "hello", Bool: false, Float64: 0,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, "PrimitivesOmitZero", tt.value)
		})
	}
}

func TestCompatibility_PrimitivesOmitEmpty(t *testing.T) {
	tests := []struct {
		name  string
		value PrimitivesOmitEmpty
	}{
		{"all zero", PrimitivesOmitEmpty{}},
		{"non-zero", PrimitivesOmitEmpty{Int: 42, String: "hello"}},
		{"empty string", PrimitivesOmitEmpty{String: ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, "PrimitivesOmitEmpty", tt.value)
		})
	}
}

func TestCompatibility_PrimitivesString(t *testing.T) {
	tests := []struct {
		name  string
		value PrimitivesString
	}{
		{"zero values", PrimitivesString{}},
		{"non-zero", PrimitivesString{Int: 42, Bool: true}},
		{"negative", PrimitivesString{Int: -100, Float64: -3.14}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, "PrimitivesString", tt.value)
		})
	}
}

// ============================================================================
// Pointer Type Tests
// ============================================================================

func TestCompatibility_Pointers(t *testing.T) {
	val := 42
	str := "test"
	b := true

	tests := []struct {
		name  string
		value Pointers
	}{
		{"all nil", Pointers{}},
		{"with values", Pointers{
			Int: &val, String: &str, Bool: &b,
		}},
		{"partial nil", Pointers{
			Int: &val, String: nil, Bool: &b,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, "Pointers", tt.value)
		})
	}
}

func TestCompatibility_PointersOmitZero(t *testing.T) {
	zeroVal := 0
	nonZeroVal := 42
	emptyStr := ""

	tests := []struct {
		name  string
		value PointersOmitZero
	}{
		{"all nil", PointersOmitZero{}},
		{"nil pointer omitted", PointersOmitZero{Int: nil}},
		{"pointer to zero included", PointersOmitZero{Int: &zeroVal}},
		{"pointer to non-zero included", PointersOmitZero{Int: &nonZeroVal}},
		{"mixed", PointersOmitZero{
			Int: &nonZeroVal, String: &emptyStr, Bool: nil,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, "PointersOmitZero", tt.value)
		})
	}
}

func TestCompatibility_MultiLevelPointers(t *testing.T) {
	val := 42

	// Test with simple pointer (more common use case)
	tests := []struct {
		name  string
		value Pointers
	}{
		{"nil", Pointers{Int: nil}},
		{"with value", Pointers{Int: &val}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, "MultiLevelPointers", tt.value)
		})
	}
}

// ============================================================================
// Collection Tests
// ============================================================================

func TestCompatibility_Collections(t *testing.T) {
	slice := []int{1, 2, 3}
	mapVal := map[string]int{"a": 1, "b": 2}
	arr := [3]int{10, 20, 30}

	tests := []struct {
		name  string
		value Collections
	}{
		{"empty", Collections{}},
		{"with values", Collections{
			Slice: slice, Map: mapVal, Array: arr,
		}},
		{"nil slice/map", Collections{
			Slice: nil, Map: nil,
		}},
		{"empty slice/map", Collections{
			Slice: []int{}, Map: map[string]int{},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, "Collections", tt.value)
		})
	}
}

func TestCompatibility_CollectionsOmitZero(t *testing.T) {
	emptySlice := []int{}
	nonEmptySlice := []int{1, 2, 3}
	emptyMap := map[string]int{}
	nonEmptyMap := map[string]int{"a": 1}

	tests := []struct {
		name  string
		value CollectionsOmitZero
	}{
		{"all nil", CollectionsOmitZero{}},
		{"nil omitted", CollectionsOmitZero{Slice: nil, Map: nil}},
		{"empty included (omitzero)", CollectionsOmitZero{
			Slice: emptySlice, Map: emptyMap,
		}},
		{"non-empty included", CollectionsOmitZero{
			Slice: nonEmptySlice, Map: nonEmptyMap,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, "CollectionsOmitZero", tt.value)
		})
	}
}

func TestCompatibility_CollectionsOmitEmpty(t *testing.T) {
	emptySlice := []int{}
	nonEmptySlice := []int{1, 2, 3}
	emptyMap := map[string]int{}
	nonEmptyMap := map[string]int{"a": 1}

	tests := []struct {
		name  string
		value CollectionsOmitEmpty
	}{
		{"all nil", CollectionsOmitEmpty{}},
		{"nil omitted (omitempty)", CollectionsOmitEmpty{
			Slice: nil, Map: nil,
		}},
		{"empty omitted (omitempty)", CollectionsOmitEmpty{
			Slice: emptySlice, Map: emptyMap,
		}},
		{"non-empty included", CollectionsOmitEmpty{
			Slice: nonEmptySlice, Map: nonEmptyMap,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, "CollectionsOmitEmpty", tt.value)
		})
	}
}

// ============================================================================
// Nested Struct Tests
// ============================================================================

func TestCompatibility_NestedStructs(t *testing.T) {
	addr := Address{
		Street: "123 Main St", City: "Springfield", State: "IL",
		ZipCode: "62701", Country: "USA",
	}
	person := Person{
		Name: "John Doe", Age: 30, Email: "john@example.com", Address: addr,
	}

	tests := []struct {
		name  string
		value interface{}
	}{
		{"Address", addr},
		{"Person", person},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, tt.name, tt.value)
		})
	}
}

func TestCompatibility_DeepNested(t *testing.T) {
	l5 := Level5{
		Level4: Level4{
			Level3: Level3{
				Level2: Level2{
					Level1: Level1{Name: "Deep"},
					Depth:  2,
				},
				Depth: 3,
			},
			Depth: 4,
		},
		Depth: 5,
		Value: "test",
	}

	compareMarshal(t, "DeepNested", l5)
	compareRoundTrip(t, "DeepNested", l5)
}

func TestCompatibility_Embedded(t *testing.T) {
	user := EmbeddedUser{
		BaseInfo: BaseInfo{
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			Version:   1,
		},
		ID:       100,
		Username: "jdoe",
	}

	compareMarshal(t, "EmbeddedUser", user)
	compareRoundTrip(t, "EmbeddedUser", user)
}

func TestCompatibility_AnonymousFields(t *testing.T) {
	outer := AnonymousOuter{
		AnonymousInner: AnonymousInner{InnerValue: "inner"},
		OuterValue:     "outer",
	}

	compareMarshal(t, "AnonymousOuter", outer)
	compareRoundTrip(t, "AnonymousOuter", outer)
}

func TestCompatibility_AnonymousTagged(t *testing.T) {
	tagged := AnonymousTagged{
		Inner: "value",
	}

	compareMarshal(t, "AnonymousTagged", tagged)
}

func TestCompatibility_CompanyEmployee(t *testing.T) {
	company := Company{
		Name: "Acme Corp",
		Employees: []Employee{
			{
				Person: Person{
					Name:  "Alice",
					Age:   25,
					Email: "alice@acme.com",
					Address: Address{
						Street: "456 Oak Ave", City: "Chicago", State: "IL",
						ZipCode: "60601", Country: "USA",
					},
				},
				Contact: Contact{
					Phone: "555-1234", Fax: "555-5678",
				},
				Department: "Engineering",
				Salary:     75000.50,
				HireDate:   time.Date(2022, 3, 15, 0, 0, 0, 0, time.UTC),
			},
			{
				Person: Person{
					Name:  "Bob",
					Age:   35,
					Email: "bob@acme.com",
					Address: Address{
						Street: "789 Pine Rd", City: "Detroit", State: "MI",
						ZipCode: "48201", Country: "USA",
					},
				},
				Contact: Contact{
					Phone: "555-9999", Fax: "555-8888",
				},
				Department: "Sales",
				Salary:     65000.00,
				HireDate:   time.Date(2021, 7, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		Location: Address{
			Street: "100 Corporate Blvd", City: "Chicago", State: "IL",
			ZipCode: "60602", Country: "USA",
		},
	}

	compareMarshal(t, "Company", company)
	compareRoundTrip(t, "Company", company)
}

// ============================================================================
// Interface Tests
// ============================================================================

func TestCompatibility_Interface(t *testing.T) {
	tests := []struct {
		name  string
		value InterfaceStruct
	}{
		{"nil", InterfaceStruct{Value: nil}},
		{"int", InterfaceStruct{Value: 42}},
		{"string", InterfaceStruct{Value: "hello"}},
		{"bool", InterfaceStruct{Value: true}},
		{"slice", InterfaceStruct{Value: []int{1, 2, 3}}},
		{"map", InterfaceStruct{Value: map[string]int{"a": 1}}},
		{"nested object", InterfaceStruct{Value: map[string]interface{}{
			"name": "test", "value": 123,
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, tt.name, tt.value)
		})
	}
}

func TestCompatibility_InterfaceWithType(t *testing.T) {
	val := 42.5
	tests := []struct {
		name  string
		value InterfaceWithType
	}{
		{"all nil", InterfaceWithType{
			Number: nil, String: nil, Bool: nil, Slice: nil, Map: nil, Null: nil,
		}},
		{"typed values", InterfaceWithType{
			Number: val, String: "test", Bool: true,
			Slice: []interface{}{1, "two", 3.0}, Map: map[string]interface{}{"key": "value"},
			Null: nil,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, tt.name, tt.value)
		})
	}
}

// ============================================================================
// Special Types Tests
// ============================================================================

func TestCompatibility_Time(t *testing.T) {
	now := time.Now()
	updated := now.Add(time.Hour)

	tests := []struct {
		name  string
		value TimeHolder
	}{
		{"with time", TimeHolder{Created: now}},
		{"with pointer time", TimeHolder{Created: now, Updated: &updated}},
		{"zero time", TimeHolder{Created: time.Time{}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, tt.name, tt.value)
		})
	}
}

func TestCompatibility_RawMessage(t *testing.T) {
	tests := []struct {
		name  string
		value RawMessageHolder
	}{
		{"empty", RawMessageHolder{Data: json.RawMessage("")}},
		{"with data", RawMessageHolder{Data: json.RawMessage(`{"key":"value"}`)}},
		{"array", RawMessageHolder{Data: json.RawMessage(`[1,2,3]`)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, tt.name, tt.value)
		})
	}
}

func TestCompatibility_Number(t *testing.T) {
	tests := []struct {
		name  string
		value NumberHolder
	}{
		{"integer string", NumberHolder{Num: json.Number("42")}},
		{"float string", NumberHolder{Num: json.Number("3.14")}},
		{"negative", NumberHolder{Num: json.Number("-100")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, tt.name, tt.value)
		})
	}
}

// ============================================================================
// OmitZero Nested Tests
// ============================================================================

func TestCompatibility_OmitZeroNested(t *testing.T) {
	tests := []struct {
		name  string
		value OmitZeroNested
	}{
		{"all zero", OmitZeroNested{}},
		{"with values", OmitZeroNested{
			Name: "John", Age: 30,
			Address: Address{Street: "123 Main", City: "Springfield"},
		}},
		{"partial zero", OmitZeroNested{
			Name: "", Age: 25,
			Address: Address{},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, tt.name, tt.value)
		})
	}
}

func TestCompatibility_OmitZeroCustomIsZero(t *testing.T) {
	tests := []struct {
		name  string
		value OmitZeroCustom
	}{
		{"zero custom", OmitZeroCustom{
			Name: "", Threshold: CustomZero{Threshold: 0.0001},
		}},
		{"non-zero custom", OmitZeroCustom{
			Name: "test", Threshold: CustomZero{Threshold: 0.5},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, tt.name, tt.value)
		})
	}
}

// ============================================================================
// Custom Marshaler Tests
// ============================================================================

func TestCompatibility_CustomMarshaler(t *testing.T) {
	items := []CustomMarshal{{Value: 1}, {Value: 2}, {Value: 3}}
	tests := []struct {
		name  string
		value WithCustomMarshal
	}{
		{"empty", WithCustomMarshal{Items: []CustomMarshal{}}},
		{"with items", WithCustomMarshal{Items: items}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, tt.name, tt.value)
		})
	}
}

func TestCompatibility_TextMarshaler(t *testing.T) {
	val := TextMarshalValue{Val: "test"}
	withText := WithTextMarshal{
		Name:  "wrapper",
		Value: val,
	}

	compareMarshal(t, "TextMarshaler", withText)
}

// ============================================================================
// Ignored and Special Field Tests
// ============================================================================

func TestCompatibility_IgnoredField(t *testing.T) {
	value := WithIgnored{
		Public:  "visible",
		Private: "hidden",
	}

	// Private should not appear in output
	stdJSON, _ := json.Marshal(value)
	goJSON, _ := gojson.Marshal(value)

	if strings.Contains(string(stdJSON), "hidden") {
		t.Error("stdlib: private field should not be serialized")
	}
	if strings.Contains(string(goJSON), "hidden") {
		t.Error("go-json: private field should not be serialized")
	}
}

func TestCompatibility_SpecialFieldNames(t *testing.T) {
	value := SpecialFieldNames{
		Dash:    "dash value",
		Empty:   "not empty",
		WithTag: "camel case",
	}

	compareMarshal(t, "SpecialFieldNames", value)
}

// ============================================================================
// Recursive Structure Tests
// ============================================================================

func TestCompatibility_Recursive(t *testing.T) {
	// Build a small tree
	root := &Node{
		Value: 1,
		Left: &Node{
			Value: 2,
			Left:  &Node{Value: 4},
			Right: &Node{Value: 5},
		},
		Right: &Node{
			Value: 3,
		},
	}

	compareMarshal(t, "Recursive", root)
	compareRoundTrip(t, "Recursive", root)
}

func TestCompatibility_MutuallyRecursive(t *testing.T) {
	tree1 := Tree{Value: 1}
	tree2 := Tree{Value: 2, Left: &tree1}
	tree3 := Tree{Value: 3, Left: &tree2}

	forest := Forest{
		Trees: []Tree{tree1, tree2, tree3},
	}

	compareMarshal(t, "MutuallyRecursive", forest)
	compareRoundTrip(t, "MutuallyRecursive", forest)
}

// ============================================================================
// Encoder/Decoder Tests
// ============================================================================

func TestCompatibility_EncoderOptions(t *testing.T) {
	value := Primitives{Int: 42, String: "<test>&"}

	// Test with HTMLEscape (default)
	stdBuf := &bytes.Buffer{}
	goBuf := &bytes.Buffer{}

	stdEnc := json.NewEncoder(stdBuf)
	goEnc := gojson.NewEncoder(goBuf)

	stdEnc.Encode(value)
	goEnc.Encode(value)

	if stdBuf.String() != goBuf.String() {
		t.Errorf("Encoder with HTML escaping mismatch:\nstdlib: %s\ngo-json: %s",
			stdBuf.String(), goBuf.String())
	}

	// Test without HTMLEscape
	stdBuf2 := &bytes.Buffer{}
	goBuf2 := &bytes.Buffer{}

	stdEnc2 := json.NewEncoder(stdBuf2)
	goEnc2 := gojson.NewEncoder(goBuf2)
	stdEnc2.SetEscapeHTML(false)
	goEnc2.SetEscapeHTML(false)

	stdEnc2.Encode(value)
	goEnc2.Encode(value)

	if stdBuf2.String() != goBuf2.String() {
		t.Errorf("Encoder without HTML escaping mismatch:\nstdlib: %s\ngo-json: %s",
			stdBuf2.String(), goBuf2.String())
	}
}

func TestCompatibility_EncoderIndent(t *testing.T) {
	value := Person{
		Name:  "John",
		Age:   30,
		Email: "john@example.com",
		Address: Address{
			Street: "123 Main St", City: "Springfield", State: "IL",
		},
	}

	compareMarshalIndent(t, "EncoderIndent", value, "", "  ")
	compareMarshalIndent(t, "EncoderIndentCustom", value, "PREFIX:", "\t")
}

func TestCompatibility_DecoderUseNumber(t *testing.T) {
	jsonData := `{"num": 42, "float": 3.14}`

	// Test with UseNumber
	type WithNumber struct {
		Num   interface{} `json:"num"`
		Float interface{} `json:"float"`
	}

	var stdResult WithNumber
	var goResult WithNumber

	stdDec := json.NewDecoder(strings.NewReader(jsonData))
	goDec := gojson.NewDecoder(strings.NewReader(jsonData))
	stdDec.UseNumber()
	goDec.UseNumber()

	stdDec.Decode(&stdResult)
	goDec.Decode(&goResult)

	// Numbers should be json.Number in both cases
	stdNum, ok := stdResult.Num.(json.Number)
	if !ok {
		t.Error("stdlib: UseNumber should return json.Number")
	}
	goNum, ok := goResult.Num.(json.Number)
	if !ok {
		t.Error("go-json: UseNumber should return json.Number")
	}

	if stdNum.String() != goNum.String() {
		t.Errorf("Number mismatch: stdlib=%s, go-json=%s", stdNum, goNum)
	}
}

func TestCompatibility_DecoderDisallowUnknownFields(t *testing.T) {
	type Simple struct {
		Name string `json:"name"`
	}

	jsonData := `{"name": "test", "unknown": "field"}`

	var stdResult Simple
	var goResult Simple

	// Use Decoder with DisallowUnknownFields
	stdDec := json.NewDecoder(strings.NewReader(jsonData))
	goDec := gojson.NewDecoder(strings.NewReader(jsonData))
	stdDec.DisallowUnknownFields()
	goDec.DisallowUnknownFields()

	stdErr := stdDec.Decode(&stdResult)
	goErr := goDec.Decode(&goResult)

	// Both should error on unknown fields
	if stdErr == nil {
		t.Error("stdlib: should error on unknown field")
	}
	if goErr == nil {
		t.Error("go-json: should error on unknown field")
	}

	// Both errors should be of similar type (UnknownFieldError)
	if stdErr != nil && goErr != nil {
		t.Logf("stdlib error: %v", stdErr)
		t.Logf("go-json error: %v", goErr)
	}
}

// ============================================================================
// Compact/Indent/Valid Tests
// ============================================================================

func TestCompatibility_Compact(t *testing.T) {
	indented := `{
  "name": "test",
  "value": 42
}`
	expected := `{"name":"test","value":42}`

	stdBuf := &bytes.Buffer{}
	goBuf := &bytes.Buffer{}

	json.Compact(stdBuf, []byte(indented))
	gojson.Compact(goBuf, []byte(indented))

	if stdBuf.String() != goBuf.String() {
		t.Errorf("Compact mismatch:\nstdlib: %s\ngo-json: %s", stdBuf.String(), goBuf.String())
	}

	// Also verify it produces the expected output
	if goBuf.String() != expected {
		t.Errorf("Compact produced wrong output: got %s, want %s", goBuf.String(), expected)
	}
}

func TestCompatibility_Indent(t *testing.T) {
	compact := `{"name":"test","value":42}`

	stdBuf := &bytes.Buffer{}
	goBuf := &bytes.Buffer{}

	json.Indent(stdBuf, []byte(compact), "", "  ")
	gojson.Indent(goBuf, []byte(compact), "", "  ")

	if stdBuf.String() != goBuf.String() {
		t.Errorf("Indent mismatch:\nstdlib: %s\ngo-json: %s", stdBuf.String(), goBuf.String())
	}
}

func TestCompatibility_Valid(t *testing.T) {
	tests := []struct {
		data string
		want bool
	}{
		{"{}", true},
		{"[]", true},
		{"null", true},
		{"true", true},
		{"42", true},
		{`"string"`, true},
		{"{invalid}", false},
		{"[1,2,", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			got := gojson.Valid([]byte(tt.data))
			if got != tt.want {
				t.Errorf("Valid(%q) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

// ============================================================================
// BigInt/BigFloat Tests
// ============================================================================

func TestCompatibility_BigInt(t *testing.T) {
	type WithBigInt struct {
		Val *big.Int `json:"val"`
	}

	tests := []struct {
		name  string
		value WithBigInt
	}{
		{"nil", WithBigInt{Val: nil}},
		{"zero", WithBigInt{Val: big.NewInt(0)}},
		{"large", WithBigInt{Val: new(big.Int).Exp(big.NewInt(2), big.NewInt(100), nil)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, tt.name, tt.value)
		})
	}
}

func TestCompatibility_BigFloat(t *testing.T) {
	type WithBigFloat struct {
		Val *big.Float `json:"val"`
	}

	tests := []struct {
		name  string
		value WithBigFloat
	}{
		{"nil", WithBigFloat{Val: nil}},
		{"zero", WithBigFloat{Val: big.NewFloat(0)}},
		{"pi", WithBigFloat{Val: big.NewFloat(math.Pi)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareMarshal(t, tt.name, tt.value)
		})
	}
}

// ============================================================================
// Complex Nested JSON Tests
// ============================================================================

func TestCompatibility_ComplexNested(t *testing.T) {
	// Simulate a realistic API response structure
	type OrderItem struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}

	type ShippingAddress struct {
		Street     string `json:"street"`
		City       string `json:"city"`
		PostalCode string `json:"postal_code"`
		Country    string `json:"country"`
	}

	type Order struct {
		ID            int             `json:"id"`
		CustomerName  string          `json:"customer_name"`
		Email         string          `json:"email"`
		Items         []OrderItem     `json:"items"`
		ShippingAddr  ShippingAddress `json:"shipping_address"`
		TotalPrice    float64         `json:"total_price"`
		Status        string          `json:"status"`
		CreatedAt     time.Time       `json:"created_at"`
		Discount      *float64        `json:"discount,omitempty"`
		Notes         string          `json:"notes,omitempty"`
	}

	order := Order{
		ID:           12345,
		CustomerName: "Jane Smith",
		Email:        "jane@example.com",
		Items: []OrderItem{
			{ID: 1, Name: "Widget A", Price: 19.99, Quantity: 2},
			{ID: 2, Name: "Widget B", Price: 29.99, Quantity: 1},
			{ID: 3, Name: "Gadget X", Price: 99.50, Quantity: 1},
		},
		ShippingAddr: ShippingAddress{
			Street:     "456 Oak Avenue",
			City:       "San Francisco",
			PostalCode: "94102",
			Country:    "USA",
		},
		TotalPrice: 169.47,
		Status:     "confirmed",
		CreatedAt:  time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		Notes:      "Please leave at door",
	}

	compareMarshal(t, "ComplexOrder", order)
	compareRoundTrip(t, "ComplexOrder", order)
}

// ============================================================================
// Unicode and Special Characters Tests
// ============================================================================

type UnicodeFields struct {
	Normal    string `json:"normal"`
	Japanese  string `json:"japanese"`
	Chinese   string `json:"chinese"`
	Emoji     string `json:"emoji"`
	Special   string `json:"special_chars"`
}

func TestCompatibility_Unicode(t *testing.T) {
	value := UnicodeFields{
		Normal:   "Hello World",
		Japanese: "こんにちは",
		Chinese:  "你好世界",
		Emoji:    "Hello 🌍 🎉",
		Special:  "Quotes: \"test\"\nNewline\nTab:\tBackslash:\\",
	}

	compareMarshal(t, "Unicode", value)
	compareRoundTrip(t, "Unicode", value)
}

// ============================================================================
// Run All Tests
// ============================================================================

func TestCompatibility_AllTypes(t *testing.T) {
	// Comprehensive test with all type categories
	// Note: Using distinct field names to avoid conflicts with embedded structs
	type Comprehensive struct {
		Prim       Primitives         `json:"prim"`
		PrimOmitZ  PrimitivesOmitZero `json:"prim_omit_z"`
		Cols       Collections        `json:"cols"`
		Addr       Address            `json:"addr"`
		Interface  InterfaceStruct    `json:"iface"`
	}

	comp := Comprehensive{
		Prim: Primitives{
			Int: 42, String: "test", Bool: true, Float64: 3.14,
		},
		PrimOmitZ: PrimitivesOmitZero{
			Int: 0, String: "", Bool: false,
		},
		Cols: Collections{
			Slice: []int{1, 2, 3}, Map: map[string]int{"a": 1},
		},
		Addr: Address{
			Street: "123 Main", City: "Springfield",
		},
		Interface: InterfaceStruct{Value: "nested interface"},
	}

	compareMarshal(t, "Comprehensive", comp)
	compareRoundTrip(t, "Comprehensive", comp)
}
