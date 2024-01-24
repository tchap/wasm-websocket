package websocketserver

import (
	"github.com/tchap/wasmcloud-websocket/bin/provider-wit-bindgen/rpc/types"
	provider "github.com/wasmCloud/provider-sdk-go"
	msgpack "github.com/wasmcloud/tinygo-msgpack"
)

type Types_MessageKind uint32

const (
	TypesMessageKind_Text Types_MessageKind = iota + 1
	TypesMessageKind_Binary
)

type Types_WebsocketMessage struct {
	Kind Types_MessageKind
	Body *types.Option[types.Bytes, *types.Bytes]
}

func (msg *Types_WebsocketMessage) MEncode(w msgpack.Writer) error {
	w.WriteMapSize(2)
	w.WriteString("kind")
	w.WriteUint32(uint32(msg.Kind))
	w.WriteString("body")
	msg.Body.MEncode(w)
	return w.CheckError()
}

func (msg *Types_WebsocketMessage) MDecode(d msgpack.Decoder) error {
	n, err := d.ReadMapSize()
	if err != nil {
		return err
	}
	for i := uint32(0); i < n; i++ {
		k, _ := d.ReadString()
		switch k {
		case "kind":
			v, err := d.ReadUint32()
			if err != nil {
				return err
			}
			msg.Kind = Types_MessageKind(v)

		case "body":
			var v types.Option[types.Bytes, *types.Bytes]
			if err := v.MDecode(d); err != nil {
				return err
			}
			msg.Body = &v
		}
	}
	return nil
}

type Types struct {
	p *provider.WasmcloudProvider
}

func NewTypes(p *provider.WasmcloudProvider) *Types {
	return &Types{p: p}
}

type Handler struct {
	p *provider.WasmcloudProvider
}

func NewHandler(p *provider.WasmcloudProvider) *Handler {
	return &Handler{p: p}
}

func (self *Handler) HandleMessage(
	actorID string,
	message *Types_WebsocketMessage,
) (*types.Option[Types_WebsocketMessage, *Types_WebsocketMessage], error) {
	// Encode the request.
	var sizer msgpack.Sizer
	sizeEncoder := &sizer
	message.MEncode(sizeEncoder)

	buf := make([]byte, sizer.Len())
	encoder := msgpack.NewEncoder(buf)
	message.MEncode(&encoder)

	if err := encoder.CheckError(); err != nil {
		return nil, err
	}

	// Send it.
	respBody, err := self.p.ToActor(actorID, buf, "WebSocketServer.HandleMessage")
	if err != nil {
		return nil, err
	}

	// Decode the response.
	var resp types.Option[Types_WebsocketMessage, *Types_WebsocketMessage]
	if err := resp.MDecode(msgpack.NewDecoder(respBody)); err != nil {
		return nil, err
	}

	return &resp, nil
}
