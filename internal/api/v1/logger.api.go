package v1

import (
	"gin-admin/internal/dto"
	"gin-admin/internal/service"
	"gin-admin/internal/types"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// Logger management
type Logger struct {
	app       types.AppContext
	LoggerSVC *service.Logger
}

var _ types.RoutableHandler = (*Logger)(nil)

func NewLogger(app types.AppContext) *Logger {
	return &Logger{
		app:       app,
		LoggerSVC: service.NewLogger(app),
	}
}

func (a *Logger) RegRoutes(group *gin.RouterGroup, engine *gin.Engine) {
	g := group.Group("loggers")
	g.Use(
		a.app.Middlewares().Auth(),
	)

	g.GET("", a.Query)
}

// @Tags LoggerAPI
// @Security ApiKeyAuth
// @Summary Query logger list
// @Param request query dto.LoggerListReq false "query params"
// @Success 200 {object} dto.ResultList[model.Logger]
// @Failure 401 {object} dto.Result[any]
// @Failure 500 {object} dto.Result[any]
// @Router /api/v1/loggers [get]
func (a *Logger) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.LoggerListReq
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
