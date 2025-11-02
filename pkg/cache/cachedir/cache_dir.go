package cachedir

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/icza/kvcache"
	"github.com/xaionaro-go/chatwebhook/pkg/cache"
)

type CacheDir struct {
	Backend kvcache.Cache
}

var _ cache.Cache = (*CacheDir)(nil)

func New(
	path string,
) (*CacheDir, error) {
	backend, err := kvcache.New(path, "")
	if err != nil {
		return nil, fmt.Errorf("unable to initialize cache dir backend: %w", err)
	}
	return &CacheDir{
		Backend: backend,
	}, nil
}

type item struct {
	UpdatedAt time.Time
	Value     []byte
}

func (c *CacheDir) Get(
	ctx context.Context,
	key string,
) (value []byte, updatedAt time.Time) {
	raw, err := c.Backend.Get(key)
	if err != nil {
		logger.Errorf(ctx, "unable to get cache item %q: %v", key, err)
		return nil, time.Time{}
	}
	if raw == nil {
		return nil, time.Time{}
	}

	var i item
	err = json.Unmarshal(raw, &i)
	if err != nil {
		logger.Errorf(ctx, "unable to unmarshal cache item %q: '%s': %v", key, raw, err)
		return nil, time.Time{}
	}

	return i.Value, i.UpdatedAt
}

func (c *CacheDir) Set(
	ctx context.Context,
	key string,
	value []byte,
) {
	i := item{
		UpdatedAt: time.Now(),
		Value:     value,
	}

	raw, err := json.Marshal(&i)
	if err != nil {
		logger.Errorf(ctx, "unable to marshal cache item %q: %v", key, err)
		return
	}

	if err := c.Backend.Put(key, raw); err != nil {
		logger.Errorf(ctx, "unable to put cache item %q: %v", key, err)
	}
}
