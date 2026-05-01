package mock

import (
	"net"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/ymz-ncnk/mok"
)

type (
	ClientLocalAddrFn          func() (addr net.Addr)
	ClientRemoteAddrFn         func() (addr net.Addr)
	ClientSetSendDeadlineFn    func(deadline time.Time) (err error)
	ClientSendFn[T any]        func(seq core.Seq, cmd core.Cmd[T]) (n int, err error)
	ClientFlushFn              func() (err error)
	ClientSetReceiveDeadlineFn func(deadline time.Time) (err error)
	ClientReceiveFn            func() (seq core.Seq, result core.Result, n int, err error)
	ClientCloseFn              func() (err error)
)

func NewClientDelegate[T any]() ClientDelegate[T] {
	return ClientDelegate[T]{Mock: mok.New("Delegate")}
}

type ClientDelegate[T any] struct {
	*mok.Mock
}

func (m ClientDelegate[T]) RegisterLocalAddr(fn ClientLocalAddrFn) ClientDelegate[T] {
	m.Register("LocalAddr", fn)
	return m
}

func (m ClientDelegate[T]) RegisterRemoteAddr(fn ClientRemoteAddrFn) ClientDelegate[T] {
	m.Register("RemoteAddr", fn)
	return m
}

func (m ClientDelegate[T]) RegisterSetSendDeadlineN(n int, fn ClientSetSendDeadlineFn) ClientDelegate[T] {
	m.RegisterN("SetSendDeadline", n, fn)
	return m
}

func (m ClientDelegate[T]) RegisterSetSendDeadline(fn ClientSetSendDeadlineFn) ClientDelegate[T] {
	m.Register("SetSendDeadline", fn)
	return m
}

func (m ClientDelegate[T]) RegisterSendN(n int, fn ClientSendFn[T]) ClientDelegate[T] {
	m.RegisterN("Send", n, fn)
	return m
}

func (m ClientDelegate[T]) RegisterSend(fn ClientSendFn[T]) ClientDelegate[T] {
	m.Register("Send", fn)
	return m
}

func (m ClientDelegate[T]) RegisterFlushN(n int, fn ClientFlushFn) ClientDelegate[T] {
	m.RegisterN("Flush", n, fn)
	return m
}

func (m ClientDelegate[T]) RegisterFlush(fn ClientFlushFn) ClientDelegate[T] {
	m.Register("Flush", fn)
	return m
}

func (m ClientDelegate[T]) RegisterSetReceiveDeadlineN(n int,
	fn ClientSetReceiveDeadlineFn,
) ClientDelegate[T] {
	m.RegisterN("SetReceiveDeadline", n, fn)
	return m
}

func (m ClientDelegate[T]) RegisterSetReceiveDeadline(fn ClientSetReceiveDeadlineFn) ClientDelegate[T] {
	m.Register("SetReceiveDeadline", fn)
	return m
}

func (m ClientDelegate[T]) RegisterReceive(fn ClientReceiveFn) ClientDelegate[T] {
	m.Register("Receive", fn)
	return m
}

func (m ClientDelegate[T]) RegisterReceiveN(n int, fn ClientReceiveFn) ClientDelegate[T] {
	m.RegisterN("Receive", n, fn)
	return m
}

func (m ClientDelegate[T]) RegisterClose(fn ClientCloseFn) ClientDelegate[T] {
	m.Register("Close", fn)
	return m
}

func (m ClientDelegate[T]) LocalAddr() (addr net.Addr) {
	vals, err := m.Call("LocalAddr")
	if err != nil {
		panic(err)
	}
	addr, _ = vals[0].(net.Addr)
	return
}

func (m ClientDelegate[T]) RemoteAddr() (addr net.Addr) {
	vals, err := m.Call("RemoteAddr")
	if err != nil {
		panic(err)
	}
	addr, _ = vals[0].(net.Addr)
	return
}

func (m ClientDelegate[T]) SetSendDeadline(deadline time.Time) (err error) {
	vals, err := m.Call("SetSendDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m ClientDelegate[T]) Send(seq core.Seq, cmd core.Cmd[T]) (n int, err error) {
	vals, err := m.Call("Send", seq, cmd)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (m ClientDelegate[T]) Flush() (err error) {
	vals, err := m.Call("Flush")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m ClientDelegate[T]) SetReceiveDeadline(deadline time.Time) (err error) {
	vals, err := m.Call("SetReceiveDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m ClientDelegate[T]) Receive() (seq core.Seq, result core.Result, n int,
	err error,
) {
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

func (m ClientDelegate[T]) Close() (err error) {
	vals, err := m.Call("Close")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}
