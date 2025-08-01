package cachex

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryCache(t *testing.T) {
	assert := assert.New(t)

	cache := NewMemoryCache(MemoryConfig{
		CleanupInterval: time.Second * 30,
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
		t.Log(key, value)
		return true
	})
	assert.Nil(err)

	err = cache.Close(ctx)
	assert.Nil(err)
}
