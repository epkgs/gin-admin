package helper

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// keyRequestBody is the key used to store request body in context
	keyRequestBody  = "req-body"
	keyResponseBody = "res-body"
)

// Get access token from header or query parameter
func GetToken(c *gin.Context) string {
	var token string
	auth := c.GetHeader("Authorization")
	prefix := "Bearer "

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = auth
	}

	if token == "" {
		token = c.Query("accessToken")
	}

	return token
}

// Get refresh token from header or query parameter
func GetRefreshToken(c *gin.Context) string {
	var token string
	auth := c.GetHeader("Authorization")
	prefix := "Bearer "

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = auth
	}

	if token == "" {
		token = c.Query("refreshToken")
	}

	return token
}

// Get request body from context
func GetRequestBody(c *gin.Context) []byte {
	if v, ok := c.Get(keyRequestBody); ok {
		if b, ok := v.([]byte); ok {
			return b
		}
	}
	return nil
}

func SetRequestBody(c *gin.Context, body []byte) {
	c.Set(keyRequestBody, body)
}

func GetResponseBody(c *gin.Context) []byte {
	if v, ok := c.Get(keyResponseBody); ok {
		if b, ok := v.([]byte); ok {
			return b
		}
	}
	return nil
}

func SetResponseBody(c *gin.Context, body []byte) {
	c.Set(keyResponseBody, body)
}
