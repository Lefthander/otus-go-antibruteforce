package errors

import (
	"errors"
)

var (
	// ErrTokenBucketInvalidFillRate appears when the parameter rate==0 is used.
	// To avoid the panic of NewTimeTicker()
	ErrTokenBucketInvalidFillRate = errors.New("invalid rate, zero value is not allowed")

	// ErrTokenBucketNotFound appears when we cannot find the bucket in the store
	ErrTokenBucketNotFound = errors.New("unable to find the specided bucket")
	// ErrTokenBucketAlreadyExists appears when we try to create a new bucket, but the bucket with such id is alredy exists
	ErrTokenBucketAlreadyExists = errors.New("such bucket already exists in the store")

	// ErrIPFilterMatchedWhiteList apprears to indicate that IP address match the White IP table
	ErrIPFilterMatchedWhiteList = errors.New("ip address is match to the white table")

	// ErrIPFilterMatchedBlackList appears to indicate that IP address match the Black IP table
	ErrIPFilterMatchedBlackList = errors.New("ip address is match to the black table")

	// ErrIPFilterNoMatch appears to indicate that there is no match in any table for provided IP
	ErrIPFilterNoMatch = errors.New("provided ip address does not match either blacl or white table")

	// ErrIPFilterNetworkAlreadyExists appears when we trying to add network which is already exists in the table
	ErrIPFilterNetworkAlreadyExists = errors.New("provided network already exists")

	// ErrIPFilterNetworkNotFound appears when there is no provided network found
	ErrIPFilterNetworkNotFound = errors.New("no such network found")
	// ErrAuthRequestLoginMissed appears when the authentication request has no login defined
	ErrAuthRequestLoginMissed = errors.New("login is missed in the request")
	// ErrAuthRequestPasswordMissed appears when the authentication request has no password defined
	ErrAuthRequestPasswordMissed = errors.New("password is missed in the request")
	// ErrAuthRequestIPMissed appears when the authentication request has no ip address defined
	ErrAuthRequestIPMissed = errors.New("ip address is missed in the request")
	// ErrNoMappingFound appears when there is no match in UUIDTable found
	ErrNoMappingFound = errors.New("no matching found between value and UUID")
)
