package adapters

import (
	"context"
	"sync"

	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/errors"
	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/interfaces"
	"github.com/google/uuid"
)

// TokenBucketMemory is a struct to store buckets
type TokenBucketMemory struct {
	TokenBuckets map[uuid.UUID]interfaces.Bucket
	mx           sync.RWMutex
}

// NewTokenBucketMemory create an instance of bucket storage
func NewTokenBucketMemory() *TokenBucketMemory {
	return &TokenBucketMemory{
		TokenBuckets: map[uuid.UUID]interfaces.Bucket{},
		mx:           sync.RWMutex{},
	}
}

// GetBucket returns bucket from the store if preset or error in case of absent
func (tb *TokenBucketMemory) GetBucket(ctx context.Context, id uuid.UUID) (interfaces.Bucket, error) {
	tb.mx.RLock()
	defer tb.mx.RUnlock()

	if b, ok := tb.TokenBuckets[id]; ok {
		return b, nil
	}

	return nil, errors.ErrTokenBucketNotFound
}

// CreateBucket in memory the instace of Bucket
func (tb *TokenBucketMemory) CreateBucket(ctx context.Context, id uuid.UUID, rate uint32,
	bucket interfaces.Bucket) error {
	_, err := tb.GetBucket(ctx, id)

	if err == nil {
		return errors.ErrTokenBucketAlreadyExists
	}

	tb.mx.Lock()
	tb.TokenBuckets[id] = bucket
	tb.mx.Unlock()

	return nil
}

// DeleteBucket removes the specified bucket from the storage
func (tb *TokenBucketMemory) DeleteBucket(ctx context.Context, id uuid.UUID) error {

	_, err := tb.GetBucket(ctx, id)

	if err != nil {
		return err
	}

	tb.mx.Lock()
	delete(tb.TokenBuckets, id)
	tb.mx.Unlock()

	return nil
}
