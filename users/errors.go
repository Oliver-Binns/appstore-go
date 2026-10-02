package users

import "errors"

// ErrNotFound is returned when no user or pending invitation exists for a given ID or email.
var ErrNotFound = errors.New("user not found")
