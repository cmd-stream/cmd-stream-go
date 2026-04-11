package core

import (
	"context"
	"net"
	"sync"
	"time"
)

// ClientDelegate helps the client to send Commands and receive Results.
type ClientDelegate[T any] interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr

	SetSendDeadline(deadline time.Time) error
	Send(seq Seq, cmd Cmd[T]) (n int, err error)
	Flush() error

	SetReceiveDeadline(deadline time.Time) error
	Receive() (seq Seq, result Result, n int, err error)

	Close() error
}

// KeepaliveDelegate defines the Keepalive method.
//
// This delegate can be used if you want the client to keepalive connection to
// the server.
type KeepaliveDelegate[T any] interface {
	ClientDelegate[T]
	Keepalive(muSn *sync.Mutex)
}

// ReconnectDelegate defines the Reconnect method.
//
// This delegate can be used if you want the client to reconnect to the server
// in case of a connection loss.
type ReconnectDelegate[T any] interface {
	ClientDelegate[T]
	Reconnect() error
}

type ClientDelegateUnwrapper[T any] interface {
	Unwrap() ClientDelegate[T]
}

// ServerDelegate helps the server handle incoming connections.
//
// Handle method should return context.Canceled error if the context is done.
type ServerDelegate interface {
	Handle(ctx context.Context, conn net.Conn) error
}
