package client

import (
	"net"

	"github.com/cmd-stream/base-go"
	bcln "github.com/cmd-stream/base-go/client"
	cser "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/delegate-go"
	dcln "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/transport-go"
	tcln "github.com/cmd-stream/transport-go/client"
)

// Default creates a new Client with the default: configuration and ServerInfo.
//
// This function is ideal for quickly initializing a Client with standard
// settings. For customized configurations, use the New constructor instead.
func Default[T any](codec Codec[T], conn net.Conn) (client *bcln.Client[T],
	err error) {
	return New[T](Conf{}, cser.DefaultServerInfo, codec, conn, nil)
}

// New creates a new Client.
//
// Parameters:
//   - conf: Client configuration.
//   - info: Must match the corresponding data provided by the server.
//   - codec: Responsible for encoding Commands and decoding Results. If the
//     server enforces a limit on Command size, Codec.Size() will be used to
//     verify whether the Command is within the allowed size.
//   - callback: Used to handle unexpected Results received from the server.
//
// Returns dcln.ErrServerInfoMismatch if the provided ServerInfo does not match
// the server's info.
func New[T any](conf Conf, info delegate.ServerInfo, codec Codec[T],
	conn net.Conn,
	callback bcln.UnexpectedResultCallback,
) (client *bcln.Client[T], err error) {
	var (
		d base.ClientDelegate[T]
		c = adaptCodec[T](conf, codec)
		t = tcln.New[T](conf.Transport, conn, c)
	)
	d, err = dcln.New[T](conf.Delegate, info, t)
	if err != nil {
		return
	}
	if conf.KeepaliveOn() {
		d = dcln.NewKeepalive[T](conf.Delegate, d)
	}
	client = bcln.New[T](d, callback)
	return
}

func adaptCodec[T any](conf Conf, codec Codec[T]) transport.Codec[base.Cmd[T], base.Result] {
	if conf.KeepaliveOn() {
		return keepaliveCodecAdapter[T]{codec}
	}
	return codecAdapter[T]{codec}
}
