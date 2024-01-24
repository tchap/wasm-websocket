package types

import msgpack "github.com/wasmcloud/tinygo-msgpack"

type Bytes struct {
	Value []byte
}

func (b Bytes) MEncode(w msgpack.Writer) error {
	w.WriteByteArray(b.Value)
	return w.CheckError()
}

func (b *Bytes) MDecode(d msgpack.Decoder) error {
	bs, err := d.ReadByteArray()
	if err != nil {
		return err
	}
	b.Value = bs
	return nil
}
