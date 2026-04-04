package transport

import (
	"io"

	"github.com/cmd-stream/cmd-stream-go/core"
)

// Codec provides methods to encode/decode data transmitted over a connection.
//
// On the client side, it encodes Commands and decodes Results. On the server
// side, it decodes Commands and encodes Results.
type Codec[T, V any] interface {
	Encode(seq core.Seq, t T, w Writer) (n int, err error)
	Decode(r Reader) (seq core.Seq, v V, n int, err error)
}

// Writer is an interface that extends io.ByteWriter, io.Writer, and
// io.StringWriter. It also includes a Flush method to ensure buffered data is
// written out.
type Writer interface {
	io.ByteWriter
	io.Writer
	io.StringWriter
	Flush() error
}

// Reader is an interface that extends io.Reader and io.ByteReader.
type Reader interface {
	io.Reader
	io.ByteReader
}
