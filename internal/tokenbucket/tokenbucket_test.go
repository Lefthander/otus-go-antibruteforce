package tokenbucket

import (
	"testing"
	"time"

	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/errors"
	"github.com/stretchr/testify/assert"
)

// Testcases:
// 1. Create a bucket and check that amount of tokens does not decreased when
//    there are no requests to the bucket.
// 2. Create a bucket and check that it's possible to get all capacity
//    without any cancelation from the bucket in the burst
// 3. Create an empty bucket and check that all followup requests will be
//    disallowed.
// 4. Create a bucket which allows to have not more then 10 request per minute
// 5. Create a bucket with significant capacity, start taking tokens from it,
//    reset bucket. Ensure that bucket returns to initial state: Capacity == currentAmount

// 1. Create a bucket and check that amount of tokens does not decreased when there are no
//    requests to the bucket.

const (
	WAITTIME3  = 3
	WAITTIME1  = 1
	WAITTIME5  = 5
	WAITTIME0  = 0
	LIFETIME3S = 3 * time.Second
)

func TestStillBucket(t *testing.T) {
	capacity := 1

	fillRate := WAITTIME1 * time.Second

	tb, err := NewTokenBucket(uint32(capacity), fillRate, LIFETIME3S)

	if err != nil {
		t.Error("Error", err)
	}

	time.Sleep(fillRate * WAITTIME3) // Whait for some time

	if tb.Capacity() != tb.Amount() {
		t.Errorf("Bucket without requests must keep the currentAmoutn=%d equals to defined capacity=%d",
			tb.Amount(), tb.Capacity())
	}
}

// 2. Create a bucket and check that it's possible to get all capacity without any cancelation from the
//    bucket in the burst
func TestFullBucket(t *testing.T) {
	var allow bool

	capacity := 5
	fillRate := WAITTIME1 * time.Second // Set a quite long period of refill

	tb, err := NewTokenBucket(uint32(capacity), fillRate, LIFETIME3S)

	if err != nil {
		t.Error("Error", err)
	}

	// Gets tokens one by one in number of capacity
	for i := 0; i < capacity-1; i++ {
		allow = tb.Allow()
		if !allow {
			t.Error("Full TokenBucket must allow all requests regading its capacity")
		}
	}
}

// 3. Create an empty bucket and check that all following requests will be disallowed.
func TestEmptyBucket(t *testing.T) {
	capacity := 3

	fillRate := WAITTIME5 * time.Second

	tb, err := NewTokenBucket(uint32(capacity), fillRate, LIFETIME3S)

	if err != nil {
		t.Error("Error", err)
	}

	// Make sure that bucket is empty
	for i := 0; i < capacity-1; i++ {
		tb.Allow()
	}

	// While bucket is empty the response must be false
	for tb.Amount() == 0 {
		if tb.Allow() {
			t.Error("Bucket must not allow while it empty")
			break
		}
	}
}

// 5. Create a bucket with significant capacity, start taking tokens from it, reset bucket.
// Ensure that bucket returns to initial state: Capacity == currentAmount
func TestResetBucket(t *testing.T) {
	capacity := 100

	fillRate := WAITTIME1 * time.Second

	tb, err := NewTokenBucket(uint32(capacity), fillRate, LIFETIME3S)

	if err != nil {
		t.Error("Error", err)
	}

	for i := 0; i < 50; i++ {
		tb.Allow()
	}

	if tb.Amount() != tb.Capacity() {
		tb.Reset()

		if tb.Amount() != tb.Capacity() {
			t.Errorf("The bucket must has the capacity=%d and amount=%d are equals after the reset!",
				tb.Capacity(), tb.Amount())
		}
	}
}

// Test Bucket with Zero Capacity, bucket allways must response Allow = false,
// and keep consistency capacity & currentAmout == 0
func TestZeroCapacityBucket(t *testing.T) {
	capacity := 0

	fillRate := WAITTIME1 * time.Millisecond

	tb, err := NewTokenBucket(uint32(capacity), fillRate, LIFETIME3S)

	if err != nil {
		t.Error("Error", err)
	}

	for i := 0; i < 50; i++ {
		if tb.Allow() {
			t.Error("Zero capacity bucket must allways response with false")
		}
	}
	// Check the consistency of Bucket

	if !(tb.Capacity() == tb.Amount()) && tb.Amount() != 0 {
		t.Errorf("Bucket is inconsistent capacity=%d, currentAmount=%d", tb.Capacity(), tb.Amount())
	}
}

// Test Bucket with Zero fillRate, bucket should return error.
func TestZeroFillRateBucket(t *testing.T) {
	t.Run("Check bucket creation with fillrate=0", func(t *testing.T) {
		capacity := 10
		fillRate := WAITTIME0 * time.Millisecond

		_, err := NewTokenBucket(uint32(capacity), fillRate, LIFETIME3S)

		assert.Equal(t, errors.ErrTokenBucketInvalidFillRate, err)
	})
}

func TestBucketCloseInactive(t *testing.T) {
	t.Run("Check that bucket will be closed correctly due to inactivity 3s", func(t *testing.T) {
		capacity := 1
		fillRate := WAITTIME1 * time.Second

		tb, err := NewTokenBucket(uint32(capacity), fillRate, LIFETIME3S)
		assert.Equal(t, nil, err)
		assert.Equal(t, true, <-tb.GetShutDownChannel())
	})
}
