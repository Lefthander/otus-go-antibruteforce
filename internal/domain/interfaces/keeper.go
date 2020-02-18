package interfaces

import (
	"context"

	"github.com/google/uuid"
)

// BucketKeeper implements the intreface to store the buckets
type BucketKeeper interface {
	CreateBucket(ctx context.Context, id uuid.UUID, rate uint32, bucket Bucket) error
	DeleteBucket(ctx context.Context, id uuid.UUID) error
	GetBucket(ctx context.Context, id uuid.UUID) (Bucket, error)
}
