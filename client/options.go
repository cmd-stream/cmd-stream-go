package client

import (
	ccln "github.com/cmd-stream/cmd-stream-go/core/cln"
	"github.com/cmd-stream/cmd-stream-go/delegate"
	dcln "github.com/cmd-stream/cmd-stream-go/delegate/cln"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
)

// Options defines the configuration settings for initializing a client.
//
// These options are composed of modular components that configure different
// layers of the client, including transport, keepalive, delegate behavior, and
// base client setup.
type Options struct {
	Info      delegate.ServerInfo
	Base      []ccln.SetOption
	Transport []tspt.SetOption
	Delegate  []dcln.SetOption
	Keepalive []dcln.SetKeepaliveOption
}

// DefaultOptions returns the default Client configuration.
func DefaultOptions() Options {
	return Options{
		Info: srv.ServerInfo,
	}
}

// SetOption defines a function for configuring Options.
type SetOption func(o *Options)

// WithServerInfo sets the ServerInfo for the client.
//
// ServerInfo helps the client identify a compatible server. Its length is
// limited to 1KB, otherwise the client will break the connection.
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
func WithTransport(ops ...tspt.SetOption) SetOption {
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

// Apply applies the given options to the Options struct.
func Apply(o *Options, opts ...SetOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
