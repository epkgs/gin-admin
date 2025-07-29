package middleware

import (
	"fmt"
	"strings"

	"gin-admin/pkg/helper"
	"gin-admin/pkg/logging"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type TraceConfig struct {
	IncludedPathPrefixes []string
	ExcludedPathPrefixes []string
	RequestHeaderKey     string
	ResponseTraceKey     string
}

var DefaultTraceConfig = TraceConfig{
	RequestHeaderKey: "X-Request-Id",
	ResponseTraceKey: "X-Trace-Id",
}

func Trace() gin.HandlerFunc {
	return TraceWithConfig(DefaultTraceConfig)
}

func TraceWithConfig(config TraceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IncludedPathPrefixes(c, config.IncludedPathPrefixes...) ||
			ExcludedPathPrefixes(c, config.ExcludedPathPrefixes...) {
			c.Next()
			return
		}

		traceID := c.GetHeader(config.RequestHeaderKey)
		if traceID == "" {
			traceID = fmt.Sprintf("TRACE-%s", strings.ToUpper(xid.New().String()))
		}

		ctx := helper.WithTraceID(c.Request.Context(), traceID)
		ctx = logging.WithTraceID(ctx, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(config.ResponseTraceKey, traceID)
		c.Next()
	}
}
