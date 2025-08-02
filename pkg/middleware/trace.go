package middleware

import (
	"fmt"
	"strings"

	"gin-admin/pkg/helper"
	"gin-admin/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type TraceConfig struct {
	RequestHeaderKey string
	ResponseTraceKey string
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

		traceID := c.GetHeader(config.RequestHeaderKey)
		if traceID == "" {
			traceID = fmt.Sprintf("TRACE-%s", strings.ToUpper(xid.New().String()))
		}

		ctx := helper.WithTraceID(c.Request.Context(), traceID)
		ctx = logger.WithTraceID(ctx, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(config.ResponseTraceKey, traceID)
		c.Next()
	}
}
