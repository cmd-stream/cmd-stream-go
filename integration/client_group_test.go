package integration_test

import (
	"net"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/core-go"

	cln "github.com/cmd-stream/cmd-stream-go/client"

	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestClientGropu(t *testing.T) {
	const addr = "127.0.0.1:9000"

	go func() {
		server := cmdstream.MakeServer(ServerCodec{}, srv.NewInvoker(struct{}{}))
		server.ListenAndServe(addr)
	}()
	time.Sleep(100 * time.Millisecond)

	var (
		wantSeq1    core.Seq = 1
		wantN1               = 2
		wantResult1          = core.AsyncResult{
			Seq:       wantSeq1,
			BytesRead: 2,
			Result:    NewResult(true),
		}
		wantClientID1 = grp.ClientID(0)
		results1      = make(chan core.AsyncResult, 1)

		factory = cln.ConnFactoryFn(func() (net.Conn, error) {
			return net.Dial("tcp", addr)
		})
	)
	group, err := cmdstream.MakeClientGroup(2, ClientCodec{}, factory)
	assertfatal.EqualError(err, nil, t)

	seq, clientID, n, err := group.Send(Cmd4{}, results1)
	assertfatal.EqualError(err, nil, t)
	asserterror.Equal(seq, wantSeq1, t)
	asserterror.Equal(clientID, wantClientID1, t)
	asserterror.Equal(n, wantN1, t)

	var (
		wantSeq2    core.Seq = 1
		wantN2               = 2
		wantResult2          = core.AsyncResult{
			Seq:       wantSeq2,
			BytesRead: 2,
			Result:    NewResult(true),
		}
		wantClientID2 = grp.ClientID(1)
		results2      = make(chan core.AsyncResult, 1)
	)
	seq, clientID, n, err = group.Send(Cmd4{}, results2)
	assertfatal.EqualError(err, nil, t)
	asserterror.Equal(seq, wantSeq2, t)
	asserterror.Equal(clientID, wantClientID2, t)
	asserterror.Equal(n, wantN2, t)

	result, err := receiveResult(results1)
	assertfatal.EqualError(err, nil, t)
	asserterror.Equal(result, wantResult1, t)

	result, err = receiveResult(results2)
	assertfatal.EqualError(err, nil, t)
	asserterror.Equal(result, wantResult2, t)
}
