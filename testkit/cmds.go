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
	ExecTime time.Duration `json:"exec_time"`
}

func (c Cmd) Exec(ctx context.Context, receiver Receiver, proxy core.Proxy) (
	err error,
) {
	time.Sleep(c.ExecTime)
	_, err = proxy.Send(Result{LastOneFlag: true})
	return
}

// -----------------------------------------------------------------------------

type MultiCmd struct {
	ResultsCount int
	ExecTime     time.Duration
}

func (c MultiCmd) Exec(ctx context.Context, receiver Receiver, proxy core.Proxy) (
	err error,
) {
	for i := range c.ResultsCount {
		time.Sleep(c.ExecTime)
		_, err = proxy.Send(Result{LastOneFlag: i == c.ResultsCount-1})
		if err != nil {
			return
		}
	}
	return
}
