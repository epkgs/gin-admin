package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
)

var dbI18n = i18n.NewBundle("database")

func init() {
	dbI18n.Load()
}

type dbErrors struct {
	Query           *Definition // 数据库查询错误
	Create          *Definition // 数据库创建错误
	Update          *Definition // 数据库更新错误
	Delete          *Definition // 数据库删除错误
	Connection      *Definition // 数据库连接错误
	Transaction     *Definition // 数据库事务错误
	QueryParamEmpty *Definition // 查询参数不能为空
	RecordNotExist  *Definition // 记录不存在
	NothingUpdate   *Definition // 未更新任何数据
}

var DB = &dbErrors{
	Query:           Define(dbI18n, 3001, "database query error", http.StatusInternalServerError),       // 数据库查询错误
	Create:          Define(dbI18n, 3002, "database create error", http.StatusInternalServerError),      // 数据库创建错误
	Update:          Define(dbI18n, 3003, "database update error", http.StatusInternalServerError),      // 数据库更新错误
	Delete:          Define(dbI18n, 3004, "database delete error", http.StatusInternalServerError),      // 数据库删除错误
	Connection:      Define(dbI18n, 3005, "database connection error", http.StatusInternalServerError),  // 数据库连接错误
	Transaction:     Define(dbI18n, 3006, "database transaction error", http.StatusInternalServerError), // 数据库事务错误
	QueryParamEmpty: Define(dbI18n, 3007, "query parameter cannot be empty", http.StatusBadRequest),     // 查询参数不能为空
	RecordNotExist:  Define(dbI18n, 3008, "record does not exist", http.StatusNotFound),                 // 记录不存在
	NothingUpdate:   Define(dbI18n, 3009, "nothing to update", http.StatusBadRequest),                   // 未更新任何数据
}
