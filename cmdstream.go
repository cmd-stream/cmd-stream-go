// Package cmdstream provides high-level APIs for creating and managing
// cmd-stream client and server.
package cmdstream

import (
	"crypto/tls"
	"net"

	cln "github.com/cmd-stream/cmd-stream-go/client"
	"github.com/cmd-stream/cmd-stream-go/core"
	ccln "github.com/cmd-stream/cmd-stream-go/core/cln"
	csrv "github.com/cmd-stream/cmd-stream-go/core/srv"
	dcln "github.com/cmd-stream/cmd-stream-go/delegate/cln"
	dsrv "github.com/cmd-stream/cmd-stream-go/delegate/srv"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	hdlr "github.com/cmd-stream/cmd-stream-go/handler"
	sndr "github.com/cmd-stream/cmd-stream-go/sender"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	tcln "github.com/cmd-stream/cmd-stream-go/transport/cln"
)

// NewClient creates and initializes a Client with the specified codec and
// connection. It simplifies the setup process by automatically assembling the
// underlying transport and delegate layers.
func NewClient[T any](codec cln.Codec[T], conn net.Conn,
	opts ...cln.SetOption,
) (client *ccln.Client[T], err error) {
	o := cln.DefaultOptions()
	cln.Apply(&o, opts...)
	var (
		delegate     core.ClientDelegate[T]
		adaptedCodec = cln.AdaptCodec(codec, o)
		transport    = tcln.New(conn, adaptedCodec, o.Transport...)
	)
	delegate, err = dcln.New(o.Info, transport, o.Delegate...)
	if err != nil {
		err = core.NewError(err)
		return
	}
	if o.Keepalive != nil {
		delegate = dcln.NewKeepalive(delegate, o.Keepalive...)
	}
	client = ccln.New(delegate, o.Base...)
	return
}

// NewReconnectClient creates a Client with automatic reconnection support. It
// uses the provided factory to establish new connections if the current one is
// lost.
func NewReconnectClient[T any](codec cln.Codec[T],
	factory cln.ConnFactory, opts ...cln.SetOption,
) (client *ccln.Client[T], err error) {
	o := cln.DefaultOptions()
	cln.Apply(&o, opts...)
	var (
		delegate         core.ClientDelegate[T]
		adaptedCodec     = cln.AdaptCodec(codec, o)
		transportFactory = cln.NewTransportFactory(adaptedCodec, factory, o.Transport...)
	)
	delegate, err = dcln.NewReconnect(o.Info, transportFactory, o.Delegate...)
	if err != nil {
		err = core.NewError(err)
		return
	}
	if o.Keepalive != nil {
		delegate = dcln.NewKeepalive(delegate, o.Keepalive...)
	}
	client = ccln.New(delegate, o.Base...)
	return
}

// NewServer creates a Server that handles incoming Commands using the provided
// receiver. It automatically configures the transport layer and command handler
// using the specified codec and options.
func NewServer[T any](receiver T, codec srv.Codec[T],
	opts ...srv.SetOption,
) (server *csrv.Server, err error) {
	return NewServerWithInvoker(srv.NewInvoker(receiver), codec, opts...)
}

// NewServerWithInvoker creates a new Server that uses the specified invoker to
// handle incoming Commands. It automatically configures the transport layer and
// command handler using the provided codec and options.
func NewServerWithInvoker[T any](invoker hdlr.Invoker[T], codec srv.Codec[T],
	opts ...srv.SetOption,
) (server *csrv.Server, err error) {
	o := srv.DefaultOptions()
	srv.Apply(&o, opts...)
	var (
		h                = hdlr.New(invoker, o.Handler...)
		adaptedCodec     = srv.AdaptCodec(codec, o)
		transportFactory = srv.NewTransportFactory(adaptedCodec, o.Transport...)
	)
	delegate, err := dsrv.New(o.Info, transportFactory, h, o.Delegate...)
	if err != nil {
		err = core.NewError(err)
		return
	}
	return csrv.New(delegate, o.Core...)
}

// NewGroup creates a Group containing the specified number of clients. It uses
// a dispatch strategy (Round-Robin by default) to distribute Commands across
// the pool of clients.
func NewGroup[T any](clientsCount int, codec cln.Codec[T],
	factory cln.ConnFactory, opts ...grp.SetOption[T],
) (group grp.Group[T], err error) {
	if clientsCount <= 0 {
		err = core.NewError(grp.ErrInvalidClientsCount)
		return
	}
	o := grp.DefaultOptions[T]()
	if err = grp.Apply(&o, opts...); err != nil {
		return
	}

	var clients []grp.Client[T]
	if o.Reconnect {
		clients, err = makeReconnectClients(clientsCount, codec, factory,
			o.ClientOpts...)
	} else {
		clients, err = makeClients(clientsCount, codec, factory, o.ClientOpts...)
	}
	if err != nil {
		for _, client := range clients {
			_ = client.Close()
		}
		err = core.NewError(err)
		return
	}
	strategy := o.Factory.New(clients)
	group = grp.New(strategy)
	return
}

// NewSender creates a high-level Sender that connects to the specified network
// address. It encapsulates a Client Group and provides a simplified interface
// for sending Commands.
func NewSender[T any](addr string, codec cln.Codec[T],
	opts ...sndr.SetMakeOption[T],
) (sender sndr.Sender[T], err error) {
	var connFactory cln.ConnFactoryFn
	o := sndr.DefaultMakeOptions[T]()
	if err = sndr.ApplyMake(&o, opts...); err != nil {
		return
	}
	if o.TLSConfig == nil {
		connFactory = func() (net.Conn, error) {
			return net.Dial("tcp", addr)
		}
	} else {
		connFactory = func() (net.Conn, error) {
			return tls.Dial("tcp", addr, o.TLSConfig)
		}
	}
	g, err := NewGroup(o.ClientsCount, codec, connFactory, o.Group...)
	if err != nil {
		return
	}
	sender = sndr.New(g, o.Sender...)
	return
}

func makeClients[T any](count int, codec cln.Codec[T],
	factory cln.ConnFactory, opts ...cln.SetOption,
) (clients []grp.Client[T], err error) {
	clients = make([]grp.Client[T], 0, count)
	var conn net.Conn
	for range count {
		conn, err = factory.New()
		if err != nil {
			return
		}
		var c *ccln.Client[T]
		c, err = NewClient(codec, conn, opts...)
		if err != nil {
			return
		}
		clients = append(clients, c)
	}
	return
}

func makeReconnectClients[T any](count int, codec cln.Codec[T],
	factory cln.ConnFactory, opts ...cln.SetOption,
) (sl []grp.Client[T], err error) {
	sl = make([]grp.Client[T], 0, count)
	for range count {
		var c *ccln.Client[T]
		c, err = NewReconnectClient(codec, factory, opts...)
		if err != nil {
			return
		}
		sl = append(sl, c)
	}
	return
}
