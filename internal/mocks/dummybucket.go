package mocks

import "context"

// DummyBucket mocks Bucket for testing purposes
type DummyBucket struct {
	ID string
}

// Allow mocks method Allow for Bucket
func (d *DummyBucket) Allow(ctx context.Context) bool {
	return true
}

// Reset Mocks method Reset for Bucket
func (d *DummyBucket) Reset(ctx context.Context) {

}
