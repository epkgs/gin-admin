package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
)

var gnI18n = i18n.NewCatalog("general")

func init() {
	gnI18n.LoadTranslations()
}

var (
	Success = Define(gnI18n, 0, "success", http.StatusOK) // 成功

	ErrUnknown          = Define(gnI18n, 5000, "unknown error", http.StatusInternalServerError)                                                         // 未知错误
	ErrInternal         = Define(gnI18n, 1000, "internal error", http.StatusInternalServerError)                                                        // 服务器内部错误
	ErrInvalidParams    = Definef[struct{ Params string }](gnI18n, 1001, "invalid parameters: {{.Params}}", http.StatusBadRequest)                      // 请求参数错误：{{.Params}}
	ErrBadRequest       = Define(gnI18n, 1002, "bad request", http.StatusBadRequest)                                                                    // 请求错误
	ErrUnauthorized     = Define(gnI18n, 1003, "unauthorized", http.StatusUnauthorized)                                                                 // 未授权
	ErrForbidden        = Define(gnI18n, 1004, "forbidden", http.StatusForbidden)                                                                       // 禁止访问
	ErrRecordNotFound   = Define(gnI18n, 1005, "record not found", http.StatusNotFound)                                                                 // 数据未找到
	ErrTimeout          = Define(gnI18n, 1006, "request timeout", http.StatusRequestTimeout)                                                            // 请求超时
	ErrTooManyRequests  = Define(gnI18n, 1007, "too many requests", http.StatusTooManyRequests)                                                         // 请求过多
	ErrAccessDenied     = Define(gnI18n, 1008, "access denied", http.StatusForbidden)                                                                   // 访问被拒绝
	ErrRouteNotFound    = Define(gnI18n, 1009, "route not found", http.StatusNotFound)                                                                  // 请求路径不存在
	ErrMethodNotAllowed = Define(gnI18n, 1010, "method not allowed", http.StatusMethodNotAllowed)                                                       // 请求函数不允许
	ErrRequestTooLarge  = Definef[struct{ Byte int64 }](gnI18n, 1011, "request body too large, limit {{.Byte}} byte", http.StatusRequestEntityTooLarge) // 请求体过大，限制 {{.Byte}} 字节
	ErrReadConfigFile   = Definef[struct{ File string }](gnI18n, 1012, "failed to read config file: {{.File}}", http.StatusInternalServerError)         // 读取配置文件失败: {{.File}}
	ErrUnmarshalConfig  = Definef[struct{ File string }](gnI18n, 1013, "failed to unmarshal config: {{.File}}", http.StatusInternalServerError)         // 解析配置文件失败: {{.File}}
	ErrGetConfigFile    = Definef[struct{ File string }](gnI18n, 1014, "failed to get config file: {{.File}}", http.StatusInternalServerError)          // 访问配置文件 {{.File}} 失败
	ErrWalkDir          = Definef[struct{ Dir string }](gnI18n, 1015, "failed to walk dir: {{.Dir}}", http.StatusInternalServerError)                   // 遍历目录 {{.Dir}} 失败
	ErrMenuNotFound     = Define(gnI18n, 1016, "menu not found", http.StatusNotFound)                                                                   // 菜单不存在
)
