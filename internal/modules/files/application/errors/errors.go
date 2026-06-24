package errors

import "errors"

var ErrInvalidInput = errors.New("invalid input")
var ErrUnauthenticated = errors.New("unauthenticated")
var ErrNotFound = errors.New("not found")
var ErrFailedPrecondition = errors.New("failed precondition")
