package logger

import (
	"context"

	"go.uber.org/zap"
)

const (
	Tag_Main     = "main"
	Tag_Recovery = "recovery"
	Tag_Request  = "request"
	Tag_Login    = "login"
	Tag_Logout   = "logout"
	Tag_System   = "system"
	Tag_Operate  = "operate"
)

const (
	key_traceID = "traceId"
	key_userID  = "userId"
	key_tag     = "tag"
	key_stack   = "stack"
)

type (
	ctxLoggerKey struct{}
	ctxValuesKey struct{}
)

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, logger)
}

func GetLogger(ctx context.Context) *zap.Logger {
	v := ctx.Value(ctxLoggerKey{})
	if v != nil {
		if vv, ok := v.(*zap.Logger); ok {
			return vv
		}
	}
	return zap.L()
}

func GetValues(ctx context.Context) map[string]any {
	v := ctx.Value(ctxValuesKey{})
	if v != nil {
		if vv, ok := v.(map[string]any); ok {
			return vv
		}
	}
	return make(map[string]any)
}

func With(ctx context.Context, key string, value any) context.Context {
	m := GetValues(ctx)

	m[key] = value
	return context.WithValue(ctx, ctxValuesKey{}, m)
}

func WithValues(ctx context.Context, values map[string]any) context.Context {
	m := GetValues(ctx)

	for k, v := range values {
		m[k] = v
	}
	return context.WithValue(ctx, ctxValuesKey{}, m)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return With(ctx, key_traceID, traceID)
}

func GetTraceID(ctx context.Context) string {
	v := GetValues(ctx)
	if v != nil {
		if s, ok := v[key_traceID].(string); ok {
			return s
		}
	}
	return ""
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return With(ctx, key_userID, userID)
}

func GetUserID(ctx context.Context) string {
	v := GetValues(ctx)
	if v != nil {
		if s, ok := v[key_userID].(string); ok {
			return s
		}
	}
	return ""
}

func WithTag(ctx context.Context, tag string) context.Context {
	return With(ctx, key_tag, tag)
}

func GetTag(ctx context.Context) string {
	v := GetValues(ctx)
	if v != nil {
		if s, ok := v[key_tag].(string); ok {
			return s
		}
	}
	return ""
}

func WithStack(ctx context.Context, stack string) context.Context {
	return With(ctx, key_stack, stack)
}

func GetStack(ctx context.Context) string {
	v := GetValues(ctx)
	if v != nil {
		if s, ok := v[key_stack].(string); ok {
			return s
		}
	}
	return ""
}

func WithStackSkip(ctx context.Context, key string, skip int) context.Context {
	return With(ctx, key, zap.StackSkip(key, skip))
}

// func Context(ctx context.Context) *zap.Logger {
// 	var fields []zap.Field
// 	values := GetValues(ctx)
// 	for k, v := range values {
// 		fields = append(fields, zap.Any(k, v))
// 	}
// 	return GetLogger(ctx).With(fields...)
// }

// type PrintLogger struct{}

// func (a *PrintLogger) Printf(format string, args ...interface{}) {
// 	zap.L().Info(fmt.Sprintf(format, args...))
// }
