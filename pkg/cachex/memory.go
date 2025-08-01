package cachex

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

type MemoryConfig struct {
	CleanupInterval time.Duration
}

func NewMemoryCache(cfg MemoryConfig, opts ...Option) Cacher {
	defaultOpts := &options{
		Delimiter: defaultDelimiter,
	}

	for _, o := range opts {
		o(defaultOpts)
	}

	return &memCache{
		opts:  defaultOpts,
		cache: cache.New(0, cfg.CleanupInterval),
	}
}

type memCache struct {
	opts  *options
	cache *cache.Cache
}

func (a *memCache) getKey(ns, key string) string {
	return fmt.Sprintf("%s%s%s", ns, a.opts.Delimiter, key)
}

func (a *memCache) Set(ctx context.Context, ns, key, value string, expiration ...time.Duration) error {
	var exp time.Duration
	if len(expiration) > 0 {
		exp = expiration[0]
	}

	a.cache.Set(a.getKey(ns, key), value, exp)
	return nil
}

func (a *memCache) Get(ctx context.Context, ns, key string) (string, error) {
	val, ok := a.cache.Get(a.getKey(ns, key))
	if !ok {
		return "", ErrNotFound
	}
	return val.(string), nil
}

func (a *memCache) Exists(ctx context.Context, ns, key string) (bool, error) {
	_, ok := a.cache.Get(a.getKey(ns, key))
	return ok, nil
}

func (a *memCache) Delete(ctx context.Context, ns, key string) error {
	a.cache.Delete(a.getKey(ns, key))
	return nil
}

func (a *memCache) GetAndDelete(ctx context.Context, ns, key string) (string, error) {
	value, err := a.Get(ctx, ns, key)
	if err != nil {
		return "", err
	}

	a.cache.Delete(a.getKey(ns, key))
	return value, nil
}

func (a *memCache) Iterator(ctx context.Context, ns string, fn func(ctx context.Context, key, value string) bool) error {
	for k, v := range a.cache.Items() {
		if strings.HasPrefix(k, a.getKey(ns, "")) {
			if !fn(ctx, strings.TrimPrefix(k, a.getKey(ns, "")), v.Object.(string)) {
				break
			}
		}
	}
	return nil
}

func (a *memCache) Close(ctx context.Context) error {
	a.cache.Flush()
	return nil
}
