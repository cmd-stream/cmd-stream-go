package sender

import "errors"

// ErrTimeout is returned when a command is sent but no result is received
// within the expected time.
var ErrTimeout = errors.New("timeout")
