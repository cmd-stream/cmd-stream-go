package transport

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	tmock "github.com/cmd-stream/cmd-stream-go/test/mock/transport"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

func LocalAddrTestCase() TransportTestCase[any, any] {
	name := "LocalAddr should return local address of the conn"

	var (
		wantAddr = &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9000}
		conn     = cmock.NewConn()
	)
	conn.RegisterLocalAddr(
		func() (addr net.Addr) { return wantAddr },
	)
	return TransportTestCase[any, any]{
		Name:  name,
		Setup: TransportSetup[any, any]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[any, any]) {
			addr := transport.LocalAddr()
			asserterror.EqualDeep(t, addr, wantAddr)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}

func RemoteAddrTestCase() TransportTestCase[any, any] {
	name := "RemoteAddr should return remote address of the conn"

	var (
		wantAddr = &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9000}
		conn     = cmock.NewConn()
	)
	conn.RegisterRemoteAddr(
		func() (addr net.Addr) { return wantAddr },
	)
	return TransportTestCase[any, any]{
		Name:  name,
		Setup: TransportSetup[any, any]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[any, any]) {
			addr := transport.RemoteAddr()
			asserterror.EqualDeep(t, addr, wantAddr)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}

func SetSendDeadlineTestCase(t *testing.T) TransportTestCase[any, any] {
	name := "Conn.SetWriteDeadline should receive same deadline as SetSendDeadline"

	var (
		wantDeadline = time.Now()
		conn         = cmock.NewConn()
	)
	conn.RegisterSetWriteDeadline(
		func(deadline time.Time) (err error) {
			asserterror.Equal(t, deadline, wantDeadline)
			return
		},
	)
	return TransportTestCase[any, any]{
		Name:  name,
		Setup: TransportSetup[any, any]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[any, any]) {
			_ = transport.SetSendDeadline(wantDeadline)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}

func SetSendDeadlineErrorTestCase() TransportTestCase[any, any] {
	name := "If Conn.SetWriteDeadline fails with an error, SetSendDeadline should return it"

	var (
		wantErr = errors.New("Conn.SetWriteDeadline error")
		conn    = cmock.NewConn()
	)
	conn.RegisterSetWriteDeadline(
		func(deadline time.Time) (err error) { return wantErr },
	)
	return TransportTestCase[any, any]{
		Name:  name,
		Setup: TransportSetup[any, any]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[any, any]) {
			err := transport.SetSendDeadline(time.Time{})
			asserterror.Equal(t, err, wantErr)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}

func SendTestCase(t *testing.T) TransportTestCase[core.Cmd[any], core.Result] {
	name := "Send should encode data with help of the Codec"

	var (
		wantSeq core.Seq = 1
		wantCmd          = cmock.NewCmd[any]()
		wantN            = 3
		wantErr error    = nil
		writer           = tmock.NewWriter()
		codec            = tmock.NewClientCodec()
	)
	codec.RegisterEncode(
		func(seq core.Seq, cmd core.Cmd[any], w tspt.Writer) (n int, err error) {
			asserterror.Equal(t, seq, wantSeq)
			asserterror.Equal[any](t, cmd, wantCmd)
			return wantN, wantErr
		},
	)
	return TransportTestCase[core.Cmd[any], core.Result]{
		Name: name,
		Setup: TransportSetup[core.Cmd[any], core.Result]{
			Writer: writer,
			Codec:  codec,
		},
		Action: func(t *testing.T, transport *tspt.CodecTransport[core.Cmd[any], core.Result]) {
			n, err := transport.Send(wantSeq, wantCmd)
			asserterror.EqualError(t, err, wantErr)
			asserterror.Equal(t, n, wantN)
		},
		Mocks: []*mok.Mock{writer.Mock, codec.Mock},
	}
}

func SendErrorTestCase() TransportTestCase[core.Cmd[any], core.Result] {
	name := "If Codec.Encode fails with an error, Send should return it"

	var (
		wantErr = errors.New("Codec.Encode error")
		codec   = tmock.NewClientCodec()
	)
	codec.RegisterEncode(
		func(seq core.Seq, cmd core.Cmd[any], w tspt.Writer) (n int, err error) {
			return 0, wantErr
		},
	)
	return TransportTestCase[core.Cmd[any], core.Result]{
		Name: name,
		Setup: TransportSetup[core.Cmd[any], core.Result]{
			Codec: codec,
		},
		Action: func(t *testing.T, transport *tspt.CodecTransport[core.Cmd[any], core.Result]) {
			_, err := transport.Send(1, nil)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{codec.Mock},
	}
}

func SetReceiveDeadlineTestCase(t *testing.T) TransportTestCase[any, any] {
	name := "Conn.SetReadDeadline should receive same deadline as SetReceiveDeadline"

	var (
		wantDeadline = time.Now()
		conn         = cmock.NewConn()
	)
	conn.RegisterSetReadDeadline(
		func(deadline time.Time) (err error) {
			asserterror.Equal(t, deadline, wantDeadline)
			return
		},
	)
	return TransportTestCase[any, any]{
		Name:  name,
		Setup: TransportSetup[any, any]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[any, any]) {
			_ = transport.SetReceiveDeadline(wantDeadline)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}

func SetReceiveDeadlineErrorTestCase() TransportTestCase[any, any] {
	name := "If Conn.SetReadDeadline fails with an error, SetReceiveDeadline should return it"

	var (
		wantErr = errors.New("Conn.SetReadDeadline error")
		conn    = cmock.NewConn()
	)
	conn.RegisterSetReadDeadline(
		func(deadline time.Time) (err error) { return wantErr },
	)
	return TransportTestCase[any, any]{
		Name:  name,
		Setup: TransportSetup[any, any]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[any, any]) {
			err := transport.SetReceiveDeadline(time.Time{})
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}

func ReceiveTestCase() TransportTestCase[core.Cmd[any], core.Result] {
	name := "Receive should decode data with help of the Codec"

	var (
		wantSeq    core.Seq = 1
		wantResult          = cmock.NewResult()
		wantN               = 3
		wantErr    error    = nil
		codec               = tmock.NewClientCodec()
	)
	codec.RegisterDecode(
		func(r tspt.Reader) (seq core.Seq, result core.Result, n int, err error) {
			return wantSeq, wantResult, wantN, wantErr
		},
	)
	return TransportTestCase[core.Cmd[any], core.Result]{
		Name: name,
		Setup: TransportSetup[core.Cmd[any], core.Result]{
			Codec: codec,
		},
		Action: func(t *testing.T, transport *tspt.CodecTransport[core.Cmd[any], core.Result]) {
			seq, result, n, err := transport.Receive()
			asserterror.EqualError(t, err, wantErr)
			asserterror.Equal(t, seq, wantSeq)
			asserterror.EqualDeep(t, result, wantResult)
			asserterror.Equal(t, n, wantN)
		},
		Mocks: []*mok.Mock{codec.Mock},
	}
}

func ReceiveErrorTestCase() TransportTestCase[core.Cmd[any], core.Result] {
	name := "If Codec.Decode fails with an error, Receive should return it"

	var (
		wantSeq    core.Seq    = 0
		wantResult core.Result = nil
		wantErr                = errors.New("Codec.Decode error")
		codec                  = tmock.NewClientCodec()
	)
	codec.RegisterDecode(
		func(r tspt.Reader) (seq core.Seq, result core.Result, n int, err error) {
			err = wantErr
			return
		},
	)
	return TransportTestCase[core.Cmd[any], core.Result]{
		Name: name,
		Setup: TransportSetup[core.Cmd[any], core.Result]{
			Codec: codec,
		},
		Action: func(t *testing.T, transport *tspt.CodecTransport[core.Cmd[any], core.Result]) {
			seq, result, _, err := transport.Receive()
			asserterror.EqualError(t, err, wantErr)
			asserterror.Equal(t, seq, wantSeq)
			asserterror.EqualDeep(t, result, wantResult)
		},
		Mocks: []*mok.Mock{codec.Mock},
	}
}

func FlushTestCase() TransportTestCase[any, any] {
	name := "Flush should flush the writer"

	var writer = tmock.NewWriter()
	writer.RegisterFlush(
		func() error { return nil },
	)
	return TransportTestCase[any, any]{
		Name:  name,
		Setup: TransportSetup[any, any]{Writer: writer},
		Action: func(t *testing.T, transport *tspt.CodecTransport[any, any]) {
			err := transport.Flush()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{writer.Mock},
	}
}

func FlushErrorTestCase() TransportTestCase[any, any] {
	name := "If Writer.Flush fails with an error, Flush should return it"

	var (
		wantErr = errors.New("Writer.Flush error")
		writer  = tmock.NewWriter()
	)
	writer.RegisterFlush(
		func() error { return wantErr },
	)
	return TransportTestCase[any, any]{
		Name:  name,
		Setup: TransportSetup[any, any]{Writer: writer},
		Action: func(t *testing.T, transport *tspt.CodecTransport[any, any]) {
			err := transport.Flush()
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{writer.Mock},
	}
}

func CloseTestCase() TransportTestCase[any, any] {
	name := "Close should close the conn"

	var conn = cmock.NewConn()
	conn.RegisterClose(
		func() (err error) { return nil },
	)
	return TransportTestCase[any, any]{
		Name:  name,
		Setup: TransportSetup[any, any]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[any, any]) {
			err := transport.Close()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}

func CloseErrorTestCase() TransportTestCase[any, any] {
	name := "If Conn.Close fails with an error, Close should return it"

	var (
		wantErr = errors.New("Conn.Close error")
		conn    = cmock.NewConn()
	)
	conn.RegisterClose(
		func() (err error) { return wantErr },
	)
	return TransportTestCase[any, any]{
		Name:  name,
		Setup: TransportSetup[any, any]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[any, any]) {
			err := transport.Close()
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}
