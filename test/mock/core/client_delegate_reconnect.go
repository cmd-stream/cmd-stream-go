package core

import (
	"net"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/ymz-ncnk/mok"
)

type ReconnectFn func() error

type ReconnectDelegate struct {
	*mok.Mock
}

func NewReconnectDelegate() ReconnectDelegate {
	return ReconnectDelegate{Mock: mok.New("ReconnectDelegate")}
}

func (m ReconnectDelegate) RegisterLocalAddr(fn ClientLocalAddrFn) ReconnectDelegate {
	m.Register("LocalAddr", fn)
	return m
}

func (m ReconnectDelegate) RegisterRemoteAddr(fn ClientRemoteAddrFn) ReconnectDelegate {
	m.Register("RemoteAddr", fn)
	return m
}

func (m ReconnectDelegate) RegisterSetSendDeadline(fn ClientSetSendDeadlineFn) ReconnectDelegate {
	m.Register("SetSendDeadline", fn)
	return m
}

func (m ReconnectDelegate) RegisterSendN(n int, fn ClientSendFn) ReconnectDelegate {
	m.RegisterN("Send", n, fn)
	return m
}

func (m ReconnectDelegate) RegisterSend(fn ClientSendFn) ReconnectDelegate {
	m.Register("Send", fn)
	return m
}

func (m ReconnectDelegate) RegisterFlush(fn ClientFlushFn) ReconnectDelegate {
	m.Register("Flush", fn)
	return m
}

func (m ReconnectDelegate) RegisterSetReceiveDeadline(fn ClientSetReceiveDeadlineFn) ReconnectDelegate {
	m.Register("SetReceiveDeadline", fn)
	return m
}

func (m ReconnectDelegate) RegisterReceive(fn ClientReceiveFn) ReconnectDelegate {
	m.Register("Receive", fn)
	return m
}

func (m ReconnectDelegate) RegisterClose(fn ClientCloseFn) ReconnectDelegate {
	m.Register("Close", fn)
	return m
}

func (m ReconnectDelegate) RegisterReconnect(fn ReconnectFn) ReconnectDelegate {
	m.Register("Reconnect", fn)
	return m
}

func (m ReconnectDelegate) LocalAddr() (addr net.Addr) {
	vals, err := m.Call("LocalAddr")
	if err != nil {
		panic(err)
	}
	addr, _ = vals[0].(net.Addr)
	return
}

func (m ReconnectDelegate) RemoteAddr() (addr net.Addr) {
	vals, err := m.Call("RemoteAddr")
	if err != nil {
		panic(err)
	}
	addr, _ = vals[0].(net.Addr)
	return
}

func (m ReconnectDelegate) SetSendDeadline(deadline time.Time) (err error) {
	vals, err := m.Call("SetSendDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m ReconnectDelegate) Send(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
	vals, err := m.Call("Send", seq, cmd)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (m ReconnectDelegate) Flush() (err error) {
	vals, err := m.Call("Flush")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m ReconnectDelegate) SetReceiveDeadline(deadline time.Time) (err error) {
	vals, err := m.Call("SetReceiveDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m ReconnectDelegate) Receive() (seq core.Seq, result core.Result, n int, err error) {
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

func (m ReconnectDelegate) Close() (err error) {
	vals, err := m.Call("Close")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m ReconnectDelegate) Reconnect() (err error) {
	vals, err := m.Call("Reconnect")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}
