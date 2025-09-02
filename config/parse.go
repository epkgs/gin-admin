package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gin-admin/errorx"
	"gin-admin/locales"
	"gin-admin/pkg/utils/util"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

type Setter func(ctx context.Context, c *Config) error

func MustLoad(ctx context.Context, file string, setters ...Setter) *Config {
	return util.Must(Load(ctx, file, setters...))
}

// Loads configuration files in various formats from a directory and parses them into
// a struct.
func Load(ctx context.Context, path string, setters ...Setter) (*Config, error) {

	cfg := new(Config)

	// Set default values
	if err := defaults.Set(cfg); err != nil {
		return nil, err
	}

	// Create a new viper instance
	v := viper.New()

	if prefix := v.GetString("ENV_PREFIX"); prefix != "" {
		v.SetEnvPrefix(prefix)
	}

	// Replace dots in config keys with underscores for environment variables
	// E.g. logger.level -> LOGGER_LEVEL
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// Enable reading configuration from environment variables
	v.AutomaticEnv()

	mode := v.GetString("APP_ENV")
	if mode == "" {
		mode = "dev"
	}

	if path == "" {
		if path = v.GetString("CONFIG_FILE"); path == "" {
			path = "config.yaml"
		}
	}

	dir, file := filepath.Split(path)

	v.AddConfigPath(".")
	v.AddConfigPath("config")
	if dir != "" {
		v.AddConfigPath(dir)
	}

	fileName, fileExt := splitFileName(file)
	if fileExt != "" {
		v.SetConfigType(fileExt)
	}

	readConfig := func(name string) (bool, error) {
		v.SetConfigName(name)

		// Try to read the config file
		if err := v.ReadInConfig(); err != nil {
			return false, err
		}

		// 获取相对于程序启动目录的相对路径
		execPath, _ := os.Executable()
		execDir := filepath.Dir(execPath)
		if relFile, err := filepath.Rel(execDir, v.ConfigFileUsed()); err == nil {
			v.SetDefault("ConfigFile", relFile)
		}
		return true, nil
	}

	var ok bool
	var err error
	if mode != "" {
		// Try to read the config file with the specified mode first
		ok, err = readConfig(fmt.Sprintf("%s.%s", fileName, mode))
	}
	if !ok {
		// If that fails, try the default file name
		ok, err = readConfig(fileName)
	}

	if !ok || err != nil {
		return nil, errorx.ErrInternalServerError.WithMsg(locales.Def.Str("failed to read config file: %s", path)).Wrap(err)
	}

	// Unmarshal the configuration into the struct
	if err := v.Unmarshal(cfg); err != nil {
		return nil, errorx.ErrInternalServerError.WithMsg(locales.Def.Str("failed to unmarshal config: %s", v.ConfigFileUsed())).Wrap(err)
	}

	cfg.preLoad()
	for _, setter := range setters {
		if err := setter(ctx, cfg); err != nil {
			return nil, err
		}
	}

	if cfg.PrintConfig {
		cfg.Print()
	}

	return cfg, nil
}

func splitFileName(fileName string) (name, ext string) {
	dotExt := filepath.Ext(fileName)
	name = fileName[:len(fileName)-len(dotExt)]
	if dotExt != "" {
		return name, dotExt[1:]
	}
	return name, ""
}
