package server_test

import (
	"testing"

	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/cmd-stream-go/test/mock"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
	tsrv "github.com/cmd-stream/cmd-stream-go/transport/srv"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestTransportFactory_New(t *testing.T) {
	var (
		wantCodec = mock.NewServerCodec[any]()
		wantConn  = mock.NewConn()
		factory   = srv.NewTransportFactory(wantCodec,
			[]tspt.SetOption{
				tspt.WithWriterBufSize(10),
				tspt.WithReaderBufSize(20),
			}...)
		transport = factory.New(wantConn).(*tsrv.ServerCodecTransport[any])
	)
	asserterror.Equal(t, transport.WriterBufSize(), 10)
	asserterror.Equal(t, transport.ReaderBufSize(), 20)
}
