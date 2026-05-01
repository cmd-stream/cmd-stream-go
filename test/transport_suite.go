package test

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/cmd-stream-go/test/mock"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type TransportSetup[T, V any] struct {
	Conn   net.Conn
	Writer tspt.Writer
	Reader tspt.Reader
	Codec  tspt.Codec[T, V]
}

type TransportTestCase[T, V any] struct {
	Name   string
	Setup  TransportSetup[T, V]
	Action func(t *testing.T, transport *tspt.CodecTransport[T, V])
	Mocks  []*mok.Mock
}

func RunTransportTestCase[T, V any](t *testing.T, tc TransportTestCase[T, V]) {
	t.Run(tc.Name, func(t *testing.T) {
		transport := tspt.New(tc.Setup.Conn, tc.Setup.Writer, tc.Setup.Reader,
			tc.Setup.Codec)
		tc.Action(t, transport)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}

type TransportSuite[T, V any] struct{}

// -----------------------------------------------------------------------------
// Test Cases
// -----------------------------------------------------------------------------

func (TransportSuite[T, V]) LocalAddr(t *testing.T) TransportTestCase[T, V] {
	name := "LocalAddr should return local address of the conn"

	var (
		wantAddr = &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9000}
		conn     = mock.NewConn()
	)
	conn.RegisterLocalAddr(
		func() (addr net.Addr) { return wantAddr },
	)
	return TransportTestCase[T, V]{
		Name:  name,
		Setup: TransportSetup[T, V]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[T, V]) {
			addr := transport.LocalAddr()
			asserterror.EqualDeep(t, addr, wantAddr)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}

func (TransportSuite[T, V]) RemoteAddr(t *testing.T) TransportTestCase[T, V] {
	name := "RemoteAddr should return remote address of the conn"

	var (
		wantAddr = &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9000}
		conn     = mock.NewConn()
	)
	conn.RegisterRemoteAddr(
		func() (addr net.Addr) { return wantAddr },
	)
	return TransportTestCase[T, V]{
		Name:  name,
		Setup: TransportSetup[T, V]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[T, V]) {
			addr := transport.RemoteAddr()
			asserterror.EqualDeep(t, addr, wantAddr)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}

func (TransportSuite[T, V]) SetSendDeadline(t *testing.T) TransportTestCase[T, V] {
	name := "Conn.SetWriteDeadline should receive same deadline as SetSendDeadline"

	var (
		wantDeadline = time.Now()
		conn         = mock.NewConn()
	)
	conn.RegisterSetWriteDeadline(
		func(deadline time.Time) (err error) {
			asserterror.Equal(t, deadline, wantDeadline)
			return
		},
	)
	return TransportTestCase[T, V]{
		Name:  name,
		Setup: TransportSetup[T, V]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[T, V]) {
			_ = transport.SetSendDeadline(wantDeadline)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}

func (TransportSuite[T, V]) SetSendDeadlineError(t *testing.T) TransportTestCase[T, V] {
	name := "If Conn.SetWriteDeadline fails with an error, SetSendDeadline should return it"

	var (
		wantErr = errors.New("Conn.SetWriteDeadline error")
		conn    = mock.NewConn()
	)
	conn.RegisterSetWriteDeadline(
		func(deadline time.Time) (err error) { return wantErr },
	)
	return TransportTestCase[T, V]{
		Name:  name,
		Setup: TransportSetup[T, V]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[T, V]) {
			err := transport.SetSendDeadline(time.Time{})
			asserterror.Equal(t, err, wantErr)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}

func (TransportSuite[T, V]) Send(t *testing.T) TransportTestCase[T, V] {
	name := "Send should encode data with help of the Codec"

	var (
		wantSeq core.Seq = 1
		wantT   T
		wantN         = 3
		wantErr error = nil
		writer        = mock.NewWriter()
		codec         = mock.NewCodec[T, V]()
	)
	codec.RegisterEncode(
		func(seq core.Seq, val T, w tspt.Writer) (n int, err error) {
			asserterror.Equal(t, seq, wantSeq)
			asserterror.EqualDeep(t, val, wantT)
			return wantN, wantErr
		},
	)
	return TransportTestCase[T, V]{
		Name: name,
		Setup: TransportSetup[T, V]{
			Writer: writer,
			Codec:  codec,
		},
		Action: func(t *testing.T, transport *tspt.CodecTransport[T, V]) {
			n, err := transport.Send(wantSeq, wantT)
			asserterror.EqualError(t, err, wantErr)
			asserterror.Equal(t, n, wantN)
		},
		Mocks: []*mok.Mock{writer.Mock, codec.Mock},
	}
}

func (TransportSuite[T, V]) SendError(t *testing.T) TransportTestCase[T, V] {
	name := "If Codec.Encode fails with an error, Send should return it"

	var (
		wantErr = errors.New("Codec.Encode error")
		codec   = mock.NewCodec[T, V]()
	)
	codec.RegisterEncode(
		func(seq core.Seq, val T, w tspt.Writer) (n int, err error) {
			return 0, wantErr
		},
	)
	return TransportTestCase[T, V]{
		Name: name,
		Setup: TransportSetup[T, V]{
			Codec: codec,
		},
		Action: func(t *testing.T, transport *tspt.CodecTransport[T, V]) {
			var val T
			_, err := transport.Send(1, val)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{codec.Mock},
	}
}

func (TransportSuite[T, V]) SetReceiveDeadline(t *testing.T) TransportTestCase[T, V] {
	name := "Conn.SetReadDeadline should receive same deadline as SetReceiveDeadline"

	var (
		wantDeadline = time.Now()
		conn         = mock.NewConn()
	)
	conn.RegisterSetReadDeadline(
		func(deadline time.Time) (err error) {
			asserterror.Equal(t, deadline, wantDeadline)
			return
		},
	)
	return TransportTestCase[T, V]{
		Name:  name,
		Setup: TransportSetup[T, V]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[T, V]) {
			_ = transport.SetReceiveDeadline(wantDeadline)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}

func (TransportSuite[T, V]) SetReceiveDeadlineError(t *testing.T) TransportTestCase[T, V] {
	name := "If Conn.SetReadDeadline fails with an error, SetReceiveDeadline should return it"

	var (
		wantErr = errors.New("Conn.SetReadDeadline error")
		conn    = mock.NewConn()
	)
	conn.RegisterSetReadDeadline(
		func(deadline time.Time) (err error) { return wantErr },
	)
	return TransportTestCase[T, V]{
		Name:  name,
		Setup: TransportSetup[T, V]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[T, V]) {
			err := transport.SetReceiveDeadline(time.Time{})
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}

func (TransportSuite[T, V]) Receive(t *testing.T) TransportTestCase[T, V] {
	name := "Receive should decode data with help of the Codec"

	var (
		wantSeq core.Seq = 1
		wantV   V
		wantN         = 3
		wantErr error = nil
		codec         = mock.NewCodec[T, V]()
	)
	codec.RegisterDecode(
		func(r tspt.Reader) (seq core.Seq, val V, n int, err error) {
			return wantSeq, wantV, wantN, wantErr
		},
	)
	return TransportTestCase[T, V]{
		Name: name,
		Setup: TransportSetup[T, V]{
			Codec: codec,
		},
		Action: func(t *testing.T, transport *tspt.CodecTransport[T, V]) {
			seq, val, n, err := transport.Receive()
			asserterror.EqualError(t, err, wantErr)
			asserterror.Equal(t, seq, wantSeq)
			asserterror.EqualDeep(t, val, wantV)
			asserterror.Equal(t, n, wantN)
		},
		Mocks: []*mok.Mock{codec.Mock},
	}
}

func (TransportSuite[T, V]) ReceiveError(t *testing.T) TransportTestCase[T, V] {
	name := "If Codec.Decode fails with an error, Receive should return it"

	var (
		wantSeq core.Seq = 0
		wantV   V
		wantErr = errors.New("Codec.Decode error")
		codec   = mock.NewCodec[T, V]()
	)
	codec.RegisterDecode(
		func(r tspt.Reader) (seq core.Seq, val V, n int, err error) {
			err = wantErr
			return
		},
	)
	return TransportTestCase[T, V]{
		Name: name,
		Setup: TransportSetup[T, V]{
			Codec: codec,
		},
		Action: func(t *testing.T, transport *tspt.CodecTransport[T, V]) {
			seq, val, _, err := transport.Receive()
			asserterror.EqualError(t, err, wantErr)
			asserterror.Equal(t, seq, wantSeq)
			asserterror.EqualDeep(t, val, wantV)
		},
		Mocks: []*mok.Mock{codec.Mock},
	}
}

func (TransportSuite[T, V]) Flush(t *testing.T) TransportTestCase[T, V] {
	name := "Flush should flush the writer"

	var writer = mock.NewWriter()
	writer.RegisterFlush(
		func() error { return nil },
	)
	return TransportTestCase[T, V]{
		Name:  name,
		Setup: TransportSetup[T, V]{Writer: writer},
		Action: func(t *testing.T, transport *tspt.CodecTransport[T, V]) {
			err := transport.Flush()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{writer.Mock},
	}
}

func (TransportSuite[T, V]) FlushError(t *testing.T) TransportTestCase[T, V] {
	name := "If Writer.Flush fails with an error, Flush should return it"

	var (
		wantErr = errors.New("Writer.Flush error")
		writer  = mock.NewWriter()
	)
	writer.RegisterFlush(
		func() error { return wantErr },
	)
	return TransportTestCase[T, V]{
		Name:  name,
		Setup: TransportSetup[T, V]{Writer: writer},
		Action: func(t *testing.T, transport *tspt.CodecTransport[T, V]) {
			err := transport.Flush()
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{writer.Mock},
	}
}

func (TransportSuite[T, V]) Close(t *testing.T) TransportTestCase[T, V] {
	name := "Close should close the conn"

	var conn = mock.NewConn()
	conn.RegisterClose(
		func() (err error) { return nil },
	)
	return TransportTestCase[T, V]{
		Name:  name,
		Setup: TransportSetup[T, V]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[T, V]) {
			err := transport.Close()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}

func (TransportSuite[T, V]) CloseError(t *testing.T) TransportTestCase[T, V] {
	name := "If Conn.Close fails with an error, Close should return it"

	var (
		wantErr = errors.New("Conn.Close error")
		conn    = mock.NewConn()
	)
	conn.RegisterClose(
		func() (err error) { return wantErr },
	)
	return TransportTestCase[T, V]{
		Name:  name,
		Setup: TransportSetup[T, V]{Conn: conn},
		Action: func(t *testing.T, transport *tspt.CodecTransport[T, V]) {
			err := transport.Close()
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{conn.Mock},
	}
}
