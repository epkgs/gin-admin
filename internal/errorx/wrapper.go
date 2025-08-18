package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
	"github.com/epkgs/i18n/errors"
	"gorm.io/gorm"
)

const (
	CodeFail    = 1
	CodeSuccess = 0
)

func Define[Args any](i18n *i18n.Bundle, code int, format string, httpStatus int) errors.Definition[Args] {
	return errors.Define[Args](i18n, format, func(e errors.I18nError) errors.I18nError {
		e.Set("code", code)
		e.Set("httpStatus", httpStatus)
		return e
	})
}

func New(i18n *i18n.Bundle, code int, format string, httpStatus int) errors.I18nError {
	e := errors.New(i18n.Sprintf(format))
	e.Set("code", code)
	e.Set("httpStatus", httpStatus)
	return e
}

func Code(err error) int {
	if err == nil {
		return CodeFail
	}

	var getter interface{ Get(key string) (any, bool) }
	if ok := errors.As(err, &getter); ok {
		if code, ok := getter.Get("code"); ok {
			if c, ok := code.(int); ok {
				return c
			}
		}
	}

	var coder interface{ Code() int }
	if ok := errors.As(err, &coder); ok {
		return coder.Code()
	}

	return CodeFail
}

func HttpStatus(err error) int {
	if err == nil {
		return http.StatusInternalServerError
	}

	var getter interface{ Get(key string) (any, bool) }
	if ok := errors.As(err, &getter); ok {
		if httpStatus, ok := getter.Get("httpStatus"); ok {
			if c, ok := httpStatus.(int); ok {
				return c
			}
		}
	}

	var httpStatuser interface{ HttpStatus() int }
	if ok := errors.As(err, &httpStatuser); ok {
		return httpStatuser.HttpStatus()
	}

	return http.StatusInternalServerError
}

func WrapGormError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(errors.I18nError); ok {
		return err
	}

	switch err {
	case gorm.ErrRecordNotFound:
		return General.RecordNotFound.Wrap(err)
	case gorm.ErrInvalidTransaction:
		return DB.Transaction.Wrap(err)
	case gorm.ErrNotImplemented:
		return General.Internal.Wrap(err)
	case gorm.ErrMissingWhereClause:
		return General.BadRequest.Wrap(err)
	case gorm.ErrUnsupportedRelation:
		return General.BadRequest.Wrap(err)
	case gorm.ErrPrimaryKeyRequired:
		return General.InvalidParams.New(struct{ Params string }{Params: "id"}).Wrap(err)
	case gorm.ErrModelValueRequired, gorm.ErrModelAccessibleFieldsRequired, gorm.ErrSubQueryRequired, gorm.ErrInvalidData, gorm.ErrUnsupportedDriver, gorm.ErrRegistered, gorm.ErrInvalidField, gorm.ErrEmptySlice, gorm.ErrDryRunModeUnsupported, gorm.ErrInvalidDB, gorm.ErrInvalidValue, gorm.ErrInvalidValueOfLength, gorm.ErrPreloadNotAllowed, gorm.ErrDuplicatedKey, gorm.ErrForeignKeyViolated, gorm.ErrCheckConstraintViolated:
		return General.Internal.Wrap(err)
	}

	return General.Internal.Wrap(err)
}
