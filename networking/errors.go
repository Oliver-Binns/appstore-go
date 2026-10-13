package networking

import "errors"

// ErrNotFound is returned when the API responds with a 404 for a requested resource.
var ErrNotFound = errors.New("resource not found")
