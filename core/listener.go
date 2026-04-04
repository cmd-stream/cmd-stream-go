package core

import (
	"net"
	"time"
)

// Listener is a server network listener.
//
// On Close, it should not close already accepted connections.
type Listener interface {
	Addr() net.Addr
	SetDeadline(time.Time) error
	Accept() (net.Conn, error)
	Close() error
}
