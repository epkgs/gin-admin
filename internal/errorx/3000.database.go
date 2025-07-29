package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
)

var dbI18n = i18n.NewCatalog("database")

func init() {
	dbI18n.LoadTranslations()
}

var (
	ErrDatabase            = Define(dbI18n, 3000, "database error", http.StatusInternalServerError)             // 数据库错误
	ErrDatabaseQuery       = Define(dbI18n, 3001, "database query error", http.StatusInternalServerError)       // 数据库查询错误
	ErrDatabaseCreate      = Define(dbI18n, 3002, "database create error", http.StatusInternalServerError)      // 数据库创建错误
	ErrDatabaseUpdate      = Define(dbI18n, 3003, "database update error", http.StatusInternalServerError)      // 数据库更新错误
	ErrDatabaseDelete      = Define(dbI18n, 3004, "database delete error", http.StatusInternalServerError)      // 数据库删除错误
	ErrDatabaseConnection  = Define(dbI18n, 3005, "database connection error", http.StatusInternalServerError)  // 数据库连接错误
	ErrDatabaseTransaction = Define(dbI18n, 3006, "database transaction error", http.StatusInternalServerError) // 数据库事务错误
	ErrQueryParamEmpty     = Define(dbI18n, 3007, "query parameter cannot be empty", http.StatusBadRequest)     // 查询参数不能为空
	ErrRecordNotExist      = Define(dbI18n, 3008, "record does not exist", http.StatusNotFound)                 // 记录不存在
	ErrNothingUpdate       = Define(dbI18n, 3009, "nothing to update", http.StatusBadRequest)                   // 未更新任何数据
)
