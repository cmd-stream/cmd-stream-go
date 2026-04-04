package srv_test

import (
	"testing"

	test "github.com/cmd-stream/cmd-stream-go/test/delegate"
)

func TestServerInfoDelegate(t *testing.T) {
	for _, tc := range []test.ServerDelegateTestCase[any]{
		test.ZeroLenServerInfoTestCase(),
		test.SendServerInfoErrorTestCase(),
		test.TransportHandleErrorTestCase(),
		test.SendServerInfoTransportDeadlineErrorTestCase(),
		test.HandleConnSuccessTestCase(t),
	} {
		test.RunServerDelegateTestCase(t, tc)
	}
}
