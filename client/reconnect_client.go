package ccln

import (
	"net"

	"github.com/cmd-stream/base-go"
	bcln "github.com/cmd-stream/base-go/client"
	cser "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/delegate-go"
	dcln "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/transport-go"
	tcln "github.com/cmd-stream/transport-go/client"
	tcom "github.com/cmd-stream/transport-go/common"
)

// NewReconnect creates a "reconnect" Client.
//
// If the Codec.Decode method encounters a network error, the client will
// attempt to reconnect automatically. In all other aspects, it functions like
// a regular client.
func NewReconnect[T any](codec Codec[T], factory ConnFactory,
	ops ...SetOption) (client *bcln.Client[T], err error) {
	options := Options{Info: cser.ServerInfo}
	Apply(ops, &options)
	var (
		d base.ClientDelegate[T]
		f = transportFactory[T]{adaptCodec[T](codec, options), factory,
			options.Transport}
	)
	d, err = dcln.NewReconnect[T](options.Info, f, options.Delegate...)
	if err != nil {
		return
	}
	if options.Keepalive != nil {
		d = dcln.NewKeepalive[T](d, options.Keepalive...)
	}
	client = bcln.New[T](d, options.Base...)
	return
}

// ConnFactory establishes a new connection to the server.
type ConnFactory interface {
	New() (net.Conn, error)
}

type transportFactory[T any] struct {
	codec   transport.Codec[base.Cmd[T], base.Result]
	factory ConnFactory
	ops     []tcom.SetOption
}

func (f transportFactory[T]) New() (transport delegate.ClienTransport[T],
	err error) {
	conn, err := f.factory.New()
	if err != nil {
		return
	}
	transport = tcln.New[T](conn, f.codec, f.ops...)
	return
}
