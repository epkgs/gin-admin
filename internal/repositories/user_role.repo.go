package repositories

import (
	"context"

	"gin-admin/internal/models"
	"gin-admin/pkg/gormx"

	"gorm.io/gorm"
)

// User roles
type UserRole struct {
	gormx.Repository[models.UserRole]
}

func NewUserRole(db *gorm.DB) *UserRole {
	return &UserRole{
		Repository: gormx.NewGenericRepo[models.UserRole](db),
	}
}

func (a *UserRole) DeleteByUserID(ctx context.Context, userID ...string) error {
	return a.Repository.DeleteBatch(ctx, gormx.WithWhere("user_id IN (?)", userID))
}

// // List user roles from the database based on the provided parameters and options.
// func (a *UserRole) List(ctx context.Context,  opts ...gormx.Option) (*dtos.List[*models.UserRoles], error) {

// 	a.Repository.List(ctx, func(db *gorm.DB) *gorm.DB {
// db = db.Table(fmt.Sprintf("%s AS a", new(models.UserRole).TableName()))
// 	if opt.JoinRole {
// 		db = db.Joins(fmt.Sprintf("left join %s b on a.role_id=b.id", new(models.Role).TableName()))
// 		db = db.Select("a.*,b.name as role_name")
// 	}

// 	if v := params.InUserIDs; len(v) > 0 {
// 		db = db.Where("a.user_id IN (?)", v)
// 	}
// 	if v := params.UserID; len(v) > 0 {
// 		db = db.Where("a.user_id = ?", v)
// 	}
// 	if v := params.RoleID; len(v) > 0 {
// 		db = db.Where("a.role_id = ?", v)
// 	}
// 	})

// 	var list models.UserRoles
// 	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	}

// 	queryResult := &models.UserRoleQueryResult{
// 		PageResult: pageResult,
// 		Data:       list,
// 	}
// 	return queryResult, nil
// }

// // Get the specified user role from the database.
// func (a *UserRole) Get(ctx context.Context, id string, opts ...models.UserRoleQueryOptions) (*models.UserRole, error) {
// 	var opt models.UserRoleQueryOptions
// 	if len(opts) > 0 {
// 		opt = opts[0]
// 	}

// 	item := new(models.UserRole)
// 	ok, err := util.FindOne(ctx, util.GetDB(ctx, a.DB).Model(new(models.UserRole)).Where("id=?", id), opt.QueryOptions, item)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	} else if !ok {
// 		return nil, nil
// 	}
// 	return item, nil
// }

// // Exist checks if the specified user role exists in the database.
// func (a *UserRole) Exists(ctx context.Context, id string) (bool, error) {
// 	ok, err := util.Exists(ctx, util.GetDB(ctx, a.DB).Model(new(models.UserRole)).Where("id=?", id))
// 	return ok, errors.WithStack(err)
// }

// // Create a new user role.
// func (a *UserRole) Create(ctx context.Context, item ...*models.UserRole) error {
// 	result := util.GetDB(ctx, a.DB).Create(item)
// 	return errors.WithStack(result.Error)
// }

// // Update the specified user role in the database.
// func (a *UserRole) Update(ctx context.Context, item *models.UserRole, selectFields ...string) error {
// 	selected := []string{"*"}
// 	if len(selectFields) > 0 {
// 		selected = selectFields
// 	}
// 	result := util.GetDB(ctx, a.DB).Where("id=?", item.ID).Select(selected).Omit("created_at").Updates(item)
// 	return errors.WithStack(result.Error)
// }

// // Delete the specified user role from the database.
// func (a *UserRole) Delete(ctx context.Context, id ...string) error {
// 	result := util.GetDB(ctx, a.DB).Where("id IN (?)", id).Delete(new(models.UserRole))
// 	return errors.WithStack(result.Error)
// }

// func (a *UserRole) DeleteByRoleID(ctx context.Context, roleID ...string) error {
// 	result := util.GetDB(ctx, a.DB).Where("role_id IN (?)", roleID).Delete(new(models.UserRole))
// 	return errors.WithStack(result.Error)
// }
