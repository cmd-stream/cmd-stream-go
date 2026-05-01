package testkit

import (
	"encoding/json"
	"fmt"

	"github.com/cmd-stream/cmd-stream-go/core"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
	"github.com/mus-format/mus-stream-go/ord"
	"github.com/mus-format/mus-stream-go/typed"
)

type ClientCodec struct{}

func (c ClientCodec) Encode(cmd core.Cmd[Receiver],
	w tspt.Writer,
) (n int, err error) {
	var (
		bs []byte
		n1 int
	)
	switch c := cmd.(type) {
	case Cmd:
		n, err = typed.DTMSer.Marshal(CmdDTM, w)
		if err != nil {
			return
		}
		bs, err = json.Marshal(c)
		if err != nil {
			return
		}
		n1, err = ord.ByteSlice.Marshal(bs, w)
		n += n1
		return
	case MultiCmd:
		n, err = typed.DTMSer.Marshal(MultiCmdDTM, w)
		if err != nil {
			return
		}
		bs, err = json.Marshal(c)
		if err != nil {
			return
		}
		n1, err = ord.ByteSlice.Marshal(bs, w)
		n += n1
		return
	default:
		panic(fmt.Sprintf("unknown cmd: %T", cmd))
	}
}

func (c ClientCodec) Decode(r tspt.Reader) (result core.Result, n int,
	err error,
) {
	bs, n, err := ord.ByteSlice.Unmarshal(r)
	if err != nil {
		return
	}
	res := Result{}
	err = json.Unmarshal(bs, &res)
	if err != nil {
		return
	}
	result = res
	return
}

// -----------------------------------------------------------------------------

type ServerCodec struct{}

func (c ServerCodec) Encode(result core.Result, w tspt.Writer) (n int, err error) {
	bs, err := json.Marshal(result)
	if err != nil {
		return
	}
	return ord.ByteSlice.Marshal(bs, w)
}

func (c ServerCodec) Decode(r tspt.Reader) (cmd core.Cmd[Receiver],
	n int, err error,
) {
	dtm, n, err := typed.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	var (
		bs []byte
		n1 int
	)
	switch dtm {
	case CmdDTM:
		bs, n1, err = ord.ByteSlice.Unmarshal(r)
		n += n1
		if err != nil {
			return
		}
		c := Cmd{}
		err = json.Unmarshal(bs, &c)
		if err != nil {
			return
		}
		cmd = c
		return
	case MultiCmdDTM:
		bs, n1, err = ord.ByteSlice.Unmarshal(r)
		n += n1
		if err != nil {
			return
		}
		c := MultiCmd{}
		err = json.Unmarshal(bs, &c)
		if err != nil {
			return
		}
		cmd = c
		return
	default:
		panic(fmt.Sprintf("unknown dtm: %d", dtm))
	}
}
