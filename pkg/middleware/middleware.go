package middleware

import (
	"github.com/gin-gonic/gin"
)

func ExcludedPathPrefixes(c *gin.Context, prefixes ...string) bool {
	if len(prefixes) == 0 {
		return false
	}

	path := c.Request.URL.Path
	pathLen := len(path)
	for _, p := range prefixes {
		if pl := len(p); pathLen >= pl && path[:pl] == p {
			return true
		}
	}
	return false
}

func IncludedPathPrefixes(c *gin.Context, prefixes ...string) bool {
	if len(prefixes) == 0 {
		return true
	}

	path := c.Request.URL.Path
	pathLen := len(path)
	for _, p := range prefixes {
		if pl := len(p); pathLen >= pl && path[:pl] == p {
			return true
		}
	}
	return false
}

func Empty() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
