package services

import (
	"context"
	"fmt"
	"strings"

	"gin-admin/internal/dtos"
	"gin-admin/internal/models"
	"gin-admin/internal/repositories"
	"gin-admin/internal/types"
	"gin-admin/pkg/gormx"

	"gorm.io/gorm"
)

// Logger management
type Logger struct {
	LoggerRepo *repositories.Logger
}

func NewLogger(app types.AppContext) *Logger {
	return &Logger{
		LoggerRepo: repositories.NewLogger(app.DB()),
	}
}

// List loggers from the data access object based on the provided parameters and options.
func (a *Logger) List(ctx context.Context, req dtos.LoggerListReq) (*dtos.List[*models.Logger], error) {
	option := func(d *gorm.DB) *gorm.DB {

		db := d.Table(fmt.Sprintf("%s AS a", new(models.Logger).TableName()))
		db = db.Joins(fmt.Sprintf("left join %s b on a.user_id=b.id", new(models.User).TableName()))
		db = db.Select("a.*, b.nick_name as nick_name, b.username as username") // 对应 models.Logger 的 NickName 和 UserName

		if v := req.Level; v != "" {
			db = db.Where("a.level = ?", v)
		}
		if v := req.Message; len(v) > 0 {
			db = db.Where("a.message LIKE ?", "%"+v+"%")
		}
		if v := req.TraceID; v != "" {
			db = db.Where("a.trace_id = ?", v)
		}
		if v := req.UserName; v != "" {
			db = db.Where("b.username LIKE ?", "%"+v+"%")
		}
		if v := req.Tag; v != "" {
			tags := strings.Split(v, ",")
			db = db.Where("a.tag IN (?)", tags)
		}
		if start := req.StartTime; start != "" {
			if end := req.EndTime; end != "" {
				db = db.Where("a.created_at BETWEEN ? AND ?", start, end)
			} else {
				db = db.Where("a.created_at >= ?", start)
			}
		} else if end := req.EndTime; end != "" {
			db = db.Where("a.created_at <= ?", end)
		}

		return db
	}

	list, err := a.LoggerRepo.Find(ctx, option, gormx.WithPage(req.Page, req.Limit))
	if err != nil {
		return nil, err
	}

	count, err := a.LoggerRepo.Count(ctx, option)
	if err != nil {
		return nil, err
	}

	return dtos.NewList(list, req.Page, req.Limit, count), nil
}
