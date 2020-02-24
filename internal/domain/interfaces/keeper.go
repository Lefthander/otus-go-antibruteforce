package interfaces

import (
	"context"
)

// BucketKeeper implements the intreface to store the buckets
type BucketKeeper interface {
	CreateBucket(ctx context.Context, id string, bucket Bucket) error
	DeleteBucket(ctx context.Context, id string) error
	GetBucket(ctx context.Context, id string) (Bucket, error)
}
