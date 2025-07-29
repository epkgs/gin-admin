package modules

import (
	"context"
	"time"

	"gin-admin/internal/types"
	"gin-admin/pkg/cachex"
)

// It returns a cachex.Cacher instance, a function to close the cache, and an error
func InitCacher(ctx context.Context, app types.AppContext) (cachex.Cacher, error) {
	cfg := app.Config().Cache

	var cache cachex.Cacher
	switch cfg.Type {
	case "redis":
		cache = cachex.NewRedisCache(cachex.RedisConfig{
			Addr:     cfg.Redis.Addr,
			DB:       cfg.Redis.DB,
			Username: cfg.Redis.Username,
			Password: cfg.Redis.Password,
		}, cachex.WithDelimiter(cfg.Delimiter))
	case "badger":
		cache = cachex.NewBadgerCache(cachex.BadgerConfig{
			Path: cfg.Badger.Path,
		}, cachex.WithDelimiter(cfg.Delimiter))
	default:
		cache = cachex.NewMemoryCache(cachex.MemoryConfig{
			CleanupInterval: time.Second * time.Duration(cfg.Memory.CleanupInterval),
		}, cachex.WithDelimiter(cfg.Delimiter))
	}

	app.AddCleaner(ctx, func() {
		_ = cache.Close(ctx)
	})

	return cache, nil
}
