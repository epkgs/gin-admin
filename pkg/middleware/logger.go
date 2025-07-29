package middleware

import (
	"fmt"
	"mime"
	"net/http"
	"time"

	"gin-admin/pkg/geo"
	"gin-admin/pkg/helper"
	"gin-admin/pkg/logging"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
)

type LoggerConfig struct {
	IncludedPathPrefixes     []string
	ExcludedPathPrefixes     []string
	MaxOutputRequestBodyLen  int
	MaxOutputResponseBodyLen int
}

var DefaultLoggerConfig = LoggerConfig{
	MaxOutputRequestBodyLen:  1024 * 1024,
	MaxOutputResponseBodyLen: 1024 * 1024,
}

// Record detailed request logs for quick troubleshooting.
func Logger() gin.HandlerFunc {
	return LoggerWithConfig(DefaultLoggerConfig)
}

func LoggerWithConfig(config LoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IncludedPathPrefixes(c, config.IncludedPathPrefixes...) ||
			ExcludedPathPrefixes(c, config.ExcludedPathPrefixes...) {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		contentType := c.Request.Header.Get("Content-Type")

		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		fields := map[string]any{
			"clientIP":      clientIP,
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"userAgent":     userAgent,
			"referer":       c.Request.Referer(),
			"uri":           c.Request.RequestURI,
			"host":          c.Request.Host,
			"remoteAddr":    c.Request.RemoteAddr,
			"proto":         c.Request.Proto,
			"contentLength": c.Request.ContentLength,
			"contentType":   contentType,
			"pragma":        c.Request.Header.Get("Pragma"),
		}

		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			mediaType, _, _ := mime.ParseMediaType(contentType)
			if mediaType == "application/json" {
				if v := helper.GetRequestBody(c); v != nil {
					if len(v) <= config.MaxOutputRequestBodyLen {
						fields["body"] = string(v)
					}
				}
			}
		}

		cost := time.Since(start).Nanoseconds() / 1e6
		fields["cost"] = cost
		fields["status"] = c.Writer.Status()
		fields["resTime"] = time.Now().Format("2006-01-02 15:04:05.999")
		fields["resSize"] = c.Writer.Size()

		if v := helper.GetResponseBody(c); v != nil {
			if len(v) <= config.MaxOutputResponseBodyLen {
				fields["resBody"] = string(v)
			}
		}

		{
			location := geo.GetCityName(clientIP, "zh-CN")
			fields["location"] = location
		}

		{
			ua := user_agent.New(userAgent)
			brw, ver := ua.Browser()
			browser := fmt.Sprintf("%s %s", brw, ver)
			fields["browser"] = browser

			system := ua.OS()
			if system == "" {
				system = ua.Platform()
			}
			fields["system"] = system
		}

		{
			ctx := c.Request.Context()
			ctx = logging.WithTag(ctx, logging.Tag_Request)
			logging.Info(
				ctx,
				fmt.Sprintf("[HTTP] %s-%s-%d (%dms)", c.Request.URL.Path, c.Request.Method, c.Writer.Status(), cost),
				fields,
			)
		}

	}
}
