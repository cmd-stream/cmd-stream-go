package hooks_test

import (
	"testing"

	test "github.com/cmd-stream/cmd-stream-go/test/sender"
)

func TestCircuitBreakerHooks(t *testing.T) {
	for _, tc := range []test.HooksCircuitBreakerTestCase{
		test.BeforeSendTestCase(t),
		test.BeforeSendErrorTestCase(),
		test.OnErrorTestCase(t),
		test.OnResultTestCase(t),
		test.OnTimeoutTestCase(t),
	} {
		test.RunHooksCircuitBreakerTestCase(t, tc)
	}
}
