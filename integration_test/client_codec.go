package intest

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	exts "github.com/mus-format/ext-mus-stream-go"
)

type ClientCodec struct{}

func (c ClientCodec) Encode(cmd base.Cmd[Receiver],
	w transport.Writer) (err error) {
	if m, ok := cmd.(exts.MarshallerTypedMUS); ok {
		_, err = m.MarshalTypedMUS(w)
		return
	}
	panic("cmd doesn't implement the ext.MarshallerTypedMUS interface")
}

func (c ClientCodec) Decode(r transport.Reader) (result base.Result,
	err error) {
	result, _, err = ResultMUS.Unmarshal(r)
	return
}

func (c ClientCodec) Size(cmd base.Cmd[Receiver]) (size int) {
	return 0
}
