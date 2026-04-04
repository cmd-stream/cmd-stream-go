package handler

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/handler"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type HandlerTestCase[T any] struct {
	Name   string
	Setup  HandlerSetup[T]
	Action func(t *testing.T, h *handler.Handler[T])
	Mocks  []*mok.Mock
}

type HandlerSetup[T any] struct {
	Invoker handler.Invoker[T]
	Opts    []handler.SetOption
}

func RunHandlerTestCase[T any](t *testing.T, tc HandlerTestCase[T]) {
	t.Run(tc.Name, func(t *testing.T) {
		h := handler.New(tc.Setup.Invoker, tc.Setup.Opts...)
		tc.Action(t, h)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}
