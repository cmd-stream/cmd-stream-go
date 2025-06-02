package group

import (
	"errors"
	"time"

	"github.com/cmd-stream/core-go"
)

// NewClientGroup creates a new ClientGroup using the provided dispatch strategy.
//
// The returned group monitors all clients and automatically closes its Done()
// channel once all clients have finished.
func NewClientGroup[T any](strategy DispatchStrategy[Client[T]]) (
	group ClientGroup[T]) {
	group = ClientGroup[T]{strategy, make(chan struct{})}
	go func() {
		sl := strategy.Slice()
		for i := range len(sl) {
			<-sl[i].Done()
		}
		close(group.done)
	}()
	return
}

// ClientGroup represents a group of clients used to communicate with the server.
//
// It distributes the load across multiple connections, which can improve
// throughput and resilience. The group selects a client for each operation
// according to the provided dispatch strategy.
type ClientGroup[T any] struct {
	strategy DispatchStrategy[Client[T]]
	done     chan struct{}
}

// Send transmits a Command to the server via one of the clients in the group,
// selected according to the dispatch strategy.
//
// The corresponding Results from the server are delivered to the provided results
// channel. If the channel lacks sufficient capacity, receiving all Results
// may block.
//
// Each Command is assigned a unique sequence number per client, starting from 1:
//   - The first Command sent by a client gets seq == 1, the second seq == 2, etc.
//   - seq == 0 is reserved for the Ping-Pong mechanism that ensures connection
//     liveness.
//
// Returns the sequence number assigned to the Command, the ClientID of the
// client it was sent through, the number of bytes written, and any error
// encountered (non-nil if the Command was not sent successfully).
func (s ClientGroup[T]) Send(cmd core.Cmd[T], results chan<- core.AsyncResult) (
	seq core.Seq, clientID ClientID, n int, err error) {
	client, index := s.strategy.Next()
	clientID = ClientID(index)
	seq, n, err = client.Send(cmd, results)
	return
}

// SendWithDeadline is like Send, but ensures that the Command is transmitted
// before the specified deadline.
func (s ClientGroup[T]) SendWithDeadline(cmd core.Cmd[T],
	results chan<- core.AsyncResult,
	deadline time.Time,
) (seq core.Seq, clientID ClientID, n int, err error) {
	client, index := s.strategy.Next()
	clientID = ClientID(index)
	seq, n, err = client.SendWithDeadline(cmd, results, deadline)
	return
}

// Has checks if the Command with the specified sequence number has been sent
// by the client and still waiting for the Result.
func (s ClientGroup[T]) Has(seq core.Seq, clientID ClientID) (ok bool) {
	return s.strategy.Slice()[int(clientID)].Has(seq)
}

// Forget makes the Client to forget about the Command which still waiting for
// the result.
//
// After calling Forget, all the results of the corresponding Command will be
// handled with UnexpectedResultCallback.
func (s ClientGroup[T]) Forget(seq core.Seq, clientID ClientID) {
	s.strategy.Slice()[int(clientID)].Forget(seq)
}

// Size returns the number of clients within this ClientGroup.
func (s ClientGroup[T]) Size() int {
	return len(s.strategy.Slice())
}

// Done returns a channel that is closed when the ClientGroup terminates.
func (s ClientGroup[T]) Done() <-chan struct{} {
	return s.done
}

// Err returns a connection error.
func (s ClientGroup[T]) Err() (err error) {
	for _, client := range s.strategy.Slice() {
		err = errors.Join(err, client.Err())
	}
	return
}

// Close terminates underlying clients.
//
// All pending Commands will receive an error (AsyncResult.Error != nil).
func (s ClientGroup[T]) Close() (err error) {
	for _, client := range s.strategy.Slice() {
		err = errors.Join(err, client.Close())
	}
	return
}
