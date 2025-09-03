package modules

import (
	"context"

	"gin-admin/pkg/jwtx"
	"gin-admin/types"

	"github.com/golang-jwt/jwt/v5"
)

func NewJWT(ctx context.Context, app types.AppContext) (jwtx.Auther, error) {
	cfg := app.Config().Jwt
	var opts []jwtx.Option
	opts = append(opts, jwtx.WithExpired(cfg.Expired))
	opts = append(opts, jwtx.WithSigningKey(cfg.SigningKey))

	var method jwt.SigningMethod
	switch cfg.SigningMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	default:
		method = jwt.SigningMethodHS512
	}
	opts = append(opts, jwtx.WithSigningMethod(method))

	auth := jwtx.New(app.Cacher(), opts...)

	app.AddCleaner(ctx, func() {
		_ = auth.Release(ctx)
	})

	return auth, nil
}
