package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"gin-admin/internal/apis"
	"gin-admin/internal/app/modules"
	"gin-admin/internal/configs"
	"gin-admin/internal/errorx"
	"gin-admin/internal/models"
	"gin-admin/internal/services"
	"gin-admin/internal/types"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/jwtx"
	"gin-admin/pkg/logger"
	"gin-admin/pkg/middleware"
	"gin-admin/pkg/response"
	"gin-admin/pkg/uploader"
	"gin-admin/pkg/utils/util"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type App struct {
	config   *configs.Config
	db       *gorm.DB
	cacher   cachex.Cacher
	jwt      jwtx.Auther
	uploader *uploader.Uploader
	casbin   types.Casbinx

	middlewares *Middlewares

	cleaners []func()
}

var _ types.AppContext = (*App)(nil)

func New(ctx context.Context, c *configs.Config) *App {

	app := &App{
		config:   c,
		cleaners: []func(){},
	}

	app.cacher = util.Must(modules.InitCacher(ctx, app))
	app.db = util.Must(modules.InitDB(ctx, app))
	app.jwt = util.Must(modules.InitJWT(ctx, app))
	app.uploader = util.Must(modules.InitUploader(ctx, app))
	app.casbin = util.Must(modules.InitCasbinx(ctx, app))

	app.middlewares = NewMiddlewares(app)

	return app
}

func (a *App) Config() *configs.Config {
	return a.config
}

func (a *App) DB() *gorm.DB {
	return a.db
}

func (a *App) Cacher() cachex.Cacher {
	return a.cacher
}

func (a *App) Jwt() jwtx.Auther {
	return a.jwt
}

func (a *App) Uploader() *uploader.Uploader {
	return a.uploader
}

func (a *App) Casbin() types.Casbinx {
	return a.casbin
}

func (a *App) Middlewares() types.Middlewares {
	return a.middlewares
}

func (a *App) AddCleaner(ctx context.Context, cleaner func()) {
	a.cleaners = append(a.cleaners, cleaner)
}

func (a *App) autoMigrate(_ context.Context) error {
	return a.db.AutoMigrate(
		new(models.Logger),
		new(models.MenuRole),
		new(models.UserRole),
		new(models.Menu),
		new(models.Role),
		new(models.User),
	)
}

func (a *App) Init(ctx context.Context) error {
	if a.Config().DB.AutoMigrate {
		if err := a.autoMigrate(ctx); err != nil {
			return err
		}

		// 插入 super 账户
		if err := services.NewUser(a).InitSuperUserIfNeed(ctx); err != nil {
			return err
		}
	}

	if err := a.Casbin().Load(ctx); err != nil {
		return err
	}

	// Init menu data
	if err := services.NewMenu(a).InitIfNeed(ctx); err != nil {
		panic(err)
	}

	return nil
}

func (a *App) InitHttp(ctx context.Context) error {
	if a.config.IsDebug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.New()
	e.GET("/health", func(c *gin.Context) {
		response.OK(c)
	})
	e.Use(middleware.RecoveryWithConfig(middleware.RecoveryConfig{
		Skip: configs.C.Middleware.Recovery.Skip,
	}))
	e.NoMethod(func(c *gin.Context) {
		response.Error(c, errorx.ErrMethodNotAllowed.New(ctx))
	})
	e.NoRoute(func(c *gin.Context) {
		response.Error(c, errorx.ErrRouteNotFound.New(ctx))
	})

	if err := apis.RegisterRouters(a, e); err != nil {
		return err
	}

	// Register swagger
	if !configs.C.Swagger.Disable {
		e.StaticFile("/openapi.json", configs.C.Swagger.StaticFile)
		e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	if dir := configs.C.Middleware.Static.Root; dir != "" {
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:                 dir,
			ExcludedPathPrefixes: configs.C.Middleware.Static.ExcludedPathPrefixes,
		}))
	}

	addr := configs.C.HTTP.Addr
	logger.Info(ctx, fmt.Sprintf("HTTP server is listening on %s", addr))
	srv := &http.Server{
		Addr:         addr,
		Handler:      e,
		ReadTimeout:  time.Second * time.Duration(configs.C.HTTP.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(configs.C.HTTP.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(configs.C.HTTP.IdleTimeout),
	}

	go func() {
		var err error
		if configs.C.HTTP.CertFile != "" && configs.C.HTTP.KeyFile != "" {
			srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = srv.ListenAndServeTLS(configs.C.HTTP.CertFile, configs.C.HTTP.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, "Failed to listen http server", err)
		}
	}()

	a.AddCleaner(ctx, func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(configs.C.HTTP.ShutdownTimeout))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error(ctx, "Failed to shutdown http server", err)
		}
	})

	return nil
}

func (a *App) Release(ctx context.Context) error {
	for _, cleaner := range a.cleaners {
		cleaner()
	}
	return nil
}
