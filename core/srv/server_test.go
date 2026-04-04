package srv_test

import (
	"testing"

	test "github.com/cmd-stream/cmd-stream-go/test/core"
)

func TestServer(t *testing.T) {
	for _, tc := range []test.ServerTestCase{
		test.ServeNoWorkersTestCase(t),
		test.ServeSeveralConnectionsTestCase(t),
		test.ServerShutdownBeforeAcceptTestCase(t),
		test.ServerShutdownAfterAcceptTestCase(t),
		test.ServerCloseBeforeAcceptTestCase(t),
		test.ServerCloseAfterAcceptTestCase(t),
		test.ListenAndServeFailOnInvalidAddrTestCase(t),
		test.ServerShutdownFailIfNotServingTestCase(t),
		test.ServerCloseFailIfNotServingTestCase(t),
		test.ServerNegativeWorkersCountTestCase(t),
	} {
		test.RunServerTestCase(t, tc)
	}
}
