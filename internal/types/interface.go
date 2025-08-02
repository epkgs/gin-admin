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

	Middlewares() Middlewares

	AddCleaner(ctx context.Context, cleaner func())
}

type Casbinx interface {
	GetEnforcer() *casbin.Enforcer
	Load(ctx context.Context) error
	Release(ctx context.Context) error
}

type Middlewares interface {
	I18n() gin.HandlerFunc
	Cors() gin.HandlerFunc
	Trace() gin.HandlerFunc
	Logger() gin.HandlerFunc
	CopyBody() gin.HandlerFunc
	Auth() gin.HandlerFunc
	RateLimiter() gin.HandlerFunc
	Casbin() gin.HandlerFunc
	Prometheus() gin.HandlerFunc
}
