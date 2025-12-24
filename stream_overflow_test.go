package json_test

import (
	"io"
	"strings"
	"testing"

	"github.com/goccy/go-json"
)

func TestStreamBufferOverflow(t *testing.T) {
	t.Run("large nested arrays don't cause panic", func(t *testing.T) {
		largeJSON := strings.Repeat("[", 10000)
		for i := 0; i < 10000; i++ {
			largeJSON += "null,"
		}
		largeJSON = largeJSON[:len(largeJSON)-1] + strings.Repeat("]", 10000)

		var v interface{}
		err := json.Unmarshal([]byte(largeJSON), &v)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("decoder with large input", func(t *testing.T) {
		largeJSON := strings.Repeat("[", 5000)
		for i := 0; i < 5000; i++ {
			largeJSON += "null,"
		}
		largeJSON = largeJSON[:len(largeJSON)-1] + strings.Repeat("]", 5000)

		var v interface{}
		decoder := json.NewDecoder(strings.NewReader(largeJSON))
		err := decoder.Decode(&v)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("large input doesn't overflow buffer", func(t *testing.T) {
		largeInput := strings.Repeat("null,", 50000)
		largeInput = "[" + largeInput[:len(largeInput)-1] + "]"

		var v interface{}
		err := json.Unmarshal([]byte(largeInput), &v)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})
}

type limitReader struct {
	r io.Reader
	n int64
}

func (l *limitReader) Read(p []byte) (n int, err error) {
	if l.n <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > l.n {
		p = p[:l.n]
	}
	n, err = l.r.Read(p)
	l.n -= int64(n)
	return
}

func TestStreamBufferLimit(t *testing.T) {
	t.Run("buffer growth stays within limit", func(t *testing.T) {
		largeStr := strings.Repeat("null,", 100000)
		largeStr = "[" + largeStr[:len(largeStr)-1] + "]"
		lr := &limitReader{r: strings.NewReader(largeStr), n: int64(len(largeStr))}
		decoder := json.NewDecoder(lr)
		var result interface{}
		err := decoder.Decode(&result)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})
}
