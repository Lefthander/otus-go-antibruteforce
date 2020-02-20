package tokenbucket

import (
	"context"
	"testing"
	"time"
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
	WAITTIME3 = 3
	WAITTIME1 = 1
	WAITTIME5 = 5
	WAITTIME0 = 0
)

func TestStillBucket(t *testing.T) {
	capacity := 5

	fillRate := WAITTIME1 * time.Second

	ctx, done := context.WithCancel(context.Background())

	defer done()

	tb, err := NewTokenBucket(ctx, uint32(capacity), fillRate)

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

	ctx, done := context.WithCancel(context.Background())

	tb, err := NewTokenBucket(ctx, uint32(capacity), fillRate)

	if err != nil {
		t.Error("Error", err)
	}

	// Gets tokens one by one in number of capacity
	for i := 0; i < capacity-1; i++ {
		allow = tb.Allow(ctx)
		if !allow {
			t.Error("Full TokenBucket must allow all requests regading its capacity")
		}
	}
	done()
}

// 3. Create an empty bucket and check that all following requests will be disallowed.
func TestEmptyBucket(t *testing.T) {
	capacity := 3

	fillRate := WAITTIME5 * time.Second

	ctx, done := context.WithCancel(context.Background())

	defer done()

	tb, err := NewTokenBucket(ctx, uint32(capacity), fillRate)

	if err != nil {
		t.Error("Error", err)
	}

	// Make sure that bucket is empty
	for i := 0; i < capacity-1; i++ {
		tb.Allow(ctx)
	}

	// While bucket is empty the response must be false
	for tb.Amount() == 0 {
		if tb.Allow(ctx) {
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

	ctx, done := context.WithCancel(context.Background())

	defer done()

	tb, err := NewTokenBucket(ctx, uint32(capacity), fillRate)

	if err != nil {
		t.Error("Error", err)
	}

	for i := 0; i < 50; i++ {
		tb.Allow(ctx)
	}

	if tb.Amount() != tb.Capacity() {
		tb.Reset(ctx)

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

	ctx, done := context.WithCancel(context.Background())

	defer done()

	tb, err := NewTokenBucket(ctx, uint32(capacity), fillRate)

	if err != nil {
		t.Error("Error", err)
	}

	for i := 0; i < 50; i++ {
		if tb.Allow(ctx) {
			t.Error("Zero capacity bucket must allways response with false")
		}
	}
	// Check the consistency of Bucket

	if !(tb.Capacity() == tb.Amount()) && tb.Amount() != 0 {
		t.Errorf("Bucket is inconsistent capacity=%d, currentAmount=%d", tb.Capacity(), tb.Amount())
	}
}

// Test Bucket with Zero fillRate, bucket must response Allow = true for all request belongs
// his capacity and after that respond Allow = false for forever
func TestZeroFillRateBucket(t *testing.T) {
	capacity := 10

	fillRate := WAITTIME0 * time.Millisecond

	ctx, done := context.WithCancel(context.Background())

	defer done()

	_, err := NewTokenBucket(ctx, uint32(capacity), fillRate)

	if err == nil {
		t.Error("Error cannot create a TimeTicker with rate ==0")
	}
}
