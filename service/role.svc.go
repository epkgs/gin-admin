package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gin-admin/errorx"
	"gin-admin/locales"
	"gin-admin/model/bo"
	"gin-admin/model/dto"
	"gin-admin/model/po"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/randx"
	"gin-admin/types"

	"github.com/epkgs/i18n/errors"
	"github.com/epkgs/object"
	"gorm.io/gen"
	"gorm.io/gorm"
)

const (
	gCacheKeyForCasbin = "sync:casbin"
	gCacheNSForRole    = "role"
)

// Role management for SYS
type Role struct {
	app    types.AppContext
	cacher cachex.Cacher
	q      *bo.Query
}

func NewRole(app types.AppContext) *Role {
	return &Role{
		app:    app,
		cacher: app.Cacher(),
		q:      bo.Use(app.DB()),
	}
}

// List roles from the data access object based on the provided parameters and options.
func (a *Role) List(ctx context.Context, req dto.RoleListReq) (*dto.List[*po.Role], error) {

	r := a.q.Role

	scope := func(d gen.Dao) gen.Dao {
		if v := req.Name; len(v) > 0 {
			d = d.Where(r.Name.Like("%" + v + "%"))
		}
		if v := req.Status; len(v) > 0 {
			d = d.Where(r.Status.Eq(v))
		}
		if req.WithMenus {
			d = d.Preload(r.Menus)
		}

		return d
	}

	list, err := r.WithContext(ctx).
		Scopes(scope, req.PageScope()).
		Order(r.Rank.Desc()).
		Order(r.CreatedAt.Desc()).
		Find()
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	count, err := r.WithContext(ctx).Scopes(scope).Count()
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	result := dto.NewList(list, &dto.Pager{
		Page:  req.Page,
		Limit: req.Limit,
		Total: count,
	})
	return result, nil
}

// Get the specified role from the data access object.s
func (a *Role) Get(ctx context.Context, id string) (*po.Role, error) {
	r := a.q.Role

	role, err := r.WithContext(ctx).Preload(r.Menus).Get(id)
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	return role, nil
}

// Create a new role in the data access object.
func (a *Role) Create(ctx context.Context, req dto.RoleCreateReq) (*po.Role, error) {
	r := a.q.Role

	if count, err := r.WithContext(ctx).Where(r.Code.Eq(req.Code)).Count(); err != nil {
		return nil, errorx.WrapGormError(err)
	} else if count > 0 {
		return nil, errorx.ErrBadRequest.WithMsg(locales.User.Str("Role code already exists"))
	}

	role := &po.Role{
		ID:        randx.NewXID(),
		CreatedAt: time.Now(),
	}

	if err := object.Assign(role, req, func(c *object.AssignConfig) {
		c.SkipKeys = []string{"Menus"}
	}); err != nil {
		return nil, errorx.ErrInternalServerError.Wrap(err)
	}

	m := a.q.Menu

	if len(req.MenuIDs) > 0 {
		menus, err := m.WithContext(ctx).Where(m.ID.In(req.MenuIDs...)).Find()
		if err != nil {
			return nil, errorx.WrapGormError(err)
		}

		role.Menus = menus
	}

	if err := r.WithContext(ctx).Create(role); err != nil {
		return nil, errorx.WrapGormError(err)
	}

	return role, nil
}

// Update the specified role in the data access object.
func (a *Role) Update(ctx context.Context, id string, req *dto.RoleUpdateReq) error {
	r := a.q.Role

	role, err := r.WithContext(ctx).Get(id)
	if err != nil {
		return errorx.WrapGormError(err)
	}

	if req.Code != nil && *req.Code != role.Code {
		if count, err := r.WithContext(ctx).Where(r.Code.Eq(*req.Code)).Count(); err != nil {
			return errorx.WrapGormError(err)
		} else if count > 0 {
			return errorx.ErrBadRequest.WithMsg(locales.User.Str("Role code already exists"))
		}
	}

	var md object.Metadata
	if err := object.Assign(role, req, func(c *object.AssignConfig) {
		c.SkipKeys = []string{"menus"}
		c.Metadata = &md
	}); err != nil {
		return err
	}

	m := a.q.Menu

	selected := md.Keys

	if req.MenuIDs != nil {
		menus, err := m.WithContext(ctx).Where(m.ID.In(*req.MenuIDs...)).Find()
		if err != nil {
			return errorx.WrapGormError(err)
		}
		role.Menus = menus
		selected = append(selected, "Menus")
	}

	role.UpdatedAt = time.Now()

	err = a.q.Transaction(func(tx *bo.Query) error {

		res := m.WithContext(ctx).UnderlyingDB().
			Omit("Menus.*").
			Select(selected).
			Updates(role)

		if err := res.Error; err != nil {
			return err
		}
		return a.RefreshUpdateTime(ctx)
	})

	return errorx.WrapGormError(err)
}

// Delete the specified role from the data access object.
func (a *Role) Delete(ctx context.Context, id string) error {

	r := a.q.Role

	role, err := r.WithContext(ctx).Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.ErrNotFound.WithMsg(locales.User.Str("Role not found"))
		}
		return errorx.WrapGormError(err)
	} else if role == nil {
		return errorx.ErrNotFound.WithMsg(locales.User.Str("Role not found"))
	}

	err = a.q.Transaction(func(tx *bo.Query) error {
		if _, err := r.WithContext(ctx).Select(r.Menus.Field(), r.Users.Field()).Delete(role); err != nil {
			return err
		}

		return a.RefreshUpdateTime(ctx)
	})

	return errorx.WrapGormError(err)
}

func (a *Role) RefreshUpdateTime(ctx context.Context) error {
	return a.cacher.Set(
		ctx,
		gCacheNSForRole,
		gCacheKeyForCasbin,
		fmt.Sprintf("%d", time.Now().Unix()), // 使用时间戳存取更方便快捷
	)
}

func (a *Role) GetUpdateTime(ctx context.Context) (int64, error) {
	val, err := a.cacher.Get(ctx, gCacheNSForRole, gCacheKeyForCasbin)
	if err != nil {
		if err == cachex.ErrNotFound {
			return 0, errorx.ErrNotFound.WithMsg(locales.DB.Str("record not found")).Wrap(err)
		}
		return 0, err
	}

	updated, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, errorx.ErrInternalServerError.Wrap(err)
	}

	return updated, nil
}
