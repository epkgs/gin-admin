package jwtx

import (
	"context"
	"time"
)

// Storer is the interface that storage the token.
type Storer interface {
	Set(ctx context.Context, tokenStr string, expiration time.Duration) error
	Delete(ctx context.Context, tokenStr string) error
	Check(ctx context.Context, tokenStr string) (bool, error)
	Close() error
}

type storeOptions struct {
	CacheNS string // default "jwt"
}

type StoreOption func(*storeOptions)

func WithCacheNS(ns string) StoreOption {
	return func(o *storeOptions) {
		o.CacheNS = ns
	}
}

type Cacher interface {
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

func NewStoreWithCache(cache Cacher, opts ...StoreOption) Storer {
	s := &storeImpl{
		c: cache,
		opts: &storeOptions{
			CacheNS: "jwt",
		},
	}
	for _, opt := range opts {
		opt(s.opts)
	}
	return s
}

type storeImpl struct {
	opts *storeOptions
	c    Cacher
}

func (s *storeImpl) Set(ctx context.Context, tokenStr string, expiration time.Duration) error {
	return s.c.Set(ctx, s.opts.CacheNS+":"+tokenStr, nil, expiration)
}

func (s *storeImpl) Delete(ctx context.Context, tokenStr string) error {
	return s.c.Delete(ctx, s.opts.CacheNS+":"+tokenStr)
}

func (s *storeImpl) Check(ctx context.Context, tokenStr string) (bool, error) {
	return s.c.Exists(ctx, s.opts.CacheNS+":"+tokenStr)
}

func (s *storeImpl) Close() error {
	return s.c.Close()
}
