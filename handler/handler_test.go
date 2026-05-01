package handler_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/test"
)

func TestHandler_Handle(t *testing.T) {
	g := test.HandlerSuite[any]{}
	for _, tc := range []test.HandlerTestCase[any]{
		g.HandleSuccess(t),
		g.SetReceiveDeadlineError(t),
		g.ReceiveError(t),
		g.InvokeError(t),
		g.OptionAt(t),
		g.OptionCmdReceiveDuration(t),
		g.CloseWhileInvokingCmds(t),
	} {
		test.RunHandlerTestCase(t, tc)
	}
}
