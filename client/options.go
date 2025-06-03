package client

import (
	ccln "github.com/cmd-stream/core-go/client"
	"github.com/cmd-stream/delegate-go"
	dcln "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/transport-go"
)

// Options defines the configuration settings for initializing a client.
//
// These options are composed of modular components that configure different
// layers of the client, including transport, keepalive, delegate behavior, and
// base client setup.
type Options struct {
	Info      delegate.ServerInfo
	Base      []ccln.SetOption
	Transport []transport.SetOption
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

// WithCore applies base-level configuration options.
func WithCore(ops ...ccln.SetOption) SetOption {
	return func(o *Options) { o.Base = ops }
}

// WithTransport applies transport-specific options.
//
// These options configure the transport layer.
func WithTransport(ops ...transport.SetOption) SetOption {
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
