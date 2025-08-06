package repositories

import (
	"context"

	"gin-admin/internal/models"
	"gin-admin/pkg/gormx"

	"gorm.io/gorm"
)

// Role management for SYS
type Role struct {
	gormx.Repository[models.Role]
}

func NewRole(db *gorm.DB) *Role {
	return &Role{
		Repository: gormx.NewGenericRepo[models.Role](db),
	}
}

func (a *Role) ExistsCode(ctx context.Context, code string) (bool, error) {
	return a.Repository.Exists(ctx, gormx.WithWhere("code=?", code))
}
