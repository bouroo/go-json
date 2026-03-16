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
// It first checks if the value implements IsZero(), then falls back to reflect.Value.IsZero().
func IsZeroForOmitZero(code *Opcode, ptr uintptr) bool {
	if ptr == 0 {
		return true
	}

	// Check if this is a pointer type by looking at PtrNum
	// If PtrNum > 0, the field is a pointer (or multi-level pointer)
	if code.PtrNum > 0 {
		// For pointer types, check if the pointer itself is nil
		// Load the pointer value from memory
		ptrVal := *(*unsafe.Pointer)(unsafe.Pointer(ptr))
		// If the pointer is nil, it's zero
		return ptrVal == nil
	}

	// For non-pointer types, use reflect to check IsZero
	rt := runtime.RType2Type(code.Type)
	rv := reflect.NewAt(rt, unsafe.Pointer(ptr)).Elem()

	// Check for custom IsZero implementation
	if rv.CanAddr() {
		if iz, ok := rv.Addr().Interface().(IsZeroer); ok {
			return iz.IsZero()
		}
	}

	return rv.IsZero()
}
