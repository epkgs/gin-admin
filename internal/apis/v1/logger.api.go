package v1

import (
	"gin-admin/internal/dtos"
	"gin-admin/internal/services"
	"gin-admin/internal/types"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// Logger management
type Logger struct {
	app       types.AppContext
	LoggerSVC *services.Logger
}

func NewLogger(app types.AppContext) *Logger {
	return &Logger{
		app:       app,
		LoggerSVC: services.NewLogger(app),
	}
}

func (a *Logger) RegisterRouter(group *gin.RouterGroup, engine *gin.Engine) {
	g := group.Group("loggers")
	g.Use(
		a.app.Middlewares().Auth(),
		a.app.Middlewares().Casbin(),
	)

	g.GET("", a.Query)
}

// @Tags LoggerAPI
// @Security ApiKeyAuth
// @Summary Query logger list
// @Param request query dtos.LoggerListReq false "query params"
// @Success 200 {object} dtos.ResultList[models.Logger]
// @Failure 401 {object} dtos.Result[any]
// @Failure 500 {object} dtos.Result[any]
// @Router /api/v1/loggers [get]
func (a *Logger) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var req dtos.LoggerListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := a.LoggerSVC.List(ctx, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.List(c, result.Items, &result.Pager)
}
