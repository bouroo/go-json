package json_test

import (
	"context"
	"testing"
	"unsafe"

	"github.com/goccy/go-json/internal/encoder"
)

func TestRuntimeContextReset(t *testing.T) {
	t.Run("reset clears all fields", func(t *testing.T) {
		ctx := encoder.TakeRuntimeContext()

		ctx.Context = context.Background()
		ctx.Buf = []byte("test")
		ctx.MarshalBuf = []byte("marshal")
		ctx.Ptrs = make([]uintptr, 10)
		ctx.KeepRefs = make([]unsafe.Pointer, 5)
		ctx.SeenPtr = make([]uintptr, 3)
		ctx.BaseIndent = 5
		ctx.Prefix = []byte("prefix")
		ctx.IndentStr = []byte("indent")

		ctx.Reset()

		if ctx.Context != nil {
			t.Error("Context should be nil after reset")
		}
		if len(ctx.Buf) != 0 {
			t.Errorf("Buf length should be 0 after reset, got %d", len(ctx.Buf))
		}
		if len(ctx.MarshalBuf) != 0 {
			t.Errorf("MarshalBuf length should be 0 after reset, got %d", len(ctx.MarshalBuf))
		}
		if len(ctx.Ptrs) != 0 {
			t.Errorf("Ptrs length should be 0 after reset, got %d", len(ctx.Ptrs))
		}
		if len(ctx.KeepRefs) != 0 {
			t.Errorf("KeepRefs length should be 0 after reset, got %d", len(ctx.KeepRefs))
		}
		if len(ctx.SeenPtr) != 0 {
			t.Errorf("SeenPtr length should be 0 after reset, got %d", len(ctx.SeenPtr))
		}
		if ctx.BaseIndent != 0 {
			t.Errorf("BaseIndent should be 0 after reset, got %d", ctx.BaseIndent)
		}
		if len(ctx.Prefix) != 0 {
			t.Errorf("Prefix length should be 0 after reset, got %d", len(ctx.Prefix))
		}
		if len(ctx.IndentStr) != 0 {
			t.Errorf("IndentStr length should be 0 after reset, got %d", len(ctx.IndentStr))
		}

		encoder.ReleaseRuntimeContext(ctx)
	})

	t.Run("release calls reset", func(t *testing.T) {
		ctx := encoder.TakeRuntimeContext()

		ctx.Buf = []byte("test")
		ctx.MarshalBuf = []byte("marshal")
		ctx.Ptrs = make([]uintptr, 10)

		encoder.ReleaseRuntimeContext(ctx)

		ctx2 := encoder.TakeRuntimeContext()

		if len(ctx2.Buf) != 0 {
			t.Errorf("Buf length should be 0 after release, got %d", len(ctx2.Buf))
		}
		if len(ctx2.MarshalBuf) != 0 {
			t.Errorf("MarshalBuf length should be 0 after release, got %d", len(ctx2.MarshalBuf))
		}
		if len(ctx2.Ptrs) != 0 {
			t.Errorf("Ptrs length should be 0 after release, got %d", len(ctx2.Ptrs))
		}

		encoder.ReleaseRuntimeContext(ctx2)
	})

	t.Run("multiple marshal cycles don't leak memory", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			ctx := encoder.TakeRuntimeContext()
			ctx.Buf = []byte("test data")
			ctx.MarshalBuf = []byte("marshal data")
			ctx.KeepRefs = make([]unsafe.Pointer, 5)
			encoder.ReleaseRuntimeContext(ctx)
		}

		ctx := encoder.TakeRuntimeContext()
		if len(ctx.Buf) != 0 {
			t.Errorf("Buf should be empty, got length %d", len(ctx.Buf))
		}
		if len(ctx.MarshalBuf) != 0 {
			t.Errorf("MarshalBuf should be empty, got length %d", len(ctx.MarshalBuf))
		}
		encoder.ReleaseRuntimeContext(ctx)
	})
}
