package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
)

var storageI18n = i18n.NewCatalog("storage")

func init() {
	storageI18n.LoadTranslations()
}

var (
	ErrFileNotFound = Define(storageI18n, 4000, "file not found", http.StatusNotFound)                  // 文件不存在
	ErrFileUpload   = Define(storageI18n, 4001, "file upload failed", http.StatusInternalServerError)   // 文件上传失败
	ErrFileDelete   = Define(storageI18n, 4002, "file deletion failed", http.StatusInternalServerError) // 文件删除失败
	ErrFileUpdate   = Define(storageI18n, 4003, "file update failed", http.StatusInternalServerError)   // 文件更新失败
	ErrFileDownload = Define(storageI18n, 4004, "file download failed", http.StatusInternalServerError) // 文件下载失败
)
