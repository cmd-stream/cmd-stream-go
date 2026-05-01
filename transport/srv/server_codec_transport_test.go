package srv_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/cmd-stream/cmd-stream-go/test/mock"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
	tsrv "github.com/cmd-stream/cmd-stream-go/transport/srv"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

func TestServerCodecTransport(t *testing.T) {
	t.Run("SendServerInfo should encode info to MUS encoding",
		func(t *testing.T) {
			var (
				wantInfo dlgt.ServerInfo = []byte("info")
				wantBs                   = infoToBs(wantInfo)
				wantErr  error           = nil
				conn                     = mock.NewConn().RegisterWrite(
					func(bs []byte) (n int, err error) {
						asserterror.EqualDeep(t, bs, wantBs)
						n = len(bs)
						return
					},
				)
				transport = tsrv.New[any](conn, nil)
				err       = transport.SendServerInfo(wantInfo)
			)
			asserterror.EqualError(t, err, wantErr)
		})

	t.Run("If Conn.Write fails with an error, SendServerInfo should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("Conn.Write error")
				conn    = mock.NewConn().RegisterWrite(
					func(b []byte) (n int, err error) {
						err = wantErr
						return
					},
				)
				transport = tsrv.New[any](conn, nil)
				err       = transport.SendServerInfo(nil)
			)
			asserterror.EqualError(t, err, wantErr)
		})

	t.Run("If MarshalServerInfo fails with an error, SendServerInfo should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("WriteByte error")
				writer  = mock.NewWriter().RegisterWriteByte(
					func(b byte) error { return wantErr },
				)
				mocks     = []*mok.Mock{writer.Mock}
				transport = &tsrv.ServerCodecTransport[any]{
					CodecTransport: tspt.New[core.Result, core.Cmd[any]](nil, writer, nil, nil),
				}
				err = transport.SendServerInfo((dlgt.ServerInfo([]byte{})))
			)
			asserterror.EqualError(t, err, wantErr)
			asserterror.EqualDeep(t, mok.CheckCalls(mocks), mok.EmptyInfomap)
		})
}

func infoToBs(info dlgt.ServerInfo) []byte {
	var (
		size = dlgt.ServerInfoValidMUS.Size(info)
		bs   = make([]byte, 0, size)
		buf  = bytes.NewBuffer(bs)
		_, _ = dlgt.ServerInfoValidMUS.Marshal(info, buf)
	)
	return buf.Bytes()
}
