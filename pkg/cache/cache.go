package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (value []byte, updatedAt time.Time)
	Set(ctx context.Context, key string, value []byte)
}
