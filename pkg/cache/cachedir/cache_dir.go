package cachedir

import (
	"context"
	"fmt"
	"time"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/chatwebhook/pkg/cache"
	"github.com/xaionaro-go/fcache"
)

type CacheDir struct {
	Backend fcache.Cache[fcache.String, fcache.Bytes32]
}

var _ cache.Cache = (*CacheDir)(nil)

func New(
	path string,
) (*CacheDir, error) {
	backend, err := fcache.Builder[fcache.String](path, 16*fcache.MiB).Build()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize cache dir backend: %w", err)
	}
	return &CacheDir{
		Backend: backend,
	}, nil
}

func (c *CacheDir) Get(
	ctx context.Context,
	key string,
) (value []byte, updatedAt time.Time) {
	value, info, err := c.Backend.Get(fcache.String(key))
	if err != nil {
		logger.Errorf(ctx, "unable to get cache item %q: %v", key, err)
		return nil, time.Time{}
	}
	if value == nil {
		return nil, time.Time{}
	}
	return value, info.Mtime
}

func (c *CacheDir) Set(
	ctx context.Context,
	key string,
	value []byte,
) {
	if _, err := c.Backend.Put(fcache.String(key), value, 0); err != nil {
		logger.Errorf(ctx, "unable to put cache item %q: %v", key, err)
	}
}
