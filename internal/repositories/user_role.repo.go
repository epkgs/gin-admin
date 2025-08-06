package repositories

import (
	"context"

	"gin-admin/internal/models"
	"gin-admin/pkg/gormx"

	"gorm.io/gorm"
)

// User roles
type UserRole struct {
	base[models.UserRole]
}

func NewUserRole(db *gorm.DB) *UserRole {
	return &UserRole{
		base: gormx.NewGenericRepo[models.UserRole](db),
	}
}

func (a *UserRole) DeleteByUserID(ctx context.Context, userID ...string) error {
	return a.base.DeleteBatch(ctx, gormx.WithWhere("user_id IN (?)", userID))
}
