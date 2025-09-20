package recovery

import (
	"context"
	"errors"
	"fmt"
	"net/http/httputil"
	"strings"
	"time"

	"gin-admin/internal/errorx"
	"gin-admin/pkg/logger"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Skip int
}

func New(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx := c.Request.Context()

				if e, ok := err.(error); ok && errors.Is(e, context.Canceled) {
					logger.Info(
						ctx,
						fmt.Sprintf("%v", err),
						"tag", "request",
					)
					return
				}

				ctx = logger.With(ctx,
					"tag", "recovery",
					"error", err,
				)

				if gin.IsDebugging() {
					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					headers := strings.Split(string(httpRequest), "\r\n")
					for idx, header := range headers {
						current := strings.Split(header, ":")
						if current[0] == "Authorization" {
							headers[idx] = current[0] + ": *"
						}
					}
					ctx = logger.With(ctx, "headers", headers)
				}

				logger.Error(
					ctx,
					fmt.Sprintf("[Recovery] %s panic recovered", time.Now().Format("2006/01/02 - 15:04:05")),
				)

				response.Error(c, errorx.ErrInternalServerError)
			}
		}()

		c.Next()
	}
}
