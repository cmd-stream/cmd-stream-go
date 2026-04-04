package integration

import (
	"net"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/cmd-stream-go/core"
	ccln "github.com/cmd-stream/cmd-stream-go/core/cln"
	"github.com/cmd-stream/cmd-stream-go/testkit"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestMultiResult(t *testing.T) {
	const addr = "127.0.0.1:9002"

	startMultiResultServer(t, addr)
	client, err := makeMultiResultClient(addr)
	assertfatal.EqualError(t, err, nil)

	exchangeMultiResult(t, client)
}

func startMultiResultServer(t *testing.T, addr string) {
	go func() {
		server, err := cmdstream.NewServer(testkit.Receiver{}, testkit.ServerCodec{})
		asserterror.EqualError(t, err, nil)
		_ = server.ListenAndServe(addr)
	}()
	time.Sleep(50 * time.Millisecond)
}

func makeMultiResultClient(addr string) (client *ccln.Client[testkit.Receiver],
	err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	return cmdstream.NewClient(testkit.ClientCodec{}, conn)
}

func exchangeMultiResult(t *testing.T, client *ccln.Client[testkit.Receiver]) {
	const resultsCount = 2
	var (
		cmd = testkit.MultiCmd{
			ResultsCount: resultsCount,
			ExecTime:     100 * time.Millisecond,
		}
		results = make(chan core.AsyncResult, resultsCount)
	)
	seq, _, err := client.Send(cmd, results)
	assertfatal.EqualError(t, err, nil)

	for i := range resultsCount {
		result := testkit.Result{LastOneFlag: i == resultsCount-1}
		receiveAndAssert(t, results, seq, result)
	}
}
