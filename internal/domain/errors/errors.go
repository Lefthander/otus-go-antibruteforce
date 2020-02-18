package errors

import (
	"errors"
)

var (
	// ErrTokenBucketInvalidFillRate appears when the rate ==0 is used. To avoid the panic of NewTimeTicker()
	ErrTokenBucketInvalidFillRate = errors.New("invalid rate, zero value is not allowed")

	// ErrTokenBucketNotFound appears when we cannot find the bucket in the store
	ErrTokenBucketNotFound = errors.New("unable to find the specided bucket")
	// ErrTokenBucketAlreadyExists appears when we try to create a new bucket, but the bucket with such id is alredy exists
	ErrTokenBucketAlreadyExists = errors.New("such bucket already exists in the store")
)
