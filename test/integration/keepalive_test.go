package integration

import (
	"net"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	"github.com/cmd-stream/cmd-stream-go/core"
	ccln "github.com/cmd-stream/cmd-stream-go/core/cln"
	dcln "github.com/cmd-stream/cmd-stream-go/delegate/cln"
	hdlr "github.com/cmd-stream/cmd-stream-go/handler"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/cmd-stream-go/testkit"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestKeepalive(t *testing.T) {
	const addr = "127.0.0.1:9004"

	startKeepaliveServer(t, addr)
	client, err := makeKeepaliveClient(addr)
	assertfatal.EqualError(t, err, nil)

	exchangeKeepalive(t, client)
}

func startKeepaliveServer(t *testing.T, addr string) {
	go func() {
		server, err := cmdstream.NewServer(testkit.Receiver{}, testkit.ServerCodec{},
			srv.WithHandler(
				hdlr.WithCmdReceiveDuration(250*time.Millisecond),
			),
		)
		asserterror.EqualError(t, err, nil)
		_ = server.ListenAndServe(addr)
	}()
	time.Sleep(50 * time.Millisecond)
}

func makeKeepaliveClient(addr string) (client *ccln.Client[testkit.Receiver],
	err error,
) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	return cmdstream.NewClient(testkit.ClientCodec{}, conn,
		cln.WithKeepalive(
			dcln.WithKeepaliveTime(100*time.Millisecond),
			dcln.WithKeepaliveIntvl(100*time.Millisecond),
		),
	)
}

func exchangeKeepalive(t *testing.T, client *ccln.Client[testkit.Receiver]) {
	var (
		cmd     = testkit.Cmd{ExecTime: 50 * time.Millisecond}
		results = make(chan core.AsyncResult, 1)
		err     error
	)
	// Send first command.
	seq, _, err := client.Send(cmd, results)
	assertfatal.EqualError(t, err, nil)
	receiveAndAssert(t, results, seq, testkit.Result{LastOneFlag: true})

	// Wait for a long duration.
	time.Sleep(500 * time.Millisecond)

	// Send second command.
	seq, _, err = client.Send(cmd, results)
	assertfatal.EqualError(t, err, nil)
	receiveAndAssert(t, results, seq, testkit.Result{LastOneFlag: true})
}
