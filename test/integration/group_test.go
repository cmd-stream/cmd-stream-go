package integration

import (
	"net"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	"github.com/cmd-stream/cmd-stream-go/core"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/cmd-stream-go/testkit"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestGroup(t *testing.T) {
	const addr = "127.0.0.1:9000"

	startGroupServer(t, addr)
	group, err := makeGroup(addr)
	assertfatal.EqualError(t, err, nil)

	exchangeGroup(t, group)
}

func startGroupServer(t *testing.T, addr string) {
	go func() {
		server, err := cmdstream.NewServer(testkit.Receiver{}, testkit.ServerCodec{})
		asserterror.EqualError(t, err, nil)
		err = server.ListenAndServe(addr)
		asserterror.EqualError(t, err, nil)
	}()
	time.Sleep(50 * time.Millisecond)
}

func makeGroup(addr string) (group grp.Group[testkit.Receiver], err error) {
	return cmdstream.NewGroup(2, testkit.ClientCodec{}, cln.ConnFactoryFn(
		func() (net.Conn, error) { return net.Dial("tcp", addr) }),
	)
}

func exchangeGroup(t *testing.T, group grp.Group[testkit.Receiver]) {
	items := []struct {
		seq            core.Seq
		clientID       grp.ClientID
		cmd            testkit.Cmd
		results        chan core.AsyncResult
		expectedResult testkit.Result
	}{
		{
			seq:            1,
			clientID:       0,
			cmd:            testkit.Cmd{},
			results:        make(chan core.AsyncResult, 1),
			expectedResult: testkit.Result{LastOneFlag: true},
		},
		{
			seq:            1,
			clientID:       1,
			cmd:            testkit.Cmd{},
			results:        make(chan core.AsyncResult, 1),
			expectedResult: testkit.Result{LastOneFlag: true},
		},
	}

	for _, item := range items {
		seq, clientID, n, err := group.Send(item.cmd, item.results)
		asserterror.Equal(t, seq, item.seq)
		asserterror.Equal(t, clientID, item.clientID)
		asserterror.Equal(t, n, testkit.CalcCmdSize(item.seq, item.cmd))
		asserterror.Equal(t, err, nil)
	}

	for _, item := range items {
		receiveAndAssert(t, item.results, item.seq, item.expectedResult)
	}
}
