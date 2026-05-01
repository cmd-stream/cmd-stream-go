package transport_test

import (
	"testing"

	tspt "github.com/cmd-stream/cmd-stream-go/transport"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestOptions(t *testing.T) {
	var (
		o                 = tspt.Options{}
		wantWriterBufSize = 1
		wantReaderBufSize = 1
	)
	tspt.Apply(&o, tspt.WithWriterBufSize(wantWriterBufSize), tspt.WithReaderBufSize(wantReaderBufSize))
	asserterror.Equal(t, o.WriterBufSize, wantWriterBufSize)
	asserterror.Equal(t, o.ReaderBufSize, wantReaderBufSize)
}
