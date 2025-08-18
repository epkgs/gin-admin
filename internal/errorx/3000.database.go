package errorx

import (
	"gin-admin/locales"
	"net/http"

	"github.com/epkgs/i18n/errors"
)

type dbErrors struct {
	Query           errors.I18nError // 数据库查询错误
	Create          errors.I18nError // 数据库创建错误
	Update          errors.I18nError // 数据库更新错误
	Delete          errors.I18nError // 数据库删除错误
	Connection      errors.I18nError // 数据库连接错误
	Transaction     errors.I18nError // 数据库事务错误
	QueryParamEmpty errors.I18nError // 查询参数不能为空
	RecordNotExist  errors.I18nError // 记录不存在
	NothingUpdate   errors.I18nError // 未更新任何数据
}

var DB = &dbErrors{
	Query:           New(locales.DB, 3001, "database query error", http.StatusInternalServerError),       // 数据库查询错误
	Create:          New(locales.DB, 3002, "database create error", http.StatusInternalServerError),      // 数据库创建错误
	Update:          New(locales.DB, 3003, "database update error", http.StatusInternalServerError),      // 数据库更新错误
	Delete:          New(locales.DB, 3004, "database delete error", http.StatusInternalServerError),      // 数据库删除错误
	Connection:      New(locales.DB, 3005, "database connection error", http.StatusInternalServerError),  // 数据库连接错误
	Transaction:     New(locales.DB, 3006, "database transaction error", http.StatusInternalServerError), // 数据库事务错误
	QueryParamEmpty: New(locales.DB, 3007, "query parameter cannot be empty", http.StatusBadRequest),     // 查询参数不能为空
	RecordNotExist:  New(locales.DB, 3008, "record does not exist", http.StatusNotFound),                 // 记录不存在
	NothingUpdate:   New(locales.DB, 3009, "nothing to update", http.StatusBadRequest),                   // 未更新任何数据
}
