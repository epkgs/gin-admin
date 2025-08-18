package errorx

import (
	"gin-admin/locales"
	"net/http"

	"github.com/epkgs/i18n/errors"
)

type fileErrors struct {
	NotFound errors.I18nError // 文件不存在
	Upload   errors.I18nError // 文件上传失败
	Delete   errors.I18nError // 文件删除失败
	Update   errors.I18nError // 文件更新失败
	Download errors.I18nError // 文件下载失败
}

var File = &fileErrors{
	NotFound: New(locales.File, 4000, "file not found", http.StatusNotFound),                  // 文件不存在
	Upload:   New(locales.File, 4001, "file upload failed", http.StatusInternalServerError),   // 文件上传失败
	Delete:   New(locales.File, 4002, "file deletion failed", http.StatusInternalServerError), // 文件删除失败
	Update:   New(locales.File, 4003, "file update failed", http.StatusInternalServerError),   // 文件更新失败
	Download: New(locales.File, 4004, "file download failed", http.StatusInternalServerError), // 文件下载失败
}
