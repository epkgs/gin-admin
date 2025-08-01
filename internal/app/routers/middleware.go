package routers

import (
	"context"
	"gin-admin/internal/services"
	"gin-admin/internal/types"
	"gin-admin/pkg/helper"
	mdw "gin-admin/pkg/middleware"
	"gin-admin/pkg/promx"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/epkgs/i18n"
	"github.com/gin-gonic/gin"
)

type middleware struct {
	excluded    map[string]struct{}
	excludedArr []string

	included    map[string]struct{}
	includedArr []string
}

var _ types.Middleware = (*middleware)(nil)

func (m *middleware) Exclude(prefixes ...string) {
	if m.excluded == nil {
		m.excluded = make(map[string]struct{})
	}
	if m.excludedArr == nil {
		m.excludedArr = make([]string, 0)
	}
	for _, pre := range prefixes {
		if _, exist := m.excluded[pre]; exist {
			continue
		}
		m.excluded[pre] = struct{}{}
		m.excludedArr = append(m.excludedArr, pre)
	}
}

func (m *middleware) GetExcluded() []string {
	return m.excludedArr
}

func (m *middleware) Include(prefixes ...string) {
	if m.included == nil {
		m.included = make(map[string]struct{})
	}
	if m.includedArr == nil {
		m.includedArr = make([]string, 0)
	}
	for _, pre := range prefixes {
		if _, exist := m.included[pre]; exist {
			continue
		}
		m.included[pre] = struct{}{}
		m.includedArr = append(m.includedArr, pre)
	}
}

func (m *middleware) GetIncluded() []string {
	return m.includedArr
}

type Middlewares struct {
	app types.AppContext

	auth        *middleware
	casbin      *middleware
	trace       *middleware
	logger      *middleware
	copyBody    *middleware
	rateLimiter *middleware
	static      *middleware
}

func NewMiddlewares(app types.AppContext) types.Middlewares {
	return &Middlewares{
		app: app,

		auth:        new(middleware),
		casbin:      new(middleware),
		trace:       new(middleware),
		logger:      new(middleware),
		copyBody:    new(middleware),
		rateLimiter: new(middleware),
		static:      new(middleware),
	}
}

func (m *Middlewares) Auth() types.Middleware {
	return m.auth
}

func (m *Middlewares) Casbin() types.Middleware {
	return m.casbin
}

func (m *Middlewares) Trace() types.Middleware {
	return m.trace
}

func (m *Middlewares) Logger() types.Middleware {
	return m.logger
}

func (m *Middlewares) CopyBody() types.Middleware {
	return m.copyBody
}

func (m *Middlewares) RateLimiter() types.Middleware {
	return m.rateLimiter
}

func (m *Middlewares) Static() types.Middleware {
	return m.static
}

func (m *Middlewares) Init(ctx context.Context, e *gin.Engine) error {

	cfg := m.app.Config()

	e.Use(i18n.GinMiddleware("zh"))

	e.Use(mdw.CORSWithConfig(mdw.CORSConfig{
		Enable:                 cfg.Middleware.CORS.Enable,
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
	}))

	e.Use(mdw.TraceWithConfig(mdw.TraceConfig{
		IncludedPathPrefixes: m.trace.GetIncluded(),
		ExcludedPathPrefixes: m.trace.GetExcluded(),
		RequestHeaderKey:     cfg.Middleware.Trace.RequestHeaderKey,
		ResponseTraceKey:     cfg.Middleware.Trace.ResponseTraceKey,
	}))

	e.Use(mdw.LoggerWithConfig(mdw.LoggerConfig{
		IncludedPathPrefixes:     m.logger.GetIncluded(),
		ExcludedPathPrefixes:     m.logger.GetExcluded(),
		MaxOutputRequestBodyLen:  cfg.Middleware.Logger.MaxOutputRequestBodyLen,
		MaxOutputResponseBodyLen: cfg.Middleware.Logger.MaxOutputResponseBodyLen,
	}))

	e.Use(mdw.CopyBodyWithConfig(mdw.CopyBodyConfig{
		IncludedPathPrefixes: m.copyBody.GetIncluded(),
		ExcludedPathPrefixes: m.copyBody.GetExcluded(),
		MaxContentLen:        cfg.Middleware.CopyBody.MaxContentLen,
	}))

	e.Use(mdw.AuthWithConfig(mdw.AuthConfig{
		IncludedPathPrefixes: m.auth.GetIncluded(),
		ExcludedPathPrefixes: m.auth.GetExcluded(),
		ParseUserID:          services.NewAuth(m.app).ParseUserID,
		RootID:               cfg.Super.ID,
	}))

	e.Use(mdw.RateLimiterWithConfig(mdw.RateLimiterConfig{
		Enable:               cfg.Middleware.RateLimiter.Enable,
		IncludedPathPrefixes: m.rateLimiter.GetIncluded(),
		ExcludedPathPrefixes: m.rateLimiter.GetExcluded(),
		Period:               cfg.Middleware.RateLimiter.Period,
		MaxRequestsPerIP:     cfg.Middleware.RateLimiter.MaxRequestsPerIP,
		MaxRequestsPerUser:   cfg.Middleware.RateLimiter.MaxRequestsPerUser,
		StoreType:            cfg.Middleware.RateLimiter.Store.Type,
		MemoryStoreConfig: mdw.RateLimiterMemoryConfig{
			Expiration:      time.Second * time.Duration(cfg.Middleware.RateLimiter.Store.Memory.Expiration),
			CleanupInterval: time.Second * time.Duration(cfg.Middleware.RateLimiter.Store.Memory.CleanupInterval),
		},
		RedisStoreConfig: mdw.RateLimiterRedisConfig{
			Addr:     cfg.Middleware.RateLimiter.Store.Redis.Addr,
			Password: cfg.Middleware.RateLimiter.Store.Redis.Password,
			DB:       cfg.Middleware.RateLimiter.Store.Redis.DB,
			Username: cfg.Middleware.RateLimiter.Store.Redis.Username,
		},
	}))

	e.Use(mdw.CasbinWithConfig(mdw.CasbinConfig{
		IncludedPathPrefixes: m.auth.GetIncluded(),
		ExcludedPathPrefixes: m.casbin.GetExcluded(),
		Skipper: func(c *gin.Context) bool {
			if cfg.Middleware.Casbin.Disable ||
				helper.GetIsRootUser(c.Request.Context()) {
				return true
			}
			return false
		},
		GetEnforcer: func(c *gin.Context) *casbin.Enforcer {
			return m.app.Casbin().GetEnforcer()
		},
		GetSubjects: func(c *gin.Context) []string {
			ctx := c.Request.Context()
			roleIDs, _ := services.NewUser(m.app).GetRoleIDsCache(ctx, helper.GetUserID(ctx))
			return roleIDs
		},
	}))

	if cfg.Prometheus.Enable {
		e.Use(promx.GinMiddleware(&promx.Config{
			Enable:         cfg.Prometheus.Enable,
			App:            cfg.AppName,
			ListenPort:     cfg.Prometheus.Port,
			BasicUserName:  cfg.Prometheus.BasicUsername,
			BasicPassword:  cfg.Prometheus.BasicPassword,
			LogApi:         cfg.Prometheus.LogApis,
			LogMethod:      cfg.Prometheus.LogMethods,
			Objectives:     map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
			DefaultCollect: cfg.Prometheus.DefaultCollect,
		}, helper.GetRequestBody))
	}

	return nil

}
