package repositories

import (
	"context"

	"gin-admin/internal/models"
	"gin-admin/pkg/gormx"

	"gorm.io/gorm"
)

// Menu management for SYS
type Menu struct {
	gormx.Repository[models.Menu]
}

func NewMenu(db *gorm.DB) *Menu {
	return &Menu{
		Repository: gormx.NewGenericRepo[models.Menu](db),
	}
}

// GetByNameAndParentID get the specified menu from the database.
func (a *Menu) GetChildByName(ctx context.Context, parentID, name string, opts ...gormx.Option) (*models.Menu, error) {
	return a.Repository.First(ctx, gormx.WithWhere("name = ? and parent_id = ?", name, parentID), func(db *gorm.DB) *gorm.DB {
		return gormx.Apply(db, opts...)
	})
}

// Updates the status of all menus whose parent path starts with the provided parent path.
func (a *Menu) UpdateStatusByParentPath(ctx context.Context, parentPath, status string) error {
	menu := &models.Menu{
		Status: status,
	}
	return a.Repository.Update(ctx, menu, gormx.WithWhere("parent_path like ?", parentPath+"%"))
}

// Updates the parent path of the specified menu.
func (a *Menu) UpdateParentPath(ctx context.Context, id, parentPath string) error {
	menu := &models.Menu{
		ParentPath: parentPath,
	}
	return a.Repository.Update(ctx, menu, gormx.WithWhere("id=?", id))
}

func (a *Menu) DeleteChildrenOfButton(ctx context.Context, parentID string) error {

	if parentID == "" {
		return nil
	}

	return a.Repository.DeleteBatch(ctx, gormx.WithWhere("parent_id = ? AND type = ?", parentID, models.MenuType_BUTTON))
}

// // Get the specified menu from the database.
// func (a *Menu) Get(ctx context.Context, id string, opts ...models.MenuQueryOptions) (*models.Menu, error) {
// 	var opt models.MenuQueryOptions
// 	if len(opts) > 0 {
// 		opt = opts[0]
// 	}

// 	item := new(models.Menu)
// 	ok, err := util.FindOne(ctx, util.GetDB(ctx, a.DB).Model(new(models.Menu)).Where("id=?", id), opt.QueryOptions, item)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	} else if !ok {
// 		return nil, nil
// 	}
// 	return item, nil
// }

// func (a *Menu) GetChildByCode(ctx context.Context, parentID, code string, opts ...models.MenuQueryOptions) (*models.Menu, error) {
// 	var opt models.MenuQueryOptions
// 	if len(opts) > 0 {
// 		opt = opts[0]
// 	}

// 	item := new(models.Menu)
// 	ok, err := util.FindOne(ctx, util.GetDB(ctx, a.DB).Model(new(models.Menu)).Where("code=? AND parent_id=?", code, parentID), opt.QueryOptions, item)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	} else if !ok {
// 		return nil, nil
// 	}
// 	return item, nil
// }

// // Checks if the specified menu exists in the database.
// func (a *Menu) Exists(ctx context.Context, id string) (bool, error) {
// 	ok, err := util.Exists(ctx, util.GetDB(ctx, a.DB).Model(new(models.Menu)).Where("id=?", id))
// 	return ok, errors.WithStack(err)
// }

// // Checks if a menu with the specified `code` exists under the specified `parentID` in the database.
// func (a *Menu) ExistsChildByCode(ctx context.Context, parentID, code string) (bool, error) {
// 	ok, err := util.Exists(ctx, util.GetDB(ctx, a.DB).Model(new(models.Menu)).Where("code=? AND parent_id=?", code, parentID))
// 	return ok, errors.WithStack(err)
// }

// // Checks if a menu with the specified `name` exists under the specified `parentID` in the database.
// func (a *Menu) ExistsChildByName(ctx context.Context, parentID, name string) (bool, error) {
// 	ok, err := util.Exists(ctx, util.GetDB(ctx, a.DB).Model(new(models.Menu)).Where("name=? AND parent_id=?", name, parentID))
// 	return ok, errors.WithStack(err)
// }

// func (a *Menu) ExistsChildByMethodPath(ctx context.Context, parentID, method, path string) (bool, error) {
// 	ok, err := util.Exists(ctx, util.GetDB(ctx, a.DB).Model(new(models.Menu)).Where("method=? AND path=? AND parent_id=?", method, path, parentID))
// 	return ok, errors.WithStack(err)
// }

// // Create a new menu.
// func (a *Menu) Create(ctx context.Context, item *models.Menu) error {
// 	result := util.GetDB(ctx, a.DB).Create(item)
// 	return errors.WithStack(result.Error)
// }

// // Update the specified menu in the database.
// func (a *Menu) Update(ctx context.Context, item *models.Menu, selectFields ...string) error {

// 	selected := []string{"*"}
// 	if len(selectFields) > 0 {
// 		selected = selectFields
// 	}

// 	result := util.GetDB(ctx, a.DB).Where("id=?", item.ID).Select(selected).Omit("created_at").Updates(item)
// 	return errors.WithStack(result.Error)
// }

// // Delete the specified menu from the database.
// func (a *Menu) Delete(ctx context.Context, id string) error {
// 	result := util.GetDB(ctx, a.DB).Where("id=?", id).Delete(new(models.Menu))
// 	return errors.WithStack(result.Error)
// }

// func (a *Menu) DeleteChildren(ctx context.Context, parentID string) error {
// 	result := util.GetDB(ctx, a.DB).Where("parent_id = ?", parentID).Delete(new(models.Menu))
// 	return errors.WithStack(result.Error)
// }
