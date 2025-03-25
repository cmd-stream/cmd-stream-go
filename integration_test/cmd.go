package intest

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	dts "github.com/mus-format/mus-stream-dts-go"
	muss "github.com/mus-format/mus-stream-go"
)

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

func (c Cmd1) MarshalMUS(w muss.Writer) (n int, err error) {
	return dts.DTMSer.Marshal(Cmd1DTM, w)
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

func (c Cmd2) MarshalMUS(w muss.Writer) (n int, err error) {
	return dts.DTMSer.Marshal(Cmd2DTM, w)
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

func (c Cmd3) MarshalMUS(w muss.Writer) (n int, err error) {
	return dts.DTMSer.Marshal(Cmd3DTM, w)
}
