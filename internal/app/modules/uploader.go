package modules

import (
	"context"

	"gin-admin/internal/types"
	"gin-admin/pkg/uploader"
)

func InitUploader(ctx context.Context, app types.AppContext) (*uploader.Uploader, error) {

	cfg := app.Config().Upload

	up := uploader.New(func(opt *uploader.Option) {
		opt.UploadPath = cfg.Path
		opt.UseDateDir = cfg.UseDateDir
	})

	return up, nil
}
