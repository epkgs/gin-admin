package jwtx

import (
	"context"
	"gin-admin/pkg/cachex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	cache := cachex.NewMemoryCache(&cachex.MemoryConfig{CleanupInterval: time.Second})

	ctx := context.Background()
	jwtAuth := New(cache)

	userID := "test"
	token, err := jwtAuth.GenerateToken(ctx, userID)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	claims, err := jwtAuth.ParseToken(ctx, token.AccessToken)
	assert.Nil(t, err)
	id := claims.UserID
	assert.Equal(t, userID, id)

	err = jwtAuth.DestroyToken(ctx, token.AccessToken)
	assert.Nil(t, err)

	_, err = jwtAuth.ParseToken(ctx, token.AccessToken)
	assert.NotNil(t, err)
	assert.EqualError(t, err, ErrInvalidToken.Error())

	err = jwtAuth.Release(ctx)
	assert.Nil(t, err)
}
