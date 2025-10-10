package modules

import (
	"context"
	"gin-admin/internal/api"
	"gin-admin/internal/errorx"
	"gin-admin/internal/middleware/recovery"
	"gin-admin/internal/types"
	"gin-admin/locales"
	"gin-admin/pkg/response"
	"log/slog"
	"net/http"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func NewHttp(ctx context.Context, app types.AppContext) (*http.Server, error) {

	cfg := app.Config()

	if cfg.IsDebug() {
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

	api.RegRoutes(app, e) // 注册 api

	// Register swagger
	if cfg.Swagger.Enable {
		g := e.Group("") // .Use(a.middlewares.Auth())
		g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	server := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      e,
		ReadTimeout:  time.Second * time.Duration(cfg.HTTP.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(cfg.HTTP.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(cfg.HTTP.IdleTimeout),
	}

	app.AddCleaner(ctx, func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.HTTP.ShutdownTimeout))
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Failed to shutdown http server", "error", err)
		}
	})

	return server, nil
}
