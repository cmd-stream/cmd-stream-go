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
// It adapts the provided codec and applies optional configuration for transport,
// delegation, and keepalive behavior.
//
// Parameters:
//   - codec: Codec used for encoding Commands / decoding Results.
//   - conn: The underlying network connection.
//   - ops: Optional configuration settings.
//
// Returns an error if setup fails at any step (e.g., delegate creation).
func MakeClient[T any](codec cln.Codec[T], conn net.Conn, ops ...cln.SetOption) (
	client *ccln.Client[T], err error,
) {
	o := cln.Options{
		Info: srv.ServerInfo,
	}
	cln.Apply(ops, &o)
	var (
		delegate     ccln.Delegate[T]
		adaptedCodec = cln.AdaptCodec(codec, o)
		transport    = tcln.New(conn, adaptedCodec, o.Transport...)
	)
	delegate, err = dcln.New(o.Info, transport, o.Delegate...)
	if err != nil {
		return
	}
	if o.Keepalive != nil {
		delegate = dcln.NewKeepalive(delegate, o.Keepalive...)
	}
	client = ccln.New(delegate, o.Base...)
	return
}

// MakeReconnectClient creates a new cmd-stream client with support for
// automatic reconnection (which occurs if the Codec.Decode method encounters a
// network error).
//
// It sets up a reconnect-aware delegate using the provided codec, connection
// factory, and optional configuration settings.
//
// Parameters:
//   - codec: Codec used for encoding Commands and decoding Results.
//   - factory: Connection factory used to establish new connections.
//   - ops: Optional client configuration settings.
//
// Returns an error if setup fails at any step (e.g., during delegate creation).
func MakeReconnectClient[T any](codec cln.Codec[T], factory cln.ConnFactory,
	ops ...cln.SetOption,
) (client *ccln.Client[T], err error) {
	o := cln.Options{
		Info: srv.ServerInfo,
	}
	cln.Apply(ops, &o)
	var (
		delegate         ccln.Delegate[T]
		adaptedCodec     = cln.AdaptCodec(codec, o)
		transportFactory = cln.NewTransportFactory(adaptedCodec, factory,
			o.Transport...)
	)
	delegate, err = dcln.NewReconnect(o.Info, transportFactory, o.Delegate...)
	if err != nil {
		return
	}
	if o.Keepalive != nil {
		delegate = dcln.NewKeepalive(delegate, o.Keepalive...)
	}
	client = ccln.New(delegate, o.Base...)
	return
}

// MakeClientGroup creates a new ClientGroup with the specified number of clients.
//
// If the Reconnect option is enabled, reconnect-capable clients are created.
// The default dispatch strategy is round-robin.
//
// Parameters:
//   - clientsCount: Number of clients to create.
//   - codec: Codec used for encoding Commands and decoding Results.
//   - factory: Connection factory for establishing client connections.
//   - ops: Optional group-level configuration (e.g., dispatch strategy,
//     reconnect, client options).
//
// If client creation fails, the function returns an error along with a group
// containing the successfully created clients.
func MakeClientGroup[T any](clientsCount int, codec cln.Codec[T],
	factory cln.ConnFactory,
	ops ...grp.SetOption[T],
) (group grp.ClientGroup[T], err error) {
	o := grp.Options[T]{
		Factory: grp.RoundRobinStrategyFactory[T]{},
	}
	grp.ApplyGroup(ops, &o)
	var clients []grp.Client[T]
	if o.Reconnect {
		clients, err = makeReconnectClients(clientsCount, codec, factory, o.ClientOps...)
	} else {
		clients, err = makeClients(clientsCount, codec, factory, o.ClientOps...)
	}
	strategy := o.Factory.New(clients)
	group = grp.NewClientGroup(strategy)
	return
}

// MakeServer creates a new cmd-stream server.
//
// It applies optional configuration to initialize transport, handler, and
// delegate components before creating the server instance.
//
// Parameters:
//   - codec: Codec used for decoding incoming Commands and encoding outgoing
//     Results.
//   - invoker: Executes the Commands.
//   - ops: Optional server configuration (e.g., transport, handler, delegate,
//     core settings).
//
// Returns a fully initialized server.
func MakeServer[T any](codec srv.Codec[T], invoker handler.Invoker[T],
	ops ...srv.SetOption,
) *csrv.Server {
	o := srv.Options{
		Info: srv.ServerInfo,
	}
	srv.Apply(ops, &o)
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
