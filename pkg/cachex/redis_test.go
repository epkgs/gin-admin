package cachex

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisCache(t *testing.T) {
	assert := assert.New(t)

	cache := NewRedisCache(RedisConfig{
		Addr: "localhost:6379",
		DB:   1,
	})

	ctx := context.Background()
	err := cache.Set(ctx, "tt", "foo", "bar")
	assert.Nil(err)

	val, err := cache.Get(ctx, "tt", "foo")
	assert.Nil(err)
	assert.Equal("bar", val)

	err = cache.Delete(ctx, "tt", "foo")
	assert.Nil(err)

	val, err = cache.Get(ctx, "tt", "foo")
	assert.Equal(ErrNotFound, err)
	assert.Equal("", val)

	tmap := make(map[string]bool)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("foo%d", i)
		err = cache.Set(ctx, "tt", key, "bar")
		assert.Nil(err)
		tmap[key] = true

		err = cache.Set(ctx, "ff", key, "bar")
		assert.Nil(err)
	}

	err = cache.Iterator(ctx, "tt", func(ctx context.Context, key, value string) bool {
		assert.True(tmap[key])
		assert.Equal("bar", value)
		return true
	})
	assert.Nil(err)

	err = cache.Close(ctx)
	assert.Nil(err)
}
