package srv_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/test"
)

func TestConnReceiver_Run(t *testing.T) {
	for _, tc := range []test.ConnReceiverTestCase{
		test.ConnReceiver.FirstSetDeadlineError(t),
		test.ConnReceiver.FirstAcceptError(t),
		test.ConnReceiver.FirstResetDeadlineError(t),
		test.ConnReceiver.SecondAcceptError(t),
		test.ConnReceiver.RunSeveralConnections(t),
		test.ConnReceiver.StopWhileAccepting(t),
		test.ConnReceiver.ShutdownWhileAccepting(t),
		test.ConnReceiver.StopWhileQueuing(t),
		test.ConnReceiver.ShutdownWhileQueuing(t),
	} {
		test.RunConnReceiverTestCase(t, tc)
	}
}
