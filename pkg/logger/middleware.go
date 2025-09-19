package logger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"

	"gin-admin/pkg/geo"
	"gin-admin/pkg/helper"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
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

// bufferPool 用于复用 bytes.Buffer 减少内存分配
var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

type GinMiddlewareConfig struct {
	MaxRequestLen  int `default:"4096"`
	MaxResponseLen int `default:"1024"`
	// 添加需要跳过日志记录的路径
	SkipPaths []string
}

// shouldSkip 检查是否应该跳过日志记录
func (cfg GinMiddlewareConfig) shouldSkip(path string) bool {
	for _, skipPath := range cfg.SkipPaths {
		if skipPath == path {
			return true
		}
	}
	return false
}

// isBufferNeeded 检查是否需要缓冲请求/响应体
func (cfg GinMiddlewareConfig) isBufferNeeded(method string, contentType string) bool {
	// 只对特定方法和内容类型进行缓冲
	if method != "POST" && method != "PUT" && method != "PATCH" {
		return false
	}
	return strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "application/xml") ||
		strings.Contains(contentType, "text/")
}

// 记录请求和响应内容的中间件
func GinMiddleware(cfg GinMiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否需要跳过日志记录
		if cfg.shouldSkip(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 获取缓冲区
		buffer := bufferPool.Get().(*bytes.Buffer)
		buffer.Reset()
		defer bufferPool.Put(buffer)

		// 保存原始的 ResponseWriter
		oldWriter := c.Writer

		// 创建新的 responseWriter 来捕获响应内容
		newWriter := &responseWriter{
			ginWriter: oldWriter,
			body:      buffer,
		}
		c.Writer = newWriter

		contentType := c.Request.Header.Get("Content-Type")

		// 读取请求体内容
		var requestBody []byte
		if c.Request.Body != nil && cfg.isBufferNeeded(c.Request.Method, contentType) && c.Request.ContentLength <= int64(cfg.MaxRequestLen) {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 恢复请求体，以便后续处理可以读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 记录开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)

		// 获取上下文
		ctx := c.Request.Context()

		resStatus := c.Writer.Status()

		// 构建基础日志属性
		ctx = WithAttrs(ctx,
			slog.String("tag", "request"),
			slog.String("trace_id", helper.GetTraceID(ctx)),
			slog.String("method", c.Request.Method),
			slog.String("uri", c.Request.RequestURI),
			slog.String("user_id", helper.GetUserID(ctx)),
			slog.String("referer", c.Request.Referer()),
			slog.String("client_ip", c.ClientIP()),
			slog.String("proto", c.Request.Proto),
			slog.String("content_type", contentType),
			slog.Duration("duration", duration),
			slog.Int("status", resStatus),
		)

		// 只有在有错误时才添加详细的用户代理信息
		if resStatus >= 400 {
			userAgent := c.Request.UserAgent()
			ua := user_agent.New(userAgent)
			brw, ver := ua.Browser()
			browser := fmt.Sprintf("%s %s", brw, ver)
			system := ua.OS()
			if system == "" {
				system = ua.Platform()
			}

			ctx = WithAttrs(ctx,
				slog.String("user_agent", userAgent),
				slog.String("browser", browser),
				slog.String("system", system),
				slog.String("location", geo.GetCityName(c.ClientIP(), "zh-CN")),
			)
		}

		parseAttrValue := func(contentType string, body []byte) any {
			if strings.Contains(contentType, "application/json") {
				data := map[string]any{}
				if err := jsoniter.Unmarshal(body, &data); err == nil {
					return data
				}
			}
			return body
		}

		// 添加请求体（如果存在且符合长度要求）
		if len(requestBody) > 0 {
			val := parseAttrValue(contentType, requestBody)
			ctx = WithAttrs(ctx, slog.Any("request", val))
		}

		// 添加响应体（如果存在且符合长度要求）
		responseBody := newWriter.body.Bytes()
		if len(responseBody) > 0 && len(responseBody) <= cfg.MaxResponseLen {
			contentType := c.Writer.Header().Get("Content-Type")
			val := parseAttrValue(contentType, responseBody)
			ctx = WithAttrs(ctx, slog.Any("response", val))
		}

		// 根据状态码选择日志级别
		message := fmt.Sprintf("%s %s [%d]",
			c.Request.Method,
			c.Request.URL.Path,
			resStatus,
		)

		switch {
		case resStatus >= 500:
			Error(ctx, message)
		case resStatus >= 400:
			Warn(ctx, message)
		case resStatus >= 300:
			Info(ctx, message)
		default:
			Info(ctx, message)
		}
	}
}
