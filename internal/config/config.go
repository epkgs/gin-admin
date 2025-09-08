package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"gin-admin/pkg/encoding/json"
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

	Cache       Cache
	DB          DB
	Captcha     Captcha
	Prometheus  Prometheus
	Swagger     Swagger
	RateLimiter RateLimiter
	Pprof       Pprof
	Menu        Menu
	Jwt         Jwt
	Casbin      Casbin
	Logger      Logger
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

	c.RuntimePath = util.Must(filepath.Abs(c.RuntimePath))

	c.Cache.Badger.Path = c.GetRuntimePath(c.Cache.Badger.Path)
	c.Logger.File.Path = c.GetRuntimePath(c.Logger.File.Path)
	c.Casbin.GenPolicyFile = c.GetRuntimePath(c.Casbin.GenPolicyFile)
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
