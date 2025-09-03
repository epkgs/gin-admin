package cachex

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

// 内存缓存实现
type memoryCache struct {
	store *cache.Cache
}

type MemoryConfig struct {
	CleanupInterval time.Duration
}

func NewMemoryCache(config *MemoryConfig) Cacher {
	c := cache.New(0, config.CleanupInterval)
	return &memoryCache{
		store: c,
	}
}

func (m *memoryCache) Get(ctx context.Context, key string) (value []byte, err error) {
	if val, found := m.store.Get(key); found {
		if data, ok := val.([]byte); ok {
			return data, nil
		}
	}
	return nil, ErrNotFound
}

func (m *memoryCache) GetObject(ctx context.Context, key string, value interface{}) error {
	raw, err := m.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, value)
}

func (m *memoryCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	m.store.Set(key, value, expiration)
	return nil
}

func (m *memoryCache) SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return m.Set(ctx, key, data, expiration)
}

func (m *memoryCache) Delete(ctx context.Context, key string) error {
	m.store.Delete(key)
	return nil
}

func (m *memoryCache) Exists(ctx context.Context, key string) (bool, error) {
	_, found := m.store.Get(key)
	return found, nil
}

func (m *memoryCache) ForEach(namespace string, fn func(key string, raw []byte) bool) error {
	items := m.store.Items()
	for key, item := range items {
		if namespace != "" && strings.HasPrefix(key, namespace+Delimiter) {
			if data, ok := item.Object.([]byte); ok {
				subKey := key[len(namespace)+len(Delimiter):]
				if ok := fn(subKey, data); !ok {
					return nil
				}
			}
		} else if namespace == "" {
			if data, ok := item.Object.([]byte); ok {
				if ok := fn(key, data); !ok {
					return nil
				}
			}
		}
	}
	return nil
}

func (m *memoryCache) Close() error {
	// go-cache 会自动清理过期项目，这里不需要手动清理
	return nil
}
