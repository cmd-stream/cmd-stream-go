package mock

import (
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/ymz-ncnk/mok"
)

type KeepaliveFn func(muSn *sync.Mutex)

func NewKeepaliveDelegate[T any]() KeepaliveDelegate[T] {
	return KeepaliveDelegate[T]{Mock: mok.New("KeepaliveDelegate")}
}

type KeepaliveDelegate[T any] struct {
	*mok.Mock
}

func (m KeepaliveDelegate[T]) RegisterLocalAddr(fn ClientLocalAddrFn) KeepaliveDelegate[T] {
	m.Register("LocalAddr", fn)
	return m
}

func (m KeepaliveDelegate[T]) RegisterRemoteAddr(fn ClientRemoteAddrFn) KeepaliveDelegate[T] {
	m.Register("RemoteAddr", fn)
	return m
}

func (m KeepaliveDelegate[T]) RegisterSetSendDeadline(fn ClientSetSendDeadlineFn) KeepaliveDelegate[T] {
	m.Register("SetSendDeadline", fn)
	return m
}

func (m KeepaliveDelegate[T]) RegisterSendN(n int, fn ClientSendFn[T]) KeepaliveDelegate[T] {
	m.RegisterN("Send", n, fn)
	return m
}

func (m KeepaliveDelegate[T]) RegisterSend(fn ClientSendFn[T]) KeepaliveDelegate[T] {
	m.Register("Send", fn)
	return m
}

func (m KeepaliveDelegate[T]) RegisterFlush(fn ClientFlushFn) KeepaliveDelegate[T] {
	m.Register("Flush", fn)
	return m
}

func (m KeepaliveDelegate[T]) RegisterSetReceiveDeadline(fn ClientSetReceiveDeadlineFn) KeepaliveDelegate[T] {
	m.Register("SetReceiveDeadline", fn)
	return m
}

func (m KeepaliveDelegate[T]) RegisterReceive(fn ClientReceiveFn) KeepaliveDelegate[T] {
	m.Register("Receive", fn)
	return m
}

func (m KeepaliveDelegate[T]) RegisterClose(fn ClientCloseFn) KeepaliveDelegate[T] {
	m.Register("Close", fn)
	return m
}

func (m KeepaliveDelegate[T]) RegisterKeepalive(fn KeepaliveFn) KeepaliveDelegate[T] {
	m.Register("Keepalive", fn)
	return m
}

func (m KeepaliveDelegate[T]) LocalAddr() (addr net.Addr) {
	vals, err := m.Call("LocalAddr")
	if err != nil {
		panic(err)
	}
	addr, _ = vals[0].(net.Addr)
	return
}

func (m KeepaliveDelegate[T]) RemoteAddr() (addr net.Addr) {
	vals, err := m.Call("RemoteAddr")
	if err != nil {
		panic(err)
	}
	addr, _ = vals[0].(net.Addr)
	return
}

func (m KeepaliveDelegate[T]) SetSendDeadline(deadline time.Time) (err error) {
	vals, err := m.Call("SetSendDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m KeepaliveDelegate[T]) Send(seq core.Seq, cmd core.Cmd[T]) (n int, err error) {
	vals, err := m.Call("Send", seq, cmd)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (m KeepaliveDelegate[T]) Flush() (err error) {
	vals, err := m.Call("Flush")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m KeepaliveDelegate[T]) SetReceiveDeadline(deadline time.Time) (err error) {
	vals, err := m.Call("SetReceiveDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m KeepaliveDelegate[T]) Receive() (seq core.Seq, result core.Result, n int, err error) {
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

func (m KeepaliveDelegate[T]) Close() (err error) {
	vals, err := m.Call("Close")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m KeepaliveDelegate[T]) Keepalive(muSn *sync.Mutex) {
	_, err := m.Call("Keepalive", muSn)
	if err != nil {
		panic(err)
	}
}
