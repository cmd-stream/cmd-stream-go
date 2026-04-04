// Package transport provides a generic, codec-driven abstraction for sending
// and receiving messages over a network connection.
package transport

import (
	"net"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
)

// CodecTransport is a common transport for both client and server.
//
// It uses a user-defined codec to encode/decode data over the connection.
type CodecTransport[T, V any] struct {
	W     Writer
	R     Reader
	conn  net.Conn
	codec Codec[T, V]
}

// New creates a new CodecTransport.
func New[T, V any](conn net.Conn, w Writer, r Reader,
	codec Codec[T, V],
) *CodecTransport[T, V] {
	return &CodecTransport[T, V]{w, r, conn, codec}
}

// LocalAddr returns the connection local network address.
func (ct *CodecTransport[T, V]) LocalAddr() net.Addr {
	return ct.conn.LocalAddr()
}

// RemoteAddr returns the connection remote network address.
func (ct *CodecTransport[T, V]) RemoteAddr() net.Addr {
	return ct.conn.RemoteAddr()
}

// SetSendDeadline sets a send deadline.
func (ct *CodecTransport[T, V]) SetSendDeadline(deadline time.Time) error {
	return ct.conn.SetWriteDeadline(deadline)
}

// Send sends data using the codec.
func (ct *CodecTransport[T, V]) Send(seq core.Seq, t T) (n int, err error) {
	return ct.codec.Encode(seq, t, ct.W)
}

// Flush flushes any buffered data.
func (ct *CodecTransport[T, V]) Flush() (err error) {
	return ct.W.Flush()
}

// SetReceiveDeadline sets a receive deadline.
func (ct *CodecTransport[T, V]) SetReceiveDeadline(deadline time.Time) error {
	return ct.conn.SetReadDeadline(deadline)
}

// Receive receives data using the codec.
func (ct *CodecTransport[T, V]) Receive() (seq core.Seq, v V, n int, err error) {
	return ct.codec.Decode(ct.R)
}

// Close closes the underlying connection.
func (ct *CodecTransport[T, V]) Close() error {
	return ct.conn.Close()
}
