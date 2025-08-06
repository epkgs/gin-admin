package repositories

import (
	"context"

	"gin-admin/internal/models"
	"gin-admin/pkg/gormx"

	"gorm.io/gorm"
)

// Menu management for SYS
type Menu struct {
	base[models.Menu]
}

func NewMenu(db *gorm.DB) *Menu {
	return &Menu{
		base: gormx.NewGenericRepo[models.Menu](db),
	}
}

// GetByNameAndParentID get the specified menu from the database.
func (a *Menu) GetChildByName(ctx context.Context, parentID, name string, opts ...gormx.Option) (*models.Menu, error) {
	return a.base.First(ctx, gormx.WithWhere("name = ? and parent_id = ?", name, parentID), func(db *gorm.DB) *gorm.DB {
		return gormx.Apply(db, opts...)
	})
}

// Updates the status of all menus whose parent path starts with the provided parent path.
func (a *Menu) UpdateStatusByParentPath(ctx context.Context, parentPath, status string) error {
	menu := &models.Menu{
		Status: status,
	}
	return a.base.Update(ctx, menu, gormx.WithWhere("parent_path like ?", parentPath+"%"))
}

// Updates the parent path of the specified menu.
func (a *Menu) UpdateParentPath(ctx context.Context, id, parentPath string) error {
	menu := &models.Menu{
		ParentPath: parentPath,
	}
	return a.base.Update(ctx, menu, gormx.WithWhere("id=?", id))
}

func (a *Menu) DeleteChildrenOfButton(ctx context.Context, parentID string) error {

	if parentID == "" {
		return nil
	}

	return a.base.DeleteBatch(ctx, gormx.WithWhere("parent_id = ? AND type = ?", parentID, models.MenuType_BUTTON))
}
