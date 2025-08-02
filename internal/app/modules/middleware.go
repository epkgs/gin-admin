package modules

import (
	"gin-admin/internal/services"
	"gin-admin/internal/types"
	"gin-admin/pkg/helper"
	"gin-admin/pkg/middleware"
	"gin-admin/pkg/promx"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/epkgs/i18n"
	"github.com/gin-gonic/gin"
)

type Middlewares struct {
	i18n        gin.HandlerFunc
	cors        gin.HandlerFunc
	trace       gin.HandlerFunc
	logger      gin.HandlerFunc
	copyBody    gin.HandlerFunc
	auth        gin.HandlerFunc
	rateLimiter gin.HandlerFunc
	casbin      gin.HandlerFunc
	prometheus  gin.HandlerFunc
}

func NewMiddlewares(app types.AppContext) *Middlewares {
	m := &Middlewares{}

	cfg := app.Config()

	m.i18n = i18n.GinMiddleware("zh")

	if cfg.Middleware.CORS.Enable {
		m.cors = middleware.CORSWithConfig(middleware.CORSConfig{
			AllowAllOrigins:        cfg.Middleware.CORS.AllowAllOrigins,
			AllowOrigins:           cfg.Middleware.CORS.AllowOrigins,
			AllowMethods:           cfg.Middleware.CORS.AllowMethods,
			AllowHeaders:           cfg.Middleware.CORS.AllowHeaders,
			AllowCredentials:       cfg.Middleware.CORS.AllowCredentials,
			ExposeHeaders:          cfg.Middleware.CORS.ExposeHeaders,
			MaxAge:                 cfg.Middleware.CORS.MaxAge,
			AllowWildcard:          cfg.Middleware.CORS.AllowWildcard,
			AllowBrowserExtensions: cfg.Middleware.CORS.AllowBrowserExtensions,
			AllowWebSockets:        cfg.Middleware.CORS.AllowWebSockets,
			AllowFiles:             cfg.Middleware.CORS.AllowFiles,
		})
	} else {
		m.cors = middleware.Empty()
	}

	m.trace = middleware.TraceWithConfig(middleware.TraceConfig{
		RequestHeaderKey: cfg.Middleware.Trace.RequestHeaderKey,
		ResponseTraceKey: cfg.Middleware.Trace.ResponseTraceKey,
	})

	m.logger = middleware.LoggerWithConfig(middleware.LoggerConfig{
		MaxOutputRequestBodyLen:  cfg.Middleware.Logger.MaxOutputRequestBodyLen,
		MaxOutputResponseBodyLen: cfg.Middleware.Logger.MaxOutputResponseBodyLen,
	})

	m.copyBody = middleware.CopyBodyWithConfig(middleware.CopyBodyConfig{
		MaxContentLen: cfg.Middleware.CopyBody.MaxContentLen,
	})

	m.auth = middleware.AuthWithConfig(middleware.AuthConfig{
		ParseUserID: services.NewAuth(app).ParseUserID,
		RootID:      cfg.Super.ID,
	})

	m.rateLimiter = middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Enable:             cfg.Middleware.RateLimiter.Enable,
		Period:             cfg.Middleware.RateLimiter.Period,
		MaxRequestsPerIP:   cfg.Middleware.RateLimiter.MaxRequestsPerIP,
		MaxRequestsPerUser: cfg.Middleware.RateLimiter.MaxRequestsPerUser,
		StoreType:          cfg.Middleware.RateLimiter.Store.Type,
		MemoryStoreConfig: middleware.RateLimiterMemoryConfig{
			Expiration:      time.Second * time.Duration(cfg.Middleware.RateLimiter.Store.Memory.Expiration),
			CleanupInterval: time.Second * time.Duration(cfg.Middleware.RateLimiter.Store.Memory.CleanupInterval),
		},
		RedisStoreConfig: middleware.RateLimiterRedisConfig{
			Addr:     cfg.Middleware.RateLimiter.Store.Redis.Addr,
			Password: cfg.Middleware.RateLimiter.Store.Redis.Password,
			DB:       cfg.Middleware.RateLimiter.Store.Redis.DB,
			Username: cfg.Middleware.RateLimiter.Store.Redis.Username,
		},
	})

	m.casbin = middleware.CasbinWithConfig(middleware.CasbinConfig{
		Skipper: func(c *gin.Context) bool {
			if cfg.Middleware.Casbin.Disable ||
				helper.GetIsRootUser(c.Request.Context()) {
				return true
			}
			return false
		},
		GetEnforcer: func(c *gin.Context) *casbin.Enforcer {
			return app.Casbin().GetEnforcer()
		},
		GetSubjects: func(c *gin.Context) []string {
			ctx := c.Request.Context()
			roleIDs, _ := services.NewUser(app).GetRoleIDsCache(ctx, helper.GetUserID(ctx))
			return roleIDs
		},
	})

	if cfg.Prometheus.Enable {
		m.prometheus = promx.GinMiddleware(&promx.Config{
			Enable:         cfg.Prometheus.Enable,
			App:            cfg.AppName,
			ListenPort:     cfg.Prometheus.Port,
			BasicUserName:  cfg.Prometheus.BasicUsername,
			BasicPassword:  cfg.Prometheus.BasicPassword,
			LogApi:         cfg.Prometheus.LogApis,
			LogMethod:      cfg.Prometheus.LogMethods,
			Objectives:     map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
			DefaultCollect: cfg.Prometheus.DefaultCollect,
		}, helper.GetRequestBody)
	} else {
		m.prometheus = middleware.Empty()
	}

	return m
}

func (m *Middlewares) I18n() gin.HandlerFunc { return m.i18n }

func (m *Middlewares) Cors() gin.HandlerFunc { return m.cors }

func (m *Middlewares) Trace() gin.HandlerFunc { return m.trace }

func (m *Middlewares) Logger() gin.HandlerFunc { return m.logger }

func (m *Middlewares) CopyBody() gin.HandlerFunc { return m.copyBody }

func (m *Middlewares) Auth() gin.HandlerFunc { return m.auth }

func (m *Middlewares) RateLimiter() gin.HandlerFunc { return m.rateLimiter }

func (m *Middlewares) Casbin() gin.HandlerFunc { return m.casbin }

func (m *Middlewares) Prometheus() gin.HandlerFunc { return m.prometheus }
