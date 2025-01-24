package client

import (
	"net"

	"github.com/cmd-stream/base-go"
	bcln "github.com/cmd-stream/base-go/client"
	"github.com/cmd-stream/delegate-go"
	dcln "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/transport-go"
	tcln "github.com/cmd-stream/transport-go/client"
	tcom "github.com/cmd-stream/transport-go/common"
)

// ConnFactory establishes a new connection to the Server.
type ConnFactory interface {
	New() (net.Conn, error)
}

type transportFactory[T any] struct {
	conf    tcom.Conf
	codec   transport.Codec[base.Cmd[T], base.Result]
	factory ConnFactory
}

func (f transportFactory[T]) New() (transport delegate.ClienTransport[T],
	err error) {
	conn, err := f.factory.New()
	if err != nil {
		return
	}
	transport = tcln.New[T](f.conf, conn, f.codec)
	return
}

// NewReconnect creates a "reconnect" Client.
//
// If the Codec.Decode method returns a network error, the client will try to
// reconnect. Otherwise works just like a regular client.
func NewReconnect[T any](conf Conf, info delegate.ServerInfo, codec Codec[T],
	factory ConnFactory,
	callback bcln.UnexpectedResultCallback,
) (client *bcln.Client[T], err error) {
	var (
		d base.ClientDelegate[T]
		c = adaptCodec[T](conf, codec)
		f = transportFactory[T]{conf.Transport, c, factory}
	)
	d, err = dcln.NewReconnect[T](conf.Delegate, info, f)
	if err != nil {
		return
	}
	if conf.KeepaliveOn() {
		d = dcln.NewKeepalive[T](conf.Delegate, d)
	}
	client = bcln.New[T](d, callback)
	return
}
