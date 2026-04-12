package delegate

import "errors"

// ErrTooLargeServerInfo happens when the received ServerInfo exceeds the
// maximum allowed length.
var ErrTooLargeServerInfo = errors.New("too large server info")
