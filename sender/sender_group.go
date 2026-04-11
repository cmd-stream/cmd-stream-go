package sender

import (
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	grp "github.com/cmd-stream/cmd-stream-go/group"
)

// Group represents a group of clients used to send commands and receive
// results.
type Group[T any] interface {
	Send(cmd core.Cmd[T], results chan<- core.AsyncResult) (
		seq core.Seq, clientID grp.ClientID, n int, err error)
	SendWithDeadline(deadline time.Time, cmd core.Cmd[T],
		results chan<- core.AsyncResult,
	) (seq core.Seq, clientID grp.ClientID, n int, err error)
	Has(seq core.Seq, clientID grp.ClientID) (ok bool)
	Forget(seq core.Seq, clientID grp.ClientID)
	Done() <-chan struct{}
	Error() (err error)
	Close() (err error)
}
