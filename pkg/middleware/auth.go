package middleware

import (
	"gin-admin/pkg/helper"
	"gin-admin/pkg/logger"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthConfig struct {
	IncludedPathPrefixes []string
	ExcludedPathPrefixes []string
	RootID               string
	Skipper              func(c *gin.Context) bool
	ParseUserID          func(c *gin.Context) (string, error)
}

func AuthWithConfig(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IncludedPathPrefixes(c, config.IncludedPathPrefixes...) ||
			ExcludedPathPrefixes(c, config.ExcludedPathPrefixes...) ||
			(config.Skipper != nil && config.Skipper(c)) {
			c.Next()
			return
		}

		userID, err := config.ParseUserID(c)
		if err != nil {
			response.Error(c, err)
			return
		}

		ctx := helper.WithUserID(c.Request.Context(), userID)
		ctx = logger.WithUserID(ctx, userID)
		if userID == config.RootID {
			ctx = helper.WithIsRootUser(ctx)
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
