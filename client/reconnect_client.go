package ccln

import (
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
	o := Options{Info: cser.ServerInfo}
	Apply(ops, &o)
	var (
		d base.ClientDelegate[T]
		f = transportFactory[T]{adaptCodec(codec, o), factory,
			o.Transport}
	)
	d, err = dcln.NewReconnect(o.Info, f, o.Delegate...)
	if err != nil {
		return
	}
	if o.Keepalive != nil {
		d = dcln.NewKeepalive(d, o.Keepalive...)
	}
	client = bcln.New(d, o.Base...)
	return
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
	transport = tcln.New(conn, f.codec, f.ops...)
	return
}
