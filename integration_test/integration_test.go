package intest

import (
	"errors"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/cmd-stream/base-go"
	bcln "github.com/cmd-stream/base-go/client"
	bser "github.com/cmd-stream/base-go/server"
	ccln "github.com/cmd-stream/cmd-stream-go/client"
	cser "github.com/cmd-stream/cmd-stream-go/server"
	dcln "github.com/cmd-stream/delegate-go/client"
)

func TestCommunication(t *testing.T) {
	const addr = "127.0.0.1:9000"

	wg := &sync.WaitGroup{}
	server, err := StartServer(addr, wg)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("We should be able to get several results from one cmd",
		func(t *testing.T) {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				t.Fatal(err)
			}
			client, err := ccln.Default[Receiver](ClientCodec{}, conn)
			if err != nil {
				t.Fatal(err)
			}
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
			if err != nil {
				t.Fatal(err)
			}
			if seq != wantSeq {
				t.Errorf("unexpected seq, want '%v' actual '%v'", wantSeq, seq)
			}
			result, err := ReceiveResult(results)
			if err != nil {
				t.Fatal(err)
			}
			if !Equal(result, wantResult1) {
				t.Errorf("unexpected result, want '%v' actual '%v'", wantResult1,
					result)
			}
			result, err = ReceiveResult(results)
			if err != nil {
				t.Fatal(err)
			}
			if !Equal(result, wantResult2) {
				t.Errorf("unexpected result, want '%v' actual '%v'", wantResult1,
					result)
			}
		})

	t.Run("We should be able to use keepalive feature", func(t *testing.T) {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		client, err := KeepaliveClient(conn)
		if err != nil {
			t.Fatal(err)
		}
		var (
			wantSeq1    base.Seq = 1
			wantResult1          = base.AsyncResult{Seq: wantSeq1, Result: Result{true}}
			results1             = make(chan base.AsyncResult, 1)
		)
		seq, err := client.Send(Cmd2{}, results1)
		if err != nil {
			t.Fatal(err)
		}
		if seq != wantSeq1 {
			t.Errorf("unexpected seq, want '%v' actual '%v'", 2, seq)
		}

		result, err := ReceiveResult(results1)
		if err != nil {
			t.Fatal(err)
		}
		if !Equal(result, wantResult1) {
			t.Errorf("unexpected result, want '%v' actual '%v'", wantResult1,
				result)
		}

		time.Sleep(5 * time.Second)

		var (
			wantSeq2    base.Seq = 2
			wantResult2          = base.AsyncResult{Seq: wantSeq2, Result: Result{true}}
			results2             = make(chan base.AsyncResult, 1)
		)
		seq, err = client.Send(Cmd3{}, results2)
		if err != nil {
			t.Fatal(err)
		}
		if seq != wantSeq2 {
			t.Errorf("unexpected seq, want '%v' actual '%v'", 2, seq)
		}

		result, err = ReceiveResult(results2)
		if err != nil {
			t.Fatal(err)
		}
		if !Equal(result, wantResult2) {
			t.Errorf("unexpected result, want '%v' actual '%v'", wantResult1,
				result)
		}
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
	server = cser.Default[Receiver](ServerCodec{}, Receiver{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Serve(l.(*net.TCPListener))
	}()
	return
}

func KeepaliveClient(conn net.Conn) (client *bcln.Client[Receiver], err error) {
	conf := ccln.Conf{
		Delegate: dcln.Conf{
			KeepaliveTime:  time.Second,
			KeepaliveIntvl: time.Second,
		},
	}
	return ccln.New[Receiver](conf, cser.DefaultServerInfo, ClientCodec{}, conn,
		nil)
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
