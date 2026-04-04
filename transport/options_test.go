package transport

import (
	"testing"

	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestOptions(t *testing.T) {
	var (
		o                 = Options{}
		wantWriterBufSize = 1
		wantReaderBufSize = 1
	)
	Apply(&o, WithWriterBufSize(wantWriterBufSize), WithReaderBufSize(wantReaderBufSize))
	asserterror.Equal(t, o.WriterBufSize, wantWriterBufSize)
	asserterror.Equal(t, o.ReaderBufSize, wantReaderBufSize)
}
