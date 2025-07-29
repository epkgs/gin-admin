package apis

import (
	"context"

	"gin-admin/internal/dtos"
	"gin-admin/internal/services"
	"gin-admin/internal/types"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// Logger management
type Logger struct {
	LoggerSVC *services.Logger
}

func NewLogger(app types.AppContext) *Logger {
	handler := &Logger{
		LoggerSVC: services.NewLogger(app),
	}

	app.Routers().GroupAPI("/api/v1/loggers", func(ctx context.Context, g *gin.RouterGroup, e *gin.Engine) error {
		g.GET("", handler.Query)
		return nil
	})

	return handler
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
