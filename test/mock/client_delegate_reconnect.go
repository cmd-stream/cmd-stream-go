package mock

import (
	"net"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/ymz-ncnk/mok"
)

type ReconnectFn func() error

func NewReconnectDelegate[T any]() ReconnectDelegate[T] {
	return ReconnectDelegate[T]{Mock: mok.New("ReconnectDelegate")}
}

type ReconnectDelegate[T any] struct {
	*mok.Mock
}

func (m ReconnectDelegate[T]) RegisterLocalAddr(fn ClientLocalAddrFn) ReconnectDelegate[T] {
	m.Register("LocalAddr", fn)
	return m
}

func (m ReconnectDelegate[T]) RegisterRemoteAddr(fn ClientRemoteAddrFn) ReconnectDelegate[T] {
	m.Register("RemoteAddr", fn)
	return m
}

func (m ReconnectDelegate[T]) RegisterSetSendDeadline(fn ClientSetSendDeadlineFn) ReconnectDelegate[T] {
	m.Register("SetSendDeadline", fn)
	return m
}

func (m ReconnectDelegate[T]) RegisterSendN(n int, fn ClientSendFn[T]) ReconnectDelegate[T] {
	m.RegisterN("Send", n, fn)
	return m
}

func (m ReconnectDelegate[T]) RegisterSend(fn ClientSendFn[T]) ReconnectDelegate[T] {
	m.Register("Send", fn)
	return m
}

func (m ReconnectDelegate[T]) RegisterFlush(fn ClientFlushFn) ReconnectDelegate[T] {
	m.Register("Flush", fn)
	return m
}

func (m ReconnectDelegate[T]) RegisterSetReceiveDeadline(fn ClientSetReceiveDeadlineFn) ReconnectDelegate[T] {
	m.Register("SetReceiveDeadline", fn)
	return m
}

func (m ReconnectDelegate[T]) RegisterReceive(fn ClientReceiveFn) ReconnectDelegate[T] {
	m.Register("Receive", fn)
	return m
}

func (m ReconnectDelegate[T]) RegisterClose(fn ClientCloseFn) ReconnectDelegate[T] {
	m.Register("Close", fn)
	return m
}

func (m ReconnectDelegate[T]) RegisterReconnect(fn ReconnectFn) ReconnectDelegate[T] {
	m.Register("Reconnect", fn)
	return m
}

func (m ReconnectDelegate[T]) LocalAddr() (addr net.Addr) {
	vals, err := m.Call("LocalAddr")
	if err != nil {
		panic(err)
	}
	addr, _ = vals[0].(net.Addr)
	return
}

func (m ReconnectDelegate[T]) RemoteAddr() (addr net.Addr) {
	vals, err := m.Call("RemoteAddr")
	if err != nil {
		panic(err)
	}
	addr, _ = vals[0].(net.Addr)
	return
}

func (m ReconnectDelegate[T]) SetSendDeadline(deadline time.Time) (err error) {
	vals, err := m.Call("SetSendDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m ReconnectDelegate[T]) Send(seq core.Seq, cmd core.Cmd[T]) (n int, err error) {
	vals, err := m.Call("Send", seq, cmd)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (m ReconnectDelegate[T]) Flush() (err error) {
	vals, err := m.Call("Flush")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m ReconnectDelegate[T]) SetReceiveDeadline(deadline time.Time) (err error) {
	vals, err := m.Call("SetReceiveDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m ReconnectDelegate[T]) Receive() (seq core.Seq, result core.Result, n int, err error) {
	vals, err := m.Call("Receive")
	if err != nil {
		panic(err)
	}
	seq = vals[0].(core.Seq)
	result, _ = vals[1].(core.Result)
	n = vals[2].(int)
	err, _ = vals[3].(error)
	return
}

func (m ReconnectDelegate[T]) Close() (err error) {
	vals, err := m.Call("Close")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m ReconnectDelegate[T]) Reconnect() (err error) {
	vals, err := m.Call("Reconnect")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}
