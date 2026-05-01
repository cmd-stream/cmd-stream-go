package cln

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	dcln "github.com/cmd-stream/cmd-stream-go/core/cln"
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

// Options returns the delegate's options.
func (d *ReconnectDelegate[T]) Options() Options {
	return d.options
}

// LocalAddr returns the local network address.
func (d *ReconnectDelegate[T]) LocalAddr() net.Addr {
	return d.Transport().LocalAddr()
}

// RemoteAddr returns the remote network address.
func (d *ReconnectDelegate[T]) RemoteAddr() net.Addr {
	return d.Transport().RemoteAddr()
}

// SetSendDeadline sets the deadline for future Send calls.
func (d *ReconnectDelegate[T]) SetSendDeadline(deadline time.Time) (err error) {
	return d.Transport().SetSendDeadline(deadline)
}

// Send transmits a command to the server.
func (d *ReconnectDelegate[T]) Send(seq core.Seq, cmd core.Cmd[T]) (n int,
	err error,
) {
	return d.Transport().Send(seq, cmd)
}

// Flush flushes the current transport's buffer.
func (d *ReconnectDelegate[T]) Flush() (err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if atomic.LoadUint32(&d.closedFlag) == 1 {
		return dcln.ErrClosed
	}
	return d.Transport().Flush()
}

// SetReceiveDeadline sets the deadline for future Receive calls.
func (d *ReconnectDelegate[T]) SetReceiveDeadline(deadline time.Time) (
	err error,
) {
	return d.Transport().SetReceiveDeadline(deadline)
}

// Receive waits for and returns the next result from the server.
func (d *ReconnectDelegate[T]) Receive() (seq core.Seq, result core.Result,
	n int, err error,
) {
	return d.Transport().Receive()
}

// Transport returns the current underlying transport.
func (d *ReconnectDelegate[T]) Transport() dlgt.ClientTransport[T] {
	return d.transport.Load().(dlgt.ClientTransport[T])
}

// Reconnect attempts to re-establish the connection.
func (d *ReconnectDelegate[T]) Reconnect() (err error) {
	var transport dlgt.ClientTransport[T]
	for {
		transport, err = d.initTransport()
		if err != nil {
			if err == dcln.ErrClosed || err == ErrServerInfoMismatch {
				return
			}
			continue
		}
		break
	}
	return d.setTransport(transport)
}

// Close closes the current transport and stops reconnection attempts.
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
		return nil, dcln.ErrClosed
	}
	transport, err = d.factory.New()
	if err != nil {
		return nil, err
	}
	if d.closed() {
		_ = transport.Close()
		return nil, dcln.ErrClosed
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
		return dcln.ErrClosed
	}
	d.transport.Store(transport)
	return nil
}

func (d *ReconnectDelegate[T]) closed() bool {
	return atomic.LoadUint32(&d.closedFlag) == 1
}
