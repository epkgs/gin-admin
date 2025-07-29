package apis

import (
	"context"

	"gin-admin/internal/dtos"
	"gin-admin/internal/services"
	"gin-admin/internal/types"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// Menu management for SYS
type Menu struct {
	MenuSVC *services.Menu
}

func NewMenu(app types.AppContext) *Menu {
	handler := &Menu{
		MenuSVC: services.NewMenu(app),
	}

	app.Routers().GroupAPI("/api/v1/menus", func(ctx context.Context, g *gin.RouterGroup, e *gin.Engine) error {
		g.GET("", handler.Query)
		g.GET(":id", handler.Get)
		g.POST("", handler.Create)
		g.PUT(":id", handler.Update)
		g.DELETE(":id", handler.Delete)
		return nil
	})

	return handler
}

// @Tags MenuAPI
// @Security ApiKeyAuth
// @Summary Query menu tree data
// @Param request query dtos.MenuListReq false "query params"
// @Success 200 {object} dtos.ResultList[models.Menu]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/menus [get]
func (a *Menu) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params dtos.MenuListReq
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
// @Success 200 {object} dtos.Result[models.Menu]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
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
// @Param body body dtos.MenuCreateReq true "Request body"
// @Success 200 {object} dtos.Result[models.Menu]
// @Failure 400 {object} dtos.Result[any]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/menus [post]
func (a *Menu) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(dtos.MenuCreateReq)
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
// @Param body body dtos.MenuUpdateReq true "Request body"
// @Success 200 {object} dtos.Result[any]
// @Failure 400 {object} dtos.Result[any]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/menus/{id} [put]
func (a *Menu) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(dtos.MenuUpdateReq)
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
// @Success 200 {object} dtos.Result[any]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
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
