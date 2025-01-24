package server

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	bser "github.com/cmd-stream/base-go/server"
	"github.com/cmd-stream/delegate-go"
	dser "github.com/cmd-stream/delegate-go/server"
	"github.com/cmd-stream/handler-go"
)

// DefaultServerInfo is the default ServerInfo.
var DefaultServerInfo = []byte("default")

// Default creates a new Server with the default: configuration (WorkersCount == 8),
// ServerInfo, ServerSettings and Invoker.
//
// This function is ideal for quickly initializing a Server with standard
// settings. For customized configurations, use the New constructor instead.
func Default[T any](codec Codec[T], receiver T) *bser.Server {
	var (
		conf                         = Conf{Base: bser.Conf{WorkersCount: 8}}
		invoker handler.InvokerFn[T] = func(ctx context.Context, at time.Time,
			seq base.Seq, cmd base.Cmd[T], proxy base.Proxy) error {
			return cmd.Exec(ctx, at, seq, receiver, proxy)
		}
	)
	return New[T](conf, DefaultServerInfo, delegate.ServerSettings{}, codec,
		invoker, nil)
}

// New creates a new Server.
//
// Parameters:
//   - conf: Configuration for the server.
//   - info: Server info, sent to the client during connection initialization.
//   - settings: Server settings, also sent to the client during connection
//     initialization.
//   - codec: Decodes Commands and encodes Results sent back to the client.
//   - invoker: Responsible for invoking the Commands.
//   - callback: Closed connections can be tracked using this callback, allowing
//     for monitoring and handling of disconnections.
func New[T any](conf Conf, info delegate.ServerInfo,
	settings delegate.ServerSettings,
	codec Codec[T],
	invoker handler.Invoker[T],
	callback bser.LostConnCallback,
) *bser.Server {
	var (
		f = TransportFactory[T]{Conf: conf.Transport, Codec: codecAdapter[T]{codec}}
		h = handler.New[T](conf.Handler, invoker)
		d = dser.New[T](conf.Delegate, info, settings, f, h)
	)
	return &bser.Server{Conf: conf.Base, Delegate: d, Callback: callback}
}
