package jwtx

import (
	"context"
	"time"
)

// storer is the interface that storage the token.
type storer interface {
	Set(ctx context.Context, tokenID string, expiration time.Duration) error
	Delete(ctx context.Context, tokenID string) error
	Check(ctx context.Context, tokenID string) (bool, error)
	Close() error
}

type storeOptions struct {
	CacheNS string // default "jwt"
}

type Cacher interface {
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

func newStore(cache Cacher) storer {
	s := &storeImpl{
		cacher: cache,
		opts: &storeOptions{
			CacheNS: "jwt",
		},
	}

	return s
}

type storeImpl struct {
	opts   *storeOptions
	cacher Cacher
}

func (a *storeImpl) Set(ctx context.Context, tokenID string, expiration time.Duration) error {
	return a.cacher.Set(ctx, a.opts.CacheNS+":"+tokenID, nil, expiration)
}

func (a *storeImpl) Delete(ctx context.Context, tokenID string) error {
	return a.cacher.Delete(ctx, a.opts.CacheNS+":"+tokenID)
}

func (a *storeImpl) Check(ctx context.Context, tokenID string) (bool, error) {
	return a.cacher.Exists(ctx, a.opts.CacheNS+":"+tokenID)
}

func (a *storeImpl) Close() error {
	return a.cacher.Close()
}
