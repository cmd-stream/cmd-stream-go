package group

import "fmt"

const errorPrefix = "cmdstream group: "

// NewGroupError creates a client group error.
func NewGroupError(cause error) error {
	return fmt.Errorf(errorPrefix+"%w", cause)
}
