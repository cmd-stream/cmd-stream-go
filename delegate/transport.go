package delegate

import (
	"context"
	"net"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
)

// Transport is a common transport for the client and server delegates.
type Transport[T, V any] interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr

	SetSendDeadline(deadline time.Time) error
	Send(seq core.Seq, t T) (n int, err error)
	Flush() error

	SetReceiveDeadline(deadline time.Time) error
	Receive() (seq core.Seq, v V, n int, err error)

	Close() error
}

// ClientTransportFactory is a factory which creates a ClientTransport.
type ClientTransportFactory[T any] interface {
	New() (ClientTransport[T], error)
}

// ClientTransport is a transport for the client delegate.
//
// It is used by the delegate to send Commands and receive Results.
type ClientTransport[T any] interface {
	Transport[core.Cmd[T], core.Result]
	ReceiveServerInfo() (info ServerInfo, err error)
}

// ServerTransportFactory is a factory which creates a ServerTransport for the
// server delegate.
type ServerTransportFactory[T any] interface {
	New(conn net.Conn) ServerTransport[T]
}

// ServerTransport is a transport for the server delegate.
//
// It is used by the delegate to receive Commands and send Results.
type ServerTransport[T any] interface {
	Transport[core.Result, core.Cmd[T]]
	SendServerInfo(info ServerInfo) error
}

// ServerTransportHandler is a handler of the ServerTransport.
type ServerTransportHandler[T any] interface {
	Handle(ctx context.Context, transport ServerTransport[T]) error
}
