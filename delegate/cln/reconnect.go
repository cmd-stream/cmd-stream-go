package cln

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/cmd-stream-go/core/cln"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
)

// ReconnectDelegate implements the core.ClientReconnectDelegate interface.
type ReconnectDelegate[T any] struct {
	info       dlgt.ServerInfo
	factory    dlgt.ClientTransportFactory[T]
	options    Options
	transport  *atomic.Value
	mu         sync.Mutex
	closedFlag uint32
}

// NewReconnect creates a new ReconnectDelegate.
func NewReconnect[T any](info dlgt.ServerInfo, factory dlgt.ClientTransportFactory[T],
	opts ...SetOption,
) (delegate *ReconnectDelegate[T], err error) {
	transport, err := factory.New()
	if err != nil {
		return
	}
	o := Options{}
	Apply(&o, opts...)
	err = checkServerInfo(o.ServerInfoReceiveDuration, transport, info)
	if err != nil {
		_ = transport.Close()
		return
	}
	delegate = &ReconnectDelegate[T]{
		info:      info,
		factory:   factory,
		options:   o,
		transport: &atomic.Value{},
	}
	err = delegate.setTransport(transport)
	return
}

// NewReconnectWithoutInfo for tests only.
func NewReconnectWithoutInfo[T any](factory dlgt.ClientTransportFactory[T],
	transport *atomic.Value,
	options Options,
) *ReconnectDelegate[T] {
	return &ReconnectDelegate[T]{
		factory:   factory,
		transport: transport,
		options:   options,
	}
}

func (d *ReconnectDelegate[T]) Options() Options {
	return d.options
}

func (d *ReconnectDelegate[T]) LocalAddr() net.Addr {
	return d.Transport().LocalAddr()
}

func (d *ReconnectDelegate[T]) RemoteAddr() net.Addr {
	return d.Transport().RemoteAddr()
}

func (d *ReconnectDelegate[T]) SetSendDeadline(deadline time.Time) (err error) {
	return d.Transport().SetSendDeadline(deadline)
}

func (d *ReconnectDelegate[T]) Send(seq core.Seq, cmd core.Cmd[T]) (n int,
	err error,
) {
	return d.Transport().Send(seq, cmd)
}

func (d *ReconnectDelegate[T]) Flush() (err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if atomic.LoadUint32(&d.closedFlag) == 1 {
		return cln.ErrClosed
	}
	return d.Transport().Flush()
}

func (d *ReconnectDelegate[T]) SetReceiveDeadline(deadline time.Time) (
	err error,
) {
	return d.Transport().SetReceiveDeadline(deadline)
}

func (d *ReconnectDelegate[T]) Receive() (seq core.Seq, result core.Result,
	n int, err error,
) {
	return d.Transport().Receive()
}

func (d *ReconnectDelegate[T]) Transport() dlgt.ClientTransport[T] {
	return d.transport.Load().(dlgt.ClientTransport[T])
}

func (d *ReconnectDelegate[T]) Reconnect() (err error) {
	var transport dlgt.ClientTransport[T]
	for {
		transport, err = d.initTransport()
		if err != nil {
			if err == cln.ErrClosed || err == ErrServerInfoMismatch {
				return
			}
			continue
		}
		break
	}
	return d.setTransport(transport)
}

func (d *ReconnectDelegate[T]) Close() (err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if atomic.LoadUint32(&d.closedFlag) == 1 {
		return
	}
	err = d.Transport().Close()
	if err == nil {
		atomic.StoreUint32(&d.closedFlag, 1)
	}
	return
}

func (d *ReconnectDelegate[T]) initTransport() (transport dlgt.ClientTransport[T],
	err error,
) {
	if d.closed() {
		return nil, cln.ErrClosed
	}
	transport, err = d.factory.New()
	if err != nil {
		return nil, err
	}
	if d.closed() {
		_ = transport.Close()
		return nil, cln.ErrClosed
	}
	err = checkServerInfo(d.options.ServerInfoReceiveDuration, transport, d.info)
	if err != nil {
		_ = transport.Close()
		return nil, err
	}
	return transport, nil
}

func (d *ReconnectDelegate[T]) setTransport(transport dlgt.ClientTransport[T]) (err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed() {
		_ = transport.Close()
		return cln.ErrClosed
	}
	d.transport.Store(transport)
	return nil
}

func (d *ReconnectDelegate[T]) closed() bool {
	return atomic.LoadUint32(&d.closedFlag) == 1
}
