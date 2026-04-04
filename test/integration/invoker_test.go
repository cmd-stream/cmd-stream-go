package integration

import (
	"context"
	"net"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/cmd-stream-go/core"
	hmock "github.com/cmd-stream/cmd-stream-go/test/mock/handler"
	"github.com/cmd-stream/cmd-stream-go/testkit"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
	"github.com/ymz-ncnk/mok"
)

func TestInvoker(t *testing.T) {
	const addr = "127.0.0.1:9007"
	var (
		receiver = testkit.Receiver{}
		invoker  = hmock.NewInvoker[testkit.Receiver]()
	)
	invoker.RegisterInvoke(func(ctx context.Context, seq core.Seq, at time.Time,
		bytesRead int, cmd core.Cmd[testkit.Receiver], proxy core.Proxy,
	) (err error) {
		return cmd.Exec(ctx, seq, at, receiver, proxy)
	})

	// Start server with custom invoker
	go func() {
		server, err := cmdstream.NewServerWithInvoker(invoker, testkit.ServerCodec{})
		asserterror.EqualError(t, err, nil)
		_ = server.ListenAndServe(addr)
	}()
	time.Sleep(50 * time.Millisecond)

	// Client sends a command
	conn, err := net.Dial("tcp", addr)
	assertfatal.EqualError(t, err, nil)
	client, err := cmdstream.NewClient(testkit.ClientCodec{}, conn)
	assertfatal.EqualError(t, err, nil)

	results := make(chan core.AsyncResult, 1)
	seq, _, err := client.Send(testkit.Cmd{}, results)
	assertfatal.EqualError(t, err, nil)

	receiveAndAssert(t, results, seq, testkit.Result{LastOneFlag: true})
	asserterror.EqualDeep(t, mok.CheckCalls([]*mok.Mock{invoker.Mock}), mok.EmptyInfomap)

	err = client.Close()
	asserterror.EqualError(t, err, nil)
}
