package errorx

import (
	"gin-admin/locales"
	"net/http"

	"github.com/epkgs/i18n/errors"
)

type generalError struct {
	Success                 errors.I18nError                           // 成功
	Unknown                 errors.I18nError                           // 未知错误
	Internal                errors.I18nError                           // 服务器内部错误
	InvalidParams           errors.Definition[struct{ Params string }] // 请求参数错误：{{.Params}}
	BadRequest              errors.I18nError                           // 请求错误
	Unauthorized            errors.I18nError                           // 未授权
	Forbidden               errors.I18nError                           // 禁止访问
	RecordNotFound          errors.I18nError                           // 数据未找到
	Timeout                 errors.I18nError                           // 请求超时
	TooManyRequests         errors.I18nError                           // 请求过多
	AccessDenied            errors.I18nError                           // 访问被拒绝
	RouteNotFound           errors.I18nError                           // 请求路径不存在
	MethodNotAllowed        errors.I18nError                           // 请求函数不允许
	RequestTooLarge         errors.Definition[struct{ Byte int64 }]    // 请求体过大，限制 {{.Byte}} 字节
	ReadConfigFile          errors.Definition[struct{ File string }]   // 读取配置文件失败: {{.File}}
	UnmarshalConfig         errors.Definition[struct{ File string }]   // 解析配置文件失败: {{.File}}
	GetConfigFile           errors.Definition[struct{ File string }]   // 访问配置文件 {{.File}} 失败
	WalkDir                 errors.Definition[struct{ Dir string }]    // 遍历目录 {{.Dir}} 失败
	MenuNotFound            errors.I18nError                           // 菜单不存在
	UnmarshalJsonFileFailed errors.Definition[string]                  // 解析JSON文件 '%s' 失败
	UnmarshalYamlFileFailed errors.Definition[string]                  // 解析YAML文件 '%s' 失败
	UnsupportedFileType     errors.Definition[string]                  // 不支持的文件类型 '%s'
}

var General = &generalError{
	Success: New(locales.General, 0, "success", http.StatusOK),

	Unknown:                 New(locales.General, 5000, "unknown error", http.StatusInternalServerError),
	Internal:                New(locales.General, 1000, "internal error", http.StatusInternalServerError),
	InvalidParams:           Define[struct{ Params string }](locales.General, 1001, "invalid parameters: {{.Params}}", http.StatusBadRequest),
	BadRequest:              New(locales.General, 1002, "bad request", http.StatusBadRequest),
	Unauthorized:            New(locales.General, 1003, "unauthorized", http.StatusUnauthorized),
	Forbidden:               New(locales.General, 1004, "forbidden", http.StatusForbidden),
	RecordNotFound:          New(locales.General, 1005, "record not found", http.StatusNotFound),
	Timeout:                 New(locales.General, 1006, "request timeout", http.StatusRequestTimeout),
	TooManyRequests:         New(locales.General, 1007, "too many requests", http.StatusTooManyRequests),
	AccessDenied:            New(locales.General, 1008, "access denied", http.StatusForbidden),
	RouteNotFound:           New(locales.General, 1009, "route not found", http.StatusNotFound),
	MethodNotAllowed:        New(locales.General, 1010, "method not allowed", http.StatusMethodNotAllowed),
	RequestTooLarge:         Define[struct{ Byte int64 }](locales.General, 1011, "request body too large, limit {{.Byte}} byte", http.StatusRequestEntityTooLarge),
	ReadConfigFile:          Define[struct{ File string }](locales.General, 1012, "failed to read config file: {{.File}}", http.StatusInternalServerError),
	UnmarshalConfig:         Define[struct{ File string }](locales.General, 1013, "failed to unmarshal config: {{.File}}", http.StatusInternalServerError),
	GetConfigFile:           Define[struct{ File string }](locales.General, 1014, "failed to get config file: {{.File}}", http.StatusInternalServerError),
	WalkDir:                 Define[struct{ Dir string }](locales.General, 1015, "failed to walk dir: {{.Dir}}", http.StatusInternalServerError),
	MenuNotFound:            New(locales.General, 1016, "menu not found", http.StatusNotFound),
	UnmarshalJsonFileFailed: Define[string](locales.General, 1017, "unmarshal JSON file '%s' failed", http.StatusInternalServerError),
	UnmarshalYamlFileFailed: Define[string](locales.General, 1018, "unmarshal YAML file '%s' failed", http.StatusInternalServerError),
	UnsupportedFileType:     Define[string](locales.General, 1019, "unsupported file type '%s'", http.StatusBadRequest),
}
