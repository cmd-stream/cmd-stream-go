package cln

import (
	"bytes"
	"net"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
)

// ClientInfoDelegate implements the core.ClientDelegate interface. It manages
// the initial handshake by expecting ServerInfo from the server immediately
// after the connection is established.
type ClientInfoDelegate[T any] struct {
	transport dlgt.ClientTransport[T]
	options   Options
}

// New creates a new ClientInfoDelegate and performs the ServerInfo handshake.
// It returns ErrServerInfoMismatch if the received ServerInfo does not match
// the expected info.
func New[T any](info dlgt.ServerInfo, transport dlgt.ClientTransport[T], opts ...SetOption) (
	d ClientInfoDelegate[T], err error,
) {
	o := Options{}
	Apply(&o, opts...)
	err = checkServerInfo(o.ServerInfoReceiveDuration, transport, info)
	if err != nil {
		return
	}
	return ClientInfoDelegate[T]{
		transport: transport,
		options:   o,
	}, nil
}

// NewWithoutInfo creates a ClientInfoDelegate without performing the handshake.
// This is primarily intended for use in tests.
func NewWithoutInfo[T any](transport dlgt.ClientTransport[T]) (d ClientInfoDelegate[T]) {
	d.transport = transport
	return
}

// Options returns the delegate's options.
func (d ClientInfoDelegate[T]) Options() Options {
	return d.options
}

// LocalAddr returns the local network address.
func (d ClientInfoDelegate[T]) LocalAddr() net.Addr {
	return d.transport.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (d ClientInfoDelegate[T]) RemoteAddr() net.Addr {
	return d.transport.RemoteAddr()
}

// SetSendDeadline sets the deadline for future Send calls.
func (d ClientInfoDelegate[T]) SetSendDeadline(deadline time.Time) error {
	return d.transport.SetSendDeadline(deadline)
}

// Send transmits a command to the server.
func (d ClientInfoDelegate[T]) Send(seq core.Seq, cmd core.Cmd[T]) (n int, err error) {
	return d.transport.Send(seq, cmd)
}

// Flush flushes the transport's buffer.
func (d ClientInfoDelegate[T]) Flush() error {
	return d.transport.Flush()
}

// SetReceiveDeadline sets the deadline for future Receive calls.
func (d ClientInfoDelegate[T]) SetReceiveDeadline(deadline time.Time) error {
	return d.transport.SetReceiveDeadline(deadline)
}

// Receive waits for and returns the next result from the server.
func (d ClientInfoDelegate[T]) Receive() (seq core.Seq, result core.Result, n int, err error) {
	return d.transport.Receive()
}

// Close closes the underlying transport.
func (d ClientInfoDelegate[T]) Close() error {
	return d.transport.Close()
}

func checkServerInfo[T any](timeout time.Duration,
	transport dlgt.ClientTransport[T],
	wantInfo dlgt.ServerInfo,
) (err error) {
	err = transport.SetReceiveDeadline(calcDeadline(timeout))
	if err != nil {
		return
	}
	info, err := transport.ReceiveServerInfo()
	if err != nil {
		return
	}
	if !bytes.Equal(info, wantInfo) {
		return ErrServerInfoMismatch
	}
	return transport.SetReceiveDeadline(time.Time{})
}

func calcDeadline(duration time.Duration) (deadline time.Time) {
	if duration != 0 {
		deadline = time.Now().Add(duration)
	}
	return
}
