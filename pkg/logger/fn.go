package logger

import (
	"context"
	"fmt"
	"log/slog"
)

type ctxLoggerKey = struct{}

func Init(configs ...func(o *Config)) (func(), error) {
	logger, clear, err := New(configs...)
	if err != nil {
		return nil, err
	}

	slog.SetDefault(logger)

	return clear, nil
}

func Enabled(ctx context.Context, level slog.Level) bool {
	return getLogger(ctx).Enabled(ctx, level)
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, logger)
}

func getLogger(ctx context.Context) *slog.Logger {
	if val := ctx.Value(ctxLoggerKey{}); val != nil {
		if logger, ok := val.(*slog.Logger); ok && logger != nil {
			return logger
		}
	}
	return slog.Default()
}

func With(ctx context.Context, args ...any) context.Context {
	return WithLogger(ctx, getLogger(ctx).With(args...))
}

func WithGroup(ctx context.Context, name string) context.Context {
	return WithLogger(ctx, getLogger(ctx).WithGroup(name))
}

func Group(key string, args ...any) slog.Attr {
	return slog.Attr{Key: key, Value: slog.GroupValue(parseAttrs(args)...)}
}

func Debug(ctx context.Context, msg string, args ...any) {
	Log(ctx, slog.LevelDebug, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	Log(ctx, slog.LevelError, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	Log(ctx, slog.LevelInfo, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	Log(ctx, slog.LevelWarn, msg, args...)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	getLogger(ctx).LogAttrs(ctx, level, msg, parseAttrs(args)...)
}

func parseAttrs(args []any) []slog.Attr {
	if len(args) == 0 {
		return nil
	}

	attrs := []slog.Attr{}

	for i := 0; i < len(args); i++ {

		switch v := args[i].(type) {
		case slog.Attr:
			attrs = append(attrs, v)
			continue
		case []slog.Attr:
			attrs = append(attrs, v...)
			continue
		case string:
			if i+1 < len(args) {
				key := v
				value := args[i+1]
				attrs = append(attrs, slog.Any(key, value))
				i++
			}
			continue
		case map[string]any:
			for key, value := range v {
				attrs = append(attrs, slog.Any(key, value))
			}
			continue
		case error:
			attrs = append(attrs, slog.String("error", v.Error()))
			continue

		default:

			if i+1 < len(args) {
				key := fmt.Sprintf("%v", v)
				value := args[i+1]
				attrs = append(attrs, slog.Any(key, value))
				i++
			}

		}
	}

	return attrs
}
