package middleware

import (
	"gin-admin/internal/errorx"
	"gin-admin/pkg/response"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type CasbinConfig struct {
	Skipper     func(c *gin.Context) bool
	GetEnforcer func(c *gin.Context) *casbin.Enforcer
	GetSubjects func(c *gin.Context) []string
}

func CasbinWithConfig(config CasbinConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Skipper != nil && config.Skipper(c) {
			c.Next()
			return
		}

		ctx := c.Request.Context()

		enforcer := config.GetEnforcer(c)
		if enforcer == nil {
			response.Error(c, errorx.ErrAccessDenied.New(ctx))
			return
		}

		for _, sub := range config.GetSubjects(c) {
			if b, err := enforcer.Enforce(sub, c.Request.URL.Path, c.Request.Method); err != nil {
				response.Error(c, err)
				return
			} else if b {
				c.Next()
				return
			}
		}
		response.Error(c, errorx.ErrAccessDenied.New(ctx))
	}
}
