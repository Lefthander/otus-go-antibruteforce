package tokenbucket

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/errors"
)

// Algorithm description
// A token is added to the bucket every 1/r seconds
// The bucket can hold at the most Capacity tokens. If tokens arrives when the
// bucket is full, it is discarded.
// when request arrives if at least one token is in the bucket, 1 token removed
// from the bucket, and the request is allowed.
// if 0 tokens in the bucket, no tokens removed and the request is considered as not allowed.

// TokenBucket contains elements to implement the algorithm
type TokenBucket struct {
	ctx           context.Context
	capacity      uint32
	currentAmount uint32
	rate          time.Duration
	ticker        *time.Ticker
	mx            sync.Mutex
}

// NewTokenBucket creates a new instance of TokenBucket
func NewTokenBucket(ctx context.Context, capacity uint32, rate time.Duration) (*TokenBucket, error) {
	// To protect for inaccurate function usage. In case of rate == 0  NewTicker will be created with a very
	// loooong period of tick ~ 290 years.
	// In other terms ticker will not run in closest time.
	if rate == 0 {
		return nil, errors.ErrTokenBucketInvalidFillRate
	}

	tb := &TokenBucket{
		ctx:           ctx,
		capacity:      capacity,
		currentAmount: capacity,
		rate:          rate,
		mx:            sync.Mutex{},
	}
	// Create a go routine with the time.Ticker to fill the bucket with desired rate.
	go func(tb *TokenBucket) {
		tb.ticker = time.NewTicker(tb.rate)

		defer tb.ticker.Stop()

		for {
			select {
			case <-tb.ctx.Done():
				tb.ticker.Stop() // Stops the Ticker
				return
			case <-tb.ticker.C:
				if tb.currentAmount == tb.capacity { //Token Bucken is full all next tokens will be discarded
					continue
				}

				atomic.AddUint32(&tb.currentAmount, 1) // Add one token to the bucket

			default: // Added to avoid blocking when we have looong time period for ticker
				continue
			}
		}
	}(tb)

	return tb, nil
}

// Allow returns true in case we have tokens in the bucket. One authorization event has a weight of one token
// When it is allowed to pass, we decrese the CurrentAmount of tokens by 1
func (tb *TokenBucket) Allow() bool {
	if tb.currentAmount > 0 { // Bucket is not empty
		atomic.AddUint32(&tb.currentAmount, ^uint32(0)) // decrease the number of tokens in the bucket
		return true
	}

	return false
}

// Reset the bucket to the initial state
func (tb *TokenBucket) Reset() {
	tb.mx.Lock()
	tb.currentAmount = tb.capacity
	tb.mx.Unlock()
}

// Capacity returns the capacity of the bucket
func (tb *TokenBucket) Capacity() uint32 {
	return tb.capacity // It's not necessary to use Mutex as it is a readonly parameter
}

// Amount returns the current amount of tokens in the bucket
func (tb *TokenBucket) Amount() uint32 {
	tb.mx.Lock()
	defer tb.mx.Unlock()

	return tb.currentAmount
}
