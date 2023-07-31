package ct

import (
	"errors"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/cmd-stream/base-go"
	base_server "github.com/cmd-stream/base-go/server"
	cs_client "github.com/cmd-stream/cmd-stream-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/delegate-go"
	delegate_client "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/handler-go"
)

const Addr = "127.0.0.1:9000"

func TestCommunication(t *testing.T) {
	wg := &sync.WaitGroup{}
	listener, err := net.Listen("tcp", Addr)
	if err != nil {
		t.Fatal(err)
	}

	server := cs_server.New[Receiver](cs_server.DefServerInfo,
		delegate.ServerSettings{},
		cs_server.Conf{
			Handler: handler.Conf{
				ReceiveTimeout: 2 * time.Second,
			},
			Base: base_server.Conf{
				WorkersCount: 2,
			},
		},
		ServerCodec{},
		Receiver{},
		nil)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Serve(listener.(*net.TCPListener)); err != base_server.ErrClosed {
			t.Error(err)
		}
	}()

	t.Run("We should be able to get several results from one cmd",
		func(t *testing.T) {
			conn, err := net.Dial("tcp", Addr)
			if err != nil {
				t.Fatal(err)
			}
			client, err := cs_client.New[Receiver](cs_server.DefServerInfo,
				cs_client.Conf{},
				ClientCodec{},
				conn,
				nil)
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
			)

			results := make(chan base.AsyncResult)
			seq, err := client.Send(Cmd1{}, results) // timeout == 0s, r1 - 0.5s, r2 - 5s
			if err != nil {
				t.Fatal(err)
			}
			if seq != wantSeq {
				t.Errorf("unexpected seq, want '%v' actual '%v'", wantSeq, seq)
			}
			result, err := receiveResult(results)
			if err != nil {
				t.Fatal(err)
			}
			if !equalOkResult(wantResult1, result) {
				t.Errorf("unexpected result, want '%v' actual '%v'", wantResult1,
					result)
			}
			result, err = receiveResult(results)
			if err != nil {
				t.Fatal(err)
			}
			if !equalOkResult(wantResult2, result) {
				t.Errorf("unexpected result, want '%v' actual '%v'", wantResult1,
					result)
			}
		})

	t.Run("We should be able to use keepalive feature", func(t *testing.T) {
		conn, err := net.Dial("tcp", Addr)
		if err != nil {
			t.Fatal(err)
		}
		client, err := cs_client.New[Receiver](cs_server.DefServerInfo,
			cs_client.Conf{
				Delegate: delegate_client.Conf{
					KeepaliveTime:  time.Second,
					KeepaliveIntvl: time.Second,
				},
			},
			ClientCodec{},
			conn,
			nil)
		if err != nil {
			t.Fatal(err)
		}
		var (
			wantSeq1    base.Seq = 1
			wantResult1          = base.AsyncResult{Seq: wantSeq1, Result: Result{true}}
			wantSeq2    base.Seq = 2
			wantResult2          = base.AsyncResult{Seq: wantSeq2, Result: Result{true}}
		)
		results1 := make(chan base.AsyncResult, 1)
		results2 := make(chan base.AsyncResult, 1)
		seq, err := client.Send(Cmd2{}, results1)
		if err != nil {
			t.Fatal(err)
		}
		if seq != wantSeq1 {
			t.Errorf("unexpected seq, want '%v' actual '%v'", 2, seq)
		}

		result, err := receiveResult(results1)
		if err != nil {
			t.Fatal(err)
		}
		if !equalOkResult(wantResult1, result) {
			t.Errorf("unexpected result, want '%v' actual '%v'", wantResult1,
				result)
		}

		time.Sleep(5 * time.Second)

		seq, err = client.Send(Cmd3{}, results2)
		if err != nil {
			t.Fatal(err)
		}
		if seq != wantSeq2 {
			t.Errorf("unexpected seq, want '%v' actual '%v'", 2, seq)
		}

		result, err = receiveResult(results2)
		if err != nil {
			t.Fatal(err)
		}
		if !equalOkResult(wantResult2, result) {
			t.Errorf("unexpected result, want '%v' actual '%v'", wantResult1,
				result)
		}
	})

	if err := server.Close(); err != nil {
		t.Fatal(err)
	}
	wg.Wait()
}

func receiveResult(results <-chan base.AsyncResult) (result base.AsyncResult,
	err error) {
	select {
	case <-time.NewTimer(time.Second).C:
		err = errors.New("test lasts too long")
	case result = <-results:
	}
	return
}

func equalOkResult(r1, r2 base.AsyncResult) bool {
	return r1.Seq == r2.Seq && reflect.DeepEqual(r1.Result, r2.Result) &&
		r1.Error == nil && r2.Error == nil
}
