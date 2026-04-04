package delegate

import (
	"net"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/ymz-ncnk/mok"
)

type ServerTransport[T any] struct {
	*mok.Mock
}

func NewServerTransport[T any]() ServerTransport[T] {
	return ServerTransport[T]{Mock: mok.New("Transport")}
}

func (m ServerTransport[T]) LocalAddr() net.Addr {
	vals, err := m.Call("LocalAddr")
	if err != nil {
		panic(err)
	}
	addr, _ := vals[0].(net.Addr)
	return addr
}

func (m ServerTransport[T]) RemoteAddr() net.Addr {
	vals, err := m.Call("RemoteAddr")
	if err != nil {
		panic(err)
	}
	addr, _ := vals[0].(net.Addr)
	return addr
}

func (m ServerTransport[T]) SetSendDeadline(deadline time.Time) error {
	vals, err := m.Call("SetSendDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return err
}

func (m ServerTransport[T]) Send(seq core.Seq, result core.Result) (n int,
	err error,
) {
	vals, err := m.Call("Send", seq, result)
	if err != nil {
		panic(err)
	}
	n, _ = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (m ServerTransport[T]) Flush() error {
	vals, err := m.Call("Flush")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return err
}

func (m ServerTransport[T]) SetReceiveDeadline(deadline time.Time) error {
	vals, err := m.Call("SetReceiveDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return err
}

func (m ServerTransport[T]) Receive() (seq core.Seq, cmd core.Cmd[T], n int,
	err error,
) {
	vals, err := m.Call("Receive")
	if err != nil {
		panic(err)
	}
	seq, _ = vals[0].(core.Seq)
	cmd, _ = vals[1].(core.Cmd[T])
	n, _ = vals[2].(int)
	err, _ = vals[3].(error)
	return
}

func (m ServerTransport[T]) Close() error {
	vals, err := m.Call("Close")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return err
}

func (m ServerTransport[T]) SendServerInfo(info dlgt.ServerInfo) error {
	vals, err := m.Call("SendServerInfo", info)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return err
}

func (m ServerTransport[T]) RegisterLocalAddr(fn func() net.Addr) ServerTransport[T] {
	m.Register("LocalAddr", fn)
	return m
}

func (m ServerTransport[T]) RegisterRemoteAddr(fn func() net.Addr) ServerTransport[T] {
	m.Register("RemoteAddr", fn)
	return m
}

func (m ServerTransport[T]) RegisterSetSendDeadline(fn func(deadline time.Time) error) ServerTransport[T] {
	m.Register("SetSendDeadline", fn)
	return m
}

func (m ServerTransport[T]) RegisterSend(fn func(seq core.Seq, result core.Result) (int, error)) ServerTransport[T] {
	m.Register("Send", fn)
	return m
}

func (m ServerTransport[T]) RegisterFlush(fn func() error) ServerTransport[T] {
	m.Register("Flush", fn)
	return m
}

func (m ServerTransport[T]) RegisterSetReceiveDeadline(fn func(deadline time.Time) error) ServerTransport[T] {
	m.Register("SetReceiveDeadline", fn)
	return m
}

func (m ServerTransport[T]) RegisterReceive(fn func() (core.Seq, core.Cmd[T], int, error)) ServerTransport[T] {
	m.Register("Receive", fn)
	return m
}

func (m ServerTransport[T]) RegisterReceiveN(n int, fn func() (core.Seq, core.Cmd[T], int, error)) ServerTransport[T] {
	m.RegisterN("Receive", n, fn)
	return m
}

func (m ServerTransport[T]) RegisterClose(fn func() error) ServerTransport[T] {
	m.Register("Close", fn)
	return m
}

func (m ServerTransport[T]) RegisterSendServerInfo(fn func(info dlgt.ServerInfo) error) ServerTransport[T] {
	m.Register("SendServerInfo", fn)
	return m
}
