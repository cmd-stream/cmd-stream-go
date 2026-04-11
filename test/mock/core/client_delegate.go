package core

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
	ClientSendFn               func(seq core.Seq, cmd core.Cmd[any]) (n int, err error)
	ClientFlushFn              func() (err error)
	ClientSetReceiveDeadlineFn func(deadline time.Time) (err error)
	ClientReceiveFn            func() (seq core.Seq, result core.Result, n int, err error)
	ClientCloseFn              func() (err error)
)

func NewClientDelegate() ClientDelegate {
	return ClientDelegate{Mock: mok.New("Delegate")}
}

type ClientDelegate struct {
	*mok.Mock
}

func (m ClientDelegate) RegisterLocalAddr(fn ClientLocalAddrFn) ClientDelegate {
	m.Register("LocalAddr", fn)
	return m
}

func (m ClientDelegate) RegisterRemoteAddr(fn ClientRemoteAddrFn) ClientDelegate {
	m.Register("RemoteAddr", fn)
	return m
}

func (m ClientDelegate) RegisterSetSendDeadlineN(n int, fn ClientSetSendDeadlineFn) ClientDelegate {
	m.RegisterN("SetSendDeadline", n, fn)
	return m
}

func (m ClientDelegate) RegisterSetSendDeadline(fn ClientSetSendDeadlineFn) ClientDelegate {
	m.Register("SetSendDeadline", fn)
	return m
}

func (m ClientDelegate) RegisterSendN(n int, fn ClientSendFn) ClientDelegate {
	m.RegisterN("Send", n, fn)
	return m
}

func (m ClientDelegate) RegisterSend(fn ClientSendFn) ClientDelegate {
	m.Register("Send", fn)
	return m
}

func (m ClientDelegate) RegisterFlushN(n int, fn ClientFlushFn) ClientDelegate {
	m.RegisterN("Flush", n, fn)
	return m
}

func (m ClientDelegate) RegisterFlush(fn ClientFlushFn) ClientDelegate {
	m.Register("Flush", fn)
	return m
}

func (m ClientDelegate) RegisterSetReceiveDeadlineN(n int,
	fn ClientSetReceiveDeadlineFn,
) ClientDelegate {
	m.RegisterN("SetReceiveDeadline", n, fn)
	return m
}

func (m ClientDelegate) RegisterSetReceiveDeadline(fn ClientSetReceiveDeadlineFn) ClientDelegate {
	m.Register("SetReceiveDeadline", fn)
	return m
}

func (m ClientDelegate) RegisterReceive(fn ClientReceiveFn) ClientDelegate {
	m.Register("Receive", fn)
	return m
}

func (m ClientDelegate) RegisterReceiveN(n int, fn ClientReceiveFn) ClientDelegate {
	m.RegisterN("Receive", n, fn)
	return m
}

func (m ClientDelegate) RegisterClose(fn ClientCloseFn) ClientDelegate {
	m.Register("Close", fn)
	return m
}

func (m ClientDelegate) LocalAddr() (addr net.Addr) {
	vals, err := m.Call("LocalAddr")
	if err != nil {
		panic(err)
	}
	addr, _ = vals[0].(net.Addr)
	return
}

func (m ClientDelegate) RemoteAddr() (addr net.Addr) {
	vals, err := m.Call("RemoteAddr")
	if err != nil {
		panic(err)
	}
	addr, _ = vals[0].(net.Addr)
	return
}

func (m ClientDelegate) SetSendDeadline(deadline time.Time) (err error) {
	vals, err := m.Call("SetSendDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m ClientDelegate) Send(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
	vals, err := m.Call("Send", seq, cmd)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (m ClientDelegate) Flush() (err error) {
	vals, err := m.Call("Flush")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m ClientDelegate) SetReceiveDeadline(deadline time.Time) (err error) {
	vals, err := m.Call("SetReceiveDeadline", deadline)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m ClientDelegate) Receive() (seq core.Seq, result core.Result, n int,
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

func (m ClientDelegate) Close() (err error) {
	vals, err := m.Call("Close")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}
