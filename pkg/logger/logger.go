package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
)

type Config struct {
	Level      string // debug/info/warn/error
	CallerSkip int
	// Console
	ConsoleEnable bool
	// File
	FileName       string // if file name not empty, will write log to file
	FileMaxSize    int    // Maximum number of backup log files
	FileMaxBackups int    // Maximum size of each log file in MB

	Database *gorm.DB // if database not nil, will write log to database
}

func New(configs ...func(o *Config)) (logger *slog.Logger, clear func(), err error) {

	cfg := Config{
		Level:          "info",
		CallerSkip:     1,
		ConsoleEnable:  true,
		FileMaxSize:    64,
		FileMaxBackups: 20,
	}

	for _, fn := range configs {
		fn(&cfg)
	}

	var fileWriter *lumberjack.Logger
	writers := []io.Writer{}

	clear = func() {
		if fileWriter != nil {
			fileWriter.Close()
			fileWriter = nil
		}
	}

	// add console output
	if cfg.ConsoleEnable {
		writers = append(writers, os.Stdout)
	}

	// add file output
	if cfg.FileName != "" {
		// create directory with permission 0750
		err = os.MkdirAll(filepath.Dir(cfg.FileName), 0750)
		if err != nil {
			// handle error without panic
			slog.Error("Failed to create log folder", "error", err, "dir", filepath.Dir(cfg.FileName))
		} else {

			fileWriter = &lumberjack.Logger{
				Filename:   cfg.FileName,
				MaxSize:    cfg.FileMaxSize,
				MaxBackups: cfg.FileMaxBackups,
				Compress:   false,
				LocalTime:  true,
			}

			writers = append(writers, fileWriter)
		}
	}

	// add database output
	if cfg.Database != nil {
		writers = append(writers, newDBWriter(cfg.Database))
	}

	level := getLevelFromString(cfg.Level)

	handler := slog.NewJSONHandler(
		io.MultiWriter(writers...),
		&slog.HandlerOptions{Level: level},
	)

	logger = slog.New(newSourceHandler(handler, cfg.CallerSkip))

	return
}

// getLevelFromString change string to slog.Level
func getLevelFromString(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
