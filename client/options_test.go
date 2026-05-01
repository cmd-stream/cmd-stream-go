package client_test

import (
	"testing"

	asserterror "github.com/ymz-ncnk/assert/error"

	cln "github.com/cmd-stream/cmd-stream-go/client"
	ccln "github.com/cmd-stream/cmd-stream-go/core/cln"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	dcln "github.com/cmd-stream/cmd-stream-go/delegate/cln"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
)

func TestDefaultOptions(t *testing.T) {
	opts := cln.DefaultOptions()
	asserterror.EqualDeep(t, opts.Info, dlgt.DefaultServerInfo)
}

func TestWithServerInfo(t *testing.T) {
	var (
		info = dlgt.ServerInfo("test-server")
		opts = cln.Options{}
	)
	cln.WithServerInfo(info)(&opts)
	asserterror.EqualDeep(t, opts.Info, info)
}

func TestWithCore(t *testing.T) {
	var (
		opts     = cln.Options{}
		coreOpts = []ccln.SetOption{func(o *ccln.Options) {}}
	)
	cln.WithCore(coreOpts...)(&opts)
	asserterror.Equal(t, len(opts.Base), len(coreOpts))
}

func TestWithTransport(t *testing.T) {
	var (
		opts     = cln.Options{}
		tsptOpts = []tspt.SetOption{tspt.WithWriterBufSize(1024)}
	)
	cln.WithTransport(tsptOpts...)(&opts)
	asserterror.Equal(t, len(opts.Transport), len(tsptOpts))
}

func TestWithDelegate(t *testing.T) {
	var (
		opts     = cln.Options{}
		dclnOpts = []dcln.SetOption{dcln.WithServerInfoReceiveDuration(100)}
	)
	cln.WithDelegate(dclnOpts...)(&opts)
	asserterror.Equal(t, len(opts.Delegate), len(dclnOpts))
}

func TestWithKeepalive(t *testing.T) {
	var (
		opts          = cln.Options{}
		keepaliveOpts = []dcln.SetKeepaliveOption{dcln.WithKeepaliveTime(1000)}
	)
	cln.WithKeepalive(keepaliveOpts...)(&opts)
	asserterror.Equal(t, len(opts.Keepalive), len(keepaliveOpts))
}

func TestApply(t *testing.T) {
	var (
		opts = cln.Options{}
		info = dlgt.ServerInfo("test-server")
	)
	cln.Apply(&opts, cln.WithServerInfo(info), nil)
	asserterror.EqualDeep(t, opts.Info, info)
}
