package handler_test

import (
	"testing"

	test "github.com/cmd-stream/cmd-stream-go/test/handler"
)

func TestHandler_Handle(t *testing.T) {
	test.RunHandlerTestCase(t, test.HandleSuccessTestCase(t))
	test.RunHandlerTestCase(t, test.SetReceiveDeadlineErrorTestCase(t))
	test.RunHandlerTestCase(t, test.ReceiveErrorTestCase(t))
	test.RunHandlerTestCase(t, test.InvokeErrorTestCase(t))
	test.RunHandlerTestCase(t, test.OptionAtTestCase(t))
	test.RunHandlerTestCase(t, test.OptionCmdReceiveDurationTestCase(t))
	test.RunHandlerTestCase(t, test.CloseWhileInvokingCmdsTestCase(t))
}
