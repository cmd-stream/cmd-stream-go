package cser

import (
	"testing"

	bmock "github.com/cmd-stream/base-go/testdata/mock"
	tcom "github.com/cmd-stream/transport-go/common"
	tser "github.com/cmd-stream/transport-go/server"
	tmock "github.com/cmd-stream/transport-go/testdata/mock"
)

func TestTransportFactory(t *testing.T) {

	t.Run("New should work correctly", func(t *testing.T) {
		var (
			wantCodec = tmock.NewServerCodec()
			wantConn  = bmock.NewConn()
			factory   = TransportFactory[any]{Codec: wantCodec, Ops: []tcom.SetOption{
				tcom.WithWriterBufSize(10),
				tcom.WithReaderBufSize(20),
			}}
			tn    = factory.New(wantConn)
			serTn = tn.(*tser.Transport[any])
		)
		if serTn.WriterBufSize() != 10 || serTn.ReaderBufSize() != 20 {
			t.Errorf("unexpected Transport.WriterBufSize(), want '%v' actual '%v'",
				10, serTn.WriterBufSize())
		}
		if serTn.ReaderBufSize() != 20 {
			t.Errorf("unexpected Transport.ReaderBufSize(), want '%v' actual '%v'",
				10, serTn.ReaderBufSize())
		}
	})

}
