package cachex

import (
	"context"
	"errors"
	"strings"
	"time"
)

var ErrNotFound = errors.New("cache: not found or expired")

const Delimiter = ":"

// Cacher 接口定义
type Cacher interface {
	Get(ctx context.Context, key string) ([]byte, error)                                          // 获取缓存，返回原始数据
	GetObject(ctx context.Context, key string, value interface{}) error                           // 获取缓存，value 必须为引用
	Set(ctx context.Context, key string, raw []byte, expiration time.Duration) error              // 设置缓存，raw 必须为原始数据
	SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error // 设置缓存， expiration 为 0 表示不过期
	Delete(ctx context.Context, key string) error                                                 // 删除缓存
	Exists(ctx context.Context, key string) (bool, error)                                         // 判断缓存是否存在
	ForEach(namespace string, fn func(key string, raw []byte) bool) error                         // 根据 namespace 迭代缓存， fn 返回false则中断迭。
	Close() error                                                                                 // 关闭缓存
}

// BuildKey 构建缓存 key
// keys: 缓存 key
// 会自动过滤空字符串
func BuildKey(keys ...string) string {
	filtered := []string{}
	for _, key := range keys {
		if key != "" {
			filtered = append(filtered, key)
		}
	}
	return strings.Join(filtered, Delimiter)
}
