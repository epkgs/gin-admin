package test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"gin-admin/internal/apis"
	"gin-admin/internal/app"
	"gin-admin/internal/configs"

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

	configs.MustLoad(context.Background(), "config.yml")

	configs.C.DB.AutoMigrate = true

	_ = os.RemoveAll(configs.C.DB.DSN)
	ctx := context.Background()
	app := app.New(ctx, configs.C)

	if err := app.Init(ctx); err != nil {
		panic(err)
	}

	engine = gin.New()
	err := apis.RegisterRouters(app, engine)
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
