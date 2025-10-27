package cmdstream

import "fmt"

const errorPrefix = "cmdstream: "

// NewMakeClientsError creates an error indicating a failure to create the
// specified number of clients.
func NewMakeClientsError(count int, cause error) error {
	return fmt.Errorf("failed to make %v clients, cause: %w", count,
		cause)
}

// NewCmdStreamError creates a cmd-stream-go module error.
func NewCmdStreamError(cause error) error {
	return fmt.Errorf(errorPrefix+"%w", cause)
}
