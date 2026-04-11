package cln

import (
	"errors"
	"fmt"
)

const errorPrefix = "cmdstream client: "

// ErrClosed happens when the Client is closed while connected to the server.
var ErrClosed = errors.New("closed")

// NewClientError wraps an error with the client-specific prefix.
func NewClientError(cause error) error {
	return fmt.Errorf(errorPrefix+"%w", cause)
}
