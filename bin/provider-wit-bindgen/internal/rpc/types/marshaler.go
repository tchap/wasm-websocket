package types

import msgpack "github.com/wasmcloud/tinygo-msgpack"

type MMarshaler interface {
	MEncode(w msgpack.Writer) error
}

type MUnmarshaler interface {
	MDecode(d *msgpack.Decoder) error
}

type MMsg[T any] interface {
	*T
	MMarshaler
	MUnmarshaler
}
