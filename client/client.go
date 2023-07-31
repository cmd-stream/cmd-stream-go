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

// NewDef creates a default Client, which uses the default ServerInfo and
// configuration.
func NewDef[T any](codec Codec[T], conn net.Conn,
	handler base_client.UnexpectedResultHandler,
) (client *base_client.Client[T], err error) {
	return New[T](cs_server.DefServerInfo, DefConf, codec, conn, handler)
}

// New creates a Client.
//
// The Client relies on the user-defined Codec. It uses the Codec.Encode method
// to send commands. And if the encoding of any command fails with an error, it
// will be returned by the Client.Send method.
// With the Codec.Decode method, the story is a little different. It is used by
// the Client in the background goroutine that receives the results. And if the
// decoding of any result fails, the Client will close.
// Also, if the Server imposes a limit on the size of the command, the Client
// will use the Codec.Size method to determine if the size of the command being
// sent is small enough.
//
// If the handler is nil, all unknown results received from the Server will be
// ignored.
//
// Returns the delegate.ErrServerInfoMismatch if the specified info does not
// match the info received from the Server.
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
