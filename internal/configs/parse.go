package configs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gin-admin/internal/errorx"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

var (
	once sync.Once
	C    = new(Config)
)

type Setter func(ctx context.Context, c *Config) error

func MustLoad(ctx context.Context, file string, setters ...Setter) {
	once.Do(func() {
		if err := Load(ctx, file, setters...); err != nil {
			panic(err)
		}
	})
}

// Loads configuration files in various formats from a directory and parses them into
// a struct.
func Load(ctx context.Context, path string, setters ...Setter) error {
	// Set default values
	if err := defaults.Set(C); err != nil {
		return err
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
	v.AddConfigPath("configs")
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
		return errorx.ErrReadConfigFile.New(ctx, struct{ File string }{path}).Wrap(err)
	}

	// Unmarshal the configuration into the struct
	if err := v.Unmarshal(C); err != nil {
		return errorx.ErrUnmarshalConfig.New(ctx, struct{ File string }{v.ConfigFileUsed()}).Wrap(err)
	}

	C.preLoad()
	for _, setter := range setters {
		if err := setter(ctx, C); err != nil {
			return err
		}
	}

	if C.PrintConfig {
		C.Print()
	}

	return nil
}

func splitFileName(fileName string) (name, ext string) {
	dotExt := filepath.Ext(fileName)
	name = fileName[:len(fileName)-len(dotExt)]
	if dotExt != "" {
		return name, dotExt[1:]
	}
	return name, ""
}
