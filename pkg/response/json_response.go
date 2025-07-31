package response

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gin-admin/internal/dtos"
	"gin-admin/internal/errorx"
	"gin-admin/pkg/helper"
	"gin-admin/pkg/logger"
	"gin-admin/pkg/validatorx"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

const (
	Code_Success = 0
	Code_Fail    = 1
)

var null any = nil

func response[T any](c *gin.Context, res *dtos.Result[T]) {

	httpStatus := http.StatusOK

	// 如果还未设置 http status
	if !c.Writer.Written() && res.HttpStatus >= 100 {
		httpStatus = res.HttpStatus
	}

	buf, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	helper.SetResponseBody(c, buf)
	c.Data(httpStatus, "application/json; charset=utf-8", buf)
	c.Abort()
}

func OK(c *gin.Context, message ...string) {
	OkData(c, null, message...)
}

func OkData[T any](c *gin.Context, data T, message ...string) {
	msg := "ok"
	if len(message) > 0 {
		msg = message[0]
	}

	response(c, dtos.NewResult(Code_Success, msg, data))
}

func Error(c *gin.Context, err error) {

	ctx := c.Request.Context()

	var res *dtos.Result[any]

	// ============== validation error ==============

	var bindingErr binding.SliceValidationError
	var validationErr validator.ValidationErrors

	if errors.As(err, &bindingErr) {
		res = handleValidationErrors(ctx, bindingErr...)
	} else if errors.As(err, &validationErr) {
		res = handleValidationErrors(ctx, validationErr)
	} else {
		res = dtos.NewResult[any](errorx.Code(err), err.Error(), nil)
		res.HttpStatus = errorx.HttpStatus(err)
	}

	if res.HttpStatus <= 0 || res.HttpStatus == 200 {
		res.HttpStatus = http.StatusInternalServerError
	}

	if res.HttpStatus >= 500 {
		ctx = logger.WithTag(ctx, logger.Tag_System)
		ctx = logger.WithStack(ctx, fmt.Sprintf("%+v", err))
		logger.Error(ctx, http.StatusText(res.HttpStatus), err)
	}

	response(c, res)
}

func handleValidationErrors(ctx context.Context, errs ...error) *dtos.Result[any] {

	res := dtos.NewResult[any](http.StatusUnprocessableEntity, "form validation failed", nil)

	faileds := map[string][]string{}
	for _, err := range errs {
		translated := err.(validator.ValidationErrors).Translate(validatorx.TranslatorDetect(ctx))
		for field, msg := range translated {
			fieldName := field[strings.Index(field, ".")+1:]
			faileds[fieldName] = append(faileds[fieldName], msg)
		}
	}
	res.HttpStatus = http.StatusUnprocessableEntity
	res.Data = faileds

	return res
}

func List[T any](c *gin.Context, items []T, pager *dtos.Pager) {

	var pg dtos.Pager
	if pager != nil {
		pg = *pager
	}

	if items == nil {
		items = make([]T, 0) // 避免返回 null
	}

	if pg.Total == 0 && len(items) > 0 {
		pg.Total = int64(len(items))
	}

	res := dtos.NewResult(Code_Success, "ok", dtos.List[T]{
		Items: items,
		Pager: pg,
	})

	response(c, res)
}
