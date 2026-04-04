package delegate

import (
	"testing"

	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/cmd-stream/cmd-stream-go/delegate/cln"
	dmock "github.com/cmd-stream/cmd-stream-go/test/mock/delegate"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type ClientDelegateTestCase[T any] struct {
	Name   string
	Setup  ClientDelegateSetup[T]
	Action func(t *testing.T, d cln.ClientInfoDelegate[T], initErr error)
	Mocks  []*mok.Mock
}

type ClientDelegateSetup[T any] struct {
	Info      dlgt.ServerInfo
	Transport dmock.ClientTransport[T]
	Opts      []cln.SetOption
}

func RunClientDelegateTestCase[T any](t *testing.T, tc ClientDelegateTestCase[T]) {
	t.Run(tc.Name, func(t *testing.T) {
		d, err := cln.New(tc.Setup.Info, tc.Setup.Transport, tc.Setup.Opts...)

		tc.Action(t, d, err)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}
