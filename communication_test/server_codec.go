package ct

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
)

type ServerCodec struct{}

func (c ServerCodec) Encode(result base.Result, w transport.Writer) (
	err error) {
	switch r := result.(type) {
	case Result:
		_, err = MarshalResultMUS(r, w)
	default:
		panic("unexpected result")
	}
	return
}

func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[Receiver],
	err error) {
	tp, _, err := UnmarshalCmdType(r)
	switch tp {
	case Cmd1CmdType:
		cmd = Cmd1{}
	case Cmd2CmdType:
		cmd = Cmd2{}
	case Cmd3CmdType:
		cmd = Cmd3{}
	}
	return
}
