package service

import (
	"context"
	"net/http"

	"gin-admin/internal/errorx"
	"gin-admin/internal/model/dto"
	"gin-admin/internal/types"
	"gin-admin/locales"

	"github.com/LyricTian/captcha"
)

// Captcha management for SYS
type Captcha struct {
	app types.AppContext
}

func NewCaptcha(app types.AppContext) *Captcha {
	return &Captcha{
		app: app,
	}
}

// This function generates a new captcha ID and returns it as a `dto.Captcha` struct. The length of
// the captcha is determined by the `config.Captcha.Length` configuration value.
func (a *Captcha) GetCaptcha(ctx context.Context) (*dto.Captcha, error) {
	return &dto.Captcha{
		CaptchaID: captcha.NewLen(a.app.Config().Captcha.Length),
	}, nil
}

// Response captcha image
func (a *Captcha) ResponseCaptcha(ctx context.Context, w http.ResponseWriter, id string, reload bool) error {
	if reload && !captcha.Reload(id) {
		return errorx.ErrBadRequest.WithMsg(locales.User.Str("Captcha id not found"))
	}

	err := captcha.WriteImage(w, id, a.app.Config().Captcha.Width, a.app.Config().Captcha.Height)
	if err != nil {
		if err == captcha.ErrNotFound {
			return errorx.ErrBadRequest.WithMsg(locales.User.Str("Captcha id not found")).Wrap(err)
		}
		return errorx.ErrInternalServerError.Wrap(err)
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")
	return nil
}
