package promx

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type WithConfig func(c *Config)

func New(withConfigs ...WithConfig) gin.HandlerFunc {

	cfg := Config{
		ListenPort:    9100,
		BasicUserName: "admin",
		BasicPassword: "admin",
	}

	for _, withConfig := range withConfigs {
		withConfig(&cfg)
	}

	prom := newPrometheusWrapper(&cfg)

	return func(c *gin.Context) {

		start := time.Now()
		// 读取请求体内容
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 恢复请求体，以便后续处理可以读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		recvBytes := len(requestBody)

		c.Next()
		latency := float64(time.Since(start).Milliseconds())
		p := c.Request.URL.Path
		for _, param := range c.Params {
			p = strings.Replace(p, param.Value, ":"+param.Key, -1)
		}
		prom.Log(p, c.Request.Method, fmt.Sprintf("%d", c.Writer.Status()), float64(c.Writer.Size()), float64(recvBytes), latency)
	}

}
