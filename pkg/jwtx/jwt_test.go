package jwtx

import (
	"context"
	"gin-admin/pkg/cachex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	cache := cachex.NewMemoryCache(cachex.MemoryConfig{CleanupInterval: time.Second})

	store := NewStoreWithCache(cache)
	ctx := context.Background()
	jwtAuth := New(store)

	userID := "test"
	token, err := jwtAuth.GenerateToken(ctx, userID)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	claims, err := jwtAuth.ParseToken(ctx, token.GetAccessToken())
	assert.Nil(t, err)
	id, err := claims.GetSubject()
	assert.Nil(t, err)
	assert.Equal(t, userID, id)

	err = jwtAuth.DestroyToken(ctx, token.GetAccessToken())
	assert.Nil(t, err)

	claims, err = jwtAuth.ParseToken(ctx, token.GetAccessToken())
	assert.Nil(t, err)
	assert.EqualError(t, err, ErrInvalidToken.Error())
	id, err = claims.GetSubject()
	assert.Nil(t, err)
	assert.Empty(t, id)

	err = jwtAuth.Release(ctx)
	assert.Nil(t, err)
}
