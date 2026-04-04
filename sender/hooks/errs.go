package hooks

import "errors"

// ErrNotAllowed indicates that sending the Command is not allowed at this
// time.
var ErrNotAllowed = errors.New("not allowed")
