package server

import (
	"reflect"
	"testing"

	base_mock "github.com/cmd-stream/base-go/testdata/mock"
	transport_common "github.com/cmd-stream/transport-go/common"
	transport_server "github.com/cmd-stream/transport-go/server"
	transport_mock "github.com/cmd-stream/transport-go/testdata/mock"
)

func TestTransportFactory(t *testing.T) {

	t.Run("New should work correctly", func(t *testing.T) {
		var (
			wantConf  = transport_common.Conf{WriterBufSize: 10, ReaderBufSize: 20}
			wantCodec = transport_mock.NewServerCodec()
			wantConn  = base_mock.NewConn()
			factory   = TransportFactory[any]{wantConf, wantCodec}

			transport     = factory.New(wantConn)
			wantTransport = transport_server.New[any](wantConf, wantConn, wantCodec)
		)
		if !reflect.DeepEqual(wantTransport, transport) {
			t.Errorf("unexpected server, want '%v' actual '%v'", wantTransport,
				transport)
		}
	})

}
