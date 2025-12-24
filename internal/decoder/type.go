package decoder

import (
	"context"
	"encoding"
	"encoding/json"
	"reflect"
	"unsafe"
)

type Decoder interface {
	Decode(ctx *RuntimeContext, cursor, length int64, p unsafe.Pointer) (int64, error)
	DecodePath(ctx *RuntimeContext, cursor, length int64) ([][]byte, int64, error)
	DecodeStream(s *Stream, cursor int64, p unsafe.Pointer) error
}

const (
	nul                   = '\000'
	maxDecodeNestingDepth = 10000
)

type unmarshalerContext interface {
	UnmarshalJSON(ctx context.Context, data []byte) error
}

var (
	unmarshalJSONType        = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	unmarshalJSONContextType = reflect.TypeOf((*unmarshalerContext)(nil)).Elem()
	unmarshalTextType        = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)
