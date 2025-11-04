//go:build integration
// +build integration

package integration

import (
	"errors"
	"net"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/core-go"
	cmds "github.com/cmd-stream/testkit-go/fixtures/cmdstream/cmds"
	helpers "github.com/cmd-stream/testkit-go/helpers"

	codecs "github.com/cmd-stream/testkit-go/fixtures/cmdstream/codecs"
	rcvr "github.com/cmd-stream/testkit-go/fixtures/cmdstream/receiver"
	results "github.com/cmd-stream/testkit-go/fixtures/cmdstream/results"

	cln "github.com/cmd-stream/cmd-stream-go/client"

	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestClientGroup(t *testing.T) {
	const addr = "127.0.0.1:9000"

	go func() {
		var (
			invoker = srv.NewInvoker(rcvr.Receiver{})
			server  = cmdstream.MakeServer(codecs.ServerCodec{}, invoker)
		)
		server.ListenAndServe(addr)
	}()
	time.Sleep(100 * time.Millisecond)

	group, err := makeClientGroup(addr)
	assertfatal.EqualError(err, nil, t)

	var (
		sendFn = func(cmd core.Cmd[rcvr.Receiver], results chan<- core.AsyncResult) (
			seq core.Seq, clientID grp.ClientID, n int, err error,
		) {
			return group.Send(cmd, results)
		}
		receiveFn = func(results <-chan core.AsyncResult) (asyncResult core.AsyncResult, err error) {
			select {
			case <-time.NewTimer(time.Second).C:
				err = errors.New("test lasts too long")
			case asyncResult = <-results:
			}
			return
		}
	)

	// First send.
	var (
		cmd1               = cmds.Cmd{}
		cmdSeq1   core.Seq = 1
		results1           = make(chan core.AsyncResult, 1)
		wantSend1          = helpers.WantSendGrp{
			Seq:      cmdSeq1,
			ClientID: grp.ClientID(0),
			N:        codecs.CmdSize(cmdSeq1, cmd1),
		}
	)
	err = helpers.SendGrp(cmd1, results1, sendFn, wantSend1)
	assertfatal.EqualError(err, nil, t)

	// Second send.
	var (
		cmd2               = cmds.Cmd{}
		cmdSeq2   core.Seq = 1
		results2           = make(chan core.AsyncResult, 1)
		wantSend2          = helpers.WantSendGrp{
			Seq:      cmdSeq2,
			ClientID: grp.ClientID(1),
			N:        codecs.CmdSize(cmdSeq2, cmd2),
		}
	)
	err = helpers.SendGrp(cmd2, results2, sendFn, wantSend2)
	assertfatal.EqualError(err, nil, t)

	// Receive from results1.
	wantReceive1 := helpers.WantReceive{
		AsyncResult: codecs.AsyncResult(cmdSeq1, results.Result{LastOneFlag: true}),
	}

	err = helpers.Receive[rcvr.Receiver](results1, receiveFn, wantReceive1)
	assertfatal.EqualError(err, nil, t)

	// Receive from results2.
	wantReceive2 := helpers.WantReceive{
		AsyncResult: codecs.AsyncResult(cmdSeq2, results.Result{LastOneFlag: true}),
	}

	err = helpers.Receive[rcvr.Receiver](results2, receiveFn, wantReceive2)
	assertfatal.EqualError(err, nil, t)
}

func makeClientGroup(addr string) (grp.ClientGroup[rcvr.Receiver], error) {
	factory := cln.ConnFactoryFn(func() (net.Conn, error) {
		return net.Dial("tcp", addr)
	})
	clientsCount := 2
	return cmdstream.MakeClientGroup(clientsCount, codecs.ClientCodec{}, factory)
}
