package ct

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
)

type ClientCodec struct{}

func (c ClientCodec) Encode(cmd base.Cmd[Receiver],
	w transport.Writer) (err error) {
	var tp CmdType
	switch cmd.(type) {
	case Cmd1:
		tp = Cmd1CmdType
	case Cmd2:
		tp = Cmd2CmdType
	case Cmd3:
		tp = Cmd3CmdType
	default:
		panic("unexpected cmd type")
	}
	_, err = MarshalCmdType(tp, w)
	return
}

func (c ClientCodec) Decode(r transport.Reader) (result base.Result,
	err error) {
	result, _, err = UnmarshalResultMUS(r)
	return
}

func (c ClientCodec) Size(cmd base.Cmd[Receiver]) (size int) {
	return 0
}
