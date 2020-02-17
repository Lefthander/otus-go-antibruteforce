package interfaces

import (
	"context"
)

// Bucket is interface which implements the bucket functionality
type Bucket interface {
	Allow(ctx context.Context) bool
    Reset(ctx context.Context)
}