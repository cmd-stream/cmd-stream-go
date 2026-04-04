package server

import (
	csrv "github.com/cmd-stream/cmd-stream-go/core/srv"
	"github.com/cmd-stream/cmd-stream-go/delegate"
	dsrv "github.com/cmd-stream/cmd-stream-go/delegate/srv"
	hdlr "github.com/cmd-stream/cmd-stream-go/handler"
	"github.com/cmd-stream/cmd-stream-go/transport"
)

// Options defines the configuration settings for initializing a server.
//
// These options are composed of modular components that configure different
// layers of the server, including transport, handler logic, delegate behavior,
// and base server setup.
type Options struct {
	Info      delegate.ServerInfo
	Core      []csrv.SetOption
	Delegate  []dsrv.SetOption
	Handler   []hdlr.SetOption
	Transport []transport.SetOption
}

func DefaultOptions() Options {
	return Options{
		Info: ServerInfo,
	}
}

type SetOption func(o *Options)

// WithServerInfo sets the ServerInfo for the server.
//
// ServerInfo helps the client identify a compatible server.
func WithServerInfo(info delegate.ServerInfo) SetOption {
	return func(o *Options) { o.Info = info }
}

// WithCore applies core-level configuration options.
func WithCore(opts ...csrv.SetOption) SetOption {
	return func(o *Options) { o.Core = append(o.Core, opts...) }
}

// WithDelegate applies delegate-specific options.
//
// These options customize the behavior of the server delegate.
func WithDelegate(opts ...dsrv.SetOption) SetOption {
	return func(o *Options) { o.Delegate = append(o.Delegate, opts...) }
}

// WithHandler sets options for the connection handler.
//
// These options customize the behavior of the connection handler.
func WithHandler(opts ...hdlr.SetOption) SetOption {
	return func(o *Options) { o.Handler = append(o.Handler, opts...) }
}

// WithTransport applies transport-specific options.
//
// These options configure the transport layer.
func WithTransport(opts ...transport.SetOption) SetOption {
	return func(o *Options) { o.Transport = append(o.Transport, opts...) }
}

func Apply(o *Options, opts ...SetOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
