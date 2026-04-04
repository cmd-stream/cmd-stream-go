package srv_test

import (
	"testing"

	test "github.com/cmd-stream/cmd-stream-go/test/core"
)

func TestWorker_Run(t *testing.T) {
	for _, tc := range []test.WorkerTestCase{
		test.WorkerRunSeveralConnectionsTestCase(t),
		test.WorkerLostConnCallbackTestCase(t),
		test.WorkerStopTestCase(t),
		test.WorkerShutdownTestCase(t),
	} {
		test.RunWorkerTestCase(t, tc)
	}
}
