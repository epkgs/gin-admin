package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-admin/internal/app/modules"
	"gin-admin/internal/config"
	"gin-admin/internal/model"
	"gin-admin/internal/service"
	_ "gin-admin/internal/swagger"
	"gin-admin/internal/types"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/jwtx"
	"gin-admin/pkg/logger"
	"gin-admin/pkg/utils/util"

	"gorm.io/gorm"
)

type App struct {
	cfg         *config.Config
	db          *gorm.DB
	cacher      cachex.Cacher
	jwt         jwtx.Auther
	casbin      types.Casbinx
	middlewares *modules.Middlewares
	http        *http.Server

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
	app.http = util.Must(modules.NewHttp(ctx, app))

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

func (a *App) Init(ctx context.Context) error {
	if a.Config().DB.AutoMigrate {
		if err := a.db.AutoMigrate(model.Models...); err != nil {
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
		return err
	}

	return nil
}

func (a *App) ListenAndServe(ctx context.Context) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		addr := a.cfg.HTTP.Addr
		logger.Info(ctx, fmt.Sprintf("HTTP server is listening on %s", addr))

		var err error
		if a.cfg.HTTP.CertFile != "" && a.cfg.HTTP.KeyFile != "" {
			a.http.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = a.http.ListenAndServeTLS(a.cfg.HTTP.CertFile, a.cfg.HTTP.KeyFile)
		} else {
			err = a.http.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, "Failed to listen http server", "error", err)
			close(sc) // 主动关闭信号通道以触发服务退出
		}
	}()

EXIT:

	for {
		sig := <-sc
		logger.Info(ctx, "Received signal",
			"signal", sig.String(),
		)

		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	a.Release(ctx)
	logger.Info(ctx, "Server exit, bye...")
	time.Sleep(time.Millisecond * 100)
	os.Exit(state)
	return nil
}

func (a *App) Release(ctx context.Context) error {
	for _, cleaner := range a.cleaners {
		cleaner()
	}
	return nil
}
