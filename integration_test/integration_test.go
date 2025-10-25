package intest

import (
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	ccln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/cmd-stream-go/integration_test/cmds"
	"github.com/cmd-stream/cmd-stream-go/integration_test/results"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/core-go"
	csrv "github.com/cmd-stream/core-go/server"
	dcln "github.com/cmd-stream/delegate-go/client"

	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestCommunication(t *testing.T) {
	const addr = "127.0.0.1:9000"

	wg := &sync.WaitGroup{}
	server, err := startServer(addr, wg)
	assertfatal.EqualError(err, nil, t)

	t.Run("We should be able to get several results from one cmd",
		func(t *testing.T) {
			conn, err := net.Dial("tcp", addr)
			assertfatal.EqualError(err, nil, t)

			client, err := cmdstream.MakeClient(ClientCodec{}, conn)
			assertfatal.EqualError(err, nil, t)
			var (
				wantSeq     core.Seq = 1
				wantResult1          = core.AsyncResult{
					Seq:       wantSeq,
					BytesRead: 2,
					Result:    results.NewResult(false),
				}
				wantResult2 = core.AsyncResult{
					Seq:       wantSeq,
					BytesRead: 2,
					Result:    results.NewResult(true),
				}
				results = make(chan core.AsyncResult)
			)

			seq, _, err := client.Send(cmds.Cmd1{}, results)
			assertfatal.EqualError(err, nil, t)
			asserterror.Equal(seq, wantSeq, t)

			result, err := receiveResult(results)
			assertfatal.EqualError(err, nil, t)
			asserterror.Equal(result, wantResult1, t)

			result, err = receiveResult(results)
			assertfatal.EqualError(err, nil, t)
			asserterror.Equal(result, wantResult2, t)
		})

	t.Run("We should be able to use keepalive feature", func(t *testing.T) {
		conn, err := net.Dial("tcp", addr)
		assertfatal.EqualError(err, nil, t)

		client, err := keepaliveClient(conn)
		assertfatal.EqualError(err, nil, t)
		var (
			wantSeq1    core.Seq = 1
			wantResult1          = core.AsyncResult{
				Seq:       wantSeq1,
				BytesRead: 2,
				Result:    results.NewResult(true),
			}
			results1 = make(chan core.AsyncResult, 1)
		)
		seq, _, err := client.Send(cmds.Cmd2{}, results1)
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
				Result:    results.NewResult(true),
			}
			results2 = make(chan core.AsyncResult, 1)
		)
		seq, _, err = client.Send(cmds.Cmd3{}, results2)
		assertfatal.EqualError(err, nil, t)
		asserterror.Equal(seq, wantSeq2, t)

		result, err = receiveResult(results2)
		assertfatal.EqualError(err, nil, t)
		asserterror.Equal(result, wantResult2, t)
	})

	t.Run("We should be able to use a client group", func(t *testing.T) {
		var (
			wantSeq1    core.Seq = 1
			wantN1      int      = 2
			wantResult1          = core.AsyncResult{
				Seq:       wantSeq1,
				BytesRead: 2,
				Result:    results.NewResult(true),
			}
			wantClientID1 = grp.ClientID(0)
			results1      = make(chan core.AsyncResult, 1)

			factory = ccln.ConnFactoryFn(func() (net.Conn, error) {
				return net.Dial("tcp", addr)
			})
		)
		group, err := cmdstream.MakeClientGroup(2, ClientCodec{}, factory)
		assertfatal.EqualError(err, nil, t)

		seq, clientID, n, err := group.Send(cmds.Cmd4{}, results1)
		assertfatal.EqualError(err, nil, t)
		asserterror.Equal(seq, wantSeq1, t)
		asserterror.Equal(clientID, wantClientID1, t)
		asserterror.Equal(n, wantN1, t)

		var (
			wantSeq2    core.Seq = 1
			wantN2      int      = 2
			wantResult2          = core.AsyncResult{
				Seq:       wantSeq2,
				BytesRead: 2,
				Result:    results.NewResult(true),
			}
			wantClientID2 = grp.ClientID(1)
			results2      = make(chan core.AsyncResult, 1)
		)
		seq, clientID, n, err = group.Send(cmds.Cmd4{}, results2)
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
	})

	if err := server.Close(); err != nil {
		t.Fatal(err)
	}
	wg.Wait()
}

func startServer(addr string, wg *sync.WaitGroup) (server *csrv.Server,
	err error,
) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	server = cmdstream.MakeServer(ServerCodec{}, srv.NewInvoker(struct{}{}))
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Serve(l.(*net.TCPListener))
	}()
	return
}

func keepaliveClient(conn net.Conn) (client grp.Client[struct{}], err error) {
	return cmdstream.MakeClient(ClientCodec{}, conn, nil, ccln.WithKeepalive(
		dcln.WithKeepaliveTime(time.Second),
		dcln.WithKeepaliveIntvl(time.Second),
	),
	)
}

func receiveResult(results <-chan core.AsyncResult) (result core.AsyncResult,
	err error,
) {
	select {
	case <-time.NewTimer(time.Second).C:
		err = errors.New("test lasts too long")
	case result = <-results:
	}
	return
}
