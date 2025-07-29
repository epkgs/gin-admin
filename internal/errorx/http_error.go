package errorx

import (
	"net/http"

	"github.com/pkg/errors"
)

const (
	CodeDefault       = 1
	HttpStatusDefault = http.StatusInternalServerError
)

// HttpError 是应用程序错误的基本结构
type HttpError struct {
	message    string
	code       int   // 错误码
	httpStatus int   // HTTP状态码
	cause      error // 原始错误
}

func NewHttpError(code int, message string, httpStatus int) *HttpError {
	return &HttpError{
		code:       code,
		message:    message,
		httpStatus: httpStatus,
	}
}

func (e *HttpError) Error() string {
	if e.message == "" {
		if e.cause != nil {
			return e.cause.Error()
		}
		return http.StatusText(e.httpStatus)
	} else {
		if e.cause != nil {
			return e.message + ": " + e.cause.Error()
		}
		return e.message
	}
}

func (e *HttpError) String() string {
	return e.message
}

// Code 获取错误码
func (e *HttpError) Code() int {
	return e.code
}

// HttpStatus 获取HTTP状态码
func (e *HttpError) HttpStatus() int {
	return e.httpStatus
}

func (e *HttpError) WithCode(code int) *HttpError {
	e.code = code
	return e
}

func (e *HttpError) WithHttpStatus(httpStatus int) *HttpError {
	e.httpStatus = httpStatus
	return e
}

func (e *HttpError) WithMessage(msg string) *HttpError {
	e.message = msg
	return e
}

// Is 检查错误是否是特定的应用程序错误
// 实现 errors.Is 接口
func (e *HttpError) Is(target error) bool {
	if target == nil {
		return false
	}

	// 如果目标错误也是 AppError，检查错误码
	if targetErr, ok := target.(*HttpError); ok {
		return e.code == targetErr.code
	}

	return false
}

// As implements the errors.As interface for HttpError.
// It attempts to match the target with the HttpError type or its cause chain.
//
// Parameters:
//   - target: a pointer to the target variable that should be set to the error value
//
// Returns:
//   - bool: true if the target was successfully matched and set, false otherwise
func (e *HttpError) As(target any) bool {
	if target == nil {
		return false
	}

	// 特定类型处理 - 最常见的情况
	if httpErr, ok := target.(**HttpError); ok {
		*httpErr = e
		return true
	}

	// 如果需要支持 error 接口类型
	if err, ok := target.(*error); ok {
		*err = e
		return true
	}

	// 如果有 cause，递归调用 errors.As
	if e.cause != nil {
		return errors.As(e.cause, target)
	}

	return false
}

// Wrap 添加原始错误
func (e *HttpError) Wrap(err error) error {
	e.cause = err
	return e.WithStack()
}

func (e *HttpError) WithStack() error {
	return errors.WithStack(e)
}

func Code(err error) int {
	if err == nil {
		return CodeDefault
	}

	var coder interface{ Code() int }
	if ok := errors.As(err, &coder); ok {
		return coder.Code()
	}

	return CodeDefault
}

func HttpStatus(err error) int {
	if err == nil {
		return HttpStatusDefault
	}

	var httpStatuser interface{ HttpStatus() int }
	if ok := errors.As(err, &httpStatuser); ok {
		return httpStatuser.HttpStatus()
	}

	return HttpStatusDefault
}

func TraceID(err error) string {
	if err == nil {
		return ""
	}

	var traceIDer interface{ TraceID() string }
	if ok := errors.As(err, &traceIDer); ok {
		return traceIDer.TraceID()
	}

	return ""
}
