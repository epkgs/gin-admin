package errorx

import (
	"context"

	"github.com/epkgs/i18n"
	i18nerr "github.com/epkgs/i18n/errorx"
	"gorm.io/gorm"
)

func httpError(code, httpStatus int) i18nerr.Wrapper[*HttpError] {
	return func(err *i18nerr.Error) *HttpError {
		return NewHttpError(code, err.Error(), httpStatus)
	}
}

func Definef[Args any](i18n *i18n.I18n, code int, format string, httpStatus int) *i18nerr.DefinitionF[*HttpError, Args] {
	return i18nerr.Definef[Args](i18n, format, httpError(code, httpStatus))
}

func Define(i18n *i18n.I18n, code int, format string, httpStatus int) *i18nerr.Definition[*HttpError] {
	return i18nerr.Define(i18n, format, httpError(code, httpStatus))
}

func WrapGormError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(*HttpError); ok {
		return err
	}

	switch err {
	case gorm.ErrRecordNotFound:
		return ErrRecordNotFound.New(ctx).Wrap(err)
	case gorm.ErrInvalidTransaction:
		return ErrDatabaseTransaction.New(ctx).Wrap(err)
	case gorm.ErrNotImplemented:
		return ErrDatabase.New(ctx).Wrap(err)
	case gorm.ErrMissingWhereClause:
		return ErrBadRequest.New(ctx).Wrap(err)
	case gorm.ErrUnsupportedRelation:
		return ErrBadRequest.New(ctx).Wrap(err)
	case gorm.ErrPrimaryKeyRequired:
		return ErrInvalidParams.New(ctx, struct{ Params string }{Params: "id"}).Wrap(err)
	case gorm.ErrModelValueRequired, gorm.ErrModelAccessibleFieldsRequired, gorm.ErrSubQueryRequired, gorm.ErrInvalidData, gorm.ErrUnsupportedDriver, gorm.ErrRegistered, gorm.ErrInvalidField, gorm.ErrEmptySlice, gorm.ErrDryRunModeUnsupported, gorm.ErrInvalidDB, gorm.ErrInvalidValue, gorm.ErrInvalidValueOfLength, gorm.ErrPreloadNotAllowed, gorm.ErrDuplicatedKey, gorm.ErrForeignKeyViolated, gorm.ErrCheckConstraintViolated:
		return ErrInternal.New(ctx).Wrap(err)
	}

	return ErrInternal.New(ctx).Wrap(err)
}
