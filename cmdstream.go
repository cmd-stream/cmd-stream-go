// Package cmdstream provides high-level factory functions for creating a
// cmd-stream client, server or client group.
//
// It integrates components from core-go, delegate-go, handler-go, and
// transport-go modules to assemble fully configured cmd-stream system with
// optional support for reconnection and keepalive.
package cmdstream

import (
	"net"

	cln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	ccln "github.com/cmd-stream/core-go/client"
	csrv "github.com/cmd-stream/core-go/server"
	dcln "github.com/cmd-stream/delegate-go/client"
	dsrv "github.com/cmd-stream/delegate-go/server"
	"github.com/cmd-stream/handler-go"
	tcln "github.com/cmd-stream/transport-go/client"
)

// MakeClient creates and initializes a new cmd-stream client.
//
// Parameters:
//   - codec: Codec used for encoding Commands / decoding Results.
//   - conn: The underlying network connection.
//   - opts: Optional configuration settings.
//
// Returns an error if setup fails at any step (e.g., delegate creation).
func MakeClient[T any](codec cln.Codec[T], conn net.Conn, opts ...cln.SetOption) (
	client *ccln.Client[T], err error,
) {
	o := cln.Options{
		Info: srv.ServerInfo,
	}
	cln.Apply(opts, &o)
	var (
		delegate     ccln.Delegate[T]
		adaptedCodec = cln.AdaptCodec(codec, o)
		transport    = tcln.New(conn, adaptedCodec, o.Transport...)
	)
	delegate, err = dcln.New(o.Info, transport, o.Delegate...)
	if err != nil {
		err = NewCmdStreamError(err)
		return
	}
	if o.Keepalive != nil {
		delegate = dcln.NewKeepalive(delegate, o.Keepalive...)
	}
	client = ccln.New(delegate, o.Base...)
	return
}

// MakeReconnectClient creates a new cmd-stream client with support for
// automatic reconnection.
//
// Parameters:
//   - codec: Codec used for encoding Commands and decoding Results.
//   - factory: Connection factory used to establish new connections.
//   - opts: Optional client configuration settings.
//
// Returns an error if setup fails at any step (e.g., during delegate creation).
func MakeReconnectClient[T any](codec cln.Codec[T], factory cln.ConnFactory,
	opts ...cln.SetOption,
) (client *ccln.Client[T], err error) {
	o := cln.Options{
		Info: srv.ServerInfo,
	}
	cln.Apply(opts, &o)
	var (
		delegate         ccln.Delegate[T]
		adaptedCodec     = cln.AdaptCodec(codec, o)
		transportFactory = cln.NewTransportFactory(adaptedCodec, factory,
			o.Transport...)
	)
	delegate, err = dcln.NewReconnect(o.Info, transportFactory, o.Delegate...)
	if err != nil {
		err = NewCmdStreamError(err)
		return
	}
	if o.Keepalive != nil {
		delegate = dcln.NewKeepalive(delegate, o.Keepalive...)
	}
	client = ccln.New(delegate, o.Base...)
	return
}

// MakeClientGroup creates a new group.ClientGroup with the specified number of
// clients.
//
// If the Reconnect option is enabled, reconnect-capable clients are created.
// The default dispatch strategy is round-robin.
//
// Parameters:
//   - clientsCount: Number of clients to create.
//   - codec: Codec used for encoding Commands and decoding Results.
//   - factory: Connection factory for establishing client connections.
//   - opts: Optional group-level configuration (e.g., dispatch strategy,
//     reconnect, client options).
//
// If client creation fails, the function returns an error along with a group
// containing the successfully created clients.
func MakeClientGroup[T any](clientsCount int, codec cln.Codec[T],
	factory cln.ConnFactory,
	opts ...grp.SetOption[T],
) (group grp.ClientGroup[T], err error) {
	o := grp.Options[T]{
		Factory: grp.RoundRobinStrategyFactory[T]{},
	}
	grp.ApplyGroup(opts, &o)
	var clients []grp.Client[T]
	if o.Reconnect {
		clients, err = makeReconnectClients(clientsCount, codec, factory, o.ClientOps...)
	} else {
		clients, err = makeClients(clientsCount, codec, factory, o.ClientOps...)
	}
	if err != nil {
		err = NewCmdStreamError(err)
	}
	strategy := o.Factory.New(clients)
	group = grp.NewClientGroup(strategy)
	return
}

// MakeServer creates a new cmd-stream server.
//
// Parameters:
//   - codec: Codec used for decoding incoming Commands and encoding outgoing
//     Results.
//   - invoker: Executes the Commands.
//   - opts: Optional server configuration (e.g., transport, handler, delegate,
//     core settings).
//
// Returns a fully initialized server.
func MakeServer[T any](codec srv.Codec[T], invoker handler.Invoker[T],
	opts ...srv.SetOption,
) *csrv.Server {
	o := srv.Options{
		Info: srv.ServerInfo,
	}
	srv.Apply(opts, &o)
	var (
		f = srv.NewTransportFactory(srv.AdaptCodec(codec, o), o.Transport...)
		h = handler.New(invoker, o.Handler...)
		d = dsrv.New(o.Info, f, h, o.Delegate...)
	)
	return csrv.New(d, o.Base...)
}

func makeClients[T any](count int, codec cln.Codec[T], factory cln.ConnFactory,
	ops ...cln.SetOption,
) (clients []grp.Client[T], err error) {
	clients = make([]grp.Client[T], 0, count)
	var conn net.Conn
	for range count {
		conn, err = factory.New()
		if err != nil {
			err = NewMakeClientsError(count, err)
			return
		}
		var client *ccln.Client[T]
		client, err = MakeClient(codec, conn, ops...)
		if err != nil {
			err = NewMakeClientsError(count, err)
			return
		}
		clients = append(clients, client)
	}
	return
}

func makeReconnectClients[T any](count int, codec cln.Codec[T],
	factory cln.ConnFactory, ops ...cln.SetOption,
) (c []grp.Client[T], err error) {
	c = make([]grp.Client[T], 0, count)
	for range count {
		var client *ccln.Client[T]
		client, err = MakeReconnectClient(codec, factory, ops...)
		if err != nil {
			return
		}
		c = append(c, client)
	}
	return
}
