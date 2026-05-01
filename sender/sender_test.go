package sender_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/test"
)

func TestSender(t *testing.T) {
	s := test.SenderSuite[any]{}
	for _, tc := range []test.SenderTestCase[any]{
		s.SendSuccess(t),
		s.SendBeforeSendError(t),
		s.SendGroupError(t),
		s.SendTimeout(t),
		s.SendWithDeadlineSuccess(t),
		s.SendWithDeadlineBeforeSendError(t),
		s.SendWithDeadlineGroupError(t),
		s.SendWithDeadlineTimeout(t),
		s.SendMultiSuccess(t),
		s.SendMultiBeforeSendError(t),
		s.SendMultiGroupError(t),
		s.SendMultiTimeout(t),
		s.SendMultiWithDeadlineSuccess(t),
		s.SendMultiWithDeadlineBeforeSendError(t),
		s.SendMultiWithDeadlineGroupError(t),
		s.SendMultiWithDeadlineTimeout(t),
	} {
		test.RunSenderTestCase(t, tc)
	}
}
