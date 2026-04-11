package sender

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/sender"
	smock "github.com/cmd-stream/cmd-stream-go/test/mock/sender"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type SenderSetup[T any] struct {
	Group   smock.Group[T]
	Options []sender.SetOption[T]
}

type SenderTestCase[T any] struct {
	Name   string
	Setup  SenderSetup[T]
	Action func(t *testing.T, s sender.Sender[T])
	Mocks  []*mok.Mock
}

func RunSenderTestCase[T any](t *testing.T, tc SenderTestCase[T]) {
	t.Run(tc.Name, func(t *testing.T) {
		s := sender.New(tc.Setup.Group, tc.Setup.Options...)
		tc.Action(t, s)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}
