package apis

import (
	"context"

	"gin-admin/internal/dtos"
	"gin-admin/internal/services"
	"gin-admin/internal/types"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// User management for SYS
type User struct {
	UserSVC *services.User
}

func NewUser(app types.AppContext) *User {
	handler := &User{
		UserSVC: services.NewUser(app),
	}

	app.Routers().GroupAPI("/api/v1/users", func(ctx context.Context, g *gin.RouterGroup, e *gin.Engine) error {
		g.GET("", handler.Query)
		g.GET(":id", handler.Get)
		g.POST("", handler.Create)
		g.PUT(":id", handler.Update)
		g.DELETE(":id", handler.Delete)
		g.PATCH(":id/reset-pwd", handler.ResetPassword)

		return nil
	})

	return handler
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Query user list
// @Param request query dtos.UserListReq false "query params"
// @Success 200 {object} dtos.ResultList[models.User]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/users [get]
func (a *User) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params dtos.UserListReq
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Error(c, err)
		return
	}

	result, err := a.UserSVC.List(ctx, params)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.List(c, result.Items, &result.Pager)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Get user record by ID
// @Param id path string true "unique id"
// @Success 200 {object} dtos.Result[models.User]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/users/{id} [get]
func (a *User) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.UserSVC.Get(ctx, c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OkData(c, item)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Create user record
// @Param body body dtos.UserCreateReq true "Request body"
// @Success 200 {object} dtos.Result[models.User]
// @Failure 400 {object} dtos.Result[any]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/users [post]
func (a *User) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(dtos.UserCreateReq)
	if err := c.ShouldBindJSON(item); err != nil {
		response.Error(c, err)
		return
	}

	result, err := a.UserSVC.Create(ctx, item)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OkData(c, result)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Update user record by ID
// @Param id path string true "unique id"
// @Param body body dtos.UserUpdateReq true "Request body"
// @Success 200 {object} dtos.Result[any]
// @Failure 400 {object} dtos.Result[any]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/users/{id} [put]
func (a *User) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(dtos.UserUpdateReq)
	if err := c.ShouldBindJSON(item); err != nil {
		response.Error(c, err)
		return
	}

	err := a.UserSVC.Update(ctx, c.Param("id"), item)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Delete user record by ID
// @Param id path string true "unique id"
// @Success 200 {object} dtos.Result[any]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/users/{id} [delete]
func (a *User) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserSVC.Delete(ctx, c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Reset user password by ID
// @Param id path string true "unique id"
// @Success 200 {object} dtos.Result[any]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/users/{id}/reset-pwd [patch]
func (a *User) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserSVC.ResetPassword(ctx, c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c)
}
