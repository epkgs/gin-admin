package auth

import (
	"gin-admin/internal/types"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

func New(app types.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := checkJWT(c, app); err != nil {
			response.Error(c, err)
			return
		}

		if err := checkCasbin(c, app); err != nil {
			response.Error(c, err)
			return
		}
		c.Next()
	}
}
