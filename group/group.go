// Package group provides client grouping and load-balancing, managing multiple
// clients as a single logical unit.
package group

import (
	"errors"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
)

// Group represents a group of clients that distributes the load across multiple
// connections using a dispatch strategy.
type Group[T any] struct {
	strategy DispatchStrategy[Client[T]]
	done     chan struct{}
}

// New creates a new Group.
func New[T any](strategy DispatchStrategy[Client[T]]) Group[T] {
	group := Group[T]{strategy, make(chan struct{})}
	go func() {
		sl := strategy.Slice()
		for i := range len(sl) {
			<-sl[i].Done()
		}
		close(group.done)
	}()
	return group
}

// Send sends a Command to the server using the dispatch strategy, delivering
// Results to the provided channel.
//
// Each Command is assigned a unique sequence number per client, starting from 1;
// 0 is reserved for Ping-Pong.
func (s Group[T]) Send(cmd core.Cmd[T], results chan<- core.AsyncResult) (
	seq core.Seq, clientID ClientID, n int, err error,
) {
	client, index := s.strategy.Next()
	clientID = ClientID(index)
	seq, n, err = client.Send(cmd, results)
	if err != nil {
		err = NewGroupError(err)
	}
	return
}

// SendWithDeadline is like Send, but with a deadline.
func (s Group[T]) SendWithDeadline(deadline time.Time, cmd core.Cmd[T],
	results chan<- core.AsyncResult,
) (seq core.Seq, clientID ClientID, n int, err error) {
	client, index := s.strategy.Next()
	clientID = ClientID(index)
	seq, n, err = client.SendWithDeadline(deadline, cmd, results)
	if err != nil {
		err = NewGroupError(err)
	}
	return
}

// Has checks if the Command with the specified sequence number was sent and is
// still waiting for a Result.
func (s Group[T]) Has(seq core.Seq, clientID ClientID) (ok bool) {
	return s.strategy.Slice()[int(clientID)].Has(seq)
}

// Forget makes the Client forget about the Command.
//
// Results of the forgotten Command are passed to UnexpectedResultCallback.
func (s Group[T]) Forget(seq core.Seq, clientID ClientID) {
	s.strategy.Slice()[int(clientID)].Forget(seq)
}

// Size returns the number of clients within this Group.
func (s Group[T]) Size() int {
	return len(s.strategy.Slice())
}

// Done returns a channel that is closed when the Group terminates.
func (s Group[T]) Done() <-chan struct{} {
	return s.done
}

// Error returns joint error of all clients.
func (s Group[T]) Error() (err error) {
	for _, client := range s.strategy.Slice() {
		err = errors.Join(err, client.Error())
	}
	return
}

// Close terminates all clients.
func (s Group[T]) Close() (err error) {
	for _, client := range s.strategy.Slice() {
		err = errors.Join(err, client.Close())
	}
	if err != nil {
		err = NewGroupError(err)
	}
	return
}
