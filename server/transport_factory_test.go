package server

import (
	"testing"

	cmocks "github.com/cmd-stream/testkit-go/mocks/core"
	tmocks "github.com/cmd-stream/testkit-go/mocks/transport"
	"github.com/cmd-stream/transport-go"
	tsrv "github.com/cmd-stream/transport-go/server"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestTransportFactory(t *testing.T) {
	t.Run("New should work correctly", func(t *testing.T) {
		var (
			wantCodec = tmocks.NewServerCodec()
			wantConn  = cmocks.NewConn()
			factory   = NewTransportFactory(
				wantCodec,
				[]transport.SetOption{
					transport.WithWriterBufSize(10),
					transport.WithReaderBufSize(20),
				}...)
			tran    = factory.New(wantConn)
			serTran = tran.(*tsrv.Transport[any])
		)
		asserterror.Equal(serTran.WriterBufSize(), 10, t)
		asserterror.Equal(serTran.ReaderBufSize(), 20, t)

		// if serTran.WriterBufSize() != 10 || serTran.ReaderBufSize() != 20 {
		// 	t.Errorf("unexpected Transport.WriterBufSize(), want '%v' actual '%v'",
		// 		10, serTran.WriterBufSize())
		// }
		// if serTran.ReaderBufSize() != 20 {
		// 	t.Errorf("unexpected Transport.ReaderBufSize(), want '%v' actual '%v'",
		// 		10, serTran.ReaderBufSize())
		// }
	})
}
