package ccln

import (
	"net"

	"github.com/cmd-stream/base-go"
	bcln "github.com/cmd-stream/base-go/client"
	cser "github.com/cmd-stream/cmd-stream-go/server"
	dcln "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/transport-go"
	tcln "github.com/cmd-stream/transport-go/client"
)

// New creates a new Client.
//
// Parameters:
//   - codec: Handles encoding of Commands and decoding of Results.
//   - conn: The network connection used for communication between the client
//     and server.
//
// Additional options (ops) can be used to configure various aspects of the
// client.
//
// Returns dcln.ErrServerInfoMismatch if the provided ServerInfo (configured via
// ops) does not match the server's expected info.
func New[T any](codec Codec[T], conn net.Conn, ops ...SetOption) (
	client *bcln.Client[T], err error) {
	o := Options{Info: cser.ServerInfo}
	Apply(ops, &o)
	var (
		d base.ClientDelegate[T]
		t = tcln.New(conn, adaptCodec(codec, o), o.Transport...)
	)
	d, err = dcln.New(o.Info, t, o.Delegate...)
	if err != nil {
		return
	}
	if o.Keepalive != nil {
		d = dcln.NewKeepalive(d, o.Keepalive...)
	}
	client = bcln.New(d, o.Base...)
	return
}

func adaptCodec[T any](codec Codec[T],
	o Options) transport.Codec[base.Cmd[T], base.Result] {
	if o.Keepalive != nil {
		return keepaliveCodecAdapter[T]{codec}
	}
	return codecAdapter[T]{codec}
}
