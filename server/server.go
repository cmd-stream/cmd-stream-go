package server

import (
	base "github.com/cmd-stream/base-go"
	base_server "github.com/cmd-stream/base-go/server"
	delegate "github.com/cmd-stream/delegate-go"
	delegate_server "github.com/cmd-stream/delegate-go/server"
	handler "github.com/cmd-stream/handler-go"
)

// DefWorkersCount defines a default workers count.
const DefWorkersCount = 8

// DefServerInfo defines a default ServerInfo.
var DefServerInfo = []byte("default")

// DefConf defines a default Server configuration.
var DefConf = Conf{
	Base: base_server.Conf{WorkersCount: DefWorkersCount},
}

// NewDef creates a default Server.
func NewDef[T any](codec Codec[T], receiver T) *base_server.Server {
	return New[T](DefServerInfo, delegate.ServerSettings{}, DefConf, codec,
		receiver,
		DefInvoker[T]{receiver})
}

// New creates a Server.
//
// Server relies on user-defined Codec. It uses Codec.Decode() to decode
// commands received from the Client and Codec.Encode() to encode the results
// sent back. If one of these methods fails, the Server closes the connection to
// the client.
//
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

// NewWith creates a Server with the delegate.
func NewWith(conf base_server.Conf,
	delegate base.ServerDelegate) *base_server.Server {
	return &base_server.Server{
		Conf:     conf,
		Delegate: delegate,
	}
}

func makeInvoker[T any](receiver T, invoker handler.Invoker[T]) handler.Invoker[T] {
	if invoker == nil {
		return DefInvoker[T]{receiver}
	}
	return invoker
}
