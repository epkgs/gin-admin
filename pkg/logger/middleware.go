package logger

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"gin-admin/pkg/geo"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
)

type ginWriter = gin.ResponseWriter

// responseWriter 是一个包装 gin.ResponseWriter 的结构体，用于捕获响应内容
type responseWriter struct {
	ginWriter
	body *bytes.Buffer
}

// Write 方法重写，用于捕获响应内容
func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ginWriter.Write(b)
}

// WriteString 方法重写，用于捕获响应内容
func (w responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ginWriter.WriteString(s)
}

// 记录请求和响应内容的中间件
func GinMiddleware(cfg Config) gin.HandlerFunc {

	return func(c *gin.Context) {
		// 保存原始的 ResponseWriter
		oldWriter := c.Writer

		// 创建新的 responseWriter 来捕获响应内容
		newWriter := &responseWriter{
			ginWriter: oldWriter,
			body:      &bytes.Buffer{},
		}
		c.Writer = newWriter

		contentType := c.Request.Header.Get("Content-Type")

		// 读取请求体内容
		var requestBody []byte
		if c.Request.Body != nil {
			if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
				if strings.Contains(contentType, "application/json") {
					requestBody, _ = io.ReadAll(c.Request.Body)
					// 恢复请求体，以便后续处理可以读取
					c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
				}
			}
		}

		// 记录开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		now := time.Now()

		// 计算处理时间
		duration := now.Sub(startTime)

		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		resSize := c.Writer.Size()

		// 构建日志字段
		var fields map[string]any
		{
			fields = map[string]any{
				"client_ip":      clientIP,
				"method":         c.Request.Method,
				"path":           c.Request.URL.Path,
				"user_agent":     userAgent,
				"referer":        c.Request.Referer(),
				"uri":            c.Request.RequestURI,
				"host":           c.Request.Host,
				"remote_addr":    c.Request.RemoteAddr,
				"proto":          c.Request.Proto,
				"content_length": c.Request.ContentLength,
				"content_type":   contentType,
				"duration":       duration,
				"status":         c.Writer.Status(),
				"res_time":       now.Format("2006-01-02 15:04:05.999"),
				"res_size":       resSize,
				"location":       geo.GetCityName(clientIP, "zh-CN"),
			}

			// 添加请求体（如果存在）
			if len(requestBody) > 0 && (cfg.MaxRequestLen > 0 && len(requestBody) <= cfg.MaxRequestLen) {
				fields["body"] = string(requestBody)
			}

			// 添加响应体（如果存在）
			responseBody := newWriter.body.Bytes()
			if len(responseBody) > 0 && (cfg.MaxResponseLen > 0 && len(responseBody) <= cfg.MaxResponseLen) {
				fields["res_body"] = string(responseBody)
			}

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

		// 记录日志
		ctx := c.Request.Context()
		ctx = WithTag(ctx, Tag_Request)
		Info(
			ctx,
			fmt.Sprintf("[HTTP] %s-%s-%d (%s)", c.Request.URL.Path, c.Request.Method, c.Writer.Status(), duration.String()),
			fields,
		)
	}
}
