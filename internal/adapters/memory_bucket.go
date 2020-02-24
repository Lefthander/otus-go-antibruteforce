package adapters

import (
	"context"
	"log"
	"sync"

	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/errors"
	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/interfaces"
)

// TokenBucketMemory is a struct to store buckets
type TokenBucketMemory struct {
	TokenBuckets map[string]interfaces.Bucket
	mx           sync.RWMutex
}

// NewTokenBucketMemory create an instance of bucket storage
func NewTokenBucketMemory() *TokenBucketMemory {
	return &TokenBucketMemory{
		TokenBuckets: map[string]interfaces.Bucket{},
		mx:           sync.RWMutex{},
	}
}

// GetBucket returns bucket from the store if preset or error in case of absent
func (tb *TokenBucketMemory) GetBucket(ctx context.Context, id string) (interfaces.Bucket, error) {
	tb.mx.RLock()
	defer tb.mx.RUnlock()

	if b, ok := tb.TokenBuckets[id]; ok {
		return b, nil
	}

	return nil, errors.ErrTokenBucketNotFound
}

// CreateBucket in memory the instace of Bucket
func (tb *TokenBucketMemory) CreateBucket(ctx context.Context, id string, bucket interfaces.Bucket) error {
	_, err := tb.GetBucket(ctx, id)

	if err == nil {
		return errors.ErrTokenBucketAlreadyExists
	}

	tb.mx.Lock()
	tb.TokenBuckets[id] = bucket
	tb.mx.Unlock()

	// Run a watchdog goroutine which takes care about buckets which are have no request for long time

	go func(ctx context.Context, id string, tb *TokenBucketMemory, shutdown chan bool) {
		<-shutdown
		err := tb.DeleteBucket(ctx, id)
		if err != nil {
			log.Println("Error deleting bucket", id, err)
		}
		log.Println("Bucket deleted due to idle timeout...", id)
	}(ctx, id, tb, bucket.GetShutDownChannel())

	return nil
}

// DeleteBucket removes the specified bucket from the storage
func (tb *TokenBucketMemory) DeleteBucket(ctx context.Context, id string) error {
	_, err := tb.GetBucket(ctx, id)

	if err != nil {
		return err
	}

	tb.mx.Lock()
	delete(tb.TokenBuckets, id)
	tb.mx.Unlock()

	return nil
}
