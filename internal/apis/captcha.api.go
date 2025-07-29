package apis

import (
	"context"

	"gin-admin/internal/dtos"
	"gin-admin/internal/services"
	"gin-admin/internal/types"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type Captcha struct {
	CaptchaSVC *services.Captcha
}

func NewCaptcha(app types.AppContext) *Captcha {
	handler := &Captcha{
		CaptchaSVC: services.NewCaptcha(app),
	}

	app.Middlewares().Casbin().Exclude("/api/v1/captcha/")
	app.Middlewares().Auth().Exclude("/api/v1/captcha/")

	app.Routers().GroupAPI("/api/v1/captcha", func(ctx context.Context, g *gin.RouterGroup, e *gin.Engine) error {
		g.GET("id", handler.GetCaptcha)
		g.GET("image", handler.Image)
		return nil
	})

	return handler
}

// @Tags CaptchaAPI
// @Summary Get captcha ID
// @Success 200 {object} dtos.Result[dtos.Captcha]
// @Router /api/v1/captcha/id [get]
func (a *Captcha) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.CaptchaSVC.GetCaptcha(ctx)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OkData(c, data)
}

// @Tags CaptchaAPI
// @Summary Response captcha image
// @Param request query dtos.CaptchaImageReq false "query params"
// @Produce image/png
// @Success 200 "Captcha image"
// @Failure 404 {object} dtos.Result[any]
// @Router /api/v1/captcha/image [get]
func (a *Captcha) Image(c *gin.Context) {
	ctx := c.Request.Context()
	var req dtos.CaptchaImageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, err)
		return
	}

	if err := a.CaptchaSVC.ResponseCaptcha(ctx, c.Writer, req.ID, req.Reload); err != nil {
		response.Error(c, err)
	}
}
