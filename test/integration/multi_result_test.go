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
	"github.com/cmd-stream/testkit-go/cmds"
	"github.com/cmd-stream/testkit-go/codecs"
	"github.com/cmd-stream/testkit-go/exch"
	rcvr "github.com/cmd-stream/testkit-go/receiver"
	"github.com/cmd-stream/testkit-go/results"

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
		_ = server.ListenAndServe(addr)
	}()
	time.Sleep(100 * time.Millisecond)

	client, err := makeClient(addr)
	assertfatal.EqualError(t, err, nil)

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
		wantSend          = exch.WantSend{
			Seq: cmdSeq,
			N:   codecs.MultiCmdSize(cmdSeq, cmd),
		}
		cmdResults = make(chan core.AsyncResult)
	)
	err = exch.Send(cmd, cmdResults, sendFn, wantSend)
	assertfatal.EqualError(t, err, nil)

	// Receive from results.
	wantReceive1 := exch.WantReceive{
		AsyncResult: codecs.AsyncResult(cmdSeq, results.Result{LastOneFlag: false}),
	}

	err = exch.Receive[rcvr.Receiver](cmdResults, receiveFn, wantReceive1)
	assertfatal.EqualError(t, err, nil)

	// Receive from results.
	wantReceive2 := exch.WantReceive{
		AsyncResult: codecs.AsyncResult(cmdSeq, results.Result{LastOneFlag: true}),
	}

	err = exch.Receive[rcvr.Receiver](cmdResults, receiveFn, wantReceive2)
	assertfatal.EqualError(t, err, nil)
}

func makeClient(addr string) (client *ccln.Client[rcvr.Receiver], err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	return cmdstream.MakeClient(codecs.ClientCodec{}, conn)
}
