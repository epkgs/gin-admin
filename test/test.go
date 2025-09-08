package test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"gin-admin/internal/api"
	"gin-admin/internal/app"
	"gin-admin/internal/config"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
)

const (
	baseAPI = "/api/v1"
)

var (
	engine *gin.Engine
)

func init() {

	ctx := context.Background()

	// 获取当前文件的绝对路径
	_, filename, _, _ := runtime.Caller(0)

	// 获取当前文件所在的目录
	currentDir := filepath.Dir(filename)

	cfg := config.MustLoad(ctx, filepath.Join(currentDir, "config.yml"))

	_ = os.RemoveAll(cfg.RuntimePath)
	if cfg.DB.Type == "sqlite3" {
		_ = os.Remove(cfg.DB.DSN)
	}
	app := app.New(ctx, cfg)

	if err := app.Init(ctx); err != nil {
		panic(err)
	}

	engine = gin.New()
	err := api.RegisterRouters(app, engine)
	if err != nil {
		panic(err)
	}
}

func ApiTester(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(engine),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}
