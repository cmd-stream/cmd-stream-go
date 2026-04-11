package integration

import (
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	"github.com/cmd-stream/cmd-stream-go/core"
	dcln "github.com/cmd-stream/cmd-stream-go/delegate/cln"
	"github.com/cmd-stream/cmd-stream-go/testkit"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestKeepaliveReconnect(t *testing.T) {
	const (
		addr               = "127.0.0.1:9003"
		clientKeepalive    = 100 * time.Millisecond
		inactivityDuration = 400 * time.Millisecond
	)

	// Start server with a receive timeout.
	startKeepaliveServer(t, addr)

	factory := &reconnectFactory{addr: addr}
	client, err := cmdstream.NewReconnectClient(testkit.ClientCodec{}, factory,
		cln.WithKeepalive(
			dcln.WithKeepaliveTime(clientKeepalive),
			dcln.WithKeepaliveIntvl(clientKeepalive),
		),
	)
	assertfatal.EqualError(t, err, nil)

	var (
		cmd     = testkit.Cmd{ExecTime: 10 * time.Millisecond}
		results = make(chan core.AsyncResult, 1)
	)

	// 1. Initial exchange.
	seq, _, err := client.Send(cmd, results)
	assertfatal.EqualError(t, err, nil)
	receiveAndAssert(t, results, seq, testkit.Result{LastOneFlag: true})

	// 2. Wait for a duration longer than the server's timeout.
	// Keepalive should keep the connection alive.
	time.Sleep(inactivityDuration)

	// 3. Verify the connection is still alive by sending another command.
	seq, _, err = client.Send(cmd, results)
	assertfatal.EqualError(t, err, nil)
	receiveAndAssert(t, results, seq, testkit.Result{LastOneFlag: true})

	// 4. Simulate disconnect.
	err = factory.conn.Close()
	assertfatal.EqualError(t, err, nil)

	// 5. Wait for client to detect disconnect.
	time.Sleep(200 * time.Millisecond)

	// 6. Final exchange - should trigger reconnect through KeepaliveDelegate.
	seq, _, err = client.Send(cmd, results)
	assertfatal.EqualError(t, err, nil)
	receiveAndAssert(t, results, seq, testkit.Result{LastOneFlag: true})
}
