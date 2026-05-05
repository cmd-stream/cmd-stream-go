package mock

import (
	"net"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/ymz-ncnk/mok"
)

type (
	ProxySendFn             func(result core.Result) (n int, err error)
	ProxySendWithDeadlineFn func(deadline time.Time, result core.Result) (
		int, err error)
	ProxyAtFn  func() time.Time
	ProxySeqFn func() core.Seq
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

func (p Proxy) RegisterSend(fn ProxySendFn) Proxy {
	p.Register("Send", fn)
	return p
}

func (p Proxy) RegisterSendWithDeadline(fn ProxySendWithDeadlineFn) Proxy {
	p.Register("SendWithDeadline", fn)
	return p
}

func (p Proxy) RegisterAt(fn ProxyAtFn) Proxy {
	p.Register("At", fn)
	return p
}

func (p Proxy) RegisterSeq(fn ProxySeqFn) Proxy {
	p.Register("Seq", fn)
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

func (p Proxy) At() (at time.Time) {
	vals, err := p.Call("At")
	if err != nil {
		panic(err)
	}
	at, _ = vals[0].(time.Time)
	return
}

func (p Proxy) Send(result core.Result) (n int, err error) {
	vals, err := p.Call("Send", result)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (p Proxy) SendWithDeadline(deadline time.Time, result core.Result) (
	n int, err error) {
	vals, err := p.Call("SendWithDeadline", deadline, result)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (p Proxy) Seq() (seq core.Seq) {
	vals, err := p.Call("Seq")
	if err != nil {
		panic(err)
	}
	seq = vals[0].(core.Seq)
	return
}
