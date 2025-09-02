package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"gin-admin/pkg/encoding/json"
	"gin-admin/pkg/logger"
	"gin-admin/pkg/utils/util"
)

type Config struct {
	AppName         string `default:"gin-starter"`
	Version         string `default:"v1.0.0"`
	AppEnv          string `default:"prod"` // dev/debug/test/prod
	ConfigFile      string
	PrintConfig     bool
	DefaultLoginPwd string `default:"6351623c8cef86fefabfa7da046fc619"` // MD5(abc-123)
	RuntimePath     string `default:"runtime"`
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

	c.RuntimePath = util.Must(filepath.Abs(c.RuntimePath))

	c.Cache.Badger.Path = c.GetRuntimePath(c.Cache.Badger.Path)
	c.Middleware.Auth.Store.Badger.Path = c.GetRuntimePath(c.Middleware.Auth.Store.Badger.Path)

	c.Logger.File.Path = c.GetRuntimePath(c.Logger.File.Path)

	c.Middleware.Casbin.GenPolicyFile = c.GetRuntimePath(c.Middleware.Casbin.GenPolicyFile)
}

func (c *Config) Print() {
	fmt.Println("// ----------------------- Load configurations start ------------------------")
	fmt.Println(c.String())
	fmt.Println("// ----------------------- Load configurations end --------------------------")
}

func (c *Config) GetRuntimePath(path string) string {

	if filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(c.RuntimePath, path)
}

func (c *Config) IsSuper(id string) bool {
	return c.Super.ID == id
}
