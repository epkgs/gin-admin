package repositories

import (
	"context"

	"gin-admin/internal/models"
	"gin-admin/pkg/gormx"

	"gorm.io/gorm"
)

// Role management for SYS
type Role struct {
	base[models.Role]
}

func NewRole(db *gorm.DB) *Role {
	return &Role{
		base: gormx.NewGenericRepo[models.Role](db),
	}
}

func (a *Role) ExistsCode(ctx context.Context, code string) (bool, error) {
	return a.base.Exists(ctx, gormx.WithWhere("code=?", code))
}
