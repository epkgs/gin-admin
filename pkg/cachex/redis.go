package cachex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis 缓存实现
type redisCache struct {
	client *redis.Client
}

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func NewRedisCache(config *RedisConfig) (Cacher, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &redisCache{client: client}, nil
}

func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return data, nil
}

func (r *redisCache) GetObject(ctx context.Context, key string, value interface{}) error {
	data, err := r.Get(ctx, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, value)
}

func (r *redisCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisCache) SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.Set(ctx, key, data, expiration)
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *redisCache) ForEach(ns string, fn func(key string, raw []byte) bool) error {
	var cursor uint64
	var err error

	ctx := context.Background()
	pattern := ns + Delimiter + "*"

	for {
		var keys []string
		keys, cursor, err = r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			data, err := r.client.Get(ctx, key).Bytes()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					continue // 键可能已过期或被删除
				}
				return err
			}

			if !fn(key, data) {
				return nil // 用户回调要求中断迭代
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (r *redisCache) Close() error {
	return r.client.Close()
}
