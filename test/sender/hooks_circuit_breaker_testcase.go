package sender

import (
	"testing"

	hks "github.com/cmd-stream/cmd-stream-go/sender/hooks"
	hksmock "github.com/cmd-stream/cmd-stream-go/test/mock/sender/hooks"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type HooksCircuitBreakerSetup struct {
	CB    hksmock.CircuitBreaker
	Hooks hksmock.Hooks[any]
}

type HooksCircuitBreakerTestCase struct {
	Name   string
	Setup  HooksCircuitBreakerSetup
	Action func(t *testing.T, h hks.Hooks[any])
	Want   HooksCircuitBreakerWant
}

type HooksCircuitBreakerWant struct {
	Mocks []*mok.Mock
}

func RunHooksCircuitBreakerTestCase(t *testing.T, tc HooksCircuitBreakerTestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		h := hks.NewCircuitBreakerHooks(tc.Setup.CB, tc.Setup.Hooks)
		tc.Action(t, h)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Want.Mocks), mok.EmptyInfomap)
	})
}
