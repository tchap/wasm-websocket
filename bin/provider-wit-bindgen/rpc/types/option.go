package types

import msgpack "github.com/wasmcloud/tinygo-msgpack"

type Option[T any, M MMsg[T]] struct {
	Value M
}

func (opt *Option[T, M]) MEncode(w msgpack.Writer) error {
	if opt.Value == nil {
		w.WriteMapSize(0)
		return w.CheckError()
	}

	w.WriteMapSize(1)
	w.WriteString("o")
	return opt.Value.MEncode(w)
}

func (opt *Option[T, M]) MDecode(d msgpack.Decoder) error {
	n, err := d.ReadMapSize()
	if err != nil {
		return err
	}

	for i := uint32(0); i < n; i++ {
		k, _ := d.ReadString()
		if k == "o" {
			var v M = new(T)
			if err := v.MDecode(d); err != nil {
				return err
			}
			opt.Value = v
		}
	}
	return nil
}
