package server

import (
	"testing"

	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	tmock "github.com/cmd-stream/cmd-stream-go/test/mock/transport"
	"github.com/cmd-stream/cmd-stream-go/transport"
	tsrv "github.com/cmd-stream/cmd-stream-go/transport/srv"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestTransportFactory_New(t *testing.T) {
	var (
		wantCodec = tmock.NewServerCodec[any]()
		wantConn  = cmock.NewConn()
		factory   = NewTransportFactory(wantCodec,
			[]transport.SetOption{
				transport.WithWriterBufSize(10),
				transport.WithReaderBufSize(20),
			}...)
		transport = factory.New(wantConn).(*tsrv.ServerCodecTransport[any])
	)
	asserterror.Equal(t, transport.WriterBufSize(), 10)
	asserterror.Equal(t, transport.ReaderBufSize(), 20)

}
