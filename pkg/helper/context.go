package helper

import (
	"context"

	"gin-admin/pkg/encoding/json"

	"gorm.io/gorm"
)

type (
	traceIDCtx    struct{}
	transCtx      struct{}
	rowLockCtx    struct{}
	userIDCtx     struct{}
	userTokenCtx  struct{}
	isRootUserCtx struct{}
	userCacheCtx  struct{}
)

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDCtx{}, traceID)
}

func GetTraceID(ctx context.Context) string {
	v := ctx.Value(traceIDCtx{})
	if v != nil {
		return v.(string)
	}
	return ""
}

func WithTrans(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, transCtx{}, db)
}

func GetTrans(ctx context.Context) (*gorm.DB, bool) {
	v := ctx.Value(transCtx{})
	if v != nil {
		return v.(*gorm.DB), true
	}
	return nil, false
}

func WithRowLock(ctx context.Context) context.Context {
	return context.WithValue(ctx, rowLockCtx{}, true)
}

func GetRowLock(ctx context.Context) bool {
	v := ctx.Value(rowLockCtx{})
	return v != nil && v.(bool)
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDCtx{}, userID)
}

func GetUserID(ctx context.Context) string {
	v := ctx.Value(userIDCtx{})
	if v != nil {
		return v.(string)
	}
	return ""
}

func WithUserToken(ctx context.Context, userToken string) context.Context {
	return context.WithValue(ctx, userTokenCtx{}, userToken)
}

func GetUserToken(ctx context.Context) string {
	v := ctx.Value(userTokenCtx{})
	if v != nil {
		return v.(string)
	}
	return ""
}

func WithIsRootUser(ctx context.Context) context.Context {
	return context.WithValue(ctx, isRootUserCtx{}, true)
}

func GetIsRootUser(ctx context.Context) bool {
	v := ctx.Value(isRootUserCtx{})
	return v != nil && v.(bool)
}

// Set user cache object
type UserCache struct {
	RoleIDs []string `json:"roles"`
}

func ParseUserCache(s string) UserCache {
	var a UserCache
	if s == "" {
		return a
	}

	_ = json.Unmarshal([]byte(s), &a)
	return a
}

func (a UserCache) String() string {
	return json.MarshalToString(a)
}

func WithUserCache(ctx context.Context, userCache UserCache) context.Context {
	return context.WithValue(ctx, userCacheCtx{}, userCache)
}

func GetUserCache(ctx context.Context) UserCache {
	v := ctx.Value(userCacheCtx{})
	if v != nil {
		return v.(UserCache)
	}
	return UserCache{}
}
