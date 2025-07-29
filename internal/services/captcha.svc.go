package services

import (
	"context"
	"net/http"

	"gin-admin/internal/configs"
	"gin-admin/internal/dtos"
	"gin-admin/internal/errorx"
	"gin-admin/internal/types"

	"github.com/LyricTian/captcha"
)

// Captcha management for SYS
type Captcha struct {
}

func NewCaptcha(app types.AppContext) *Captcha {
	return &Captcha{}
}

// This function generates a new captcha ID and returns it as a `dtos.Captcha` struct. The length of
// the captcha is determined by the `configs.C.Captcha.Length` configuration value.
func (a *Captcha) GetCaptcha(ctx context.Context) (*dtos.Captcha, error) {
	return &dtos.Captcha{
		CaptchaID: captcha.NewLen(configs.C.Captcha.Length),
	}, nil
}

// Response captcha image
func (a *Captcha) ResponseCaptcha(ctx context.Context, w http.ResponseWriter, id string, reload bool) error {
	if reload && !captcha.Reload(id) {
		return errorx.ErrCaptchaIDNotFound.New(ctx)
	}

	err := captcha.WriteImage(w, id, configs.C.Captcha.Width, configs.C.Captcha.Height)
	if err != nil {
		if err == captcha.ErrNotFound {
			return errorx.ErrCaptchaIDNotFound.New(ctx).Wrap(err)
		}
		return errorx.ErrInternal.New(ctx).Wrap(err)
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")
	return nil
}
