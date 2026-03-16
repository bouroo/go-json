package encoder

import (
	"reflect"
	"unsafe"

	"github.com/goccy/go-json/internal/runtime"
)

// IsZeroer is the interface used to check custom zero values.
type IsZeroer interface {
	IsZero() bool
}

// IsZeroForOmitZero determines if a value should be omitted when using the omitzero tag.
// It checks if a value is zero (default/empty for its type).
// NOTE: This function uses unsafe memory operations. Direct pointer dereferences
// are avoided to minimize issues with the race detector, though some operations
// may still require unsafe memory access for performance reasons.
// The race detector is disabled for this function because it operates on unsafe pointers
// that are valid in the encoder's context but would trigger false positives.
//
//go:nocheckptr
func IsZeroForOmitZero(code *Opcode, ptr uintptr) bool {
	if ptr == 0 {
		return true
	}

	kind := code.Type.Kind()

	// Implement zero checking for each type kind
	switch kind {
	// Pointer types: zero only if nil
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Interface:
		// Read the pointer/interface value safely
		// These types use a pointer as their first field
		p := *(*unsafe.Pointer)(unsafe.Pointer(ptr))
		return p == nil

	// Boolean type
	case reflect.Bool:
		return *(*bool)(unsafe.Pointer(ptr)) == false

	// Signed integer types
	case reflect.Int:
		return *(*int)(unsafe.Pointer(ptr)) == 0
	case reflect.Int8:
		return *(*int8)(unsafe.Pointer(ptr)) == 0
	case reflect.Int16:
		return *(*int16)(unsafe.Pointer(ptr)) == 0
	case reflect.Int32:
		return *(*int32)(unsafe.Pointer(ptr)) == 0
	case reflect.Int64:
		return *(*int64)(unsafe.Pointer(ptr)) == 0

	// Unsigned integer types
	case reflect.Uint:
		return *(*uint)(unsafe.Pointer(ptr)) == 0
	case reflect.Uint8:
		return *(*uint8)(unsafe.Pointer(ptr)) == 0
	case reflect.Uint16:
		return *(*uint16)(unsafe.Pointer(ptr)) == 0
	case reflect.Uint32:
		return *(*uint32)(unsafe.Pointer(ptr)) == 0
	case reflect.Uint64:
		return *(*uint64)(unsafe.Pointer(ptr)) == 0
	case reflect.Uintptr:
		return *(*uintptr)(unsafe.Pointer(ptr)) == 0

	// Floating point types
	case reflect.Float32:
		return *(*float32)(unsafe.Pointer(ptr)) == 0
	case reflect.Float64:
		return *(*float64)(unsafe.Pointer(ptr)) == 0

	// Complex types
	case reflect.Complex64:
		return *(*complex64)(unsafe.Pointer(ptr)) == 0
	case reflect.Complex128:
		return *(*complex128)(unsafe.Pointer(ptr)) == 0

	// String type
	case reflect.String:
		// String is three words: pointer, length, capacity (on 64-bit)
		// Length is zero for empty string
		s := *(*string)(unsafe.Pointer(ptr))
		return len(s) == 0

	// For complex types (Array, Struct), we need to check each field
	// Use reflect for these but attempt to minimize allocations
	default:
		rt := runtime.RType2Type(code.Type)
		// Note: This may trigger race detector warnings with -race,
		// but it's unavoidable for complex types
		rv := reflect.NewAt(rt, unsafe.Pointer(ptr)).Elem()
		return isZeroValue(rv)
	}
}

// isZeroValue checks if a reflect.Value is zero, including support for custom IsZero methods
func isZeroValue(rv reflect.Value) bool {
	// Check for custom IsZero implementation first
	if rv.CanInterface() {
		if iz, ok := rv.Interface().(IsZeroer); ok {
			return iz.IsZero()
		}
	}

	// Use reflect's built-in IsZero
	// This correctly handles:
	// - nil pointers (IsZero() == true)
	// - non-nil pointers to zero values (IsZero() == false, since pointer is not nil)
	// - nil slices/maps/channels (IsZero() == true)
	// - empty non-nil slices/maps (IsZero() == false, since they're not nil)
	// - zero primitives like 0, "", false (IsZero() == true)
	return rv.IsZero()
}
