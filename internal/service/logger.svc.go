package service

import (
	"context"
	"strings"

	"gin-admin/internal/model/bo"
	"gin-admin/internal/model/dto"
	"gin-admin/internal/model/po"
	"gin-admin/internal/types"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

// Logger management
type Logger struct {
	app types.AppContext
	db  *gorm.DB
	q   *bo.Query
}

func NewLogger(app types.AppContext) *Logger {
	return &Logger{
		app: app,
		db:  app.DB(),
		q:   bo.Use(app.DB()),
	}
}

// List loggers from the data access object based on the provided parameters and options.
func (a *Logger) List(ctx context.Context, req dto.LoggerListReq) (*dto.List[*po.Logger], error) {

	l := a.q.Logger
	u := a.q.User
	db := a.db

	l.WithContext(ctx).Scopes()

	scope := func(d gen.Dao) gen.Dao {

		d = d.Select(l.ALL, u.NickName, u.Username).
			Join(u)

		if req.Level != "" {
			d = d.Where(l.Level.Eq(req.Level))
		}
		if len(req.Message) > 0 {
			d = d.Where(l.Message.Like("%" + req.Message + "%"))
		}
		if req.TraceID != "" {
			d = d.Where(l.TraceID.Eq(req.TraceID))
		}
		if req.UserName != "" {
			d = d.Where(u.Username.Like("%" + req.UserName + "%"))
		}
		if req.Tag != "" {
			tags := strings.Split(req.Tag, ",")
			d = d.Where(l.Tag.In(tags...))
		}

		if req.StartTime != "" {
			d = d.Where(field.CompareSubQuery(field.GteOp, l.CreatedAt, db.Raw(req.StartTime)))
		}

		if req.EndTime != "" {
			d = d.Where(field.CompareSubQuery(field.LteOp, l.CreatedAt, db.Raw(req.EndTime)))
		}

		// if startTime, err := dateparse.ParseAny(req.StartTime); err == nil {
		// 	if endTime, err := dateparse.ParseAny(req.EndTime); err == nil {
		// 		d = d.Where(l.CreatedAt.Between(startTime, endTime))
		// 	} else {
		// 		d = d.Where(l.CreatedAt.Gte(startTime))
		// 	}
		// } else if endTime, err := dateparse.ParseAny(req.EndTime); err == nil {
		// 	d = d.Where(l.CreatedAt.Lte(endTime))
		// }

		return d
	}

	list, err := l.WithContext(ctx).Scopes(scope, req.PageScope()).Find()
	if err != nil {
		return nil, err
	}

	count, err := l.WithContext(ctx).Scopes(scope).Count()
	if err != nil {
		return nil, err
	}

	return dto.NewList(list, &dto.Pager{
		Page:  req.Page,
		Limit: req.Limit,
		Total: count,
	}), nil
}
