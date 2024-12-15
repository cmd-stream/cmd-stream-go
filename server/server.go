package server

import (
	base "github.com/cmd-stream/base-go"
	base_server "github.com/cmd-stream/base-go/server"
	delegate "github.com/cmd-stream/delegate-go"
	delegate_server "github.com/cmd-stream/delegate-go/server"
	handler "github.com/cmd-stream/handler-go"
)

var (
	// ServerInfo defines a default ServerInfo.
	DefServerInfo = []byte("default")
	// Conf defines a default Server configuration.
	DefConf = Conf{
		Base: base_server.Conf{WorkersCount: 8},
	}
)

// NewDef creates a Server with default ServerInfo and configuration.
//
// The server will be able to handle only 8 connections at a time.
func NewDef[T any](codec Codec[T], receiver T) *base_server.Server {
	return New[T](DefServerInfo, delegate.ServerSettings{}, DefConf, codec,
		receiver,
		Invoker[T]{receiver})
}

// New creates a Server.
//
// Server relies on user-defined Codec. It uses Codec.Decode() to decode
// commands received from the Client and Codec.Encode() to encode the results.
// If one of these methods fails, the connection to the client will be closed.
// All closed connections could be tracked with Conf.Base.LostConnCallback.
// If the invoker parameter is nil, the default value is used.
func New[T any](info delegate.ServerInfo, settings delegate.ServerSettings,
	conf Conf,
	codec Codec[T],
	receiver T,
	invoker handler.Invoker[T],
) *base_server.Server {
	var (
		factory = TransportFactory[T]{
			Conf:  conf.Transport,
			Codec: codecAdapter[T]{codec},
		}
		handler  = handler.New[T](conf.Handler, makeInvoker[T](receiver, invoker))
		delegate = delegate_server.New[T](conf.Delegate, info, settings,
			factory,
			handler)
	)
	return NewWith(conf.Base, delegate)
}

// NewWith creates a Server with the specified delegate.
func NewWith(conf base_server.Conf,
	delegate base.ServerDelegate) *base_server.Server {
	return &base_server.Server{
		Conf:     conf,
		Delegate: delegate,
	}
}

func makeInvoker[T any](receiver T, invoker handler.Invoker[T]) handler.Invoker[T] {
	if invoker == nil {
		return Invoker[T]{receiver}
	}
	return invoker
}
