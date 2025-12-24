package json_test

import (
	"testing"

	"github.com/goccy/go-json/internal/encoder"
)

func TestMapContextReset(t *testing.T) {
	t.Run("reset clears all fields", func(t *testing.T) {
		ctx := encoder.NewMapContext(10, false)

		ctx.Start = 5
		ctx.First = 3
		ctx.Idx = 7
		ctx.Buf = []byte("test buffer")
		ctx.Len = 10

		ctx.Reset()

		if ctx.Start != 0 {
			t.Errorf("Start should be 0 after reset, got %d", ctx.Start)
		}
		if ctx.First != 0 {
			t.Errorf("First should be 0 after reset, got %d", ctx.First)
		}
		if ctx.Idx != 0 {
			t.Errorf("Idx should be 0 after reset, got %d", ctx.Idx)
		}
		if len(ctx.Buf) != 0 {
			t.Errorf("Buf length should be 0 after reset, got %d", len(ctx.Buf))
		}
		if ctx.Len != 0 {
			t.Errorf("Len should be 0 after reset, got %d", ctx.Len)
		}

		encoder.ReleaseMapContext(ctx)
	})

	t.Run("release calls reset", func(t *testing.T) {
		ctx := encoder.NewMapContext(10, false)
		ctx.Buf = []byte("test buffer")
		ctx.Start = 5
		ctx.Idx = 7

		encoder.ReleaseMapContext(ctx)

		ctx2 := encoder.NewMapContext(10, false)

		if len(ctx2.Buf) != 0 {
			t.Errorf("Buf length should be 0 after release, got %d", len(ctx2.Buf))
		}
		if ctx2.Start != 0 {
			t.Errorf("Start should be 0 after release, got %d", ctx2.Start)
		}
		if ctx2.Idx != 0 {
			t.Errorf("Idx should be 0 after release, got %d", ctx2.Idx)
		}

		encoder.ReleaseMapContext(ctx2)
	})

	t.Run("multiple map encoding cycles don't leak memory", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			ctx := encoder.NewMapContext(10, false)
			ctx.Buf = []byte("test data for map encoding")
			ctx.Start = 5
			ctx.Idx = 7
			encoder.ReleaseMapContext(ctx)
		}

		ctx := encoder.NewMapContext(10, false)
		if len(ctx.Buf) != 0 {
			t.Errorf("Buf should be empty, got length %d", len(ctx.Buf))
		}
		if ctx.Start != 0 {
			t.Errorf("Start should be 0, got %d", ctx.Start)
		}
		encoder.ReleaseMapContext(ctx)
	})

	t.Run("unordered map context reset", func(t *testing.T) {
		ctx := encoder.NewMapContext(10, true)
		ctx.Buf = []byte("unordered buffer")

		ctx.Reset()

		if len(ctx.Buf) != 0 {
			t.Errorf("Buf length should be 0 after reset, got %d", len(ctx.Buf))
		}

		encoder.ReleaseMapContext(ctx)
	})
}
