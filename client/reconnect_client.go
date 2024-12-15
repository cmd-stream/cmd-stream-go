package client

import (
	"net"

	"github.com/cmd-stream/base-go"
	base_client "github.com/cmd-stream/base-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	delegate "github.com/cmd-stream/delegate-go"
	delegate_client "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/transport-go"
	transport_client "github.com/cmd-stream/transport-go/client"
	transport_common "github.com/cmd-stream/transport-go/common"
)

// ConnFactory creates a connection to the Server.
type ConnFactory interface {
	New() (net.Conn, error)
}

type transportFactory[T any] struct {
	conf    transport_common.Conf
	codec   transport.Codec[base.Cmd[T], base.Result]
	factory ConnFactory
}

func (f transportFactory[T]) New() (transport delegate.ClienTransport[T],
	err error) {
	conn, err := f.factory.New()
	if err != nil {
		return
	}
	transport = transport_client.New[T](f.conf, conn, f.codec)
	return
}

// NewDefReconnect creates a "reconnect" Client with default ServerInfo and
// configuration.
func NewDefReconnect[T any](codec Codec[T], factory ConnFactory,
	callback base_client.UnexpectedResultHandler,
) (client *base_client.Client[T], err error) {
	return NewReconnect[T](cs_server.DefServerInfo, DefConf, codec, factory,
		callback)
}

// NewReconnect creates a "reconnect" Client.
//
// If the Codec.Decode method returns a network error, the Client will try to
// reconnect. Otherwise works just like a regular Client.
func NewReconnect[T any](info delegate.ServerInfo, conf Conf,
	codec Codec[T],
	factory ConnFactory,
	callback base_client.UnexpectedResultHandler,
) (client *base_client.Client[T], err error) {
	var (
		d base.ClientDelegate[T]
		c = adaptCodec[T](conf, codec)
		f = transportFactory[T]{conf.Transport, c, factory}
	)
	d, err = delegate_client.NewReconnect[T](conf.Delegate, info, f)
	if err != nil {
		return
	}
	if conf.KeepaliveOn() {
		d = delegate_client.NewKeepalive[T](conf.Delegate, d)
	}
	return NewWith[T](d, callback), nil
}
