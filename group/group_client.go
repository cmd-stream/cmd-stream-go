package group

import (
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
)

// ClientID identifies a specific client within a Group.
type ClientID int

// Client represents a client used by the Group for sending commands and
// receiving results.
type Client[T any] interface {
	Send(cmd core.Cmd[T], results chan<- core.AsyncResult) (seq core.Seq, n int,
		err error)
	SendWithDeadline(deadline time.Time, cmd core.Cmd[T],
		results chan<- core.AsyncResult) (seq core.Seq, n int, err error)
	Has(seq core.Seq) bool
	Forget(seq core.Seq)
	Error() (err error)
	Close() (err error)
	Done() <-chan struct{}
}
