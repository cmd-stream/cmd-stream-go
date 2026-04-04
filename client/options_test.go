package client

import (
	"testing"

	asserterror "github.com/ymz-ncnk/assert/error"

	ccln "github.com/cmd-stream/cmd-stream-go/core/cln"
	"github.com/cmd-stream/cmd-stream-go/delegate"
	dcln "github.com/cmd-stream/cmd-stream-go/delegate/cln"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	asserterror.EqualDeep(t, opts.Info, delegate.ServerInfo(srv.ServerInfo))
}

func TestWithServerInfo(t *testing.T) {
	var (
		info = delegate.ServerInfo("test-server")
		opts = Options{}
	)
	WithServerInfo(info)(&opts)
	asserterror.EqualDeep(t, opts.Info, info)
}

func TestWithCore(t *testing.T) {
	var (
		opts     = Options{}
		coreOpts = []ccln.SetOption{func(o *ccln.Options) {}}
	)
	WithCore(coreOpts...)(&opts)
	asserterror.Equal(t, len(opts.Base), len(coreOpts))
}

func TestWithTransport(t *testing.T) {
	var (
		opts     = Options{}
		tsptOpts = []tspt.SetOption{tspt.WithWriterBufSize(1024)}
	)
	WithTransport(tsptOpts...)(&opts)
	asserterror.Equal(t, len(opts.Transport), len(tsptOpts))
}

func TestWithDelegate(t *testing.T) {
	var (
		opts     = Options{}
		dclnOpts = []dcln.SetOption{dcln.WithServerInfoReceiveDuration(100)}
	)
	WithDelegate(dclnOpts...)(&opts)
	asserterror.Equal(t, len(opts.Delegate), len(dclnOpts))
}

func TestWithKeepalive(t *testing.T) {
	var (
		opts          = Options{}
		keepaliveOpts = []dcln.SetKeepaliveOption{dcln.WithKeepaliveTime(1000)}
	)
	WithKeepalive(keepaliveOpts...)(&opts)
	asserterror.Equal(t, len(opts.Keepalive), len(keepaliveOpts))
}

func TestApply(t *testing.T) {
	var (
		opts = Options{}
		info = delegate.ServerInfo("test-server")
	)
	Apply(&opts, WithServerInfo(info), nil)
	asserterror.EqualDeep(t, opts.Info, info)
}
