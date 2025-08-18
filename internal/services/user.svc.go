package services

import (
	"context"
	"encoding/json"
	"time"

	"gin-admin/internal/configs"
	"gin-admin/internal/dtos"
	"gin-admin/internal/errorx"
	"gin-admin/internal/models"
	"gin-admin/internal/repositories"
	"gin-admin/internal/types"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/crypto/hash"
	"gin-admin/pkg/gormx"
	"gin-admin/pkg/randx"

	"github.com/epkgs/i18n/errors"
	"github.com/epkgs/object"
	"gorm.io/gorm"
)

const (
	gCacheNSForUserRoles = "user_roles"
)

// User management for SYS
type User struct {
	Cacher       cachex.Cacher
	UserRepo     *repositories.User
	RoleRepo     *repositories.Role
	UserRoleRepo *repositories.UserRole
}

func NewUser(app types.AppContext) *User {
	return &User{
		Cacher:       app.Cacher(),
		UserRepo:     repositories.NewUser(app.DB()),
		RoleRepo:     repositories.NewRole(app.DB()),
		UserRoleRepo: repositories.NewUserRole(app.DB()),
	}
}

// List users from the data access object based on the provided parameters and options.
func (a *User) List(ctx context.Context, req dtos.UserListReq) (*dtos.List[*models.User], error) {
	option := func(db *gorm.DB) *gorm.DB {

		if v := req.LikeUsername; len(v) > 0 {
			db = db.Where("username LIKE ?", "%"+v+"%")
		}
		if v := req.LikeName; len(v) > 0 {
			db = db.Where("name LIKE ?", "%"+v+"%")
		}
		if v := req.Status; len(v) > 0 {
			db = db.Where("status = ?", v)
		}
		if req.WithRoles {
			db = db.Preload("Roles")
		}
		return db
	}

	list, err := a.UserRepo.Find(ctx, option, gormx.WithPage(req.Page, req.Limit))
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	count, err := a.UserRepo.Count(ctx, option)
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	return dtos.NewList(list, &dtos.Pager{
		Page:  req.Page,
		Limit: req.Limit,
		Total: count,
	}), nil
}

// Get the specified user from the data access object.
func (a *User) Get(ctx context.Context, id string) (*models.User, error) {
	user, err := a.UserRepo.Get(ctx, id, gormx.WithPreload("Roles"))
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	return user, nil
}

// Create a new user in the data access object.
func (a *User) Create(ctx context.Context, req *dtos.UserCreateReq) (*models.User, error) {

	if req.Username == configs.C.Super.Username {
		return nil, errorx.User.ModifySuperUser // 超级管理员不允许修改
	}

	existsUsername, err := a.UserRepo.ExistsUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	} else if existsUsername {
		return nil, errorx.User.Exists.New(struct{ Name string }{Name: req.Username}) // 用户名已存在
	}

	user := &models.User{
		ID:        randx.NewXID(),
		CreatedAt: time.Now(),
	}

	if req.Password == "" {
		req.Password = configs.C.DefaultLoginPwd
	}

	if err := object.Assign(user, req); err != nil {
		return nil, errorx.General.Internal.Wrap(err)
	}

	if pass := req.Password; pass != "" {
		hashPass, err := hash.GeneratePassword(pass)
		if err != nil {
			return nil, errorx.User.PasswordEncrypt.Wrap(err)
		}
		user.Password = hashPass
	}

	roles, err := a.RoleRepo.Find(ctx, gormx.WithWhere("id IN ?", req.RoleIDs))
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	user.Roles = roles
	if err := a.UserRepo.Create(ctx, user); err != nil {
		return nil, errorx.WrapGormError(err)
	}

	return user, nil
}

// Update the specified user in the data access object.
func (a *User) Update(ctx context.Context, id string, req *dtos.UserUpdateReq) error {

	if id == configs.C.Super.ID {
		return errorx.User.ModifySuperUser // 超级管理员不允许修改
	}

	user, err := a.UserRepo.Get(ctx, id)
	if err != nil {
		return errorx.WrapGormError(err)
	}

	if req.Username != nil && user.Username != *req.Username {
		existsUsername, err := a.UserRepo.ExistsUsername(ctx, *req.Username)
		if err != nil {
			return errorx.WrapGormError(err)
		} else if existsUsername {
			return errorx.User.Exists.New(struct{ Name string }{Name: *req.Username}) // 用户名已存在
		}
	}

	var md object.Metadata
	if err := object.Assign(user, req, func(c *object.AssignConfig) {
		c.Metadata = &md
	}); err != nil {
		return errorx.General.Internal.Wrap(err)
	}

	selected := md.Keys

	if req.Password != nil {
		pass := *req.Password
		hashPass, err := hash.GeneratePassword(pass)
		if err != nil {
			return errorx.User.PasswordEncrypt.Wrap(err)
		}
		user.Password = hashPass
	}

	if req.RoleIDs != nil {
		roles, err := a.RoleRepo.Find(ctx, gormx.WithWhere("id IN ?", req.RoleIDs))
		if err != nil {
			return errorx.WrapGormError(err)
		}
		user.Roles = roles
		selected = append(selected, "Roles")
	}

	user.UpdatedAt = time.Now()

	if err := a.UserRepo.Update(ctx, user, gormx.WithSelect(selected)); err != nil {
		return errorx.WrapGormError(err)
	}

	return nil
}

// Delete the specified user from the data access object.
func (a *User) Delete(ctx context.Context, id string) error {

	if id == configs.C.Super.ID {
		return errorx.User.ModifySuperUser // 超级管理员不允许修改
	}

	exists, err := a.UserRepo.Exists(ctx, gormx.WithWhere("id = ?", id))
	if err != nil {
		return errorx.WrapGormError(err)
	} else if !exists {
		return errorx.User.NotFound
	}

	err = a.UserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := a.UserRepo.Delete(ctx, id); err != nil {
			return err
		}
		if err := a.UserRoleRepo.DeleteByUserID(ctx, id); err != nil {
			return err
		}
		return a.DeleteRoleIDsCache(ctx, id)
	})

	return errorx.WrapGormError(err)
}

func (a *User) ResetPassword(ctx context.Context, id string) error {
	if id == configs.C.Super.ID {
		return errorx.User.ModifySuperUser // 超级管理员不允许修改
	}

	exists, err := a.UserRepo.Exists(ctx, gormx.WithWhere("id=?", id))
	if err != nil {
		return errorx.WrapGormError(err)
	} else if !exists {
		return errorx.User.NotFound
	}

	hashPass, err := hash.GeneratePassword(configs.C.DefaultLoginPwd)
	if err != nil {
		return errorx.User.PasswordEncrypt.Wrap(err)
	}

	err = a.UserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := a.UserRepo.UpdatePassword(ctx, id, hashPass); err != nil {
			return err
		}
		return nil
	})

	return errorx.WrapGormError(err)
}

func (a *User) GetRoleIDs(ctx context.Context, id string) ([]string, error) {
	userRoles, err := a.UserRoleRepo.Find(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", id)
	})
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	return models.UserRoles(userRoles).ToRoleIDs(), nil
}

func (a *User) SetRoleIDsCache(ctx context.Context, userID string, roleIDs []string, expiration ...time.Duration) error {
	byt, err := json.Marshal(roleIDs)
	if err != nil {
		return errorx.General.Internal.Wrap(err)
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
			return nil, errorx.General.RecordNotFound.Wrap(err)
		}
		return nil, errorx.General.Internal.Wrap(err)
	}

	var roleIDs []string
	if err := json.Unmarshal([]byte(val), &roleIDs); err != nil {
		return nil, errorx.General.Internal.Wrap(err)
	}

	return roleIDs, nil
}

func (a *User) InitSuperUserIfNeed(ctx context.Context) error {

	err := a.UserRepo.Transaction(ctx, func(tx *gorm.DB) error {

		user, err := a.UserRepo.Get(ctx, configs.C.Super.ID)
		if user == nil || errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有 root 账户，则插入数据库
			hashedPass, err := hash.GeneratePassword(configs.C.Super.Password)
			if err != nil {
				return err
			}
			user := &models.User{
				ID:       configs.C.Super.ID,
				Username: configs.C.Super.Username,
				NickName: configs.C.Super.NickName,
				Password: hashedPass,
				Status:   models.UserStatus_Activated,
			}
			return a.UserRepo.Create(ctx, user)
		}
		if err != nil {
			return err
		}

		if user.Username != configs.C.Super.Username || hash.CompareHashAndPassword(user.Password, configs.C.Super.Password) != nil || user.NickName != configs.C.Super.NickName {
			// 如果root账户信息有误，则更新数据库
			hashedPass, err := hash.GeneratePassword(configs.C.Super.Password)
			if err != nil {
				return err
			}

			user.NickName = configs.C.Super.NickName
			user.Password = hashedPass
			user.Username = configs.C.Super.Username
			return a.UserRepo.Update(ctx, user, gormx.WithSelect("NickName", "Password", "Username"))
		}

		return nil

	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return errorx.WrapGormError(err)
}
