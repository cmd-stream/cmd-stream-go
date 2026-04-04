package transport

import (
	"net"
	"testing"

	tspt "github.com/cmd-stream/cmd-stream-go/transport"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type TransportSetup[T, V any] struct {
	Conn   net.Conn
	Writer tspt.Writer
	Reader tspt.Reader
	Codec  tspt.Codec[T, V]
}

type TransportTestCase[T, V any] struct {
	Name   string
	Setup  TransportSetup[T, V]
	Action func(t *testing.T, transport *tspt.CodecTransport[T, V])
	Mocks  []*mok.Mock
}

func RunTransportTestCase[T, V any](t *testing.T, tc TransportTestCase[T, V]) {
	t.Run(tc.Name, func(t *testing.T) {
		transport := tspt.New(tc.Setup.Conn, tc.Setup.Writer, tc.Setup.Reader,
			tc.Setup.Codec)
		tc.Action(t, transport)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}
