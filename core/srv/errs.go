package srv

import (
	"errors"
	"fmt"
)

const errorPrefix = "cmdstream server: "

// ErrNoWorkers happens when the server is configured with 0 workers.
var ErrNoWorkers = errors.New("not positive Conf.WorkersCount")

// ErrNotServing happens when the server is not serving and is closed or
// shutdown.
var ErrNotServing = errors.New("server is not serving")

// ErrShutdown happens when the server is shutdown while serving.
var ErrShutdown = errors.New("shutdown")

// ErrClosed happens when the server is closed while serving.
var ErrClosed = errors.New("closed")

// NewServerError wraps an error with the server-specific prefix.
func NewServerError(err error) error {
	return fmt.Errorf(errorPrefix+"%w", err)
}
