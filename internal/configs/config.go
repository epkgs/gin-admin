package configs

import (
	"fmt"
	"strings"

	"gin-admin/pkg/encoding/json"
	"gin-admin/pkg/logger"
)

type Config struct {
	AppName         string `default:"gin-starter"`
	Version         string `default:"v1.0.0"`
	AppEnv          string `default:"dev"` // dev/debug/test/prod
	ConfigFile      string
	PrintConfig     bool
	DefaultLoginPwd string `default:"6351623c8cef86fefabfa7da046fc619"` // MD5(abc-123)
	Super           struct {
		ID       string `default:"super"`
		Username string `default:"super"`
		Password string
		NickName string `default:"Super Admin"`
	}

	HTTP struct {
		Addr            string `default:":8080"`
		ShutdownTimeout int    `default:"10"` // seconds
		ReadTimeout     int    `default:"60"` // seconds
		WriteTimeout    int    `default:"60"` // seconds
		IdleTimeout     int    `default:"10"` // seconds
		CertFile        string
		KeyFile         string
	}

	Cache      Cache
	DB         DB
	Upload     Upload
	Captcha    Captcha
	Prometheus Prometheus
	Swagger    Swagger
	Pprof      Pprof
	Menu       Menu

	Logger     logger.Config
	Middleware Middleware
}

func (c *Config) IsDebug() bool {
	mode := strings.ToLower(c.AppEnv)
	return mode == "dev" || mode == "debug" || mode == "test"
}

func (c *Config) String() string {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic("Failed to marshal config: " + err.Error())
	}
	return string(b)
}

func (c *Config) preLoad() {
	if addr := c.Cache.Redis.Addr; addr != "" {
		username := c.Cache.Redis.Username
		password := c.Cache.Redis.Password
		if c.Captcha.CacheType == "redis" &&
			c.Captcha.Redis.Addr == "" {
			c.Captcha.Redis.Addr = addr
			c.Captcha.Redis.Username = username
			c.Captcha.Redis.Password = password
		}
		if c.Middleware.RateLimiter.Store.Type == "redis" &&
			c.Middleware.RateLimiter.Store.Redis.Addr == "" {
			c.Middleware.RateLimiter.Store.Redis.Addr = addr
			c.Middleware.RateLimiter.Store.Redis.Username = username
			c.Middleware.RateLimiter.Store.Redis.Password = password
		}
		if c.Middleware.Auth.Store.Type == "redis" &&
			c.Middleware.Auth.Store.Redis.Addr == "" {
			c.Middleware.Auth.Store.Redis.Addr = addr
			c.Middleware.Auth.Store.Redis.Username = username
			c.Middleware.Auth.Store.Redis.Password = password
		}
	}
}

func (c *Config) Print() {
	fmt.Println("// ----------------------- Load configurations start ------------------------")
	fmt.Println(c.String())
	fmt.Println("// ----------------------- Load configurations end --------------------------")
}

func (c *Config) FormatTableName(name string) string {
	return c.DB.TablePrefix + name
}

func (c *Config) IsSuper(id string) bool {
	return c.Super.ID == id
}
