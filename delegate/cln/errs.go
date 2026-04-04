package cln

import "errors"

// ErrServerInfoMismatch happens when ServerInfo of the client and server
// does not match.
var ErrServerInfoMismatch = errors.New("server info mismatch")
