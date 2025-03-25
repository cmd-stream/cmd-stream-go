package cser

import (
	bser "github.com/cmd-stream/base-go/server"
	"github.com/cmd-stream/delegate-go"
	dser "github.com/cmd-stream/delegate-go/server"
	"github.com/cmd-stream/handler-go"
	tcom "github.com/cmd-stream/transport-go/common"
)

type Options struct {
	Info      delegate.ServerInfo
	Base      []bser.SetOption
	Delegate  []dser.SetOption
	Handler   []handler.SetOption
	Transport []tcom.SetOption
}

type SetOption func(o *Options)

// WithServerInfo sets the ServerInfo for the server.
//
// ServerInfo helps the client identify a compatible server.
func WithServerInfo(info delegate.ServerInfo) SetOption {
	return func(o *Options) { o.Info = info }
}

// WithBase applies base-level configuration options.
func WithBase(ops ...bser.SetOption) SetOption {
	return func(o *Options) { o.Base = ops }
}

// WithDelegate applies delegate-specific options.
//
// These options customize the behavior of the server delegate.
func WithDelegate(ops ...dser.SetOption) SetOption {
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
func WithTransport(ops ...tcom.SetOption) SetOption {
	return func(o *Options) { o.Transport = ops }
}

func Apply(ops []SetOption, o *Options) {
	for i := range ops {
		if ops[i] != nil {
			ops[i](o)
		}
	}
}
