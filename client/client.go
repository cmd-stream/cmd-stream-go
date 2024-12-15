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

// NewDef creates a Client with default ServerInfo and configuration.
func NewDef[T any](codec Codec[T], conn net.Conn,
	handler base_client.UnexpectedResultHandler,
) (client *base_client.Client[T], err error) {
	return New[T](cs_server.DefServerInfo, DefConf, codec, conn, handler)
}

// New creates a Client.
//
// Client relies on user-defined Codec - Codec.Encode() is used to encode
// commands, Codec.Decode() to decode results received from the Server. If
// decoding fails, the Client will close.
// If the Server impose a limit on the command size, Codec.Size() will be used
// to determine whether the command being sent is small enough.
// If the handler parameter is nil, all unknown results received from the Server
// will be ignored.
//
// Returns client.ErrServerInfoMismatch (from the delegate module) if the
// specified info does not match the info received from the Server.
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
