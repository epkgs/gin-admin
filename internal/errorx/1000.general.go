package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
)

var gnI18n = i18n.NewBundle("general")

func init() {
	gnI18n.Load()
}

type generalError struct {
	Success          *Definition                           // 成功
	Unknown          *Definition                           // 未知错误
	Internal         *Definition                           // 服务器内部错误
	InvalidParams    *DefinitionF[struct{ Params string }] // 请求参数错误：{{.Params}}
	BadRequest       *Definition                           // 请求错误
	Unauthorized     *Definition                           // 未授权
	Forbidden        *Definition                           // 禁止访问
	RecordNotFound   *Definition                           // 数据未找到
	Timeout          *Definition                           // 请求超时
	TooManyRequests  *Definition                           // 请求过多
	AccessDenied     *Definition                           // 访问被拒绝
	RouteNotFound    *Definition                           // 请求路径不存在
	MethodNotAllowed *Definition                           // 请求函数不允许
	RequestTooLarge  *DefinitionF[struct{ Byte int64 }]    // 请求体过大，限制 {{.Byte}} 字节
	ReadConfigFile   *DefinitionF[struct{ File string }]   // 读取配置文件失败: {{.File}}
	UnmarshalConfig  *DefinitionF[struct{ File string }]   // 解析配置文件失败: {{.File}}
	GetConfigFile    *DefinitionF[struct{ File string }]   // 访问配置文件 {{.File}} 失败
	WalkDir          *DefinitionF[struct{ Dir string }]    // 遍历目录 {{.Dir}} 失败
	MenuNotFound     *Definition                           // 菜单不存在
}

var General = &generalError{
	Success: Define(gnI18n, 0, "success", http.StatusOK),

	Unknown:          Define(gnI18n, 5000, "unknown error", http.StatusInternalServerError),
	Internal:         Define(gnI18n, 1000, "internal error", http.StatusInternalServerError),
	InvalidParams:    Definef[struct{ Params string }](gnI18n, 1001, "invalid parameters: {{.Params}}", http.StatusBadRequest),
	BadRequest:       Define(gnI18n, 1002, "bad request", http.StatusBadRequest),
	Unauthorized:     Define(gnI18n, 1003, "unauthorized", http.StatusUnauthorized),
	Forbidden:        Define(gnI18n, 1004, "forbidden", http.StatusForbidden),
	RecordNotFound:   Define(gnI18n, 1005, "record not found", http.StatusNotFound),
	Timeout:          Define(gnI18n, 1006, "request timeout", http.StatusRequestTimeout),
	TooManyRequests:  Define(gnI18n, 1007, "too many requests", http.StatusTooManyRequests),
	AccessDenied:     Define(gnI18n, 1008, "access denied", http.StatusForbidden),
	RouteNotFound:    Define(gnI18n, 1009, "route not found", http.StatusNotFound),
	MethodNotAllowed: Define(gnI18n, 1010, "method not allowed", http.StatusMethodNotAllowed),
	RequestTooLarge:  Definef[struct{ Byte int64 }](gnI18n, 1011, "request body too large, limit {{.Byte}} byte", http.StatusRequestEntityTooLarge),
	ReadConfigFile:   Definef[struct{ File string }](gnI18n, 1012, "failed to read config file: {{.File}}", http.StatusInternalServerError),
	UnmarshalConfig:  Definef[struct{ File string }](gnI18n, 1013, "failed to unmarshal config: {{.File}}", http.StatusInternalServerError),
	GetConfigFile:    Definef[struct{ File string }](gnI18n, 1014, "failed to get config file: {{.File}}", http.StatusInternalServerError),
	WalkDir:          Definef[struct{ Dir string }](gnI18n, 1015, "failed to walk dir: {{.Dir}}", http.StatusInternalServerError),
	MenuNotFound:     Define(gnI18n, 1016, "menu not found", http.StatusNotFound),
}
