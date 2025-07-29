package promx

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GinMiddleware(cfg *Config, getRequestBody func(c *gin.Context) []byte) gin.HandlerFunc {
	prom := NewPrometheusWrapper(cfg)

	return func(c *gin.Context) {
		if !cfg.Enable {
			c.Next()
			return
		}

		start := time.Now()
		recvBytes := 0
		if v := getRequestBody(c); v != nil {
			recvBytes = len(v)
		}
		c.Next()
		latency := float64(time.Since(start).Milliseconds())
		p := c.Request.URL.Path
		for _, param := range c.Params {
			p = strings.Replace(p, param.Value, ":"+param.Key, -1)
		}
		prom.Log(p, c.Request.Method, fmt.Sprintf("%d", c.Writer.Status()), float64(c.Writer.Size()), float64(recvBytes), latency)
	}

}
