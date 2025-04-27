package intest

import (
	"errors"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/cmd-stream/base-go"
	bser "github.com/cmd-stream/base-go/server"
	ccln "github.com/cmd-stream/cmd-stream-go/client"
	cser "github.com/cmd-stream/cmd-stream-go/server"
	dcln "github.com/cmd-stream/delegate-go/client"

	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestCommunication(t *testing.T) {
	const addr = "127.0.0.1:9000"

	wg := &sync.WaitGroup{}
	server, err := StartServer(addr, wg)
	assertfatal.EqualError(err, nil, t)

	t.Run("We should be able to get several results from one cmd",
		func(t *testing.T) {
			conn, err := net.Dial("tcp", addr)
			assertfatal.EqualError(err, nil, t)

			client, err := ccln.New[Receiver](ClientCodec{}, conn)
			assertfatal.EqualError(err, nil, t)
			var (
				wantSeq     base.Seq = 1
				wantResult1          = base.AsyncResult{
					Seq:    wantSeq,
					Result: Result{false},
				}
				wantResult2 = base.AsyncResult{
					Seq:    wantSeq,
					Result: Result{true},
				}
				results = make(chan base.AsyncResult)
			)

			seq, err := client.Send(Cmd1{}, results)
			assertfatal.EqualError(err, nil, t)
			asserterror.Equal(seq, wantSeq, t)

			result, err := ReceiveResult(results)
			assertfatal.EqualError(err, nil, t)
			asserterror.Equal(result, wantResult1, t)

			result, err = ReceiveResult(results)
			assertfatal.EqualError(err, nil, t)
			asserterror.Equal(result, wantResult2, t)
		})

	t.Run("We should be able to use keepalive feature", func(t *testing.T) {
		conn, err := net.Dial("tcp", addr)
		assertfatal.EqualError(err, nil, t)

		client, err := KeepaliveClient(conn)
		assertfatal.EqualError(err, nil, t)
		var (
			wantSeq1    base.Seq = 1
			wantResult1          = base.AsyncResult{Seq: wantSeq1, Result: Result{true}}
			results1             = make(chan base.AsyncResult, 1)
		)
		seq, err := client.Send(Cmd2{}, results1)
		assertfatal.EqualError(err, nil, t)
		asserterror.Equal(seq, wantSeq1, t)

		result, err := ReceiveResult(results1)
		assertfatal.EqualError(err, nil, t)
		asserterror.Equal(result, wantResult1, t)

		time.Sleep(5 * time.Second)

		var (
			wantSeq2    base.Seq = 2
			wantResult2          = base.AsyncResult{Seq: wantSeq2, Result: Result{true}}
			results2             = make(chan base.AsyncResult, 1)
		)
		seq, err = client.Send(Cmd3{}, results2)
		assertfatal.EqualError(err, nil, t)
		asserterror.Equal(seq, wantSeq2, t)

		result, err = ReceiveResult(results2)
		assertfatal.EqualError(err, nil, t)
		asserterror.Equal(result, wantResult2, t)
	})

	t.Run("We should be able to use Streamline", func(t *testing.T) {
		var (
			wantSeq1      base.Seq = 1
			wantResult1            = base.AsyncResult{Seq: wantSeq1, Result: Result{true}}
			wantClientID1          = ccln.ClientID(0)
			results1               = make(chan base.AsyncResult, 1)

			factory = ccln.ConnFactoryFn(func() (net.Conn, error) {
				return net.Dial("tcp", addr)
			})
			clients = ccln.MustMakeClients(2, ClientCodec{}, factory)
			grp     = ccln.NewGroup(ccln.NewRoundRobinStrategy(clients))
		)
		seq, clientID, err := grp.Send(Cmd4{}, results1)
		assertfatal.EqualError(err, nil, t)
		asserterror.Equal(seq, wantSeq1, t)
		asserterror.Equal(clientID, wantClientID1, t)

		var (
			wantSeq2      base.Seq = 1
			wantResult2            = base.AsyncResult{Seq: wantSeq2, Result: Result{true}}
			wantClientID2          = ccln.ClientID(1)
			results2               = make(chan base.AsyncResult, 1)
		)
		seq, clientID, err = grp.Send(Cmd4{}, results2)
		assertfatal.EqualError(err, nil, t)
		asserterror.Equal(seq, wantSeq1, t)
		asserterror.Equal(clientID, wantClientID2, t)

		result, err := ReceiveResult(results1)
		assertfatal.EqualError(err, nil, t)
		asserterror.Equal(result, wantResult1, t)

		result, err = ReceiveResult(results2)
		assertfatal.EqualError(err, nil, t)
		asserterror.Equal(result, wantResult2, t)
	})

	if err := server.Close(); err != nil {
		t.Fatal(err)
	}
	wg.Wait()
}

func StartServer(addr string, wg *sync.WaitGroup) (server *bser.Server,
	err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	server = cser.New[Receiver](ServerCodec{}, cser.NewInvoker(Receiver{}))
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Serve(l.(*net.TCPListener))
	}()
	return
}

func KeepaliveClient(conn net.Conn) (client ccln.Client[Receiver], err error) {
	return ccln.New[Receiver](ClientCodec{}, conn, nil, ccln.WithKeepalive(
		dcln.WithKeepaliveTime(time.Second),
		dcln.WithKeepaliveIntvl(time.Second),
	),
	)
}

func ReceiveResult(results <-chan base.AsyncResult) (result base.AsyncResult,
	err error) {
	select {
	case <-time.NewTimer(time.Second).C:
		err = errors.New("test lasts too long")
	case result = <-results:
	}
	return
}

func Equal(r1, r2 base.AsyncResult) bool {
	return reflect.DeepEqual(r1.Result, r2.Result) &&
		r1.Seq == r2.Seq &&
		r1.Error == nil &&
		r2.Error == nil
}
