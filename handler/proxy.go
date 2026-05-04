package handler

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
)

// Proxy implemets the core.Proxy interface.
type Proxy[T any] struct {
	transport dlgt.ServerTransport[T]
	flushFlag *uint32
	mu        *sync.Mutex
	seq       core.Seq
	at        time.Time
}

// ReceivedAt returns the time the Proxy received the command from the network.
func (p Proxy[T]) ReceivedAt() time.Time {
	return p.at
}

// Seq returns the sequence number of the Command.
func (p Proxy[T]) Seq() core.Seq {
	return p.seq
}

// LocalAddr returns the local network address.
func (p Proxy[T]) LocalAddr() net.Addr {
	return p.transport.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (p Proxy[T]) RemoteAddr() net.Addr {
	return p.transport.RemoteAddr()
}

// Send transmits a result to the client.
func (p Proxy[T]) Send(result core.Result) (n int, err error) {
	p.mu.Lock()
	n, err = p.transport.Send(p.seq, result)
	p.mu.Unlock()
	if err != nil {
		return
	}
	return n, p.flush()
}

// SendWithDeadline transmits a result with a specified deadline.
func (p Proxy[T]) SendWithDeadline(deadline time.Time, result core.Result) (
	n int, err error,
) {
	p.mu.Lock()
	err = p.transport.SetSendDeadline(deadline)
	if err != nil {
		p.mu.Unlock()
		return
	}
	n, err = p.transport.Send(p.seq, result)
	if err != nil {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()
	return n, p.flush()
}

func (p Proxy[T]) flush() (err error) {
	if swapped := atomic.CompareAndSwapUint32(p.flushFlag, 0, 1); swapped {
		p.mu.Lock()
		atomic.CompareAndSwapUint32(p.flushFlag, 1, 0)
		err = p.transport.Flush()
		p.mu.Unlock()
	}
	return
}
