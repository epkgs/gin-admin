package modules

import (
	"fmt"
	"gin-admin/internal/errorx"
	"gin-admin/internal/service"
	"gin-admin/internal/types"
	"gin-admin/locales"
	"gin-admin/pkg/helper"
	"gin-admin/pkg/jwtx"
	"gin-admin/pkg/logger"
	"gin-admin/pkg/middleware/promx"
	"gin-admin/pkg/middleware/ratelimiter"
	"gin-admin/pkg/response"
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
	casbin      mdw
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

		m.logger.handler = logger.GinMiddleware(cfg)
	})

	return m.logger.handler
}

func (m *Middlewares) Auth() gin.HandlerFunc {

	m.auth.once.Do(func() {
		m.auth.handler = func(c *gin.Context) {

			ctx := c.Request.Context()

			var token string
			{

				auth := c.GetHeader("Authorization")
				prefix := "Bearer "

				if auth != "" && strings.HasPrefix(auth, prefix) {
					token = auth[len(prefix):]
				} else {
					token = auth
				}

				if token == "" {
					token = c.Query("token")
				}
			}

			if token == "" {
				response.Error(c, errorx.ErrUnauthorized.WithMsg(locales.User.Str("invalid token")))
				return
			}

			ctx = helper.WithToken(ctx, token)

			claims, err := m.app.Jwt().ParseToken(ctx, token)
			if err != nil {
				if err == jwtx.ErrInvalidToken {
					response.Error(c, errorx.ErrUnauthorized.WithMsg(locales.User.Str("invalid token")))
					return
				}
				response.Error(c, errorx.ErrInternalServerError.Wrap(err))
				return
			}

			userID := claims.UserID

			ctx = helper.WithUserID(ctx, userID)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		}
	})

	return m.auth.handler
}

func (m *Middlewares) RoutePermission() gin.HandlerFunc {
	m.casbin.once.Do(func() {
		cfg := m.app.Config()

		m.casbin.handler = func(c *gin.Context) {
			ctx := c.Request.Context()

			userID := helper.GetUserID(ctx)
			if cfg.IsSuper(userID) {
				c.Next()
				return
			}

			enforcer := m.app.Casbin().GetEnforcer()
			if enforcer == nil {
				response.Error(c, errorx.ErrForbidden)
				return
			}

			userSVC := service.NewUser(m.app)

			roleIDs, err := userSVC.GetRoleIDsCache(ctx, userID)
			if err != nil {
				response.Error(c, errorx.ErrForbidden.Wrap(err))
				return
			}

			for _, roleID := range roleIDs {
				if ok, err := enforcer.Enforce(roleID, c.Request.URL.Path, c.Request.Method); err != nil {
					response.Error(c, errorx.ErrInternalServerError.Wrap(err))
					return
				} else if ok {
					c.Next()
					return
				}
			}
			response.Error(c, errorx.ErrForbidden)
		}

	})

	return m.casbin.handler
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
