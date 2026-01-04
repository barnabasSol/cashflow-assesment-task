package common

import "errors"

var ErrInternal = errors.New("internal error")

var ErrNotFound = errors.New("not found")

var ErrBadRequest = errors.New("invalid input")
