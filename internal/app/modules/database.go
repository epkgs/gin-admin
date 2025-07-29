package modules

import (
	"context"

	"gin-admin/internal/types"
	"gin-admin/pkg/gormx"

	"gorm.io/gorm"
)

// It creates a new database connection, and returns a function that closes the connection
func InitDB(ctx context.Context, app types.AppContext) (*gorm.DB, error) {
	cfg := app.Config().DB

	resolver := make([]gormx.ResolverConfig, len(cfg.Resolver))
	for i, v := range cfg.Resolver {
		resolver[i] = gormx.ResolverConfig{
			DBType:   v.DBType,
			Sources:  v.Sources,
			Replicas: v.Replicas,
			Tables:   v.Tables,
		}
	}

	db, err := gormx.New(gormx.Config{
		Debug:        cfg.Debug,
		PrepareStmt:  cfg.PrepareStmt,
		DBType:       cfg.Type,
		DSN:          cfg.DSN,
		MaxLifetime:  cfg.MaxLifetime,
		MaxIdleTime:  cfg.MaxIdleTime,
		MaxOpenConns: cfg.MaxOpenConns,
		MaxIdleConns: cfg.MaxIdleConns,
		TablePrefix:  cfg.TablePrefix,
		Resolver:     resolver,
	})
	if err != nil {
		return nil, err
	}

	app.AddCleaner(ctx, func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	return db, nil
}
