package srv

import (
	csrv "github.com/cmd-stream/core-go/server"
	"github.com/cmd-stream/delegate-go"
	dsrv "github.com/cmd-stream/delegate-go/server"
	"github.com/cmd-stream/handler-go"
	"github.com/cmd-stream/transport-go"
)

// Options defines the configuration settings for initializing a server.
//
// These options are composed of modular components that configure different
// layers of the server, including transport, handler logic, delegate behavior,
// and base server setup.
type Options struct {
	Info      delegate.ServerInfo
	Base      []csrv.SetOption
	Delegate  []dsrv.SetOption
	Handler   []handler.SetOption
	Transport []transport.SetOption
}

type SetOption func(o *Options)

// WithServerInfo sets the ServerInfo for the server.
//
// ServerInfo helps the client identify a compatible server.
func WithServerInfo(info delegate.ServerInfo) SetOption {
	return func(o *Options) { o.Info = info }
}

// WithCore applies core-level configuration options.
func WithCore(ops ...csrv.SetOption) SetOption {
	return func(o *Options) { o.Base = ops }
}

// WithDelegate applies delegate-specific options.
//
// These options customize the behavior of the server delegate.
func WithDelegate(ops ...dsrv.SetOption) SetOption {
	return func(o *Options) { o.Delegate = ops }
}

// WithHandler sets options for the connection handler.
//
// These options customize the behavior of the connection handler.
func WithHandler(ops ...handler.SetOption) SetOption {
	return func(o *Options) { o.Handler = ops }
}

// WithTransport applies transport-specific options.
//
// These options configure the transport layer.
func WithTransport(ops ...transport.SetOption) SetOption {
	return func(o *Options) { o.Transport = ops }
}

func Apply(ops []SetOption, o *Options) {
	for i := range ops {
		if ops[i] != nil {
			ops[i](o)
		}
	}
}
