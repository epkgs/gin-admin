package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"gin-admin/internal/api"
	"gin-admin/internal/app/modules"
	"gin-admin/internal/config"
	"gin-admin/internal/errorx"
	"gin-admin/internal/middleware/recovery"
	"gin-admin/internal/model"
	"gin-admin/internal/service"
	_ "gin-admin/internal/swagger"
	"gin-admin/internal/types"
	"gin-admin/locales"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/jwtx"
	"gin-admin/pkg/response"
	"gin-admin/pkg/utils/util"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type App struct {
	cfg    *config.Config
	db     *gorm.DB
	cacher cachex.Cacher
	jwt    jwtx.Auther
	casbin types.Casbinx

	middlewares *modules.Middlewares

	cleaners []func()
}

var _ types.AppContext = (*App)(nil)

func New(ctx context.Context, c *config.Config) *App {

	app := &App{
		cfg:      c,
		cleaners: []func(){},
	}

	app.cacher = util.Must(modules.NewCacher(ctx, app))
	app.db = util.Must(modules.NewDB(ctx, app))
	app.jwt = util.Must(modules.NewJWT(ctx, app))
	app.casbin = util.Must(modules.NewCasbinx(ctx, app))

	app.middlewares = modules.NewMiddlewares(app)

	return app
}

func (a *App) Config() *config.Config {
	return a.cfg
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
	return a.db.AutoMigrate(model.Models...)
}

func (a *App) Init(ctx context.Context) error {
	if a.Config().DB.AutoMigrate {
		if err := a.autoMigrate(ctx); err != nil {
			return err
		}

		// 插入 super 账户
		if err := service.NewUser(a).InitSuperUserIfNeed(ctx); err != nil {
			return err
		}
	}

	if err := a.Casbin().Load(ctx); err != nil {
		return err
	}

	// Init menu data
	if err := service.NewMenu(a).InitIfNeed(ctx); err != nil {
		panic(err)
	}

	return nil
}

func (a *App) InitHttp(ctx context.Context) error {
	if a.cfg.IsDebug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.New()
	e.GET("/health", func(c *gin.Context) {
		response.OK(c)
	})
	e.Use(recovery.New(recovery.Config{Skip: 3}))
	e.NoMethod(func(c *gin.Context) {
		response.Error(c, errorx.ErrMethodNotAllowed)
	})
	e.NoRoute(func(c *gin.Context) {
		response.Error(c, errorx.ErrNotFound.WithMsg(locales.Def.Str("Route not found")))
	})

	if err := api.RegisterRouters(a, e); err != nil {
		return err
	}

	// Register swagger
	if a.cfg.Swagger.Enable {
		g := e.Group("") // .Use(a.middlewares.Auth())
		g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	addr := a.cfg.HTTP.Addr
	slog.Info(fmt.Sprintf("HTTP server is listening on %s", addr))
	srv := &http.Server{
		Addr:         addr,
		Handler:      e,
		ReadTimeout:  time.Second * time.Duration(a.cfg.HTTP.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(a.cfg.HTTP.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(a.cfg.HTTP.IdleTimeout),
	}

	go func() {
		var err error
		if a.cfg.HTTP.CertFile != "" && a.cfg.HTTP.KeyFile != "" {
			srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = srv.ListenAndServeTLS(a.cfg.HTTP.CertFile, a.cfg.HTTP.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to listen http server", "error", err)
		}
	}()

	a.AddCleaner(ctx, func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(a.cfg.HTTP.ShutdownTimeout))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Failed to shutdown http server", "error", err)
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
