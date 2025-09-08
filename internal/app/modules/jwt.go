package modules

import (
	"context"

	"gin-admin/internal/types"
	"gin-admin/pkg/jwtx"

	"github.com/golang-jwt/jwt/v5"
)

func NewJWT(ctx context.Context, app types.AppContext) (jwtx.Auther, error) {
	cfg := app.Config()
	var method jwt.SigningMethod
	switch cfg.Jwt.SigningMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	default:
		method = jwt.SigningMethodHS512
	}

	auth := jwtx.New(
		jwtx.WithExpired(cfg.Jwt.Expired),
		jwtx.WithSigningKey(cfg.Jwt.SigningKey),
		jwtx.WithCachePath(cfg.GetRuntimePath("jwt")),
		jwtx.WithSigningMethod(method),
	)

	app.AddCleaner(ctx, func() {
		_ = auth.Release(ctx)
	})

	return auth, nil
}
