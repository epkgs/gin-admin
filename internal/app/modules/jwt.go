package modules

import (
	"context"
	"time"

	"gin-admin/internal/types"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/jwtx"

	"github.com/golang-jwt/jwt/v5"
)

func InitJWT(ctx context.Context, app types.AppContext) (jwtx.Auther, error) {
	cfg := app.Config().Middleware.Auth
	var opts []jwtx.Option
	opts = append(opts, jwtx.SetExpired(cfg.Expired))
	opts = append(opts, jwtx.SetSigningKey(cfg.SigningKey, cfg.OldSigningKey))
	opts = append(opts, jwtx.SetRefreshKey(cfg.RefreshKey))

	var method jwt.SigningMethod
	switch cfg.SigningMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	default:
		method = jwt.SigningMethodHS512
	}
	opts = append(opts, jwtx.SetSigningMethod(method))

	var cache cachex.Cacher
	switch cfg.Store.Type {
	case "redis":
		cache = cachex.NewRedisCache(cachex.RedisConfig{
			Addr:     cfg.Store.Redis.Addr,
			DB:       cfg.Store.Redis.DB,
			Username: cfg.Store.Redis.Username,
			Password: cfg.Store.Redis.Password,
		}, cachex.WithDelimiter(cfg.Store.Delimiter))
	case "badger":
		cache = cachex.NewBadgerCache(cachex.BadgerConfig{
			Path: cfg.Store.Badger.Path,
		}, cachex.WithDelimiter(cfg.Store.Delimiter))
	default:
		cache = cachex.NewMemoryCache(cachex.MemoryConfig{
			CleanupInterval: time.Second * time.Duration(cfg.Store.Memory.CleanupInterval),
		}, cachex.WithDelimiter(cfg.Store.Delimiter))
	}

	auth := jwtx.New(jwtx.NewStoreWithCache(cache), opts...)

	app.AddCleaner(ctx, func() {
		_ = auth.Release(ctx)
	})

	return auth, nil
}
