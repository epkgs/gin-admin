package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http/httputil"
	"strings"
	"time"

	"gin-admin/internal/errorx"
	"gin-admin/pkg/logging"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RecoveryConfig struct {
	Skip int // default: 3
}

var DefaultRecoveryConfig = RecoveryConfig{
	Skip: 3,
}

// Recovery from any panics and writes a 500 if there was one.
func Recovery() gin.HandlerFunc {
	return RecoveryWithConfig(DefaultRecoveryConfig)
}

func RecoveryWithConfig(config RecoveryConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx := c.Request.Context()

				if e, ok := err.(error); ok && errors.Is(e, context.Canceled) {
					ctx = logging.WithTag(ctx, logging.Tag_Request)
					logging.Info(
						ctx,
						fmt.Sprintf("%v", err),
					)
					return
				}

				ctx = logging.WithTag(ctx, logging.Tag_Recovery)

				values := map[string]any{
					"stack": zap.StackSkip("stack", config.Skip),
				}

				if gin.IsDebugging() {
					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					headers := strings.Split(string(httpRequest), "\r\n")
					for idx, header := range headers {
						current := strings.Split(header, ":")
						if current[0] == "Authorization" {
							headers[idx] = current[0] + ": *"
						}
					}

					values["headers"] = headers
				}

				logging.Error(
					ctx,
					fmt.Sprintf("[Recovery] %s panic recovered", time.Now().Format("2006/01/02 - 15:04:05")),
					fmt.Errorf("%v", err),
					values,
				)

				response.Error(c, errorx.ErrInternal.New(ctx))
			}
		}()

		c.Next()
	}
}
