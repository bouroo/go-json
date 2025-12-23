package encoder

import (
	"reflect"

	"github.com/goccy/go-json/internal/runtime"
)

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
