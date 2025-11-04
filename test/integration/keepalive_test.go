//go:build integration
// +build integration

package integration

import (
	"errors"
	"net"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/core-go"
	ccln "github.com/cmd-stream/core-go/client"
	dcln "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/testkit-go/fixtures/cmdstream/cmds"
	"github.com/cmd-stream/testkit-go/fixtures/cmdstream/codecs"
	rcvr "github.com/cmd-stream/testkit-go/fixtures/cmdstream/receiver"
	"github.com/cmd-stream/testkit-go/fixtures/cmdstream/results"
	"github.com/cmd-stream/testkit-go/helpers"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestKeepalive(t *testing.T) {
	const addr = "127.0.0.1:9001"

	go func() {
		var (
			invoker = srv.NewInvoker(rcvr.Receiver{})
			// TODO Use option to disconnect after some time.
			server = cmdstream.MakeServer(codecs.ServerCodec{}, invoker)
		)
		server.ListenAndServe(addr)
	}()
	time.Sleep(100 * time.Millisecond)

	client, err := makeKeepaliveClient(addr)
	assertfatal.EqualError(err, nil, t)

	var (
		sendFn = func(cmd core.Cmd[rcvr.Receiver], results chan<- core.AsyncResult) (
			seq core.Seq, n int, err error,
		) {
			return client.Send(cmd, results)
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

	// Send first Command.
	var (
		cmd1              = cmds.Cmd{ExecTime: 500 * time.Millisecond}
		cmdSeq1  core.Seq = 1
		wantSend          = helpers.WantSend{
			Seq: cmdSeq1,
			N:   codecs.CmdSize(cmdSeq1, cmd1),
		}
		results1 = make(chan core.AsyncResult, 1)
	)
	err = helpers.Send(cmd1, results1, sendFn, wantSend)
	assertfatal.EqualError(err, nil, t)

	// Receive from results1.
	wantReceive1 := helpers.WantReceive{
		AsyncResult: codecs.AsyncResult(cmdSeq1, results.Result{LastOneFlag: true}),
	}

	err = helpers.Receive[rcvr.Receiver](results1, receiveFn, wantReceive1)
	assertfatal.EqualError(err, nil, t)

	// Wait for a long duration.
	time.Sleep(5 * time.Second)

	// Send second Command.
	var (
		cmd2               = cmds.Cmd{ExecTime: 500 * time.Millisecond}
		cmdSeq2   core.Seq = 2
		wantSend2          = helpers.WantSend{
			Seq: cmdSeq2,
			N:   codecs.CmdSize(cmdSeq2, cmd2),
		}
		results2 = make(chan core.AsyncResult, 1)
	)
	err = helpers.Send(cmd2, results2, sendFn, wantSend2)
	assertfatal.EqualError(err, nil, t)

	// Receive from results2.
	wantReceive2 := helpers.WantReceive{
		AsyncResult: codecs.AsyncResult(cmdSeq2, results.Result{LastOneFlag: true}),
	}

	err = helpers.Receive[rcvr.Receiver](results2, receiveFn, wantReceive2)
	assertfatal.EqualError(err, nil, t)
}

func makeKeepaliveClient(addr string) (client *ccln.Client[rcvr.Receiver],
	err error,
) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	return cmdstream.MakeClient(codecs.ClientCodec{}, conn, nil,
		cln.WithKeepalive(
			dcln.WithKeepaliveTime(time.Second),
			dcln.WithKeepaliveIntvl(time.Second),
		),
	)
}
