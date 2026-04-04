package delegate

import (
	"testing"

	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/cmd-stream/cmd-stream-go/delegate/srv"
	dmock "github.com/cmd-stream/cmd-stream-go/test/mock/delegate"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type ServerDelegateTestCase[T any] struct {
	Name   string
	Setup  ServerDelegateSetup[T]
	Action func(t *testing.T, d srv.ServerInfoDelegate[T], initErr error)
	Mocks  []*mok.Mock
}

type ServerDelegateSetup[T any] struct {
	Info    dlgt.ServerInfo
	Factory dmock.ServerTransportFactory[T]
	Handler dmock.ServerTransportHandler[T]
	Opts    []srv.SetOption
}

func RunServerDelegateTestCase[T any](t *testing.T, tc ServerDelegateTestCase[T]) {
	t.Run(tc.Name, func(t *testing.T) {
		d, err := srv.New(tc.Setup.Info, tc.Setup.Factory, tc.Setup.Handler,
			tc.Setup.Opts...)

		tc.Action(t, d, err)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}
