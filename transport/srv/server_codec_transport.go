package srv

import (
	"bufio"
	"net"

	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
)

// New creates a new ServerCodecTransport.
func New[T any](conn net.Conn, codec tspt.Codec[core.Result, core.Cmd[T]],
	opts ...tspt.SetOption,
) *ServerCodecTransport[T] {
	o := tspt.Options{}
	tspt.Apply(&o, opts...)
	var (
		w = bufio.NewWriterSize(conn, o.WriterBufSize)
		r = bufio.NewReaderSize(conn, o.ReaderBufSize)
	)
	return &ServerCodecTransport[T]{tspt.New(conn, w, r, codec), w}
}

// ServerCodecTransport implements the delegate.ServerTransport interface.
type ServerCodecTransport[T any] struct {
	*tspt.CodecTransport[core.Result, core.Cmd[T]]
	w tspt.Writer
}

// SendServerInfo transmits the server info to the client.
func (ct *ServerCodecTransport[T]) SendServerInfo(info dlgt.ServerInfo) (
	err error,
) {
	_, err = dlgt.ServerInfoMUS.Marshal(info, ct.w)
	if err != nil {
		return
	}
	return ct.Flush()
}

// WriterBufSize returns the transport's writer buffer size.
func (ct *ServerCodecTransport[T]) WriterBufSize() int {
	return ct.CodecTransport.W.(*bufio.Writer).Size()
}

// ReaderBufSize returns the transport's reader buffer size.
func (ct *ServerCodecTransport[T]) ReaderBufSize() int {
	return ct.CodecTransport.R.(*bufio.Reader).Size()
}
