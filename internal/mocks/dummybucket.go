package mocks

// DummyBucket mocks Bucket for testing purposes
type DummyBucket struct {
	ID string
}

// Allow mocks method Allow for Bucket
func (d *DummyBucket) Allow() bool {
	return true
}

// Reset Mocks method Reset for Bucket
func (d *DummyBucket) Reset() {

}

// GetShutDownChannel Mocks method GetShutDownChannel for Bucket
func (d *DummyBucket) GetShutDownChannel() chan bool {
	return make(chan bool)
}
