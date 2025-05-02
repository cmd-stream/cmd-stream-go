package ccln

import (
	"errors"
	"time"

	"github.com/cmd-stream/base-go"
)

// ClientID identifies a specific client within a Group.
type ClientID int

// Client represents a client used by the Group.
type Client[T any] interface {
	Send(cmd base.Cmd[T], results chan<- base.AsyncResult) (seq base.Seq, err error)
	SendWithDeadline(deadline time.Time, cmd base.Cmd[T],
		results chan<- base.AsyncResult) (seq base.Seq, err error)
	Has(seq base.Seq) bool
	Forget(seq base.Seq)
	Err() (err error)
	Close() (err error)
	Done() <-chan struct{}
}

// NewGroup creates a new Group.
func NewGroup[T any](strategy DispatchStrategy[Client[T]]) (grp Group[T]) {
	grp = Group[T]{strategy, make(chan struct{})}
	go func() {
		sl := strategy.Slice()
		for i := range len(sl) {
			<-sl[i].Done()
		}
		close(grp.done)
	}()
	return
}

// Group represents a group of clients which are used to communicate with the
// server. It can be used to increase the communication speed by spreading the
// load across multiple connections.
type Group[T any] struct {
	strategy DispatchStrategy[Client[T]]
	done     chan struct{}
}

// Send transmits a Command to the server.
//
// Results received from the server are delivered to the provided results channel.
// If the channel does not have sufficient capacity, waiting for all results may
// hang.
//
// Each Command is assigned a unique sequence number per client, starting from 1:
//   - The first Command is sent with seq == 1, the second with seq == 2, and so
//     on.
//   - seq == 0 is reserved for the Ping-Pong mechanism, which maintains
//     connection liveness.
//
// Returns the sequence number of the Command, the ClientID it was sent through,
// and any error encountered (non-nil if the Command was not sent successfully).
func (s Group[T]) Send(cmd base.Cmd[T], results chan<- base.AsyncResult) (
	seq base.Seq, clientID ClientID, err error) {
	client, index := s.strategy.Next()
	clientID = ClientID(index)
	seq, err = client.Send(cmd, results)
	return
}

// SendWithDeadline behaves like Send, but ensures the Command is transmitted
// before the specified deadline.
func (s Group[T]) SendWithDeadline(deadline time.Time, cmd base.Cmd[T],
	results chan<- base.AsyncResult,
) (seq base.Seq, clientID ClientID, err error) {
	client, index := s.strategy.Next()
	clientID = ClientID(index)
	seq, err = client.SendWithDeadline(deadline, cmd, results)
	return
}

// Has checks if the Command with the specified sequence number has been sent
// by the client and still waiting for the Result.
func (s Group[T]) Has(seq base.Seq, clientID ClientID) (ok bool) {
	return s.strategy.Slice()[int(clientID)].Has(seq)
}

// Forget makes the Client to forget about the Command which still waiting for
// the result.
//
// After calling Forget, all the results of the corresponding Command will be
// handled with UnexpectedResultCallback.
func (s Group[T]) Forget(seq base.Seq, clientID ClientID) {
	s.strategy.Slice()[int(clientID)].Forget(seq)
}

// Done returns a channel that is closed when the Group terminates.
func (s Group[T]) Done() <-chan struct{} {
	return s.done
}

// Err returns a connection error.
func (s Group[T]) Err() (err error) {
	for _, client := range s.strategy.Slice() {
		err = errors.Join(err, client.Err())
	}
	return
}

// Close terminates underlying clients.
//
// All pending Commands will receive an error (AsyncResult.Error != nil).
func (s Group[T]) Close() (err error) {
	for _, client := range s.strategy.Slice() {
		err = errors.Join(err, client.Close())
	}
	return
}
