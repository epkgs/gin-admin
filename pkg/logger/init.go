package logger

import (
	"context"
	"gin-admin/pkg/gormx"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Debug      bool
	Level      string // debug/info/warn/error/dpanic/panic/fatal
	CallerSkip int
	File       struct {
		Enable     bool
		Path       string
		MaxSize    int
		MaxBackups int
	}
	Database struct {
		Enable       bool
		Level        string
		Type         string // Database type: sqlite3/mysql/postgres
		DSN          string // Database connection string
		TablePrefix  string // Table prefix for database tables
		MaxBuffer    int
		MaxThread    int
		MaxOpenConns int // Maximum open connections
		MaxIdleConns int // Maximum idle connections
		MaxLifetime  int // Maximum connection lifetime in seconds
		MaxIdleTime  int // Maximum connection idle time in seconds
	}
}

type HookHandlerFunc func(ctx context.Context, cfg *Config) (*Hook, error)

func InitWithConfig(ctx context.Context, cfg *Config, hooks ...HookHandlerFunc) (func(), error) {
	var zconfig zap.Config
	if cfg.Debug {
		cfg.Level = "debug"
		zconfig = zap.NewDevelopmentConfig()
	} else {
		zconfig = zap.NewProductionConfig()
	}

	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	zconfig.Level.SetLevel(level)

	var (
		logger   *zap.Logger
		cleanFns []func()
	)

	if cfg.File.Enable {
		filename := cfg.File.Path
		_ = os.MkdirAll(filepath.Dir(filename), 0777)
		fileWriter := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    cfg.File.MaxSize,
			MaxBackups: cfg.File.MaxBackups,
			Compress:   false,
			LocalTime:  true,
		}

		cleanFns = append(cleanFns, func() {
			_ = fileWriter.Close()
		})

		zc := zapcore.NewCore(
			zapcore.NewJSONEncoder(zconfig.EncoderConfig),
			zapcore.AddSync(fileWriter),
			zconfig.Level,
		)
		logger = zap.New(zc)
	} else {
		ilogger, err := zconfig.Build()
		if err != nil {
			return nil, err
		}
		logger = ilogger
	}

	skip := cfg.CallerSkip
	if skip <= 0 {
		skip = 2
	}

	logger = logger.WithOptions(
		zap.WithCaller(true),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCallerSkip(skip),
	)

	handleHook := func(writer *Hook) *zap.Logger {
		cleanFns = append(cleanFns, func() {
			writer.Flush()
		})

		hookLevel := zap.NewAtomicLevel()
		hookLevel.SetLevel(level)

		hookEncoder := zap.NewProductionEncoderConfig()
		hookEncoder.EncodeTime = zapcore.EpochMillisTimeEncoder
		hookEncoder.EncodeDuration = zapcore.MillisDurationEncoder
		hookCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(hookEncoder),
			zapcore.AddSync(writer),
			hookLevel,
		)

		return logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, hookCore)
		}))
	}

	if cfg.Database.Enable {
		db, err := gormx.New(gormx.Config{
			Debug:        strings.ToUpper(cfg.Level) == "DEBUG",
			DBType:       cfg.Database.Type,
			DSN:          cfg.Database.DSN,
			MaxLifetime:  cfg.Database.MaxLifetime,
			MaxIdleTime:  cfg.Database.MaxIdleTime,
			MaxOpenConns: cfg.Database.MaxOpenConns,
			MaxIdleConns: cfg.Database.MaxIdleConns,
			TablePrefix:  cfg.Database.TablePrefix,
		})

		if err != nil {
			return nil, err
		}

		hook := NewHook(
			NewGormHook(db),
			SetHookMaxJobs(cfg.Database.MaxBuffer),
			SetHookMaxWorkers(cfg.Database.MaxThread),
		)

		logger = handleHook(hook)
	}

	for _, hook := range hooks {
		writer, err := hook(ctx, cfg)
		if err != nil {
			return nil, err
		} else if writer == nil {
			continue
		}

		logger = handleHook(writer)
	}

	zap.ReplaceGlobals(logger)
	return func() {
		for _, fn := range cleanFns {
			fn()
		}
	}, nil
}
