package vm_indent

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/goccy/go-json/internal/encoder"
	"github.com/goccy/go-json/internal/runtime"
)

const uintptrSize = 4 << (^uintptr(0) >> 63)

var (
	appendInt           = encoder.AppendInt
	appendUint          = encoder.AppendUint
	appendFloat32       = encoder.AppendFloat32
	appendFloat64       = encoder.AppendFloat64
	appendString        = encoder.AppendString
	appendByteSlice     = encoder.AppendByteSlice
	appendNumber        = encoder.AppendNumber
	appendStructEnd     = encoder.AppendStructEndIndent
	appendIndent        = encoder.AppendIndent
	errUnsupportedValue = encoder.ErrUnsupportedValue
	errUnsupportedFloat = encoder.ErrUnsupportedFloat
	mapiterinit         = encoder.MapIterInit
	mapiterkey          = encoder.MapIterKey
	mapitervalue        = encoder.MapIterValue
	mapiternext         = encoder.MapIterNext
	maplen              = encoder.MapLen
)

type emptyInterface struct {
	typ *runtime.Type
	ptr unsafe.Pointer
}

type nonEmptyInterface struct {
	itab *struct {
		ityp *runtime.Type // static interface type
		typ  *runtime.Type // dynamic concrete type
		// unused fields...
	}
	ptr unsafe.Pointer
}

func errUnimplementedOp(op encoder.OpType) error {
	return fmt.Errorf("encoder (indent): opcode %s has not been implemented", op)
}

func load(base uintptr, idx uint32) uintptr {
	addr := base + uintptr(idx)
	return **(**uintptr)(unsafe.Pointer(&addr))
}

func store(base uintptr, idx uint32, p uintptr) {
	addr := base + uintptr(idx)
	**(**uintptr)(unsafe.Pointer(&addr)) = p
}

func loadNPtr(base uintptr, idx uint32, ptrNum uint8) uintptr {
	addr := base + uintptr(idx)
	p := **(**uintptr)(unsafe.Pointer(&addr))
	for i := uint8(0); i < ptrNum; i++ {
		if p == 0 {
			return 0
		}
		p = ptrToPtr(p)
	}
	return p
}

func ptrToUint64(p uintptr, bitSize uint8) uint64 {
	switch bitSize {
	case 8:
		return (uint64)(**(**uint8)(unsafe.Pointer(&p)))
	case 16:
		return (uint64)(**(**uint16)(unsafe.Pointer(&p)))
	case 32:
		return (uint64)(**(**uint32)(unsafe.Pointer(&p)))
	case 64:
		return **(**uint64)(unsafe.Pointer(&p))
	}
	return 0
}
func ptrToFloat32(p uintptr) float32            { return **(**float32)(unsafe.Pointer(&p)) }
func ptrToFloat64(p uintptr) float64            { return **(**float64)(unsafe.Pointer(&p)) }
func ptrToBool(p uintptr) bool                  { return **(**bool)(unsafe.Pointer(&p)) }
func ptrToBytes(p uintptr) []byte               { return **(**[]byte)(unsafe.Pointer(&p)) }
func ptrToNumber(p uintptr) json.Number         { return **(**json.Number)(unsafe.Pointer(&p)) }
func ptrToString(p uintptr) string              { return **(**string)(unsafe.Pointer(&p)) }
func ptrToSlice(p uintptr) *runtime.SliceHeader { return *(**runtime.SliceHeader)(unsafe.Pointer(&p)) }
func ptrToPtr(p uintptr) uintptr {
	return uintptr(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
}
func ptrToNPtr(p uintptr, ptrNum uint8) uintptr {
	for i := uint8(0); i < ptrNum; i++ {
		if p == 0 {
			return 0
		}
		p = ptrToPtr(p)
	}
	return p
}

func ptrToUnsafePtr(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}
func ptrToInterface(code *encoder.Opcode, p uintptr) interface{} {
	return *(*interface{})(unsafe.Pointer(&emptyInterface{
		typ: code.Type,
		ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
	}))
}

func appendBool(_ *encoder.RuntimeContext, b []byte, v bool) []byte {
	if v {
		return append(b, "true"...)
	}
	return append(b, "false"...)
}

func appendNull(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, "null"...)
}

func appendComma(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, ',', '\n')
}

func appendNullComma(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, "null,\n"...)
}

func appendColon(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b[:len(b)-2], ':', ' ')
}

func appendMapKeyValue(ctx *encoder.RuntimeContext, code *encoder.Opcode, b, key, value []byte) []byte {
	b = appendIndent(ctx, b, code.Indent+1)
	b = append(b, key...)
	b[len(b)-2] = ':'
	b[len(b)-1] = ' '
	return append(b, value...)
}

func appendMapEnd(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	b = b[:len(b)-2]
	b = append(b, '\n')
	b = appendIndent(ctx, b, code.Indent)
	return append(b, '}', ',', '\n')
}

func appendArrayHead(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	b = append(b, '[', '\n')
	return appendIndent(ctx, b, code.Indent+1)
}

func appendArrayEnd(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	b = b[:len(b)-2]
	b = append(b, '\n')
	b = appendIndent(ctx, b, code.Indent)
	return append(b, ']', ',', '\n')
}

func appendEmptyArray(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, '[', ']', ',', '\n')
}

func appendEmptyObject(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, '{', '}', ',', '\n')
}

func appendObjectEnd(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	last := len(b) - 1
	// replace comma to newline
	b[last-1] = '\n'
	b = appendIndent(ctx, b[:last], code.Indent)
	return append(b, '}', ',', '\n')
}

func appendMarshalJSON(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte, v interface{}) ([]byte, error) {
	return encoder.AppendMarshalJSONIndent(ctx, code, b, v)
}

func appendMarshalText(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte, v interface{}) ([]byte, error) {
	return encoder.AppendMarshalTextIndent(ctx, code, b, v)
}

func appendStructHead(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, '{', '\n')
}

func appendStructKey(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	b = appendIndent(ctx, b, code.Indent)
	b = append(b, code.Key...)
	return append(b, ' ')
}

func appendStructEndSkipLast(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	last := len(b) - 1
	if b[last-1] == '{' {
		b[last] = '}'
	} else {
		if b[last] == '\n' {
			// to remove ',' and '\n' characters
			b = b[:len(b)-2]
		}
		b = append(b, '\n')
		b = appendIndent(ctx, b, code.Indent-1)
		b = append(b, '}')
	}
	return appendComma(ctx, b)
}

func restoreIndent(ctx *encoder.RuntimeContext, code *encoder.Opcode, ctxptr uintptr) {
	ctx.BaseIndent = uint32(load(ctxptr, code.Length))
}

func storeIndent(ctxptr uintptr, code *encoder.Opcode, indent uintptr) {
	store(ctxptr, code.Length, indent)
}

func appendArrayElemIndent(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	return appendIndent(ctx, b, code.Indent+1)
}

func appendMapKeyIndent(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	return appendIndent(ctx, b, code.Indent)
}

// callIsZeroMethod calls the IsZero() method on a value if it exists.
// Used for omitzero tag support when a custom IsZero() method is present.
func callIsZeroMethod(code *encoder.Opcode, ptr uintptr) bool {
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
	v := ptrToInterface(code, ptr)
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

// hasIsZeroMethod checks if the type has a custom IsZero() bool method.
func hasIsZeroMethod(typ *runtime.Type) bool {
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

// isStructZero checks if a struct value is zero-valued.
// Uses custom IsZero() method if available, otherwise falls back to reflect.Value.IsZero().
func isStructZero(typ *runtime.Type, ptr uintptr) bool {
	if typ == nil || ptr == 0 {
		return false
	}

	// Create an opcode temporarily to use ptrToInterface
	code := &encoder.Opcode{Type: typ}

	// First check if the type has a custom IsZero method
	if hasIsZeroMethod(typ) {
		return callIsZeroMethod(code, ptr)
	}

	// Fall back to reflection-based zero detection
	v := ptrToInterface(code, ptr)
	reflectValue := reflect.ValueOf(v)
	return reflectValue.IsZero()
}
