package integration_test

import (
	"net"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	ccln "github.com/cmd-stream/cmd-stream-go/client"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/core-go"
	dcln "github.com/cmd-stream/delegate-go/client"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestKeepalive(t *testing.T) {
	const addr = "127.0.0.1:9001"

	go func() {
		server := cmdstream.MakeServer(ServerCodec{}, srv.NewInvoker(struct{}{}))
		server.ListenAndServe(addr)
	}()
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", addr)
	assertfatal.EqualError(err, nil, t)
	client, err := cmdstream.MakeClient(ClientCodec{}, conn, nil,
		ccln.WithKeepalive(
			dcln.WithKeepaliveTime(time.Second),
			dcln.WithKeepaliveIntvl(time.Second),
		),
	)
	assertfatal.EqualError(err, nil, t)

	var (
		wantSeq1    core.Seq = 1
		wantResult1          = core.AsyncResult{
			Seq:       wantSeq1,
			BytesRead: 2,
			Result:    NewResult(true),
		}
		results1 = make(chan core.AsyncResult, 1)
	)
	seq, _, err := client.Send(Cmd2{}, results1)
	assertfatal.EqualError(err, nil, t)
	asserterror.Equal(seq, wantSeq1, t)

	result, err := receiveResult(results1)
	assertfatal.EqualError(err, nil, t)
	asserterror.Equal(result, wantResult1, t)

	time.Sleep(5 * time.Second)

	var (
		wantSeq2    core.Seq = 2
		wantResult2          = core.AsyncResult{
			Seq:       wantSeq2,
			BytesRead: 2,
			Result:    NewResult(true),
		}
		results2 = make(chan core.AsyncResult, 1)
	)
	seq, _, err = client.Send(Cmd3{}, results2)
	assertfatal.EqualError(err, nil, t)
	asserterror.Equal(seq, wantSeq2, t)

	result, err = receiveResult(results2)
	assertfatal.EqualError(err, nil, t)
	asserterror.Equal(result, wantResult2, t)
}
