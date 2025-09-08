package v1

import (
	"gin-admin/internal/dto"
	"gin-admin/internal/service"
	"gin-admin/internal/types"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// Menu management for SYS
type Menu struct {
	app     types.AppContext
	MenuSVC *service.Menu
}

func NewMenu(app types.AppContext) *Menu {
	return &Menu{
		app:     app,
		MenuSVC: service.NewMenu(app),
	}
}

func (a *Menu) RegisterRouter(group *gin.RouterGroup, engine *gin.Engine) {

	g := group.Group("menus")
	g.Use(
		a.app.Middlewares().Auth(),
	)

	g.GET("", a.Query)
	g.GET(":id", a.Get)
	g.POST("", a.Create)
	g.PUT(":id", a.Update)
	g.DELETE(":id", a.Delete)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Query menu tree data
// @Param request query dto.MenuListReq false "query params"
// @Success 200 {object} dto.ResultList[model.Menu]
// @Failure 401 {object} dto.Result[any]
// @Failure 500 {object} dto.Result[any]
// @Router /api/v1/menus [get]
func (a *Menu) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params dto.MenuListReq
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Error(c, err)
		return
	}

	result, err := a.MenuSVC.List(ctx, params)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.List(c, result.Items, &result.Pager)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Get menu record by ID
// @Param id path string true "unique id"
// @Success 200 {object} dto.Result[model.Menu]
// @Failure 401 {object} dto.Result[any]
// @Failure 500 {object} dto.Result[any]
// @Router /api/v1/menus/{id} [get]
func (a *Menu) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.MenuSVC.Get(ctx, c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OkData(c, item)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Create menu record
// @Param body body dto.MenuCreateReq true "Request body"
// @Success 200 {object} dto.Result[model.Menu]
// @Failure 400 {object} dto.Result[any]
// @Failure 401 {object} dto.Result[any]
// @Failure 500 {object} dto.Result[any]
// @Router /api/v1/menus [post]
func (a *Menu) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(dto.MenuCreateReq)
	if err := c.ShouldBindJSON(item); err != nil {
		response.Error(c, err)
		return
	}

	result, err := a.MenuSVC.Create(ctx, item)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OkData(c, result)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Update menu record by ID
// @Param id path string true "unique id"
// @Param body body dto.MenuUpdateReq true "Request body"
// @Success 200 {object} dto.Result[any]
// @Failure 400 {object} dto.Result[any]
// @Failure 401 {object} dto.Result[any]
// @Failure 500 {object} dto.Result[any]
// @Router /api/v1/menus/{id} [put]
func (a *Menu) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(dto.MenuUpdateReq)
	if err := c.ShouldBindJSON(item); err != nil {
		response.Error(c, err)
		return
	}

	err := a.MenuSVC.Update(ctx, c.Param("id"), item)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c)
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Delete menu record by ID
// @Param id path string true "unique id"
// @Success 200 {object} dto.Result[any]
// @Failure 401 {object} dto.Result[any]
// @Failure 500 {object} dto.Result[any]
// @Router /api/v1/menus/{id} [delete]
func (a *Menu) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuSVC.Delete(ctx, c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c)
}
