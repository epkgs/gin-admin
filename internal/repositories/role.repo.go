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

// // List roles from the database based on the provided parameters and options.
// func (a *Role) List(ctx context.Context, params models.RoleQueryParam, opts ...models.RoleQueryOptions) (*models.RoleQueryResult, error) {
// 	var opt models.RoleQueryOptions
// 	if len(opts) > 0 {
// 		opt = opts[0]
// 	}

// 	db := util.GetDB(ctx, a.DB).Model(new(models.Role))
// 	if v := params.InIDs; len(v) > 0 {
// 		db = db.Where("id IN (?)", v)
// 	}
// 	if v := params.LikeName; len(v) > 0 {
// 		db = db.Where("name LIKE ?", "%"+v+"%")
// 	}
// 	if v := params.Status; len(v) > 0 {
// 		db = db.Where("status = ?", v)
// 	}
// 	if v := params.GtUpdatedAt; v != nil {
// 		db = db.Where("updated_at > ?", v)
// 	}
// 	if params.WithMenuIDs {
// 		db = db.Preload("RoleMenus")
// 	}

// 	var list models.Roles
// 	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	}

// 	if params.WithMenuIDs {
// 		for _, item := range list {
// 			for _, roleMenu := range item.RoleMenus {
// 				item.MenuIDs = append(item.MenuIDs, roleMenu.MenuID)
// 			}
// 		}
// 	}

// 	queryResult := &models.RoleQueryResult{
// 		PageResult: pageResult,
// 		Data:       list,
// 	}

// 	return queryResult, nil
// }

// // Get the specified role from the database.
// func (a *Role) Get(ctx context.Context, id string, opts ...models.RoleQueryOptions) (*models.Role, error) {
// 	var opt models.RoleQueryOptions
// 	if len(opts) > 0 {
// 		opt = opts[0]
// 	}

// 	db := util.GetDB(ctx, a.DB).Model(new(models.Role)).Where("id=?", id)

// 	item := new(models.Role)
// 	ok, err := util.FindOne(ctx, db, opt.QueryOptions, item)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	} else if !ok {
// 		return nil, nil
// 	}

// 	for _, roleMenu := range item.RoleMenus {
// 		item.MenuIDs = append(item.MenuIDs, roleMenu.MenuID)
// 	}

// 	return item, nil
// }

// // Exist checks if the specified role exists in the database.
// func (a *Role) Exists(ctx context.Context, id string) (bool, error) {
// 	ok, err := util.Exists(ctx, util.GetDB(ctx, a.DB).Model(new(models.Role)).Where("id=?", id))
// 	return ok, errors.WithStack(err)
// }

// // Create a new role.
// func (a *Role) Create(ctx context.Context, item ...*models.Role) error {
// 	result := util.GetDB(ctx, a.DB).Create(item)
// 	return errors.WithStack(result.Error)
// }

// // Update the specified role in the database.
// func (a *Role) Update(ctx context.Context, item *models.Role, selectFields ...string) error {
// 	selected := []string{"*"}
// 	if len(selectFields) > 0 {
// 		selected = selectFields
// 	}
// 	result := util.GetDB(ctx, a.DB).Where("id=?", item.ID).Select(selected).Omit("created_at").Updates(item)
// 	return errors.WithStack(result.Error)
// }

// // Delete the specified role from the database.
// func (a *Role) Delete(ctx context.Context, id ...string) error {
// 	result := util.GetDB(ctx, a.DB).Where("id IN (?)", id).Delete(new(models.Role))
// 	return errors.WithStack(result.Error)
// }
