package cln_test

import (
	"bytes"
	"errors"
	"testing"

	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/cmd-stream/cmd-stream-go/test/mock"
	tcln "github.com/cmd-stream/cmd-stream-go/transport/cln"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestTransport(t *testing.T) {
	t.Run("ReceiveServerInfo should decode info from MUS encoding",
		func(t *testing.T) {
			var (
				wantInfo dlgt.ServerInfo = []byte("info")
				wantErr  error           = nil
				bs                       = infoToBs(wantInfo)
				conn                     = mock.NewConn().RegisterRead(
					func(b []byte) (n int, err error) {
						n = copy(b, bs)
						return
					},
				)
				transport = tcln.New[any](conn, nil)
				info, err = transport.ReceiveServerInfo()
			)
			asserterror.EqualDeep(t, info, wantInfo)
			asserterror.EqualError(t, err, wantErr)
		})

	t.Run("If decoding fails with an error, ReceiveServerInfo should return this error",
		func(t *testing.T) {
			var (
				wantInfo dlgt.ServerInfo = nil
				wantErr                  = errors.New("Read error")
				conn                     = mock.NewConn().RegisterRead(
					func(b []byte) (n int, err error) {
						return 0, wantErr
					},
				)
				transport = tcln.New[any](conn, nil)
				info, err = transport.ReceiveServerInfo()
			)
			asserterror.EqualDeep(t, info, wantInfo)
			asserterror.EqualError(t, err, wantErr)
		})

	t.Run("If ServerInfo length exceeds the limit, ReceiveServerInfo should return an error",
		func(t *testing.T) {
			var (
				wantInfo dlgt.ServerInfo = nil
				wantErr                  = dlgt.ErrTooLargeServerInfo
				infoBs                   = make([]byte, dlgt.DefaultServerInfoMaxLength+1)
				bs                       = infoToBs(infoBs)
				conn                     = mock.NewConn().RegisterRead(
					func(b []byte) (n int, err error) {
						n = copy(b, bs)
						return
					},
				)
				transport = tcln.New[any](conn, nil)
				info, err = transport.ReceiveServerInfo()
			)
			asserterror.EqualDeep(t, info, wantInfo)
			asserterror.EqualError(t, err, wantErr)
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
