package response

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gin-admin/internal/dto"
	"gin-admin/pkg/logger"
	"gin-admin/pkg/validatorx"

	"github.com/epkgs/i18n/errors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

const (
	Code_Success = 0
	Code_Fail    = 1
)

var null any = nil

func response[T any](c *gin.Context, res *dto.Result[T]) {

	httpStatus := http.StatusOK

	// 如果还未设置 http status
	if !c.Writer.Written() && res.HttpStatus >= 100 {
		httpStatus = res.HttpStatus
	}

	buf, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

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

	response(c, dto.NewResult(Code_Success, msg, data))
}

func Error(c *gin.Context, err error) {

	ctx := c.Request.Context()

	var res *dto.Result[any]

	// ============== validation error ==============

	var bindingErr binding.SliceValidationError
	var validationErr validator.ValidationErrors

	if errors.As(err, &bindingErr) {
		res = handleValidationErrors(ctx, bindingErr...)
	} else if errors.As(err, &validationErr) {
		res = handleValidationErrors(ctx, validationErr)
	} else {

		var msg string
		if tran, ok := err.(interface {
			T(ctx context.Context) string
		}); ok {
			msg = tran.T(ctx)
		} else {
			msg = err.Error()
		}

		res = dto.NewResult[any](errors.Code(err), msg, nil)
		res.HttpStatus = errors.HttpStatus(err)
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

func handleValidationErrors(ctx context.Context, errs ...error) *dto.Result[any] {

	res := dto.NewResult[any](http.StatusUnprocessableEntity, "form validation failed", nil)

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

func List[T any](c *gin.Context, items []T, pager *dto.Pager) {

	if items == nil {
		items = make([]T, 0) // 避免返回 null
	}

	res := dto.NewResult(Code_Success, "ok", dto.NewList(items, pager))

	response(c, res)
}
