package core

import (
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/ymz-ncnk/mok"
)

type KeepaliveFn func(muSn *sync.Mutex)

type KeepaliveDelegate struct {
	*mok.Mock
}

func NewKeepaliveDelegate() KeepaliveDelegate {
	return KeepaliveDelegate{Mock: mok.New("KeepaliveDelegate")}
}

func (m KeepaliveDelegate) RegisterLocalAddr(fn ClientLocalAddrFn) KeepaliveDelegate {
	m.Register("LocalAddr", fn)
	return m
}

func (m KeepaliveDelegate) RegisterRemoteAddr(fn ClientRemoteAddrFn) KeepaliveDelegate {
	m.Register("RemoteAddr", fn)
	return m
}

func (m KeepaliveDelegate) RegisterSetSendDeadline(fn ClientSetSendDeadlineFn) KeepaliveDelegate {
	m.Register("SetSendDeadline", fn)
	return m
}

func (m KeepaliveDelegate) RegisterSendN(n int, fn ClientSendFn) KeepaliveDelegate {
	m.RegisterN("Send", n, fn)
	return m
}

func (m KeepaliveDelegate) RegisterSend(fn ClientSendFn) KeepaliveDelegate {
	m.Register("Send", fn)
	return m
}

func (m KeepaliveDelegate) RegisterFlush(fn ClientFlushFn) KeepaliveDelegate {
	m.Register("Flush", fn)
	return m
}

func (m KeepaliveDelegate) RegisterSetReceiveDeadline(fn ClientSetReceiveDeadlineFn) KeepaliveDelegate {
	m.Register("SetReceiveDeadline", fn)
	return m
}

func (m KeepaliveDelegate) RegisterReceive(fn ClientReceiveFn) KeepaliveDelegate {
	m.Register("Receive", fn)
	return m
}

func (m KeepaliveDelegate) RegisterClose(fn ClientCloseFn) KeepaliveDelegate {
	m.Register("Close", fn)
	return m
}

func (m KeepaliveDelegate) RegisterKeepalive(fn KeepaliveFn) KeepaliveDelegate {
	m.Register("Keepalive", fn)
	return m
}

func (m KeepaliveDelegate) LocalAddr() (addr net.Addr) {
	vals, err := m.Call("LocalAddr")
	if err != nil {
		panic(err)
	}
	addr, _ = vals[0].(net.Addr)
	return
}

func (m KeepaliveDelegate) RemoteAddr() (addr net.Addr) {
	vals, err := m.Call("RemoteAddr")
	if err != nil {
		panic(err)
	}
	addr, _ = vals[0].(net.Addr)
	return
}

func (m KeepaliveDelegate) SetSendDeadline(deadline time.Time) (err error) {
	vals, err := m.Call("SetSendDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m KeepaliveDelegate) Send(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
	vals, err := m.Call("Send", seq, mok.SafeVal[core.Cmd[any]](cmd))
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (m KeepaliveDelegate) Flush() (err error) {
	vals, err := m.Call("Flush")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m KeepaliveDelegate) SetReceiveDeadline(deadline time.Time) (err error) {
	vals, err := m.Call("SetReceiveDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m KeepaliveDelegate) Receive() (seq core.Seq, result core.Result, n int, err error) {
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

func (m KeepaliveDelegate) Close() (err error) {
	vals, err := m.Call("Close")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m KeepaliveDelegate) Keepalive(muSn *sync.Mutex) {
	_, err := m.Call("Keepalive", muSn)
	if err != nil {
		panic(err)
	}
}
