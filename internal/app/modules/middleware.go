package modules

import (
	"fmt"
	"gin-admin/internal/middleware/auth"
	"gin-admin/internal/middleware/logger"
	"gin-admin/internal/middleware/promx"
	"gin-admin/internal/middleware/ratelimiter"
	"gin-admin/internal/types"
	"gin-admin/pkg/helper"
	"strings"
	"sync"
	"time"

	"github.com/epkgs/i18n"
	"github.com/gin-contrib/cors"
	"github.com/rs/xid"

	"github.com/gin-gonic/gin"
)

type mdw struct {
	once    sync.Once
	handler gin.HandlerFunc
}

type Middlewares struct {
	app types.AppContext

	i18n        mdw
	cors        mdw
	trace       mdw
	logger      mdw
	auth        mdw
	rateLimiter mdw
	prometheus  mdw
}

func NewMiddlewares(app types.AppContext) *Middlewares {
	return &Middlewares{
		app: app,
	}
}

func (m *Middlewares) I18n() gin.HandlerFunc {
	m.i18n.once.Do(func() {
		m.i18n.handler = i18n.GinMiddleware("zh")
	})
	return m.i18n.handler
}

func (m *Middlewares) Cors() gin.HandlerFunc {
	m.cors.once.Do(func() {
		m.cors.handler = cors.New(cors.Config{
			AllowOrigins:           []string{"*"},
			AllowMethods:           []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:           []string{"*"},
			AllowCredentials:       true,
			ExposeHeaders:          []string{"Content-Disposition"},
			MaxAge:                 86400 * time.Second,
			AllowWildcard:          true,
			AllowBrowserExtensions: false,
			AllowWebSockets:        true,
			AllowFiles:             true,
		})

	})

	return m.cors.handler
}

func (m *Middlewares) Trace() gin.HandlerFunc {
	m.trace.once.Do(func() {

		requestHeaderKey := "X-Request-Id"

		m.trace.handler = func(c *gin.Context) {

			traceID := c.GetHeader(requestHeaderKey)
			if traceID == "" {
				traceID = fmt.Sprintf("TRACE-%s", strings.ToUpper(xid.New().String()))
			}

			ctx := helper.WithTraceID(c.Request.Context(), traceID)
			c.Request = c.Request.WithContext(ctx)
			c.Writer.Header().Set(requestHeaderKey, traceID)
			c.Next()
		}
	})

	return m.trace.handler
}

func (m *Middlewares) Logger() gin.HandlerFunc {
	m.logger.once.Do(func() {
		cfg := m.app.Config().Logger

		m.logger.handler = logger.New(logger.Config{
			MaxRequestLen:  cfg.MaxRequestLen,
			MaxResponseLen: cfg.MaxResponseLen,
		})
	})

	return m.logger.handler
}

func (m *Middlewares) Auth() gin.HandlerFunc {

	m.auth.once.Do(func() {
		m.auth.handler = auth.New(m.app)
	})

	return m.auth.handler
}

func (m *Middlewares) RateLimiter() gin.HandlerFunc {

	m.rateLimiter.once.Do(func() {
		cfg := m.app.Config()

		lmcfg := ratelimiter.Config{
			Period:             cfg.RateLimiter.Period,
			MaxRequestsPerIP:   cfg.RateLimiter.MaxRequestsPerIP,
			MaxRequestsPerUser: cfg.RateLimiter.MaxRequestsPerUser,
		}

		if cfg.Cache.Type == "redis" {
			lmcfg.RedisConfig = &ratelimiter.RedisConfig{
				Addr:     cfg.Cache.Redis.Addr,
				DB:       cfg.Cache.Redis.DB,
				Username: cfg.Cache.Redis.Username,
				Password: cfg.Cache.Redis.Password,
			}
		}

		m.rateLimiter.handler = ratelimiter.New(lmcfg)
	})

	return m.rateLimiter.handler
}

func (m *Middlewares) Prometheus() gin.HandlerFunc {
	m.prometheus.once.Do(func() {
		cfg := m.app.Config()

		m.prometheus.handler = promx.New(func(c *promx.Config) {
			c.App = cfg.AppName
			c.ListenPort = cfg.Prometheus.Port
			c.BasicUserName = cfg.Prometheus.BasicUsername
			c.BasicPassword = cfg.Prometheus.BasicPassword
			c.LogApi = cfg.Prometheus.LogApis
			c.LogMethod = cfg.Prometheus.LogMethods
			c.DefaultCollect = cfg.Prometheus.DefaultCollect
			c.Objectives = map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001}
		})

	})

	return m.prometheus.handler
}
