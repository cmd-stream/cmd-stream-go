package cln

import (
	"bufio"
	"net"

	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
)

// ClientCodecTransport implements the delegate.ClientTransport interface.
type ClientCodecTransport[T any] struct {
	*tspt.CodecTransport[core.Cmd[T], core.Result]
}

// New creates a new Transport.
func New[T any](conn net.Conn, codec tspt.Codec[core.Cmd[T], core.Result],
	opts ...tspt.SetOption,
) *ClientCodecTransport[T] {
	o := tspt.Options{}
	tspt.Apply(&o, opts...)
	var (
		w = bufio.NewWriterSize(conn, o.WriterBufSize)
		r = bufio.NewReaderSize(conn, o.ReaderBufSize)
	)
	return &ClientCodecTransport[T]{tspt.New(conn, w, r, codec)}
}

// ReceiveServerInfo waits for and returns the server info.
func (t *ClientCodecTransport[T]) ReceiveServerInfo() (info dlgt.ServerInfo,
	err error,
) {
	info, _, err = dlgt.ServerInfoMUS.Unmarshal(t.R)
	return
}

// WriterBufSize returns the transport's writer buffer size.
func (t *ClientCodecTransport[T]) WriterBufSize() int {
	return t.CodecTransport.W.(*bufio.Writer).Size()
}

// ReaderBufSize returns the transport's reader buffer size.
func (t *ClientCodecTransport[T]) ReaderBufSize() int {
	return t.CodecTransport.R.(*bufio.Reader).Size()
}
