package v1

import (
	"gin-admin/internal/dto"
	"gin-admin/internal/service"
	"gin-admin/internal/types"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// User management for SYS
type User struct {
	app     types.AppContext
	UserSVC *service.User
}

func NewUser(app types.AppContext) *User {
	return &User{
		app:     app,
		UserSVC: service.NewUser(app),
	}
}

func (a *User) RegisterRouter(group *gin.RouterGroup, engine *gin.Engine) {

	g := group.Group("users")
	g.Use(
		a.app.Middlewares().Auth(),
		a.app.Middlewares().RoutePermission(),
	)

	g.GET("", a.Query)
	g.GET(":id", a.Get)
	g.POST("", a.Create)
	g.PUT(":id", a.Update)
	g.DELETE(":id", a.Delete)
	g.PATCH(":id/reset-pwd", a.ResetPassword)
}

// @Tags UserAPI
// @Security ApiKeyAuth
// @Summary Query user list
// @Param request query dto.UserListReq false "query params"
// @Success 200 {object} dto.ResultList[model.User]
// @Failure 401 {object} dto.Result[any]
// @Failure 500 {object} dto.Result[any]
// @Router /api/v1/users [get]
func (a *User) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params dto.UserListReq
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
// @Success 200 {object} dto.Result[model.User]
// @Failure 401 {object} dto.Result[any]
// @Failure 500 {object} dto.Result[any]
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
// @Param body body dto.UserCreateReq true "Request body"
// @Success 200 {object} dto.Result[model.User]
// @Failure 400 {object} dto.Result[any]
// @Failure 401 {object} dto.Result[any]
// @Failure 500 {object} dto.Result[any]
// @Router /api/v1/users [post]
func (a *User) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(dto.UserCreateReq)
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
// @Param body body dto.UserUpdateReq true "Request body"
// @Success 200 {object} dto.Result[any]
// @Failure 400 {object} dto.Result[any]
// @Failure 401 {object} dto.Result[any]
// @Failure 500 {object} dto.Result[any]
// @Router /api/v1/users/{id} [put]
func (a *User) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(dto.UserUpdateReq)
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
// @Success 200 {object} dto.Result[any]
// @Failure 401 {object} dto.Result[any]
// @Failure 500 {object} dto.Result[any]
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
// @Success 200 {object} dto.Result[any]
// @Failure 401 {object} dto.Result[any]
// @Failure 500 {object} dto.Result[any]
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
