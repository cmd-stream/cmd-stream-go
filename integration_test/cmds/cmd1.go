package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/cmd-stream-go/integration_test/results"
	"github.com/cmd-stream/core-go"
	"github.com/mus-format/dts-stream-go"
	"github.com/mus-format/mus-stream-go"
)

type Cmd1 struct{}

func (c Cmd1) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver struct{},
	proxy core.Proxy,
) (err error) {
	t := 500 * time.Millisecond
	time.Sleep(t)
	_, err = proxy.Send(seq, results.NewResult(false))
	if err != nil {
		return
	}
	time.Sleep(t)
	_, err = proxy.Send(seq, results.NewResult(true))
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
