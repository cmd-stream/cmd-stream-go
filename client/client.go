package client

import (
	"net"

	"github.com/cmd-stream/base-go"
	base_client "github.com/cmd-stream/base-go/client"
	cs_server "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/delegate-go"
	delegate_client "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/transport-go"
	transport_client "github.com/cmd-stream/transport-go/client"
)

// DefConf is a default Client configuration.
var DefConf = Conf{}

// NewDef creates a default Client, which uses default ServerInfo and the
// default configuration.
func NewDef[T any](codec Codec[T], conn net.Conn,
	handler base_client.UnexpectedResultHandler,
) (client *base_client.Client[T], err error) {
	return New[T](cs_server.DefServerInfo, DefConf, codec, conn, handler)
}

// New creates a Client.
//
// Client relies on user-defined Codec. It uses Codec.Encode to encode
// commands for Server and Codec.Decode to decode received results. If the
// last one method fails Client will be closed.
// Also, if Server imposes a limit on the size of a command, Client will use
// Codec.Size to determine if the command being sent is small enough.
// If the handler parameter is nil, all unknown results received from Server
// will be ignored.
//
// Returns delegate.ErrServerInfoMismatch if the specified info does not
// match the info received from Server.
func New[T any](info delegate.ServerInfo, conf Conf, codec Codec[T],
	conn net.Conn,
	handler base_client.UnexpectedResultHandler,
) (client *base_client.Client[T], err error) {
	var (
		d base.ClientDelegate[T]
		c = adaptCodec[T](conf, codec)
		t = transport_client.New[T](conf.Transport, conn, c)
	)
	d, err = delegate_client.New[T](conf.Delegate, info, t)
	if err != nil {
		return
	}
	if conf.KeepaliveOn() {
		d = delegate_client.NewKeepalive[T](conf.Delegate, d)
	}
	return NewWith[T](d, handler), nil
}

// NewWith creates a Client with the specified delegate.
func NewWith[T any](delegate base.ClientDelegate[T],
	handler base_client.UnexpectedResultHandler) *base_client.Client[T] {
	return base_client.New[T](delegate, handler)
}

func adaptCodec[T any](conf Conf, codec Codec[T]) transport.Codec[base.Cmd[T], base.Result] {
	if conf.KeepaliveOn() {
		return keepaliveCodecAdapter[T]{codec}
	}
	return codecAdapter[T]{codec}
}
