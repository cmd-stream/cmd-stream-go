package integration

import (
	"net"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	"github.com/cmd-stream/cmd-stream-go/core"
	csrv "github.com/cmd-stream/cmd-stream-go/core/srv"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/cmd-stream-go/testkit"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestGroup(t *testing.T) {
	const addr = "127.0.0.1:9001"

	server := startGroupServer(t, addr)
	defer server.Close()

	group, err := makeGroup(addr)
	assertfatal.EqualError(t, err, nil)
	defer group.Close()

	exchangeGroup(t, group)
}

func startGroupServer(t *testing.T, addr string) *csrv.Server {
	server, err := cmdstream.NewServer(testkit.Receiver{}, testkit.ServerCodec{})
	asserterror.EqualError(t, err, nil)
	go func() {
		_ = server.ListenAndServe(addr)
	}()
	time.Sleep(50 * time.Millisecond)
	return server
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
