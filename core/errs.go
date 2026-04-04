package core

import "fmt"

const errorPrefix = "cmdstream: "

// NewError creates a cmd-stream error.
func NewError(cause error) error {
	return fmt.Errorf(errorPrefix+"%w", cause)
}
