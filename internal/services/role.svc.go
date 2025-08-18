package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gin-admin/internal/dtos"
	"gin-admin/internal/errorx"
	"gin-admin/internal/models"
	"gin-admin/internal/repositories"
	"gin-admin/internal/types"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/gormx"
	"gin-admin/pkg/randx"

	"github.com/epkgs/object"
	"gorm.io/gorm"
)

const (
	gCacheKeyForCasbin = "sync:casbin"
	gCacheNSForRole    = "role"
)

// Role management for SYS
type Role struct {
	Cacher       cachex.Cacher
	RoleRepo     *repositories.Role
	MenuRepo     *repositories.Menu
	UserRoleRepo *repositories.UserRole
}

func NewRole(app types.AppContext) *Role {
	return &Role{
		Cacher:       app.Cacher(),
		RoleRepo:     repositories.NewRole(app.DB()),
		MenuRepo:     repositories.NewMenu(app.DB()),
		UserRoleRepo: repositories.NewUserRole(app.DB()),
	}
}

// List roles from the data access object based on the provided parameters and options.
func (a *Role) List(ctx context.Context, req dtos.RoleListReq) (*dtos.List[*models.Role], error) {

	option := func(db *gorm.DB) *gorm.DB {

		if v := req.Name; len(v) > 0 {
			db = db.Where("name LIKE ?", "%"+v+"%")
		}
		if v := req.Status; len(v) > 0 {
			db = db.Where("status = ?", v)
		}
		if req.WithMenus {
			db = db.Preload("Menus")
		}

		return db
	}

	list, err := a.RoleRepo.Find(ctx, option, gormx.WithOrder("rank", "desc"), gormx.WithOrder("created_at", "desc"), gormx.WithPage(req.Page, req.Limit))
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	count, err := a.RoleRepo.Count(ctx, option)
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	result := dtos.NewList(list, &dtos.Pager{
		Page:  req.Page,
		Limit: req.Limit,
		Total: count,
	})
	return result, nil
}

// Get the specified role from the data access object.s
func (a *Role) Get(ctx context.Context, id string) (*models.Role, error) {
	role, err := a.RoleRepo.Get(ctx, id, gormx.WithPreload("Menus"))
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	return role, nil
}

// Create a new role in the data access object.
func (a *Role) Create(ctx context.Context, req dtos.RoleCreateReq) (*models.Role, error) {
	if exists, err := a.RoleRepo.ExistsCode(ctx, req.Code); err != nil {
		return nil, errorx.WrapGormError(err)
	} else if exists {
		return nil, errorx.User.RoleCodeExists
	}

	role := &models.Role{
		ID:        randx.NewXID(),
		CreatedAt: time.Now(),
	}

	if err := object.Assign(role, req, func(c *object.AssignConfig) {
		c.SkipKeys = []string{"Menus"}
	}); err != nil {
		return nil, errorx.General.Internal.Wrap(err)
	}

	if len(req.MenuIDs) > 0 {
		menus, err := a.MenuRepo.Find(ctx, gormx.WithWhere("id IN (?)", req.MenuIDs))
		if err != nil {
			return nil, errorx.WrapGormError(err)
		}

		role.Menus = menus
	}

	if err := a.RoleRepo.Create(ctx, role); err != nil {
		return nil, errorx.WrapGormError(err)
	}

	return role, nil
}

// Update the specified role in the data access object.
func (a *Role) Update(ctx context.Context, id string, req *dtos.RoleUpdateReq) error {
	role, err := a.RoleRepo.Get(ctx, id)
	if err != nil {
		return errorx.WrapGormError(err)
	}

	if req.Code != nil && *req.Code != role.Code {
		if exists, err := a.RoleRepo.ExistsCode(ctx, *req.Code); err != nil {
			return errorx.WrapGormError(err)
		} else if exists {
			return errorx.User.RoleCodeExists
		}
	}

	var md object.Metadata
	if err := object.Assign(role, req, func(c *object.AssignConfig) {
		c.SkipKeys = []string{"menus"}
		c.Metadata = &md
	}); err != nil {
		return err
	}

	selected := md.Keys

	if req.MenuIDs != nil {
		menus, err := a.MenuRepo.Find(ctx, gormx.WithWhere("id IN (?)", *req.MenuIDs))
		if err != nil {
			return errorx.WrapGormError(err)
		}
		role.Menus = menus
		selected = append(selected, "Menus")
	}

	role.UpdatedAt = time.Now()

	err = a.RoleRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := a.RoleRepo.Update(ctx, role, gormx.WithOmit("Menus.*"), gormx.WithSelect(selected)); err != nil {
			return err
		}

		return a.RefreshUpdateTime(ctx)
	})

	return errorx.WrapGormError(err)
}

// Delete the specified role from the data access object.
func (a *Role) Delete(ctx context.Context, id string) error {
	exists, err := a.RoleRepo.Exists(ctx, gormx.WithWhere("id = ?", id))
	if err != nil {
		return errorx.WrapGormError(err)
	} else if !exists {
		return errorx.User.RoleNotFount
	}

	err = a.RoleRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := a.RoleRepo.Delete(ctx, id, gormx.WithSelect("Menus", "Users")); err != nil {
			return err
		}
		return a.RefreshUpdateTime(ctx)
	})

	return errorx.WrapGormError(err)
}

func (a *Role) RefreshUpdateTime(ctx context.Context) error {
	return a.Cacher.Set(ctx, gCacheNSForRole, gCacheKeyForCasbin, fmt.Sprintf("%d", time.Now().Unix()))
}

func (a *Role) GetUpdateTime(ctx context.Context) (int64, error) {
	val, err := a.Cacher.Get(ctx, gCacheNSForRole, gCacheKeyForCasbin)
	if err != nil {
		if err == cachex.ErrNotFound {
			return 0, errorx.General.RecordNotFound.Wrap(err)
		}
		return 0, err
	}

	updated, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, errorx.General.Internal.Wrap(err)
	}

	return updated, nil
}
