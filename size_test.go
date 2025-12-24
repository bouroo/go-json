package json

import (
	"testing"
	"unsafe"

	"github.com/goccy/go-json/internal/encoder"
)

func TestOpcodeSize(t *testing.T) {
	const uintptrSize = 4 << (^uintptr(0) >> 63)
	if uintptrSize == 8 {
		size := unsafe.Sizeof(encoder.Opcode{})
		// Size is 152 bytes after adding IsZeroMethodFunc (interface{} = 16 bytes)
		// and IsZeroMethodNeedsPtr (bool = 1 byte, padded) to Opcode
		if size != 152 {
			t.Fatalf("unexpected opcode size: expected 152bytes but got %dbytes", size)
		}
	}
}
