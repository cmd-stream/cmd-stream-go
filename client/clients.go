package ccln

import (
	"fmt"
	"net"
)

// MustMakeClients creates a specified number of clients using the provided codec
// and connection factory. If the client creation fails, it panics.
func MustMakeClients[T any](count int, codec Codec[T], factory ConnFactory,
	ops ...SetOption) (clients []Client[T]) {
	clients, err := MakeClients(count, codec, factory, ops...)
	if err != nil {
		panic(err)
	}
	return
}

// MakeClients creates a specified number of clients using the provided codec
// and connection factory. If any of the client creations fail, it returns an
// error with a description of how many clients were successfully created.
// how many clients were successfully created.
func MakeClients[T any](count int, codec Codec[T], factory ConnFactory,
	ops ...SetOption) (clients []Client[T], err error) {
	clients = make([]Client[T], 0, count)
	for i := range count {
		var conn net.Conn
		conn, err = factory.New()
		if err != nil {
			err = fmt.Errorf("only %v clients were created, cause: %w", i, err)
			return
		}
		var client Client[T]
		client, err = New(codec, conn, ops...)
		if err != nil {
			err = fmt.Errorf("only %v clients were created, cause: %w", i, err)
			return
		}
		clients = append(clients, client)
	}
	return
}

// MakeReconnectClients creates a specified number of reconnect clients using the
// provided codec and connection factory. If any of the client creations fail,
// it returns an error with a description of how many clients were successfully
// created.
func MakeReconnectClients[T any](count int, codec Codec[T],
	factory ConnFactory, ops ...SetOption) (c []Client[T], err error) {
	c = make([]Client[T], 0, count)
	for range count {
		var client Client[T]
		client, err = NewReconnect(codec, factory, ops...)
		if err != nil {
			return
		}
		c = append(c, client)
	}
	return
}
