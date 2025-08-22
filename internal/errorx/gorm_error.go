package errorx

import (
	"gin-admin/locales"

	"github.com/epkgs/i18n/errors"

	"gorm.io/gorm"
)

func WrapGormError(err error) errors.Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(errors.Error); ok {
		return e
	}

	switch err {
	case gorm.ErrRecordNotFound:
		return ErrNotFound.WithMsg(locales.Def.Str("record not found")).Wrap(err)
	case gorm.ErrInvalidTransaction:
		return ErrInternalServerError.WithMsg(locales.DB.Str("database transaction error")).Wrap(err)
	case gorm.ErrNotImplemented:
		return ErrInternalServerError.Wrap(err)
	case gorm.ErrMissingWhereClause:
		return ErrBadRequest.Wrap(err)
	case gorm.ErrUnsupportedRelation:
		return ErrBadRequest.Wrap(err)
	case gorm.ErrPrimaryKeyRequired:
		return ErrBadRequest.WithMsg(locales.Def.Str("invalid parameters")).Wrap(err)
	case gorm.ErrModelValueRequired, gorm.ErrModelAccessibleFieldsRequired, gorm.ErrSubQueryRequired, gorm.ErrInvalidData, gorm.ErrUnsupportedDriver, gorm.ErrRegistered, gorm.ErrInvalidField, gorm.ErrEmptySlice, gorm.ErrDryRunModeUnsupported, gorm.ErrInvalidDB, gorm.ErrInvalidValue, gorm.ErrInvalidValueOfLength, gorm.ErrPreloadNotAllowed, gorm.ErrDuplicatedKey, gorm.ErrForeignKeyViolated, gorm.ErrCheckConstraintViolated:
		return ErrInternalServerError.Wrap(err)
	}

	return ErrInternalServerError.Wrap(err)
}
