package logging

import (
	"context"

	"go.uber.org/zap"
)

func logger(ctx context.Context) *zap.Logger {
	var fields []zap.Field
	values := GetValues(ctx)
	for k, v := range values {
		fields = append(fields, zap.Any(k, v))
	}
	return GetLogger(ctx).With(fields...)
}

func toFields(fields []map[string]any) []zap.Field {
	var zfields []zap.Field
	for _, m := range fields {
		for k, v := range m {

			if f, ok := v.(zap.Field); ok {
				zfields = append(zfields, f)
			} else {
				zfields = append(zfields, zap.Any(k, v))
			}
		}
	}
	return zfields
}

func Info(ctx context.Context, msg string, fields ...map[string]any) {
	logger(ctx).Info(msg, toFields(fields)...)
}

func Debug(ctx context.Context, msg string, fields ...map[string]any) {
	logger(ctx).Debug(msg, toFields(fields)...)
}

func Warn(ctx context.Context, msg string, fields ...map[string]any) {
	logger(ctx).Warn(msg, toFields(fields)...)
}

func Error(ctx context.Context, msg string, err error, fields ...map[string]any) {
	zfields := []zap.Field{zap.Error(err)}
	zfields = append(zfields, toFields(fields)...)
	logger(ctx).Error(msg, zfields...)
}

func DPanic(ctx context.Context, msg string, fields ...map[string]any) {
	logger(ctx).DPanic(msg, toFields(fields)...)
}

func Panic(ctx context.Context, msg string, fields ...map[string]any) {
	logger(ctx).Panic(msg, toFields(fields)...)
}

func Fatal(ctx context.Context, msg string, fields ...map[string]any) {
	logger(ctx).Fatal(msg, toFields(fields)...)
}

// Sync calls the underlying Core's Sync method, flushing any buffered log entries. Applications should take care to call Sync before exiting.
func Sync(ctx context.Context) error {
	return GetLogger(ctx).Sync()
}
