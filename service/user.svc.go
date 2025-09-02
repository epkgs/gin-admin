package service

import (
	"context"
	"encoding/json"
	"time"

	"gin-admin/errorx"
	"gin-admin/locales"
	"gin-admin/model/bo"
	"gin-admin/model/dto"
	"gin-admin/model/po"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/crypto/hash"
	"gin-admin/pkg/randx"
	"gin-admin/types"

	"github.com/epkgs/i18n/errors"

	"github.com/epkgs/object"
	"gorm.io/gen"
	"gorm.io/gorm"
)

const (
	gCacheNSForUserRoles = "user_roles"
)

// User management for SYS
type User struct {
	app    types.AppContext
	Cacher cachex.Cacher
	q      *bo.Query
}

func NewUser(app types.AppContext) *User {
	return &User{
		app:    app,
		Cacher: app.Cacher(),
		q:      bo.Use(app.DB()),
	}
}

// List users from the data access object based on the provided parameters and options.
func (a *User) List(ctx context.Context, req dto.UserListReq) (*dto.List[*po.User], error) {

	u := a.q.User

	scope := func(d gen.Dao) gen.Dao {
		if v := req.LikeUsername; len(v) > 0 {
			d = d.Where(u.Username.Like("%" + v + "%"))
		}
		if v := req.LikeName; len(v) > 0 {
			d = d.Where(u.NickName.Like("%" + v + "%"))
		}
		if v := req.Status; len(v) > 0 {
			d = d.Where(u.Status.Eq(v))
		}
		if req.WithRoles {
			d = d.Preload(u.Roles)
		}
		return d
	}

	list, err := u.WithContext(ctx).Scopes(scope, req.PageScope()).Find()
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	count, err := u.WithContext(ctx).Scopes(scope).Count()
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	return dto.NewList(list, &dto.Pager{
		Page:  req.Page,
		Limit: req.Limit,
		Total: count,
	}), nil
}

// Get the specified user from the data access object.
func (a *User) Get(ctx context.Context, id string) (*po.User, error) {

	u := a.q.User

	user, err := u.WithContext(ctx).Preload(u.Roles).Get(id)
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	return user, nil
}

// Create a new user in the data access object.
func (a *User) Create(ctx context.Context, req *dto.UserCreateReq) (*po.User, error) {

	if req.Username == a.app.Config().Super.Username {
		return nil, errorx.ErrForbidden.WithMsg(locales.User.Str("Super user can not modify")) // 超级管理员不允许修改
	}

	u := a.q.User

	count, err := u.WithContext(ctx).Where(u.Username.Eq(req.Username)).Count()
	if err != nil {
		return nil, errorx.WrapGormError(err)
	} else if count > 0 {
		return nil, errorx.ErrConflict.WithMsg(locales.User.Str("User %s already exists", req.Username)) // 用户名已存在
	}

	user := &po.User{
		ID:        randx.NewXID(),
		CreatedAt: time.Now(),
	}

	if req.Password == "" {
		req.Password = a.app.Config().DefaultLoginPwd
	}

	if err := object.Assign(user, req); err != nil {
		return nil, errorx.ErrInternalServerError.Wrap(err)
	}

	if pass := req.Password; pass != "" {
		hashPass, err := hash.GeneratePassword(pass)
		if err != nil {
			return nil, errorx.ErrInternalServerError.WithMsg(locales.User.Str("Password encrypt failed")).Wrap(err)
		}
		user.Password = hashPass
	}

	r := a.q.Role

	roles, err := r.WithContext(ctx).Where(r.ID.In(req.RoleIDs...)).Find()
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	user.Roles = roles
	if err := u.WithContext(ctx).Create(user); err != nil {
		return nil, errorx.WrapGormError(err)
	}

	return user, nil
}

// Update the specified user in the data access object.
func (a *User) Update(ctx context.Context, id string, req *dto.UserUpdateReq) error {

	if id == a.app.Config().Super.ID {
		return errorx.ErrForbidden.WithMsg(locales.User.Str("Super user can not modify")) // 超级管理员不允许修改
	}

	u := a.q.User

	user, err := u.WithContext(ctx).Get(id)
	if err != nil {
		return errorx.WrapGormError(err)
	}

	if req.Username != nil && user.Username != *req.Username {
		count, err := u.WithContext(ctx).Where(u.Username.Eq(*req.Username)).Count()
		if err != nil {
			return errorx.WrapGormError(err)
		} else if count > 0 {
			return errorx.ErrConflict.WithMsg(locales.User.Str("User %s already exists", *req.Username)) // 用户名已存在
		}
	}

	var md object.Metadata
	if err := object.Assign(user, req, func(c *object.AssignConfig) {
		c.Metadata = &md
	}); err != nil {
		return errorx.ErrInternalServerError.Wrap(err)
	}

	selected := md.Keys

	if req.Password != nil {
		pass := *req.Password
		hashPass, err := hash.GeneratePassword(pass)
		if err != nil {
			return errorx.ErrInternalServerError.WithMsg(locales.User.Str("Password encrypt failed")).Wrap(err)
		}
		user.Password = hashPass
	}

	r := a.q.Role

	if req.RoleIDs != nil {
		roles, err := r.WithContext(ctx).Where(r.ID.In(*req.RoleIDs...)).Find()
		if err != nil {
			return errorx.WrapGormError(err)
		}
		user.Roles = roles
		selected = append(selected, "Roles")
	}

	user.UpdatedAt = time.Now()

	if err := u.WithContext(ctx).UnderlyingDB().Select(selected).Updates(user).Error; err != nil {
		return errorx.WrapGormError(err)
	}

	return nil
}

// Delete the specified user from the data access object.
func (a *User) Delete(ctx context.Context, id string) error {

	if id == a.app.Config().Super.ID {
		return errorx.ErrForbidden.WithMsg(locales.User.Str("Super user can not modify")) // 超级管理员不允许修改
	}

	u := a.q.User

	user, err := u.WithContext(ctx).Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.ErrNotFound.WithMsg(locales.User.Str("User not found"))
		}
		return errorx.WrapGormError(err)
	}

	err = a.q.Transaction(func(tx *bo.Query) error {
		if _, err := u.WithContext(ctx).Delete(user); err != nil {
			return err
		}

		sql := "DELETE FROM r_role_users WHERE user_id= ?"
		if err := u.WithContext(ctx).UnderlyingDB().Exec(sql, user.ID).Error; err != nil {
			return err
		}
		return a.DeleteRoleIDsCache(ctx, id)
	})

	return errorx.WrapGormError(err)
}

func (a *User) ResetPassword(ctx context.Context, id string) error {
	if id == a.app.Config().Super.ID {
		return errorx.ErrForbidden.WithMsg(locales.User.Str("Super user can not modify")) // 超级管理员不允许修改
	}

	u := a.q.User

	_, err := u.WithContext(ctx).Select(u.ID).Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.ErrNotFound.WithMsg(locales.User.Str("User not found"))
		}
		return errorx.WrapGormError(err)
	}

	hashPass, err := hash.GeneratePassword(a.app.Config().DefaultLoginPwd)
	if err != nil {
		return errorx.ErrInternalServerError.WithMsg(locales.User.Str("Password encrypt failed")).Wrap(err)
	}

	err = a.q.Transaction(func(tx *bo.Query) error {
		if _, err := u.WithContext(ctx).Where(u.ID.Eq(id)).Update(u.Password, hashPass); err != nil {
			return err
		}
		return nil
	})

	return errorx.WrapGormError(err)
}

func (a *User) GetRoleIDs(ctx context.Context, id string) ([]string, error) {
	var roles []*po.Role

	u := a.q.User

	err := u.WithContext(ctx).UnderlyingDB().
		Where(u.ID.Eq(id)).
		Association(u.Roles.Name()).
		Find(&roles)

	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	ids := make([]string, 0, len(roles))
	for i, role := range roles {
		ids[i] = role.ID
	}

	return ids, nil
}

func (a *User) SetRoleIDsCache(ctx context.Context, userID string, roleIDs []string, expiration ...time.Duration) error {
	byt, err := json.Marshal(roleIDs)
	if err != nil {
		return errorx.ErrInternalServerError.Wrap(err)
	}
	return a.Cacher.Set(ctx, gCacheNSForUserRoles, userID, string(byt), expiration...)
}

func (a *User) DeleteRoleIDsCache(ctx context.Context, userID string) error {
	return a.Cacher.Delete(ctx, gCacheNSForUserRoles, userID)
}

func (a *User) GetRoleIDsCache(ctx context.Context, userID string) ([]string, error) {
	val, err := a.Cacher.Get(ctx, gCacheNSForUserRoles, userID)
	if err != nil {
		if err == cachex.ErrNotFound {
			return nil, errorx.ErrNotFound.WithMsg("record not found").Wrap(err)
		}
		return nil, errorx.ErrInternalServerError.Wrap(err)
	}

	var roleIDs []string
	if err := json.Unmarshal([]byte(val), &roleIDs); err != nil {
		return nil, errorx.ErrInternalServerError.Wrap(err)
	}

	return roleIDs, nil
}

func (a *User) InitSuperUserIfNeed(ctx context.Context) error {

	u := a.q.User
	cfg := a.app.Config()

	err := a.q.Transaction(func(tx *bo.Query) error {

		user, err := u.WithContext(ctx).Get(cfg.Super.ID)
		if user == nil || errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有 root 账户，则插入数据库
			hashedPass, err := hash.GeneratePassword(cfg.Super.Password)
			if err != nil {
				return err
			}
			user := &po.User{
				ID:       cfg.Super.ID,
				Username: cfg.Super.Username,
				NickName: cfg.Super.NickName,
				Password: hashedPass,
				Status:   po.UserStatus_Activated,
			}
			return u.WithContext(ctx).Create(user)
		}
		if err != nil {
			return err
		}

		if user.Username != cfg.Super.Username ||
			hash.CompareHashAndPassword(user.Password, cfg.Super.Password) != nil ||
			user.NickName != cfg.Super.NickName {
			// 如果root账户信息有误，则更新数据库
			hashedPass, err := hash.GeneratePassword(cfg.Super.Password)
			if err != nil {
				return err
			}

			user.NickName = cfg.Super.NickName
			user.Password = hashedPass
			user.Username = cfg.Super.Username

			_, err = u.WithContext(ctx).Select(u.NickName, u.Password, u.Username).Updates(user)
			return err
		}

		return nil
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return errorx.WrapGormError(err)
}
