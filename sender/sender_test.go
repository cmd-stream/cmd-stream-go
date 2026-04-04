package sender_test

import (
	"testing"

	test "github.com/cmd-stream/cmd-stream-go/test/sender"
)

func TestSender(t *testing.T) {
	for _, tc := range []test.SenderTestCase[any]{
		test.SendSuccessTestCase(t),
		test.SendBeforeSendErrorTestCase(),
		test.SendGroupErrorTestCase(),
		test.SendTimeoutTestCase(),
		test.SendWithDeadlineSuccessTestCase(t),
		test.SendWithDeadlineBeforeSendErrorTestCase(),
		test.SendWithDeadlineGroupErrorTestCase(),
		test.SendWithDeadlineTimeoutTestCase(),
		test.SendMultiSuccessTestCase(t),
		test.SendMultiBeforeSendErrorTestCase(),
		test.SendMultiGroupErrorTestCase(),
		test.SendMultiTimeoutTestCase(t),
		test.SendMultiWithDeadlineSuccessTestCase(t),
		test.SendMultiWithDeadlineBeforeSendErrorTestCase(),
		test.SendMultiWithDeadlineGroupErrorTestCase(),
		test.SendMultiWithDeadlineTimeoutTestCase(t),
	} {
		test.RunSenderTestCase(t, tc)
	}
}
