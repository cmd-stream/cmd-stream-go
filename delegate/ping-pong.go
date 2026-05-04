package delegate

import (
	"context"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
)

// PingCmd represents a keepalive ping command.
type PingCmd[T any] struct{}

// Exec sends a pong result back to the server.
func (c PingCmd[T]) Exec(_ context.Context, _ T, proxy core.Proxy) (err error) {
	_, err = proxy.SendWithDeadline(time.Time{}, PongResult{})
	return
}

// PongResult represents a keepalive pong result.
type PongResult struct{}

// LastOne indicates that this is the final result for the ping command.
func (r PongResult) LastOne() bool {
	return true
}
