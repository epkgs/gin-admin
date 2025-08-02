package v1

import (
	"gin-admin/internal/dtos"
	"gin-admin/internal/errorx"
	"gin-admin/internal/services"
	"gin-admin/internal/types"
	"gin-admin/pkg/helper"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	app     types.AppContext
	AuthSVC *services.Auth
}

func NewAuth(app types.AppContext) *Auth {
	return &Auth{
		app:     app,
		AuthSVC: services.NewAuth(app),
	}
}

func (a *Auth) RegisterRouter(group *gin.RouterGroup, engine *gin.Engine) {

	g := group.Group("auth")

	g.POST("login", a.Login)
	g.POST("refresh-token", a.RefreshToken)
	g.GET("user", a.app.Middlewares().Auth(), a.GetUserInfo)
	g.GET("menus", a.app.Middlewares().Auth(), a.QueryMenus)
	g.PUT("password", a.app.Middlewares().Auth(), a.UpdatePassword)
	g.PUT("user", a.app.Middlewares().Auth(), a.UpdateUser)
	g.POST("logout", a.app.Middlewares().Auth(), a.Logout)
}

// @Tags AuthAPI
// @Summary Login system with username and password
// @Param body body dtos.Login true "Request body"
// @Success 200 {object} dtos.Result[dtos.LoginToken]
// @Failure 400 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/auth/login [post]
func (a *Auth) Login(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(dtos.Login)
	if err := c.ShouldBindJSON(item); err != nil {
		response.Error(c, err)
		return
	}

	data, err := a.AuthSVC.Login(ctx, item.Trim())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OkData(c, data)
}

// @Tags AuthAPI
// @Security ApiKeyAuth
// @Summary Logout system
// @Success 200 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/auth/logout [post]
func (a *Auth) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.AuthSVC.Logout(ctx)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c)
}

// @Tags AuthAPI
// @Security ApiKeyAuth
// @Summary Refresh current access token
// @Success 200 {object} dtos.Result[dtos.LoginToken]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/auth/refresh-token [post]
func (a *Auth) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	refreshToken := helper.GetToken(c)
	if refreshToken == "" {
		response.Error(c, errorx.ErrInvalidToken.New(ctx))
		return
	}

	data, err := a.AuthSVC.RefreshToken(ctx, refreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OkData(c, data)
}

// @Tags AuthAPI
// @Security ApiKeyAuth
// @Summary Get current user info
// @Success 200 {object} dtos.Result[models.User]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/auth/user [get]
func (a *Auth) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.AuthSVC.GetUserInfo(ctx)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OkData(c, data)
}

// @Tags AuthAPI
// @Security ApiKeyAuth
// @Summary Change current user password
// @Param body body dtos.AuthUpdatePasswordReq true "Request body"
// @Success 200 {object} dtos.Result[any]
// @Failure 400 {object} dtos.Result[any]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/auth/password [put]
func (a *Auth) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(dtos.AuthUpdatePasswordReq)
	if err := c.ShouldBindJSON(item); err != nil {
		response.Error(c, err)
		return
	}

	err := a.AuthSVC.UpdatePassword(ctx, item)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c)
}

// @Tags AuthAPI
// @Security ApiKeyAuth
// @Summary Query current user menus based on the current user role
// @Success 200 {object} dtos.Result[models.Menus]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/auth/menus [get]
func (a *Auth) QueryMenus(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.AuthSVC.QueryMenus(ctx)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OkData(c, data)
}

// @Tags AuthAPI
// @Security ApiKeyAuth
// @Summary Update current user info
// @Param body body dtos.AuthUpdateUserReq true "Request body"
// @Success 200 {object} dtos.Result[any]
// @Failure 400 {object} dtos.Result[any]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/auth/user [put]
func (a *Auth) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(dtos.AuthUpdateUserReq)
	if err := c.ShouldBindJSON(item); err != nil {
		response.Error(c, err)
		return
	}

	err := a.AuthSVC.UpdateUser(ctx, item)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c)
}
