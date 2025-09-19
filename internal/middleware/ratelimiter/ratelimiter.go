package ratelimiter

import (
	"context"
	"time"

	"gin-admin/internal/errorx"
	"gin-admin/pkg/helper"
	"gin-admin/pkg/logger"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

type Config struct {
	Period             int
	MaxRequestsPerIP   int
	MaxRequestsPerUser int
	RedisConfig        *RedisConfig
}

const prefix = "gin-admin-ratelimiter:"

func New(config Config) gin.HandlerFunc {

	var store Storer
	if config.RedisConfig != nil {
		store = newRedisStore(*config.RedisConfig)
	} else {
		store = newMemoryStore(time.Second*3600, time.Second*30)
	}

	return func(c *gin.Context) {

		var (
			allowed bool
			err     error
		)

		ctx := c.Request.Context()
		if userID := helper.GetUserID(ctx); userID != "" {
			allowed, err = store.Allow(ctx, userID, time.Second*time.Duration(config.Period), config.MaxRequestsPerUser)
		} else {
			allowed, err = store.Allow(ctx, c.ClientIP(), time.Second*time.Duration(config.Period), config.MaxRequestsPerIP)
		}

		if err != nil {
			logger.Error(ctx, "Rate limiter middleware error",
				"error", err,
			)
			response.Error(c, errorx.ErrInternalServerError)
		} else if allowed {
			c.Next()
		} else {
			response.Error(c, errorx.ErrTooManyRequests)
		}
	}
}

type Storer interface {
	Allow(ctx context.Context, identifier string, period time.Duration, maxRequests int) (bool, error)
}

func newMemoryStore(expiration, cleanupInterval time.Duration) Storer {
	return &memoryStore{
		cache: cache.New(expiration, cleanupInterval),
	}
}

type memoryStore struct {
	cache *cache.Cache
}

func (s *memoryStore) Allow(ctx context.Context, identifier string, period time.Duration, maxRequests int) (bool, error) {
	if period.Seconds() <= 0 || maxRequests <= 0 {
		return true, nil
	}

	if limiter, exists := s.cache.Get(identifier); exists {
		isAllow := limiter.(*rate.Limiter).Allow()
		s.cache.SetDefault(identifier, limiter)
		return isAllow, nil
	}

	limiter := rate.NewLimiter(rate.Every(period), maxRequests)
	limiter.Allow()
	s.cache.SetDefault(identifier, limiter)

	return true, nil
}

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func newRedisStore(config RedisConfig) Storer {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	})

	return &redisStore{
		limiter: redis_rate.NewLimiter(rdb),
	}
}

type redisStore struct {
	limiter *redis_rate.Limiter
}

func (s *redisStore) Allow(ctx context.Context, identifier string, period time.Duration, maxRequests int) (bool, error) {
	if period.Seconds() <= 0 || maxRequests <= 0 {
		return true, nil
	}

	result, err := s.limiter.Allow(ctx, prefix+identifier, redis_rate.PerSecond(maxRequests/int(period.Seconds())))
	if err != nil {
		return false, err
	}
	return result.Allowed > 0, nil
}
