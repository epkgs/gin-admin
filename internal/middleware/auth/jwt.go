package auth

import (
	"gin-admin/internal/errorx"
	"gin-admin/internal/types"
	"gin-admin/locales"
	"gin-admin/pkg/helper"
	"gin-admin/pkg/jwtx"
	"strings"

	"github.com/gin-gonic/gin"
)

func checkJWT(c *gin.Context, app types.AppContext) error {
	ctx := c.Request.Context()

	var token string
	{

		auth := c.GetHeader("Authorization")
		prefix := "Bearer "

		if auth != "" && strings.HasPrefix(auth, prefix) {
			token = auth[len(prefix):]
		} else {
			token = auth
		}

		if token == "" {
			token = c.Query("token")
		}
	}

	if token == "" {
		return errorx.ErrUnauthorized.WithMsg(locales.User.Str("Invalid token"))
	}

	ctx = helper.WithToken(ctx, token)

	claims, err := app.Jwt().ParseToken(ctx, token)
	if err != nil {
		if err == jwtx.ErrInvalidToken {
			return errorx.ErrUnauthorized.WithMsg(locales.User.Str("Invalid token"))
		}
		return errorx.ErrInternalServerError.Wrap(err)
	}

	userID := claims.UserID

	ctx = helper.WithUserID(ctx, userID)
	c.Request = c.Request.WithContext(ctx)
	return nil
}
