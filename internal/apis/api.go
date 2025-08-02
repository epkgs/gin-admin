package apis

import (
	v1 "gin-admin/internal/apis/v1"
	"gin-admin/internal/types"

	"github.com/gin-gonic/gin"
)

func RegisterRouters(app types.AppContext, e *gin.Engine) error {
	apiV1 := e.Group("/api/v1")

	apiV1.Use(
		app.Middlewares().I18n(),
		app.Middlewares().Cors(),
		app.Middlewares().Trace(),
		app.Middlewares().Logger(),
		app.Middlewares().CopyBody(),
		// app.Middlewares().Auth(),
		app.Middlewares().RateLimiter(),
		// app.Middlewares().Casbin(),
		app.Middlewares().Prometheus(),
	)

	registerRouters(apiV1, e,
		v1.NewAuth(app),
		v1.NewCaptcha(app),
		v1.NewLogger(app),
		v1.NewMenu(app),
		v1.NewRole(app),
		v1.NewUser(app),
	)

	return nil
}

type routerRegister interface {
	RegisterRouter(group *gin.RouterGroup, engine *gin.Engine)
}

func registerRouters(group *gin.RouterGroup, engine *gin.Engine, registers ...routerRegister) {
	for _, register := range registers {
		register.RegisterRouter(group, engine)
	}
}
