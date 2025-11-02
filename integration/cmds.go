package integration_test

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
	com "github.com/mus-format/common-go"
	"github.com/mus-format/dts-stream-go"
	"github.com/mus-format/mus-stream-go"
)

const (
	Cmd1DTM com.DTM = iota + 1
	Cmd2DTM
	Cmd3DTM
	Cmd4DTM
)

type Cmd1 struct{}

func (c Cmd1) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver struct{},
	proxy core.Proxy,
) (err error) {
	t := 500 * time.Millisecond
	time.Sleep(t)
	_, err = proxy.Send(seq, NewResult(false))
	if err != nil {
		return
	}
	time.Sleep(t)
	_, err = proxy.Send(seq, NewResult(true))
	if err != nil {
		return
	}
	return
}

func (c Cmd1) MarshalTypedMUS(w mus.Writer) (n int, err error) {
	return dts.DTMSer.Marshal(Cmd1DTM, w)
}

func (c Cmd1) SizeTypedMUS() (size int) {
	return dts.DTMSer.Size(Cmd1DTM)
}

type Cmd2 struct{}

func (c Cmd2) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver struct{},
	proxy core.Proxy,
) (err error) {
	time.Sleep(500 * time.Millisecond)
	_, err = proxy.Send(seq, NewResult(true))
	return
}

func (c Cmd2) MarshalTypedMUS(w mus.Writer) (n int, err error) {
	return dts.DTMSer.Marshal(Cmd2DTM, w)
}

func (c Cmd2) SizeTypedMUS() (size int) {
	return dts.DTMSer.Size(Cmd2DTM)
}

type Cmd3 struct{}

func (c Cmd3) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver struct{},
	proxy core.Proxy,
) (err error) {
	time.Sleep(500 * time.Millisecond)
	_, err = proxy.Send(seq, NewResult(true))
	return
}

func (c Cmd3) MarshalTypedMUS(w mus.Writer) (n int, err error) {
	return dts.DTMSer.Marshal(Cmd3DTM, w)
}

func (c Cmd3) SizeTypedMUS() (size int) {
	return dts.DTMSer.Size(Cmd3DTM)
}

type Cmd4 struct{}

func (c Cmd4) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver struct{},
	proxy core.Proxy,
) (err error) {
	_, err = proxy.Send(seq, NewResult(true))
	return
}

func (c Cmd4) MarshalTypedMUS(w mus.Writer) (n int, err error) {
	return dts.DTMSer.Marshal(Cmd4DTM, w)
}

func (c Cmd4) SizeTypedMUS() (size int) {
	return dts.DTMSer.Size(Cmd4DTM)
}
