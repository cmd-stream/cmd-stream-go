package core

import (
	"net"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/ymz-ncnk/mok"
)

type (
	SendFn             func(seq core.Seq, result core.Result) (n int, err error)
	SendWithDeadlineFn func(deadline time.Time, seq core.Seq, result core.Result) (
		int, err error)
)

type Proxy struct {
	*mok.Mock
}

func NewProxy() Proxy {
	return Proxy{mok.New("Proxy")}
}

func (p Proxy) RegisterLocalAddr(fn LocalAddrFn) Proxy {
	p.Register("LocalAddr", fn)
	return p
}

func (p Proxy) RegisterRemoteAddr(fn RemoteAddrFn) Proxy {
	p.Register("RemoteAddr", fn)
	return p
}

func (p Proxy) RegisterSend(fn SendFn) Proxy {
	p.Register("Send", fn)
	return p
}

func (p Proxy) RegisterSendWithDeadline(fn SendWithDeadlineFn) Proxy {
	p.Register("SendWithDeadline", fn)
	return p
}

func (p Proxy) LocalAddr() (addr net.Addr) {
	vals, err := p.Call("LocalAddr")
	if err != nil {
		panic(err)
	}
	addr, _ = vals[0].(net.Addr)
	return
}

func (p Proxy) RemoteAddr() (addr net.Addr) {
	vals, err := p.Call("RemoteAddr")
	if err != nil {
		panic(err)
	}
	addr, _ = vals[0].(net.Addr)
	return
}

func (p Proxy) Send(seq core.Seq, result core.Result) (n int, err error) {
	vals, err := p.Call("Send", seq, result)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (p Proxy) SendWithDeadline(deadline time.Time, seq core.Seq, result core.Result,
) (n int, err error) {
	vals, err := p.Call("SendWithDeadline", deadline, seq, result)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}
