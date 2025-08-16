package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
)

var fileI18n = i18n.NewBundle("file")

func init() {
	fileI18n.Load()
}

type fileErrors struct {
	NotFound *Definition // 文件不存在
	Upload   *Definition // 文件上传失败
	Delete   *Definition // 文件删除失败
	Update   *Definition // 文件更新失败
	Download *Definition // 文件下载失败
}

var File = &fileErrors{
	NotFound: Define(fileI18n, 4000, "file not found", http.StatusNotFound),                  // 文件不存在
	Upload:   Define(fileI18n, 4001, "file upload failed", http.StatusInternalServerError),   // 文件上传失败
	Delete:   Define(fileI18n, 4002, "file deletion failed", http.StatusInternalServerError), // 文件删除失败
	Update:   Define(fileI18n, 4003, "file update failed", http.StatusInternalServerError),   // 文件更新失败
	Download: Define(fileI18n, 4004, "file download failed", http.StatusInternalServerError), // 文件下载失败
}
