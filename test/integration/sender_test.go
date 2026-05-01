package integration

import (
	"context"
	"sync"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	csrv "github.com/cmd-stream/cmd-stream-go/core/srv"
	sndr "github.com/cmd-stream/cmd-stream-go/sender"
	"github.com/cmd-stream/cmd-stream-go/testkit"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestSender(t *testing.T) {
	const addr = "127.0.0.1:9007"

	server := startSenderServer(t, addr)
	defer server.Close()

	sender, err := makeSender(addr)
	assertfatal.EqualError(t, err, nil)
	defer func() {
		err := sender.Close()
		asserterror.EqualError(t, err, nil)
	}()
	exchangeSender(t, sender)
}

func TestConcurrentSender(t *testing.T) {
	const addr = "127.0.0.1:9008"

	server := startSenderServer(t, addr)
	defer server.Close()
	sender, err := makeSender(addr)
	assertfatal.EqualError(t, err, nil)
	defer func() {
		err := sender.Close()
		asserterror.EqualError(t, err, nil)
	}()

	exchangeConcurrentSender(t, sender)
}

func startSenderServer(t *testing.T, addr string) *csrv.Server {
	server, err := cmdstream.NewServer(testkit.Receiver{}, testkit.ServerCodec{})
	asserterror.EqualError(t, err, nil)
	go func() {
		_ = server.ListenAndServe(addr)
	}()
	time.Sleep(50 * time.Millisecond)
	return server
}

func makeSender(addr string) (sender sndr.Sender[testkit.Receiver], err error) {
	return cmdstream.NewSender(addr, testkit.ClientCodec{},
		sndr.WithClientsCount[testkit.Receiver](2),
	)
}

func exchangeSender(t *testing.T, sender sndr.Sender[testkit.Receiver]) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cmd := testkit.Cmd{ExecTime: 10 * time.Millisecond}
	expectedResult := testkit.Result{LastOneFlag: true}

	// Send multiple commands to verify the sender works.
	for range 3 {
		result, err := sender.Send(ctx, cmd)
		assertfatal.EqualError(t, err, nil)
		asserterror.EqualDeep(t, result, expectedResult)
	}
}

func exchangeConcurrentSender(t *testing.T, sender sndr.Sender[testkit.Receiver]) {
	var (
		wg            sync.WaitGroup
		clientsCount  = 10
		cmdsPerClient = 5
	)

	wg.Add(clientsCount)
	for range clientsCount {
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			cmd := testkit.Cmd{ExecTime: 10 * time.Millisecond}
			expectedResult := testkit.Result{LastOneFlag: true}

			for range cmdsPerClient {
				result, err := sender.Send(ctx, cmd)
				asserterror.EqualError(t, err, nil)
				asserterror.EqualDeep(t, result, expectedResult)
			}
		}()
	}
	wg.Wait()
}
