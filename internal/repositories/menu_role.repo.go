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

// // Query role menus from the database based on the provided parameters and options.
// func (a *MenuRole) Query(ctx context.Context, params models.RoleMenuQueryParam, opts ...models.RoleMenuQueryOptions) (*models.RoleMenuQueryResult, error) {
// 	var opt models.RoleMenuQueryOptions
// 	if len(opts) > 0 {
// 		opt = opts[0]
// 	}

// 	db := util.GetDB(ctx, a.DB).Model(new(models.MenuRole))
// 	if v := params.RoleIDs; len(v) > 0 {
// 		db = db.Where("role_id in (?)", v)
// 	}

// 	var list models.RoleMenus
// 	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	}

// 	queryResult := &models.RoleMenuQueryResult{
// 		PageResult: pageResult,
// 		Data:       list,
// 	}
// 	return queryResult, nil
// }

// // Get the specified role menu from the database.
// func (a *MenuRole) Get(ctx context.Context, id string, opts ...models.RoleMenuQueryOptions) (*models.MenuRole, error) {
// 	var opt models.RoleMenuQueryOptions
// 	if len(opts) > 0 {
// 		opt = opts[0]
// 	}

// 	item := new(models.MenuRole)
// 	ok, err := util.FindOne(ctx, util.GetDB(ctx, a.DB).Model(new(models.MenuRole)).Where("id=?", id), opt.QueryOptions, item)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	} else if !ok {
// 		return nil, nil
// 	}
// 	return item, nil
// }

// // Exist checks if the specified role menu exists in the database.
// func (a *MenuRole) Exists(ctx context.Context, id string) (bool, error) {
// 	ok, err := util.Exists(ctx, util.GetDB(ctx, a.DB).Model(new(models.MenuRole)).Where("id=?", id))
// 	return ok, errors.WithStack(err)
// }

// // Create a new role menu.
// func (a *MenuRole) Create(ctx context.Context, item ...*models.MenuRole) error {
// 	result := util.GetDB(ctx, a.DB).Create(item)
// 	return errors.WithStack(result.Error)
// }

// // Update the specified role menu in the database.
// func (a *MenuRole) Update(ctx context.Context, item *models.MenuRole, selectFields ...string) error {
// 	selected := []string{"*"}
// 	if len(selectFields) > 0 {
// 		selected = selectFields
// 	}
// 	result := util.GetDB(ctx, a.DB).Where("id=?", item.ID).Select(selected).Omit("created_at").Updates(item)
// 	return errors.WithStack(result.Error)
// }

// func (a *MenuRole) Save(ctx context.Context, items models.RoleMenus) error {
// 	result := util.GetDB(ctx, a.DB).Save(&items)
// 	return errors.WithStack(result.Error)
// }

// // Delete the specified role menu from the database.
// func (a *MenuRole) Delete(ctx context.Context, id ...string) error {
// 	result := util.GetDB(ctx, a.DB).Where("id IN (?)", id).Delete(new(models.MenuRole))
// 	return errors.WithStack(result.Error)
// }

// // Deletes role menus by role id.
// func (a *MenuRole) DeleteByRoleID(ctx context.Context, roleID ...string) error {
// 	result := util.GetDB(ctx, a.DB).Where("role_id IN (?)", roleID).Delete(new(models.MenuRole))
// 	return errors.WithStack(result.Error)
// }

// func (a *MenuRole) DeleteByRoleMenuIDs(ctx context.Context, roleIDs []string, menuIDs []string) error {
// 	result := util.GetDB(ctx, a.DB).Where("role_id IN (?) and menu_id IN (?)", roleIDs, menuIDs).Delete(new(models.MenuRole))
// 	return errors.WithStack(result.Error)
// }
