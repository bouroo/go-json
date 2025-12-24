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

// CacheIsZeroMethod caches the IsZero() method for optimized calls.
// Returns the reflect.Value of the method's function and whether it requires a pointer receiver.
// This should be called during compilation to cache the method for runtime use.
func CacheIsZeroMethod(typ *runtime.Type) (interface{}, bool) {
	// First check if the type itself has IsZero() method
	method, found := typ.MethodByName("IsZero")
	if found {
		// Verify signature: no parameters, one bool return
		if method.Type.NumIn() == 1 && // receiver only
			method.Type.NumOut() == 1 &&
			method.Type.Out(0).Kind() == reflect.Bool {
			return method.Func, false
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
				return ptrMethod.Func, true
			}
		}
	}
	return nil, false
}

// CallIsZeroMethod calls the IsZero() method on a value.
// Used for omitzero tag support when a custom IsZero() method is present.
// The caller must have verified that the type has an IsZero() method
// (via Opcode.HasIsZeroMethod) before calling this function.
func CallIsZeroMethod(code *Opcode, ptr uintptr) bool {
	if code.Type == nil || ptr == 0 {
		return false
	}

	// Fast path: use cached method function if available
	if code.IsZeroMethodFunc != nil {
		// Convert the runtime.Type pointer value to an interface{}
		v := PtrToInterface(code, ptr)
		reflectValue := reflect.ValueOf(v)
		// Type assert to get the reflect.Value back
		method, ok := code.IsZeroMethodFunc.(reflect.Value)
		if !ok || !method.IsValid() {
			return false
		}
		results := method.Call([]reflect.Value{reflectValue})
		if len(results) > 0 && results[0].Kind() == reflect.Bool {
			return results[0].Bool()
		}
		return false
	}

	// Fallback to reflection for cases where caching wasn't done
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

// isStructZeroByReflection checks if a struct value is zero-valued using reflection.
// This function assumes the type does NOT have a custom IsZero() method.
// The caller must verify this before calling this function.
func isStructZeroByReflection(typ *runtime.Type, ptr uintptr) bool {
	if typ == nil || ptr == 0 {
		return false
	}

	code := &Opcode{Type: typ}
	v := PtrToInterface(code, ptr)
	reflectValue := reflect.ValueOf(v)
	return reflectValue.IsZero()
}

// IsStructZero checks if a struct value is zero-valued.
// Uses custom IsZero() method if available, otherwise falls back to reflect.Value.IsZero().
func IsStructZero(typ *runtime.Type, ptr uintptr) bool {
	if typ == nil || ptr == 0 {
		return false
	}

	// Check if the type has a custom IsZero method (for external callers)
	if HasIsZeroMethod(typ) {
		code := &Opcode{Type: typ}
		return CallIsZeroMethod(code, ptr)
	}

	// Use pure reflection path
	return isStructZeroByReflection(typ, ptr)
}
