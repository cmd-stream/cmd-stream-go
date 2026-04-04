package srv_test

import (
	"testing"

	test "github.com/cmd-stream/cmd-stream-go/test/core"
)

func TestConnReceiver_Run(t *testing.T) {
	for _, tc := range []test.ConnReceiverTestCase{
		test.FirstSetDeadlineErrorTestCase(t),
		test.FirstAcceptErrorTestCase(t),
		test.FirstResetDeadlineErrorTestCase(t),
		test.SecondAcceptErrorTestCase(t),
		test.RunSeveralConnectionsTestCase(t),
		test.StopWhileAcceptingTestCase(t),
		test.ShutdownWhileAcceptingTestCase(t),
		test.StopWhileQueuingTestCase(t),
		test.ShutdownWhileQueuingTestCase(t),
	} {
		test.RunConnReceiverTestCase(t, tc)
	}
}
