package srv

import (
	"context"
	"net"
	"time"

	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
)

// ServerInfoDelegate implements the core.ServerDelegate interface.
//
// It initializes the connection by sending ServerInfo to the client.
type ServerInfoDelegate[T any] struct {
	info    dlgt.ServerInfo
	factory dlgt.ServerTransportFactory[T]
	handler dlgt.ServerTransportHandler[T]
	options Options
}

// New creates a new ServerInfoDelegate.
func New[T any](info dlgt.ServerInfo, factory dlgt.ServerTransportFactory[T],
	handler dlgt.ServerTransportHandler[T],
	opts ...SetOption,
) (delegate ServerInfoDelegate[T], err error) {
	if len(info) == 0 {
		err = ErrEmptyInfo
		return
	}
	o := Options{}
	Apply(&o, opts...)
	return ServerInfoDelegate[T]{
		info:    info,
		factory: factory,
		handler: handler,
		options: o,
	}, nil
}

// Handle performs the initial handshake by sending ServerInfo to the client.
func (d ServerInfoDelegate[T]) Handle(ctx context.Context, conn net.Conn) (err error) {
	transport := d.factory.New(conn)
	err = d.sendServerInfo(transport)
	if err != nil {
		_ = transport.Close()
		return
	}
	return d.handler.Handle(ctx, transport)
}

func (d ServerInfoDelegate[T]) sendServerInfo(transport dlgt.ServerTransport[T]) (err error) {
	if d.options.ServerInfoSendDuration != 0 {
		deadline := time.Now().Add(d.options.ServerInfoSendDuration)
		if err = transport.SetSendDeadline(deadline); err != nil {
			return
		}
	}
	return transport.SendServerInfo(d.info)
}
