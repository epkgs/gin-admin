package services

import (
	"context"
	"sort"
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
	"gin-admin/pkg/helper"
	"gin-admin/pkg/jwtx"
	"gin-admin/pkg/logger"

	"github.com/epkgs/i18n/errors"
	"github.com/epkgs/object"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Auth management for SYS
type Auth struct {
	Cacher       cachex.Cacher
	Jwt          jwtx.Auther
	UserRepo     *repositories.User
	UserRoleRepo *repositories.UserRole
	MenuRepo     *repositories.Menu
	UserSvc      *User
	MenuSvc      *Menu
}

func NewAuth(app types.AppContext) *Auth {
	return &Auth{
		Cacher:       app.Cacher(),
		Jwt:          app.Jwt(),
		UserRepo:     repositories.NewUser(app.DB()),
		UserRoleRepo: repositories.NewUserRole(app.DB()),
		MenuRepo:     repositories.NewMenu(app.DB()),
		UserSvc:      NewUser(app),
		MenuSvc:      NewMenu(app),
	}
}

func (a *Auth) ParseUserID(c *gin.Context) (string, error) {
	ctx := c.Request.Context()

	rootID := configs.C.Super.ID
	if configs.C.Middleware.Auth.Disable {
		return rootID, nil
	}

	token := helper.GetToken(c)
	if token == "" {
		return "", errorx.User.InvalidToken
	}

	ctx = helper.WithUserToken(ctx, token)

	claims, err := a.Jwt.ParseToken(ctx, token)
	if err != nil {
		if err == jwtx.ErrInvalidToken {
			return "", errorx.User.InvalidToken
		}
		return "", errorx.General.Internal.Wrap(err)
	}

	userID, _ := claims.GetSubject()

	if userID == rootID {
		c.Request = c.Request.WithContext(helper.WithIsRootUser(ctx))
		return userID, nil
	}

	_, err = a.UserSvc.GetRoleIDsCache(ctx, userID)
	if err != nil {
		if errors.Is(err, cachex.ErrNotFound) {

			// Check user status, if not activated, force to logout
			user, err := a.UserRepo.Get(ctx, userID, gormx.WithSelect("status"))
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return "", errorx.User.InvalidToken
				}
				return "", errorx.WrapGormError(err)
			}

			if user == nil || user.Status != models.UserStatus_Activated {
				return "", errorx.User.InvalidToken
			}

			roleIDs, err := a.UserSvc.GetRoleIDs(ctx, userID)
			if err != nil {
				return "", err
			}

			err = a.UserSvc.SetRoleIDsCache(ctx, userID, roleIDs)
			if err != nil {
				return "", err
			}
			return userID, nil
		}
		return "", err
	}

	return userID, nil
}

func (a *Auth) Login(ctx context.Context, req *dtos.Login) (*dtos.LoginToken, error) {
	// verify captcha
	// if !captcha.VerifyString(req.CaptchaID, req.CaptchaCode) {
	// 	return nil, errors.BadRequest("Incorrect captcha")
	// }

	ctx = logger.WithTag(ctx, logger.Tag_Login)

	// get user info
	user, err := a.UserRepo.GetByUsername(ctx, req.Username, gormx.WithSelect("id", "password", "status"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.User.UsernamePassword
		}
		return nil, errorx.WrapGormError(err)
	}

	if user.Status != models.UserStatus_Activated {
		return nil, errorx.User.Disabled.New(struct{ Name string }{req.Username})
	}

	// check password
	if err := hash.CompareHashAndPassword(user.Password, req.Password); err != nil {
		return nil, errorx.User.UsernamePassword
	}

	userID := user.ID
	ctx = logger.WithUserID(ctx, userID)

	// set user cache with role ids
	roleIDs, err := a.UserSvc.GetRoleIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = a.UserSvc.SetRoleIDsCache(ctx, userID, roleIDs, time.Duration(configs.C.Cache.Expiration.User)*time.Hour)
	if err != nil {
		logger.Error(ctx, "Failed to set cache", err)
	}

	// generate token
	token, err := a.Jwt.GenerateToken(ctx, userID)
	if err != nil {
		return nil, errorx.General.Internal.Wrap(err)
	}

	loginToken := &dtos.LoginToken{
		AccessToken:  token.GetAccessToken(),
		RefreshToken: token.GetRefreshToken(),
		TokenType:    token.GetTokenType(),
		Expires:      token.GetExpires(),
	}

	logger.Info(ctx, "Login success",

		map[string]any{
			"username":     req.Username,
			"accessToken":  loginToken.AccessToken,
			"refreshToken": loginToken.RefreshToken,
			"tokenType":    loginToken.TokenType,
			"expires":      loginToken.Expires,
		},
	)

	return loginToken, nil
}

func (a *Auth) RefreshToken(ctx context.Context, refreshToken string) (*dtos.LoginToken, error) {

	ctx = logger.WithTag(ctx, logger.Tag_Login)

	claims, err := a.Jwt.ParseRefreshToken(ctx, refreshToken)
	if err != nil {
		if err == jwtx.ErrInvalidToken {
			return nil, errorx.User.InvalidToken.Wrap(err)
		}
		return nil, err
	}

	userID, _ := claims.GetSubject()

	user, err := a.UserRepo.Get(ctx, userID, gormx.WithSelect("status", "username"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.User.Incorrect
		}
		return nil, errorx.WrapGormError(err)
	}

	if user.Status != models.UserStatus_Activated {
		return nil, errorx.User.Disabled.New(struct{ Name string }{user.NickName})
	}

	ctx = logger.WithUserID(ctx, userID)

	token, err := a.Jwt.GenerateToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	loginToken := &dtos.LoginToken{
		AccessToken:  token.GetAccessToken(),
		RefreshToken: token.GetRefreshToken(),
		TokenType:    token.GetTokenType(),
		Expires:      token.GetExpires(),
	}

	logger.Info(ctx, "Login success",
		map[string]any{
			"username":     user.Username,
			"accessToken":  loginToken.AccessToken,
			"refreshToken": loginToken.RefreshToken,
			"tokenType":    loginToken.TokenType,
			"expires":      loginToken.Expires,
		},
	)

	return loginToken, nil
}

func (a *Auth) Logout(ctx context.Context) error {
	userToken := helper.GetUserToken(ctx)
	if userToken == "" {
		return nil
	}

	ctx = logger.WithTag(ctx, logger.Tag_Logout)
	if err := a.Jwt.DestroyToken(ctx, userToken); err != nil {
		return err
	}

	userID := helper.GetUserID(ctx)
	err := a.UserSvc.DeleteRoleIDsCache(ctx, userID)
	if err != nil {
		logger.Error(ctx, "Failed to delete user cache", err)
	}
	logger.Info(ctx, "Logout success")

	return nil
}

// Get user info
func (a *Auth) GetUserInfo(ctx context.Context) (*models.User, error) {

	userID := helper.GetUserID(ctx)
	user, err := a.UserRepo.Get(ctx, userID, gormx.WithPreload("Roles"), gormx.WithOmit("password"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.User.NotLogin
		}
		return nil, errorx.WrapGormError(err)
	}

	return user, nil
}

// Change login password
func (a *Auth) UpdatePassword(ctx context.Context, req *dtos.AuthUpdatePasswordReq) error {
	if helper.GetIsRootUser(ctx) {
		return errorx.User.ModifySuperUser
	}

	userID := helper.GetUserID(ctx)
	user, err := a.UserRepo.Get(ctx, userID, gormx.WithSelect("password"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.User.NotLogin
		}
		return errorx.WrapGormError(err)
	}

	// check old password
	if err := hash.CompareHashAndPassword(user.Password, req.OldPassword); err != nil {
		return errorx.User.OldPassword.Wrap(err)
	}

	// update password
	newPassword, err := hash.GeneratePassword(req.NewPassword)
	if err != nil {
		return errorx.General.Internal.Wrap(err)
	}
	return a.UserRepo.UpdatePassword(ctx, userID, newPassword)
}

// Query menus based on user permissions
func (a *Auth) QueryMenus(ctx context.Context) (models.Menus, error) {
	req := dtos.MenuListReq{
		Status: models.MenuStatus_ENABLED,
		Pager: dtos.Pager{
			Page: -1,
		},
	}

	isRoot := helper.GetIsRootUser(ctx)
	if !isRoot {
		req.UserID = helper.GetUserID(ctx)
	}
	list, err := a.MenuSvc.List(ctx, req)
	if err != nil {
		return nil, err
	}

	menus := models.Menus(list.Items)

	if isRoot {
		return menus.ToTree(), nil
	}

	// fill parent menus
	if parentIDs := menus.ParentIDs(); len(parentIDs) > 0 {
		var missMenusIDs []string
		menuIDMapper := menus.ToIDMapper()
		for _, parentID := range parentIDs {
			if _, ok := menuIDMapper[parentID]; !ok {
				missMenusIDs = append(missMenusIDs, parentID)
			}
		}
		if len(missMenusIDs) > 0 {
			res, err := a.MenuRepo.Find(ctx, gormx.WithWhere("id in (?)", missMenusIDs))
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errorx.WrapGormError(err)
			}

			if len(res) > 0 {
				menus = append(menus, res...)
				sort.Sort(menus)
			}
		}
	}

	return menus, nil
}

// Update current user info
func (a *Auth) UpdateUser(ctx context.Context, req *dtos.AuthUpdateUserReq) error {
	// if util.GetIsRootUser(ctx) {
	// 	return errors.BadRequest("Super user cannot update")
	// }

	userID := helper.GetUserID(ctx)
	user, err := a.UserRepo.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.User.NotLogin
		}
		return errorx.WrapGormError(err)
	}

	var md object.Metadata
	err = object.Assign(user, req, func(c *object.AssignConfig) {
		c.SkipSameValues = true
		c.Metadata = &md
	})
	if err != nil {
		return errorx.General.Internal.Wrap(err)
	}

	if len(md.Keys) == 0 {
		return errorx.DB.NothingUpdate
	}

	return a.UserRepo.Update(ctx, user, gormx.WithSelect(md.Keys))
}
