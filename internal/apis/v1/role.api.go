package v1

import (
	"gin-admin/internal/dtos"
	"gin-admin/internal/services"
	"gin-admin/internal/types"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// Role management for SYS
type Role struct {
	app     types.AppContext
	RoleSVC *services.Role
}

func NewRole(app types.AppContext) *Role {
	return &Role{
		app:     app,
		RoleSVC: services.NewRole(app),
	}
}

func (a *Role) RegisterRouter(group *gin.RouterGroup, engine *gin.Engine) {

	g := group.Group("roles")
	g.Use(
		a.app.Middlewares().Auth(),
		a.app.Middlewares().Casbin(),
	)

	g.GET("", a.Query)
	g.GET(":id", a.Get)
	g.POST("", a.Create)
	g.PUT(":id", a.Update)
	g.DELETE(":id", a.Delete)
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Query role list
// @Param request query dtos.RoleListReq false "query params"
// @Success 200 {object} dtos.ResultList[models.Role]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/roles [get]
func (a *Role) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params dtos.RoleListReq
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Error(c, err)
		return
	}

	result, err := a.RoleSVC.List(ctx, params)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.List(c, result.Items, &result.Pager)
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Get role record by ID
// @Param id path string true "unique id"
// @Success 200 {object} dtos.Result[models.Role]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/roles/{id} [get]
func (a *Role) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.RoleSVC.Get(ctx, c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OkData(c, item)
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Create role record
// @Param body body dtos.RoleCreateReq true "Request body"
// @Success 200 {object} dtos.Result[models.Role]
// @Failure 400 {object} dtos.Result[any]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/roles [post]
func (a *Role) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(dtos.RoleCreateReq)
	if err := c.ShouldBindJSON(item); err != nil {
		response.Error(c, err)
		return
	}

	result, err := a.RoleSVC.Create(ctx, *item)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OkData(c, result)
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Update role record by ID
// @Param id path string true "unique id"
// @Param body body dtos.RoleUpdateReq true "Request body"
// @Success 200 {object} dtos.Result[any]
// @Failure 400 {object} dtos.Result[any]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/roles/{id} [put]
func (a *Role) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(dtos.RoleUpdateReq)
	if err := c.ShouldBindJSON(item); err != nil {
		response.Error(c, err)
		return
	}

	err := a.RoleSVC.Update(ctx, c.Param("id"), item)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c)
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Delete role record by ID
// @Param id path string true "unique id"
// @Success 200 {object} dtos.Result[any]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/roles/{id} [delete]
func (a *Role) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleSVC.Delete(ctx, c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c)
}
