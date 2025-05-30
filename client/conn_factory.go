package cln

import "net"

// ConnFactory establishes a new connection to the server.
type ConnFactory interface {
	New() (net.Conn, error)
}

// ConnFactoryFn is a function implementation of the ConnFactory.
type ConnFactoryFn func() (net.Conn, error)

func (f ConnFactoryFn) New() (net.Conn, error) {
	return f()
}
