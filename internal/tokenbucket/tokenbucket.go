package tokenbucket

import (
	"log"
	"sync"
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
	// ctx            context.Context
	capacity       uint32        // Amount of tokens in the bucket, to survyve the bursts for our case = 1
	currentAmount  uint32        // Current amount of tokens in the bucket
	rate           time.Duration // rate to generate new tokens in the bucket
	lastAccessTime time.Time     // time of lastest request to the bucket
	lifeTime       time.Duration // how long bucket will live without any requests
	ticker         *time.Ticker  // ticker to generate enents with desired rate
	shutDown       chan bool     // shutDown channel to close the bucket by inactivity
	mx             sync.RWMutex  // mutex to protect concurency access
}

// NewTokenBucket creates a new instance of TokenBucket
func NewTokenBucket(capacity uint32, rate time.Duration, lifetime time.Duration) (*TokenBucket, error) {
	// To protect for inaccurate function usage. In case of rate == 0 returns - error.
	if rate == 0 {
		return nil, errors.ErrTokenBucketInvalidFillRate
	}

	tb := &TokenBucket{
		//ctx:            ctx,
		capacity:       capacity,
		currentAmount:  capacity,
		rate:           rate,
		lifeTime:       lifetime,
		lastAccessTime: time.Now(),
		shutDown:       make(chan bool, 1),
		mx:             sync.RWMutex{},
	}
	// Create a go routine with the time.Ticker to fill the bucket with desired rate.
	go func(tb *TokenBucket) {
		tb.ticker = time.NewTicker(tb.rate)

		defer tb.ticker.Stop()

		for {
			select {
			case <-tb.ticker.C:
				tb.mx.RLock()
				if time.Now().Sub(tb.lastAccessTime) > tb.lifeTime {
					// Inactive timeout exeded - initiate the closure procedure of the bucket.
					log.Println("Idle timeout exiting...")
					tb.shutDown <- true
					tb.mx.RUnlock()
					return
				}
				tb.mx.RUnlock()

				if tb.currentAmount == tb.capacity { //Token Bucken is full all next tokens will be discarded
					continue
				}
				tb.mx.Lock()
				tb.currentAmount = tb.currentAmount + 1 // Add one token to the bucket
				tb.mx.Unlock()
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
	tb.mx.Lock()
	defer tb.mx.Unlock()

	tb.lastAccessTime = time.Now()

	if tb.currentAmount > 0 { // Bucket is not empty
		tb.currentAmount = tb.currentAmount - 1 // decrease the number of tokens in the bucket
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

// GetShutDownChannel returns the shutDown channel for provide it to the watchdog goroutine.
func (tb *TokenBucket) GetShutDownChannel() chan bool {
	tb.mx.Lock()
	defer tb.mx.Unlock()
	return tb.shutDown
}

// Amount returns the current amount of tokens in the bucket
func (tb *TokenBucket) Amount() uint32 {
	tb.mx.Lock()
	defer tb.mx.Unlock()

	return tb.currentAmount
}
