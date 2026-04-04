package testkit

import (
	"context"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	com "github.com/mus-format/common-go"
)

const (
	CmdDTM com.DTM = iota + 1
	MultiCmdDTM
)

const (
	CmdSize = 17
)

// -----------------------------------------------------------------------------

type Cmd struct {
	ExecTime time.Duration `json:"execTime"`
}

func (c Cmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver Receiver,
	proxy core.Proxy,
) (err error) {
	time.Sleep(c.ExecTime)
	_, err = proxy.Send(seq, Result{LastOneFlag: true})
	return
}

// -----------------------------------------------------------------------------

type MultiCmd struct {
	ResultsCount int
	ExecTime     time.Duration
}

func (c MultiCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver Receiver,
	proxy core.Proxy,
) (err error) {
	for i := range c.ResultsCount {
		time.Sleep(c.ExecTime)
		_, err = proxy.Send(seq, Result{LastOneFlag: i == c.ResultsCount-1})
		if err != nil {
			return
		}
	}
	return
}
