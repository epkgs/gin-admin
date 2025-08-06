package repositories

import (
	"context"

	"gin-admin/internal/models"
	"gin-admin/pkg/gormx"

	"gorm.io/gorm"
)

// Menu role permissions
type MenuRole struct {
	gormx.Repository[models.MenuRole]
}

func NewMenuRole(db *gorm.DB) *MenuRole {
	return &MenuRole{
		Repository: gormx.NewGenericRepo[models.MenuRole](db),
	}
}

// Deletes role menus by menu id.
func (a *MenuRole) DeleteByMenuID(ctx context.Context, menuID ...string) error {
	return a.Repository.DeleteBatch(ctx, gormx.WithWhere("menu_id IN (?)", menuID))
}
