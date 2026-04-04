package integration

import (
	"net"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/cmd-stream-go/core"
	ccln "github.com/cmd-stream/cmd-stream-go/core/cln"
	"github.com/cmd-stream/cmd-stream-go/testkit"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestReconnect(t *testing.T) {
	const addr = "127.0.0.1:9004"

	startReconnectServer(addr)
	factory := &reconnectFactory{addr: addr}
	client, err := cmdstream.NewReconnectClient(testkit.ClientCodec{}, factory)
	assertfatal.EqualError(t, err, nil)

	exchangeReconnect(t, client, factory)
}

type reconnectFactory struct {
	addr string
	conn net.Conn
}

func (f *reconnectFactory) New() (net.Conn, error) {
	conn, err := net.Dial("tcp", f.addr)
	if err != nil {
		return nil, err
	}
	f.conn = conn
	return conn, nil
}

func startReconnectServer(addr string) {
	go func() {
		server, err := cmdstream.NewServer(testkit.Receiver{}, testkit.ServerCodec{})
		if err != nil {
			return
		}
		_ = server.ListenAndServe(addr)
	}()
	time.Sleep(50 * time.Millisecond)
}

func exchangeReconnect(t *testing.T, client *ccln.Client[testkit.Receiver],
	factory *reconnectFactory) {
	var (
		cmd     = testkit.Cmd{ExecTime: 10 * time.Millisecond}
		results = make(chan core.AsyncResult, 1)
	)

	// Send first command.
	seq, _, err := client.Send(cmd, results)
	assertfatal.EqualError(t, err, nil)
	receiveAndAssert(t, results, seq, testkit.Result{LastOneFlag: true})

	// Simulate disconnect.
	err = factory.conn.Close()
	assertfatal.EqualError(t, err, nil)
	// Wait for client to detect disconnect.
	time.Sleep(100 * time.Millisecond)

	// Send second command - should trigger reconnect.
	seq, _, err = client.Send(cmd, results)
	assertfatal.EqualError(t, err, nil)
	receiveAndAssert(t, results, seq, testkit.Result{LastOneFlag: true})
}
