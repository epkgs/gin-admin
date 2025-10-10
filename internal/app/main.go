package app

import (
	"context"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"os"

	"gin-admin/internal/config"
	"gin-admin/pkg/gormx"
	"gin-admin/pkg/logger"
)

// The Run function initializes and starts a service with configuration and logger, and handles
// cleanup upon exit.
func Run(ctx context.Context, configFile string) error {

	// Load configuration.
	cfg := config.MustLoad(ctx, configFile)

	// Initialize logger.
	cleanLoggerFn, err := logger.Init(func(o *logger.Config) {

		o.Level = cfg.Logger.Level
		o.ConsoleEnable = cfg.Logger.Console.Enable

		if cfg.Logger.File.Enable {
			o.FileName = cfg.Logger.File.Path
			o.FileMaxSize = cfg.Logger.File.MaxSize
			o.FileMaxBackups = cfg.Logger.File.MaxBackups
		}

		if cfg.Logger.Database.Enable {
			db, err := gormx.New(gormx.Config{
				DBType:       cfg.DB.Type,
				DSN:          cfg.DB.DSN,
				MaxLifetime:  cfg.DB.MaxLifetime,
				MaxIdleTime:  cfg.DB.MaxIdleTime,
				MaxOpenConns: cfg.DB.MaxOpenConns,
				MaxIdleConns: cfg.DB.MaxIdleConns,
			})
			if err == nil {
				o.Database = db
			}
		}

	})
	if err != nil {
		return err
	}

	ctx = logger.With(ctx, "tag", "main")

	logger.Info(ctx, "starting service ...",
		"version", cfg.Version,
		"pid", os.Getpid(),
		"config", cfg.ConfigFile,
		"env", cfg.AppEnv,
	)

	// Start pprof server.
	if addr := cfg.Pprof.Addr; addr != "" {
		logger.Info(ctx, "pprof server is listening on "+addr)
		go func() {
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				logger.Error(ctx, "failed to listen pprof server",
					"error", err,
				)
			}
		}()
	}

	app := New(ctx, cfg)
	app.AddCleaner(ctx, cleanLoggerFn)

	if err := app.Init(ctx); err != nil {
		return err
	}

	return app.ListenAndServe(ctx)
}
