package modules

import (
	"context"
	"time"

	"gin-admin/pkg/cachex"
	"gin-admin/types"
)

// It returns a cachex.Cacher instance, a function to close the cache, and an error
func NewCacher(ctx context.Context, app types.AppContext) (cachex.Cacher, error) {
	cfg := app.Config().Cache

	var cache cachex.Cacher
	var err error
	switch cfg.Type {
	case "redis":
		cache, err = cachex.NewRedisCache(&cachex.RedisConfig{
			Addr:     cfg.Redis.Addr,
			DB:       cfg.Redis.DB,
			Username: cfg.Redis.Username,
			Password: cfg.Redis.Password,
		})
		if err != nil {
			return nil, err
		}
	case "badger":
		cache, err = cachex.NewBadgerCache(&cachex.BadgerConfig{
			Path: cfg.Badger.Path,
		})
		if err != nil {
			return nil, err
		}
	default:
		cache = cachex.NewMemoryCache(&cachex.MemoryConfig{
			CleanupInterval: time.Second * time.Duration(cfg.Memory.CleanupInterval),
		})
	}

	app.AddCleaner(ctx, func() {
		_ = cache.Close()
	})

	return cache, nil
}
