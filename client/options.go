package ccln

import (
	bcln "github.com/cmd-stream/base-go/client"
	"github.com/cmd-stream/delegate-go"
	dcln "github.com/cmd-stream/delegate-go/client"
	tcom "github.com/cmd-stream/transport-go/common"
)

type Options struct {
	Info      delegate.ServerInfo
	Base      []bcln.SetOption
	Transport []tcom.SetOption
	Delegate  []dcln.SetOption
	Keepalive []dcln.SetKeepaliveOption
}

type SetOption func(o *Options)

// WithServerInfo sets the ServerInfo for the client.
//
// ServerInfo helps the client identify a compatible server.
func WithServerInfo(info delegate.ServerInfo) SetOption {
	return func(o *Options) { o.Info = info }
}

// WithBase applies base-level configuration options.
func WithBase(ops ...bcln.SetOption) SetOption {
	return func(o *Options) { o.Base = ops }
}

// WithTransport applies transport-specific options.
//
// These options configure the transport layer.
func WithTransport(ops ...tcom.SetOption) SetOption {
	return func(o *Options) { o.Transport = ops }
}

// WithDelegate applies delegate-specific options.
//
// These options customize the behavior of the client delegate.
func WithDelegate(ops ...dcln.SetOption) SetOption {
	return func(o *Options) { o.Delegate = ops }
}

// WithKeepalive applies keepalive-specific options.
//
// These options configure how the client maintains an active connection,
// including Ping intervals and timeout settings.
func WithKeepalive(ops ...dcln.SetKeepaliveOption) SetOption {
	return func(o *Options) { o.Keepalive = ops }
}

func Apply(ops []SetOption, o *Options) {
	for i := range ops {
		if ops[i] != nil {
			ops[i](o)
		}
	}
}
