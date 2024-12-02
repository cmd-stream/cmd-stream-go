package ct

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
)

type ClientCodec struct{}

func (c ClientCodec) Encode(cmd base.Cmd[Receiver],
	w transport.Writer) (err error) {
	if m, ok := cmd.(MarshallerMUS); ok {
		_, err = m.MarshalMUS(w)
		return
	}
	panic("cmd doesn't implement the MarshallerMUS interface")
}

func (c ClientCodec) Decode(r transport.Reader) (result base.Result,
	err error) {
	result, _, err = UnmarshalResultMUS(r)
	return
}

func (c ClientCodec) Size(cmd base.Cmd[Receiver]) (size int) {
	return 0
}
