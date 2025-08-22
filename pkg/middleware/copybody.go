package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"

	"gin-admin/internal/errorx"
	"gin-admin/locales"
	"gin-admin/pkg/helper"
	"gin-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type CopyBodyConfig struct {
	MaxContentLen int64
}

var DefaultCopyBodyConfig = CopyBodyConfig{
	MaxContentLen: 32 << 20, // 32MB
}

func CopyBody() gin.HandlerFunc {
	return CopyBodyWithConfig(DefaultCopyBodyConfig)
}

func CopyBodyWithConfig(config CopyBodyConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil {
			c.Next()
			return
		}

		var (
			requestBody []byte
			err         error
		)

		isGzip := false
		safe := http.MaxBytesReader(c.Writer, c.Request.Body, config.MaxContentLen)
		if c.GetHeader("Content-Encoding") == "gzip" {
			if reader, ierr := gzip.NewReader(safe); ierr == nil {
				isGzip = true
				requestBody, err = io.ReadAll(reader)
			}
		}

		if !isGzip {
			requestBody, err = io.ReadAll(safe)
		}

		if err != nil {
			response.Error(c, errorx.ErrRequestEntityTooLarge.WithMsg(locales.Def.Str("request body too large, limit %s byte", config.MaxContentLen)))
			return
		}

		c.Request.Body.Close()
		bf := bytes.NewBuffer(requestBody)
		c.Request.Body = io.NopCloser(bf)
		helper.SetRequestBody(c, requestBody)
		c.Next()
	}
}
