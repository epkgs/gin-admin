package api

import (
	v1 "gin-admin/internal/api/v1"
	"gin-admin/internal/types"

	"github.com/gin-gonic/gin"
)

func RegRoutes(app types.AppContext, e *gin.Engine) {
	apiV1 := e.Group("/api/v1")

	apiV1.Use(
		app.Middlewares().I18n(),
		app.Middlewares().Cors(),
		app.Middlewares().Trace(),
		app.Middlewares().Logger(),
		// app.Middlewares().Auth(),
		app.Middlewares().RateLimiter(),
		app.Middlewares().Prometheus(),
	)

	v1.NewAuth(app).RegRoutes(apiV1, e)
	v1.NewCaptcha(app).RegRoutes(apiV1, e)
	v1.NewLogger(app).RegRoutes(apiV1, e)
	v1.NewMenu(app).RegRoutes(apiV1, e)
	v1.NewRole(app).RegRoutes(apiV1, e)
	v1.NewUser(app).RegRoutes(apiV1, e)
	v1.NewAPI(app).RegRoutes(apiV1, e)
}
