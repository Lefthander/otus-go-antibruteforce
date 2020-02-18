package adapters

import (
	"context"
	"testing"

	"github.com/Lefthander/otus-go-antibruteforce/internal/domain/errors"
	"github.com/Lefthander/otus-go-antibruteforce/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTockenBucketMemory(t *testing.T) {

	tb := NewTokenBucketMemory()

	ctx := context.Background()

	teststruuid := "aaaaaaaa-1111-2222-3333-ffffffffffff"

	testID, err := uuid.Parse(teststruuid)

	if err != nil {
		t.Error("Error create a test uuid", err)
	}

	t.Run("Check empty", func(t *testing.T) {
		assert.Equal(t, 0, len(tb.TokenBuckets))
	})

	t.Run("Verify CreateBucket()", func(t *testing.T) {
		err := tb.CreateBucket(ctx, testID, 1, &mocks.DummyBucket{})
		if err != nil {
			t.Error("Cannot create a bucket in the store", err)
		}
		assert.Equal(t, 1, len(tb.TokenBuckets))
	})

	t.Run("Verify Create an existing bucket", func(t *testing.T) {
		err := tb.CreateBucket(ctx, testID, 1, &mocks.DummyBucket{})
		assert.Equal(t, errors.ErrTokenBucketAlreadyExists, err)
	})

	t.Run("Verify GetBucket()", func(t *testing.T) {
		b, err := tb.GetBucket(ctx, testID)
		if err != nil {
			t.Error("Cannot get a bucket from the store", err)
		}
		assert.NotNil(t, b)
	})

	t.Run("Verify DeleteBucket()", func(t *testing.T) {
		err := tb.DeleteBucket(ctx, testID)
		if err != nil {
			t.Error("Unable to delete bucket from the store", err)
		}
		assert.Equal(t, 0, len(tb.TokenBuckets))
	})

	t.Run("Verify delete the non-exist bucket", func(t *testing.T) {
		err := tb.DeleteBucket(ctx, testID)
		assert.Equal(t, errors.ErrTokenBucketNotFound, err)
	})

}
