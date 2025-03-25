package cser

import (
	bser "github.com/cmd-stream/base-go/server"
	dser "github.com/cmd-stream/delegate-go/server"
	"github.com/cmd-stream/handler-go"
)

// ServerInfo is the default ServerInfo.
var ServerInfo = []byte("default")

// New creates a new server.
//
// Parameters:
//   - codec: Handles encoding of Results and decoding of Commands from clients.
//   - invoker: Executes incoming Commands.
//
// Additional options (ops) can be used to configure various aspects of the
// server.
func New[T any](codec Codec[T], invoker handler.Invoker[T],
	ops ...SetOption) *bser.Server {
	options := Options{Info: ServerInfo}
	Apply(ops, &options)
	var (
		f = TransportFactory[T]{
			Codec: codecAdapter[T]{codec},
			Ops:   options.Transport,
		}
		h = handler.New[T](invoker, options.Handler...)
		d = dser.New[T](options.Info, f, h, options.Delegate...)
	)
	return bser.New(d, options.Base...)
}
