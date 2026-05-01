package srv_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/test"
)

func TestServer(t *testing.T) {
	for _, tc := range []test.ServerTestCase{
		test.Server.ServeNoWorkers(t),
		test.Server.ServeSeveralConnections(t),
		test.Server.ShutdownBeforeAccept(t),
		test.Server.ShutdownAfterAccept(t),
		test.Server.CloseBeforeAccept(t),
		test.Server.CloseAfterAccept(t),
		test.Server.ListenAndServeFailOnInvalidAddr(t),
		test.Server.ShutdownFailIfNotServing(t),
		test.Server.CloseFailIfNotServing(t),
		test.Server.NegativeWorkersCount(t),
	} {
		test.RunServerTestCase(t, tc)
	}
}
