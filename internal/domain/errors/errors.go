package errors

import (
	"errors"
)

var (
	// ErrTokenBucketInvalidFillRate appears when the rate ==0 is used. To avoid the panic of NewTimeTicker()
	ErrTokenBucketInvalidFillRate = errors.New("invalid rate, zero value is not allowed")
)
