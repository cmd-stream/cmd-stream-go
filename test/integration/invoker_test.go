package integration

import (
	"context"
	"net"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/cmd-stream-go/test/mock"
	"github.com/cmd-stream/cmd-stream-go/testkit"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
	"github.com/ymz-ncnk/mok"
)

func TestInvoker(t *testing.T) {
	const addr = "127.0.0.1:9002"
	var (
		receiver = testkit.Receiver{}
		invoker  = mock.NewInvoker[testkit.Receiver]()
	)
	invoker.RegisterInvoke(func(ctx context.Context, seq core.Seq, at time.Time,
		bytesRead int, cmd core.Cmd[testkit.Receiver], proxy core.Proxy,
	) (err error) {
		return cmd.Exec(ctx, seq, at, receiver, proxy)
	})

	// Start server with custom invoker
	server, err := cmdstream.NewServerWithInvoker(invoker, testkit.ServerCodec{})
	assertfatal.EqualError(t, err, nil)
	go func() {
		_ = server.ListenAndServe(addr)
	}()
	defer server.Close()
	time.Sleep(50 * time.Millisecond)

	// Client sends a command
	conn, err := net.Dial("tcp", addr)
	assertfatal.EqualError(t, err, nil)
	client, err := cmdstream.NewClient(testkit.ClientCodec{}, conn)
	assertfatal.EqualError(t, err, nil)
	defer client.Close()

	results := make(chan core.AsyncResult, 1)
	seq, _, err := client.Send(testkit.Cmd{}, results)
	assertfatal.EqualError(t, err, nil)

	receiveAndAssert(t, results, seq, testkit.Result{LastOneFlag: true})
	asserterror.EqualDeep(t, mok.CheckCalls([]*mok.Mock{invoker.Mock}), mok.EmptyInfomap)
}
