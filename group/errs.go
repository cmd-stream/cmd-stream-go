package group

import (
	"errors"
	"fmt"
)

const errorPrefix = "cmdstream group: "

// NewGroupError creates a client group error.
func NewGroupError(cause error) error {
	return fmt.Errorf(errorPrefix+"%w", cause)
}

// ErrInvalidClientsCount is returned when the number of clients is less than or
// equal to 0.
var ErrInvalidClientsCount = errors.New(errorPrefix + "invalid clients count")
