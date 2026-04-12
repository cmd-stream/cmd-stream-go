package cln

import (
	"bytes"
	"errors"
	"testing"

	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestTransport(t *testing.T) {
	t.Run("ReceiveServerInfo should decode info from MUS encoding",
		func(t *testing.T) {
			var (
				wantInfo dlgt.ServerInfo = []byte("info")
				wantErr  error           = nil
				bs                       = infoToBs(wantInfo)
				conn                     = cmock.NewConn().RegisterRead(
					func(b []byte) (n int, err error) {
						n = copy(b, bs)
						return
					},
				)
				transport = New[any](conn, nil)
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
				conn                     = cmock.NewConn().RegisterRead(
					func(b []byte) (n int, err error) {
						return 0, wantErr
					},
				)
				transport = New[any](conn, nil)
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
				conn                     = cmock.NewConn().RegisterRead(
					func(b []byte) (n int, err error) {
						n = copy(b, bs)
						return
					},
				)
				transport = New[any](conn, nil)
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
