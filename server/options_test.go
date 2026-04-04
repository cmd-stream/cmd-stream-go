package server

import (
	"testing"
	"time"

	csrv "github.com/cmd-stream/cmd-stream-go/core/srv"
	"github.com/cmd-stream/cmd-stream-go/delegate"
	dsrv "github.com/cmd-stream/cmd-stream-go/delegate/srv"
	hdlr "github.com/cmd-stream/cmd-stream-go/handler"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestOptions(t *testing.T) {
	var (
		o             = Options{}
		wantInfo      = delegate.ServerInfo("info")
		wantCore      = []csrv.SetOption{csrv.WithWorkersCount(10)}
		wantDelegate  = []dsrv.SetOption{dsrv.WithServerInfoSendDuration(time.Second)}
		wantHandler   = []hdlr.SetOption{hdlr.WithAt()}
		wantTransport = []tspt.SetOption{tspt.WithWriterBufSize(1024)}
	)
	Apply(&o,
		WithServerInfo(wantInfo),
		WithCore(wantCore...),
		WithDelegate(wantDelegate...),
		WithHandler(wantHandler...),
		WithTransport(wantTransport...),
	)
	asserterror.EqualDeep(t, o.Info, wantInfo)
	asserterror.Equal(t, len(o.Core), len(wantCore))
	asserterror.Equal(t, len(o.Delegate), len(wantDelegate))
	asserterror.Equal(t, len(o.Handler), len(wantHandler))
	asserterror.Equal(t, len(o.Transport), len(wantTransport))
}
