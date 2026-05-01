package mock

import (
	"net"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/ymz-ncnk/mok"
)

type ClientTransport[T any] struct {
	*mok.Mock
}

func NewClientTransport[T any]() ClientTransport[T] {
	return ClientTransport[T]{Mock: mok.New("Transport")}
}

func (m ClientTransport[T]) LocalAddr() net.Addr {
	res, err := m.Call("LocalAddr")
	if err != nil {
		panic(err)
	}
	addr, _ := res[0].(net.Addr)
	return addr
}

func (m ClientTransport[T]) RemoteAddr() net.Addr {
	res, err := m.Call("RemoteAddr")
	if err != nil {
		panic(err)
	}
	addr, _ := res[0].(net.Addr)
	return addr
}

func (m ClientTransport[T]) SetSendDeadline(deadline time.Time) error {
	res, err := m.Call("SetSendDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = res[0].(error)
	return err
}

func (m ClientTransport[T]) Send(seq core.Seq, cmd core.Cmd[T]) (n int, err error) {
	res, err := m.Call("Send", seq, cmd)
	if err != nil {
		panic(err)
	}
	n, _ = res[0].(int)
	err, _ = res[1].(error)
	return
}

func (m ClientTransport[T]) Flush() error {
	res, err := m.Call("Flush")
	if err != nil {
		panic(err)
	}
	err, _ = res[0].(error)
	return err
}

func (m ClientTransport[T]) SetReceiveDeadline(deadline time.Time) error {
	res, err := m.Call("SetReceiveDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = res[0].(error)
	return err
}

func (m ClientTransport[T]) Receive() (seq core.Seq, result core.Result, n int,
	err error,
) {
	res, err := m.Call("Receive")
	if err != nil {
		panic(err)
	}
	seq, _ = res[0].(core.Seq)
	result, _ = res[1].(core.Result)
	n, _ = res[2].(int)
	err, _ = res[3].(error)
	return
}

func (m ClientTransport[T]) Close() error {
	res, err := m.Call("Close")
	if err != nil {
		panic(err)
	}
	err, _ = res[0].(error)
	return err
}

func (m ClientTransport[T]) ReceiveServerInfo() (info dlgt.ServerInfo, err error) {
	res, err := m.Call("ReceiveServerInfo")
	if err != nil {
		panic(err)
	}
	info, _ = res[0].(dlgt.ServerInfo)
	err, _ = res[1].(error)
	return
}

func (m ClientTransport[T]) RegisterLocalAddr(fn func() net.Addr) ClientTransport[T] {
	m.Register("LocalAddr", fn)
	return m
}

func (m ClientTransport[T]) RegisterRemoteAddr(fn func() net.Addr) ClientTransport[T] {
	m.Register("RemoteAddr", fn)
	return m
}

func (m ClientTransport[T]) RegisterSetSendDeadline(fn func(deadline time.Time) error) ClientTransport[T] {
	m.Register("SetSendDeadline", fn)
	return m
}

func (m ClientTransport[T]) RegisterSend(fn func(seq core.Seq, cmd core.Cmd[T]) (int, error)) ClientTransport[T] {
	m.Register("Send", fn)
	return m
}

func (m ClientTransport[T]) RegisterFlush(fn func() error) ClientTransport[T] {
	m.Register("Flush", fn)
	return m
}

func (m ClientTransport[T]) RegisterSetReceiveDeadline(fn func(deadline time.Time) error) ClientTransport[T] {
	m.Register("SetReceiveDeadline", fn)
	return m
}

func (m ClientTransport[T]) RegisterReceive(fn func() (core.Seq, core.Result, int, error)) ClientTransport[T] {
	m.Register("Receive", fn)
	return m
}

func (m ClientTransport[T]) RegisterClose(fn func() error) ClientTransport[T] {
	m.Register("Close", fn)
	return m
}

func (m ClientTransport[T]) RegisterReceiveServerInfo(fn func() (dlgt.ServerInfo, error)) ClientTransport[T] {
	m.Register("ReceiveServerInfo", fn)
	return m
}
