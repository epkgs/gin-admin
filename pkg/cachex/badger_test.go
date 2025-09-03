package cachex

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadgerCache(t *testing.T) {
	assert := assert.New(t)

	cache, err := NewBadgerCache(&BadgerConfig{
		Path: "./tmp/badger",
	})

	assert.Nil(err)

	ctx := context.Background()
	err = cache.SetObject(ctx, BuildKey("tt", "foo"), "bar", 0)
	assert.Nil(err)

	var val string
	err = cache.GetObject(ctx, BuildKey("tt", "foo"), &val)
	assert.Nil(err)
	assert.Equal("bar", val)

	err = cache.Delete(ctx, BuildKey("tt", "foo"))
	assert.Nil(err)

	var val2 string
	err = cache.GetObject(ctx, BuildKey("tt", "foo"), &val2)
	assert.Equal(ErrNotFound, err)
	assert.Equal("", val2)

	tmap := make(map[string]bool)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("foo%d", i)
		err = cache.SetObject(ctx, BuildKey("tt", key), "bar", 0)
		assert.Nil(err)
		tmap[key] = true

		err = cache.SetObject(ctx, BuildKey("ff", key), "bar", 0)
		assert.Nil(err)
	}

	err = cache.ForEach("tt", func(key string, raw []byte) bool {
		assert.True(tmap[key])
		var value string

		if err := json.Unmarshal(raw, &value); err != nil {
			return false
		}

		assert.Equal("bar", value)
		t.Log(key, value)
		return true
	})
	assert.Nil(err)

	err = cache.Close()
	assert.Nil(err)
}
