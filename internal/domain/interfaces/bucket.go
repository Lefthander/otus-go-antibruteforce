package interfaces

// Bucket is interface which implements the bucket functionality
type Bucket interface {
	Allow() bool
	Reset()
	GetShutDownChannel() chan bool
}
