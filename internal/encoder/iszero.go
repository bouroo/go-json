package encoder

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-json/internal/runtime"
)

// UintptrSize is the size of a uintptr in bytes.
const UintptrSize = 4 << (^uintptr(0) >> 63)

// ErrUnimplementedOp creates an error for an unimplemented opcode.
// The variant parameter specifies which VM variant (e.g., "encoder", "encoder (indent)").
func ErrUnimplementedOp(op OpType, variant string) error {
	if variant == "" {
		variant = "encoder"
	}
	return fmt.Errorf("%s: opcode %s has not been implemented", variant, op)
}

// HasIsZeroMethod checks if the type has a custom IsZero() bool method.
// This function is used by both the compiler and the VM packages for omitzero tag support.
func HasIsZeroMethod(typ *runtime.Type) bool {
	// Check if type has IsZero() bool method
	method, found := typ.MethodByName("IsZero")
	if found {
		// Verify signature: no parameters, one bool return
		if method.Type.NumIn() == 1 && // receiver only
			method.Type.NumOut() == 1 &&
			method.Type.Out(0).Kind() == reflect.Bool {
			return true
		}
	}
	// Check pointer type if original type is not a pointer
	if typ.Kind() != reflect.Ptr {
		ptrTyp := runtime.PtrTo(typ)
		ptrMethod, found := ptrTyp.MethodByName("IsZero")
		if found {
			if ptrMethod.Type.NumIn() == 1 && // receiver only
				ptrMethod.Type.NumOut() == 1 &&
				ptrMethod.Type.Out(0).Kind() == reflect.Bool {
				return true
			}
		}
	}
	return false
}

// CallIsZeroMethod calls the IsZero() method on a value if it exists.
// Used for omitzero tag support when a custom IsZero() method is present.
func CallIsZeroMethod(code *Opcode, ptr uintptr) bool {
	if code.Type == nil || ptr == 0 {
		return false
	}

	rtType := code.Type

	// Check if the type (or its pointer) has IsZero method
	_, found := rtType.MethodByName("IsZero")
	if !found {
		// Try pointer type if the original type doesn't have the method
		if rtType.Kind() != reflect.Ptr {
			ptrType := runtime.PtrTo(rtType)
			_, found = ptrType.MethodByName("IsZero")
			if !found {
				return false
			}
		} else {
			return false
		}
	}

	// Convert the runtime.Type pointer value to an interface{}
	// This allows us to use reflect.ValueOf to get a reflect.Value
	v := PtrToInterface(code, ptr)
	reflectValue := reflect.ValueOf(v)

	// Call the IsZero() method through reflection
	method := reflectValue.MethodByName("IsZero")
	if !method.IsValid() {
		return false
	}

	results := method.Call([]reflect.Value{})
	if len(results) > 0 && results[0].Kind() == reflect.Bool {
		return results[0].Bool()
	}
	return false
}

// IsStructZero checks if a struct value is zero-valued.
// Uses custom IsZero() method if available, otherwise falls back to reflect.Value.IsZero().
func IsStructZero(typ *runtime.Type, ptr uintptr) bool {
	if typ == nil || ptr == 0 {
		return false
	}

	// Create an opcode temporarily to use PtrToInterface
	code := &Opcode{Type: typ}

	// First check if the type has a custom IsZero method
	if HasIsZeroMethod(typ) {
		return CallIsZeroMethod(code, ptr)
	}

	// Fall back to reflection-based zero detection
	v := PtrToInterface(code, ptr)
	reflectValue := reflect.ValueOf(v)
	return reflectValue.IsZero()
}
