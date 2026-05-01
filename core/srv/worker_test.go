package srv_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/test"
)

func TestWorker_Run(t *testing.T) {
	for _, tc := range []test.WorkerTestCase{
		test.Worker.RunSeveralConnections(t),
		test.Worker.LostConnCallback(t),
		test.Worker.Stop(t),
		test.Worker.Shutdown(t),
	} {
		test.RunWorkerTestCase(t, tc)
	}
}
