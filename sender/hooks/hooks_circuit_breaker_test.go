package hooks_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/test"
)

func TestCircuitBreakerHooks(t *testing.T) {
	s := test.SenderHooksCircuitBreaker[any]{}
	for _, tc := range []test.HooksCircuitBreakerTestCase[any]{
		s.BeforeSend(t),
		s.BeforeSendError(t),
		s.OnError(t),
		s.OnResult(t),
		s.OnTimeout(t),
	} {
		test.RunHooksCircuitBreakerTestCase(t, tc)
	}
}
