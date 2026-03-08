package server

import (
	"testing"

	cmock "github.com/cmd-stream/core-go/test/mock"
	"github.com/cmd-stream/transport-go"
	tsrv "github.com/cmd-stream/transport-go/server"
	tsrvmock "github.com/cmd-stream/transport-go/test/mock/server"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestTransportFactory(t *testing.T) {
	t.Run("New should work correctly", func(t *testing.T) {
		var (
			wantCodec = tsrvmock.NewServerCodec()
			wantConn  = cmock.NewConn()
			factory   = NewTransportFactory(
				wantCodec,
				[]transport.SetOption{
					transport.WithWriterBufSize(10),
					transport.WithReaderBufSize(20),
				}...)
			tran    = factory.New(wantConn)
			serTran = tran.(*tsrv.Transport[any])
		)
		asserterror.Equal(t, serTran.WriterBufSize(), 10)
		asserterror.Equal(t, serTran.ReaderBufSize(), 20)

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
