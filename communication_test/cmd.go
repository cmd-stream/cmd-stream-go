package ct

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/raw"
)

type CmdType byte

const (
	Cmd1CmdType CmdType = iota + 1
	Cmd2CmdType
	Cmd3CmdType
)

func MarshalCmdType(tp CmdType, w muss.Writer) (n int, err error) {
	return raw.MarshalByte(byte(tp), w)
}

func UnmarshalCmdType(r muss.Reader) (tp CmdType, n int, err error) {
	b, n, err := raw.UnmarshalByte(r)
	tp = CmdType(b)
	return
}

func SizeCmdType(tp CmdType) (size int) {
	return raw.SizeByte(byte(tp))
}

// -----------------------------------------------------------------------------
type Cmd1 struct{}

func (c Cmd1) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Receiver,
	proxy base.Proxy,
) (err error) {
	t := 500 * time.Millisecond
	time.Sleep(t)
	err = proxy.Send(seq, Result{false})
	if err != nil {
		return
	}
	time.Sleep(t)
	err = proxy.Send(seq, Result{true})
	if err != nil {
		return
	}
	return
}

// -----------------------------------------------------------------------------
type Cmd2 struct{}

func (c Cmd2) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Receiver,
	proxy base.Proxy,
) (err error) {
	time.Sleep(500 * time.Millisecond)
	return proxy.Send(seq, Result{true})
}

// -----------------------------------------------------------------------------
type Cmd3 struct{}

func (c Cmd3) Exec(ctx context.Context, at time.Time, seq base.Seq,
	receiver Receiver,
	proxy base.Proxy,
) (err error) {
	time.Sleep(500 * time.Millisecond)
	return proxy.Send(seq, Result{true})
}
