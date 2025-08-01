package types

import (
	"context"

	"gin-admin/internal/configs"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/jwtx"
	"gin-admin/pkg/uploader"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppContext interface {
	Config() *configs.Config
	DB() *gorm.DB
	Cacher() cachex.Cacher
	Jwt() jwtx.Auther
	Casbin() Casbinx
	Uploader() *uploader.Uploader

	Routers() Routers
	Middlewares() Middlewares

	AddCleaner(ctx context.Context, cleaner func())
}

type Casbinx interface {
	GetEnforcer() *casbin.Enforcer
	Load(ctx context.Context) error
	Release(ctx context.Context) error
}

type Initializer interface {
	Init(ctx context.Context) error
}

type RouteRegister func(ctx context.Context, g *gin.RouterGroup, e *gin.Engine) error

type Routers interface {
	ApiGroup(prefix string, register RouteRegister)
	Init(ctx context.Context, e *gin.Engine) error
}

type Middleware interface {
	Exclude(prefixes ...string)
	GetExcluded() []string
	Include(prefixes ...string)
	GetIncluded() []string
}

type Middlewares interface {
	Auth() Middleware
	Casbin() Middleware
	Trace() Middleware
	Logger() Middleware
	CopyBody() Middleware
	RateLimiter() Middleware
	Static() Middleware

	Init(ctx context.Context, e *gin.Engine) error
}
