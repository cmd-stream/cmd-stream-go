package srv

import "errors"

// ErrEmptyInfo happens when ServerInfo is empty during Delegate creation.
var ErrEmptyInfo = errors.New("empty info")
