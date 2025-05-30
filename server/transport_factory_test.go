package csrv

import (
	"testing"

	bmock "github.com/cmd-stream/base-go/testdata/mock"
	"github.com/cmd-stream/transport-go"
	tser "github.com/cmd-stream/transport-go/server"
	tmock "github.com/cmd-stream/transport-go/testdata/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestTransportFactory(t *testing.T) {

	t.Run("New should work correctly", func(t *testing.T) {
		var (
			wantCodec = tmock.NewServerCodec()
			wantConn  = bmock.NewConn()
			factory   = NewTransportFactory[any](
				wantCodec,
				[]transport.SetOption{
					transport.WithWriterBufSize(10),
					transport.WithReaderBufSize(20),
				}...)
			tran    = factory.New(wantConn)
			serTran = tran.(*tser.Transport[any])
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
