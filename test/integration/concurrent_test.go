package integration

import (
	"net"
	"sync"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/cmd-stream-go/core"
	ccln "github.com/cmd-stream/cmd-stream-go/core/cln"
	"github.com/cmd-stream/cmd-stream-go/testkit"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestConcurrent(t *testing.T) {
	const addr = "127.0.0.1:9003"

	startConcurrentServer(t, addr)
	client, err := makeConcurrentClient(addr)
	assertfatal.EqualError(t, err, nil)

	exchangeConcurrent(t, client)
}

func startConcurrentServer(t *testing.T, addr string) {
	go func() {
		server, err := cmdstream.NewServer(testkit.Receiver{}, testkit.ServerCodec{})
		asserterror.EqualError(t, err, nil)
		_ = server.ListenAndServe(addr)
	}()
	time.Sleep(50 * time.Millisecond)
}

func makeConcurrentClient(addr string) (client *ccln.Client[testkit.Receiver],
	err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	return cmdstream.NewClient(testkit.ClientCodec{}, conn)
}

func exchangeConcurrent(t *testing.T, client *ccln.Client[testkit.Receiver]) {
	var (
		wg            sync.WaitGroup
		clientsCount  = 10
		cmdsPerClient = 5
	)

	wg.Add(clientsCount)
	for range clientsCount {
		go func() {
			defer wg.Done()
			for range cmdsPerClient {
				var (
					cmd     = testkit.Cmd{ExecTime: 10 * time.Millisecond}
					results = make(chan core.AsyncResult, 1)
				)
				seq, _, err := client.Send(cmd, results)
				asserterror.EqualError(t, err, nil)

				receiveAndAssert(t, results, seq, testkit.Result{LastOneFlag: true})
			}
		}()
	}
	wg.Wait()
}
