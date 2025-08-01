package cachex

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("cache: not found")

// Cacher is the interface that wraps the basic Get, Set, and Delete methods.
type Cacher interface {
	Set(ctx context.Context, ns, key, value string, expiration ...time.Duration) error
	Get(ctx context.Context, ns, key string) (string, error)
	GetAndDelete(ctx context.Context, ns, key string) (string, error)
	Exists(ctx context.Context, ns, key string) (bool, error)
	Delete(ctx context.Context, ns, key string) error
	Iterator(ctx context.Context, ns string, fn func(ctx context.Context, key, value string) bool) error
	Close(ctx context.Context) error
}

var defaultDelimiter = ":"

type options struct {
	Delimiter string
}

type Option func(*options)

func WithDelimiter(delimiter string) Option {
	return func(o *options) {
		o.Delimiter = delimiter
	}
}
