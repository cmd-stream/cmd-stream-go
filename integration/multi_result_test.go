package integration_test

import (
	"net"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/core-go"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestMultiResult(t *testing.T) {
	const addr = "127.0.0.1:9002"

	go func() {
		server := cmdstream.MakeServer(ServerCodec{}, srv.NewInvoker(struct{}{}))
		server.ListenAndServe(addr)
	}()
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", addr)
	assertfatal.EqualError(err, nil, t)
	client, err := cmdstream.MakeClient(ClientCodec{}, conn)
	assertfatal.EqualError(err, nil, t)

	var (
		wantSeq     core.Seq = 1
		wantResult1          = core.AsyncResult{
			Seq:       wantSeq,
			BytesRead: 2,
			Result:    NewResult(false),
		}
		wantResult2 = core.AsyncResult{
			Seq:       wantSeq,
			BytesRead: 2,
			Result:    NewResult(true),
		}
		results = make(chan core.AsyncResult)
	)

	seq, _, err := client.Send(Cmd1{}, results)
	assertfatal.EqualError(err, nil, t)
	asserterror.Equal(seq, wantSeq, t)

	result, err := receiveResult(results)
	assertfatal.EqualError(err, nil, t)
	asserterror.Equal(result, wantResult1, t)

	result, err = receiveResult(results)
	assertfatal.EqualError(err, nil, t)
	asserterror.Equal(result, wantResult2, t)
}
