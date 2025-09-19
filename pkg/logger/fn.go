package logger

import (
	"context"
	"log/slog"
)

type ctxLoggerKey = struct{}
type ctxValuesKey = struct{}

func Init(configs ...func(o *Config)) (func(), error) {
	logger, clear, err := New(configs...)
	if err != nil {
		return nil, err
	}

	slog.SetDefault(logger)

	return clear, nil
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, logger)
}

func getLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(ctxLoggerKey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

func WithAttrs(ctx context.Context, attrs ...any) context.Context {
	args, ok := ctx.Value(ctxValuesKey{}).([]any)
	if ok {
		args = append(args, attrs...)
	} else {
		args = attrs
	}

	return context.WithValue(ctx, ctxValuesKey{}, args)
}

func getAttrs(ctx context.Context) []any {
	attrs, ok := ctx.Value(ctxValuesKey{}).([]any)
	if ok {
		return attrs
	}
	return []any{}
}

func Debug(ctx context.Context, msg string, args ...any) {
	attrs := append(getAttrs(ctx), args...)
	getLogger(ctx).DebugContext(ctx, msg, attrs...)
}

func Enabled(ctx context.Context, level slog.Level) bool {
	return getLogger(ctx).Enabled(ctx, level)
}

func Error(ctx context.Context, msg string, args ...any) {
	attrs := append(getAttrs(ctx), args...)
	getLogger(ctx).ErrorContext(ctx, msg, attrs...)
}

func Info(ctx context.Context, msg string, args ...any) {
	attrs := append(getAttrs(ctx), args...)
	getLogger(ctx).InfoContext(ctx, msg, attrs...)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	attrs := append(getAttrs(ctx), args...)
	getLogger(ctx).Log(ctx, level, msg, attrs...)
}

func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	getLogger(ctx).LogAttrs(ctx, level, msg, attrs...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	attrs := append(getAttrs(ctx), args...)
	getLogger(ctx).WarnContext(ctx, msg, attrs...)
}

func WithGroup(ctx context.Context, name string) context.Context {
	attrs := getAttrs(ctx)
	logger := getLogger(ctx)
	ctx = context.WithValue(ctx, ctxValuesKey{}, []any{}) // 清空上下文的 attrs
	return WithLogger(ctx, logger.With(attrs...).WithGroup(name))
}

func Group(key string, args ...any) slog.Attr {
	return slog.Group(key, args...)
}
