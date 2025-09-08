package auth

import (
	"gin-admin/internal/errorx"
	"gin-admin/internal/service"
	"gin-admin/internal/types"
	"gin-admin/pkg/helper"

	"github.com/gin-gonic/gin"
)

func checkCasbin(c *gin.Context, app types.AppContext) error {
	ctx := c.Request.Context()

	cfg := app.Config()

	userID := helper.GetUserID(ctx)
	if cfg.IsSuper(userID) {
		return nil
	}

	enforcer := app.Casbin().GetEnforcer()
	if enforcer == nil {
		return errorx.ErrForbidden
	}

	userSVC := service.NewUser(app)

	roleIDs, err := userSVC.GetRoleIDsCache(ctx, userID)
	if err != nil {
		return errorx.ErrForbidden.Wrap(err)
	}

	for _, roleID := range roleIDs {
		if ok, err := enforcer.Enforce(roleID, c.Request.URL.Path, c.Request.Method); err != nil {
			return errorx.ErrInternalServerError.Wrap(err)
		} else if ok {
			return nil
		}
	}

	return errorx.ErrForbidden
}
