//go:build integration
// +build integration

package integration

import (
	"errors"
	"net"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/testkit-go/fixtures/cmdstream/cmds"
	"github.com/cmd-stream/testkit-go/fixtures/cmdstream/codecs"
	rcvr "github.com/cmd-stream/testkit-go/fixtures/cmdstream/receiver"
	"github.com/cmd-stream/testkit-go/fixtures/cmdstream/results"
	"github.com/cmd-stream/testkit-go/helpers"

	ccln "github.com/cmd-stream/core-go/client"

	"github.com/cmd-stream/core-go"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestMultiResult(t *testing.T) {
	const addr = "127.0.0.1:9002"

	go func() {
		var (
			invoker = srv.NewInvoker(rcvr.Receiver{})
			server  = cmdstream.MakeServer(codecs.ServerCodec{}, invoker)
		)
		server.ListenAndServe(addr)
	}()
	time.Sleep(100 * time.Millisecond)

	client, err := makeClient(addr)
	assertfatal.EqualError(err, nil, t)

	var (
		sendFn = func(cmd core.Cmd[rcvr.Receiver], cmdResults chan<- core.AsyncResult) (
			seq core.Seq, n int, err error,
		) {
			return client.Send(cmd, cmdResults)
		}
		receiveFn = func(results <-chan core.AsyncResult) (asyncResult core.AsyncResult,
			err error,
		) {
			select {
			case <-time.NewTimer(time.Second).C:
				err = errors.New("test lasts too long")
			case asyncResult = <-results:
			}
			return
		}
	)

	// Send multi-result Command.
	var (
		cmd = cmds.MultiCmd{
			ResultsCount: 2,
			ExecTime:     500 * time.Millisecond,
		}
		cmdSeq   core.Seq = 1
		wantSend          = helpers.WantSend{
			Seq: cmdSeq,
			N:   codecs.MultiCmdSize(cmdSeq, cmd),
		}
		cmdResults = make(chan core.AsyncResult)
	)
	err = helpers.Send(cmd, cmdResults, sendFn, wantSend)
	assertfatal.EqualError(err, nil, t)

	// Receive from results.
	wantReceive1 := helpers.WantReceive{
		AsyncResult: codecs.AsyncResult(cmdSeq, results.Result{LastOneFlag: false}),
	}

	err = helpers.Receive[rcvr.Receiver](cmdResults, receiveFn, wantReceive1)
	assertfatal.EqualError(err, nil, t)

	// Receive from results.
	wantReceive2 := helpers.WantReceive{
		AsyncResult: codecs.AsyncResult(cmdSeq, results.Result{LastOneFlag: true}),
	}

	err = helpers.Receive[rcvr.Receiver](cmdResults, receiveFn, wantReceive2)
	assertfatal.EqualError(err, nil, t)
}

func makeClient(addr string) (client *ccln.Client[rcvr.Receiver], err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	return cmdstream.MakeClient(codecs.ClientCodec{}, conn)
}
