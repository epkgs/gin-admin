package service

import (
	"context"
	"sort"
	"time"

	"gin-admin/internal/dao"
	"gin-admin/internal/dto"
	"gin-admin/internal/errorx"
	"gin-admin/internal/model"
	"gin-admin/internal/types"
	"gin-admin/locales"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/crypto/hash"
	"gin-admin/pkg/helper"
	"gin-admin/pkg/jwtx"
	"gin-admin/pkg/logger"

	"github.com/epkgs/i18n/errors"

	"gorm.io/gorm"
)

var ErrInvalidToken = errorx.ErrUnauthorized.WithMsg(locales.User.Str("Invalid token"))

// Auth management for SYS
type Auth struct {
	app     types.AppContext
	cacher  cachex.Cacher
	jwt     jwtx.Auther
	userSvc *User
	menuSvc *Menu
	q       *dao.Query
}

func NewAuth(app types.AppContext) *Auth {
	return &Auth{
		app:     app,
		cacher:  app.Cacher(),
		jwt:     app.Jwt(),
		userSvc: NewUser(app),
		menuSvc: NewMenu(app),
		q:       dao.Use(app.DB()),
	}
}

func (a *Auth) Login(ctx context.Context, req *dto.Login) (*dto.LoginToken, error) {
	// verify captcha
	// if !captcha.VerifyString(req.CaptchaID, req.CaptchaCode) {
	// 	return nil, errors.BadRequest("Incorrect captcha")
	// }

	ctx = logger.WithAttrs(ctx, "tag", "login")

	u := a.q.User

	// get user info
	user, err := u.WithContext(ctx).Select(u.ID, u.Password, u.Status).Where(u.Username.Eq(req.Username)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.ErrUnauthorized.WithMsg(locales.User.Str("Incorrect username or password"))
		}
		return nil, errorx.WrapGormError(err)
	}

	if user.Status != model.UserStatus_Activated {
		return nil, errorx.ErrForbidden.WithMsg(locales.User.Str("User %s is disabled", req.Username))
	}

	// check password
	if err := hash.CompareHashAndPassword(user.Password, req.Password); err != nil {
		return nil, errorx.ErrUnauthorized.WithMsg(locales.User.Str("Incorrect username or password"))
	}

	userID := user.ID
	ctx = helper.WithUserID(ctx, userID)

	// set user cache with role ids
	roleIDs, err := a.userSvc.GetRoleIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = a.userSvc.SetRoleIDsCache(ctx, userID, roleIDs, time.Duration(a.app.Config().Cache.Expiration.User)*time.Hour)
	if err != nil {
		logger.Error(ctx, "Failed to set cache",
			"error", err,
		)
	}

	// generate token
	token, err := a.jwt.GenerateToken(ctx, userID)
	if err != nil {
		return nil, errorx.ErrInternalServerError.Wrap(err)
	}

	loginToken := &dto.LoginToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expires:      token.Expires,
	}

	logger.Info(ctx, "Login success",
		"username", req.Username,
		"accessToken", loginToken.AccessToken,
		"refreshToken", loginToken.RefreshToken,
		"tokenType", loginToken.TokenType,
		"expires", loginToken.Expires,
	)

	return loginToken, nil
}

func (a *Auth) RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginToken, error) {

	ctx = logger.WithAttrs(ctx, "tag", "login")

	claims, err := a.jwt.ParseToken(ctx, refreshToken)
	if err != nil {
		if err == jwtx.ErrInvalidToken {
			return nil, ErrInvalidToken.Wrap(err)
		}
		return nil, err
	}

	userID, _ := claims.GetSubject()

	u := a.q.User

	user, err := u.WithContext(ctx).Select(u.Status, u.Username).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.ErrBadRequest.WithMsg(locales.User.Str("Incorrect user"))
		}
		return nil, errorx.WrapGormError(err)
	}

	if user.Status != model.UserStatus_Activated {
		return nil, errorx.ErrForbidden.WithMsg(locales.User.Str("User %s is disabled", user.NickName))
	}

	ctx = helper.WithUserID(ctx, userID)

	token, err := a.jwt.GenerateToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	loginToken := &dto.LoginToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expires:      token.Expires,
	}

	logger.Info(ctx, "Login success",
		"username", user.Username,
		"accessToken", loginToken.AccessToken,
		"refreshToken", loginToken.RefreshToken,
		"tokenType", loginToken.TokenType,
		"expires", loginToken.Expires,
	)

	return loginToken, nil
}

func (a *Auth) Logout(ctx context.Context) error {
	userToken := helper.GetToken(ctx)
	if userToken == "" {
		return nil
	}

	ctx = logger.WithAttrs(ctx, "tag", "logout")
	if err := a.jwt.DestroyToken(ctx, userToken); err != nil {
		return err
	}

	userID := helper.GetUserID(ctx)
	err := a.userSvc.DeleteRoleIDsCache(ctx, userID)
	if err != nil {
		logger.Error(ctx, "Failed to delete user cache",
			"error", err,
		)
	}
	logger.Info(ctx, "Logout success")

	return nil
}

// Get user info
func (a *Auth) GetUserInfo(ctx context.Context) (*model.User, error) {

	userID := helper.GetUserID(ctx)

	u := a.q.User
	user, err := u.WithContext(ctx).Where(u.ID.Eq(userID)).Preload(u.Roles).Omit(u.Password).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.ErrUnauthorized.WithMsg(locales.User.Str("User is not logged in"))
		}
		return nil, errorx.WrapGormError(err)
	}

	return user, nil
}

// Change login password
func (a *Auth) UpdatePassword(ctx context.Context, req *dto.AuthUpdatePasswordReq) error {
	if a.app.Config().IsSuper(helper.GetUserID(ctx)) {
		return errorx.ErrForbidden.WithMsg(locales.User.Str("Super user can not modify"))
	}

	userID := helper.GetUserID(ctx)

	u := a.q.User
	user, err := u.WithContext(ctx).Select(u.Password).Where(u.ID.Eq(userID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.ErrUnauthorized.WithMsg(locales.User.Str("User is not logged in"))
		}
		return errorx.WrapGormError(err)
	}

	// check old password
	if err := hash.CompareHashAndPassword(user.Password, req.OldPassword); err != nil {
		return errorx.ErrBadRequest.WithMsg(locales.User.Str("Old password incorrect")).Wrap(err)
	}

	// update password
	newPassword, err := hash.GeneratePassword(req.NewPassword)
	if err != nil {
		return errorx.ErrInternalServerError.Wrap(err)
	}

	_, err = u.WithContext(ctx).Where(u.ID.Eq(userID)).Update(u.Password, newPassword)
	return err
}

// Query menus based on user permissions
func (a *Auth) QueryMenus(ctx context.Context) (model.Menus, error) {
	req := dto.MenuListReq{
		Status: model.MenuStatus_ENABLED,
		Pager: dto.Pager{
			Page: -1,
		},
	}

	isRoot := a.app.Config().IsSuper(helper.GetUserID(ctx))
	if !isRoot {
		req.UserID = helper.GetUserID(ctx)
	}
	list, err := a.menuSvc.List(ctx, req)
	if err != nil {
		return nil, err
	}

	menus := model.Menus(list.Items)

	if isRoot {
		return menus.ToTree(), nil
	}

	m := a.q.Menu

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
			res, err := m.WithContext(ctx).Where(m.ID.In(missMenusIDs...)).Find()
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
func (a *Auth) UpdateUser(ctx context.Context, req *dto.AuthUpdateUserReq) error {
	if a.app.Config().IsSuper(helper.GetUserID(ctx)) {
		return errorx.ErrForbidden.WithMsg(locales.User.Str("Super user can not modify"))
	}

	userID := helper.GetUserID(ctx)

	u := a.q.User
	user, err := u.WithContext(ctx).Where(u.ID.Eq(userID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.ErrUnauthorized.WithMsg(locales.User.Str("User is not logged in"))
		}
		return errorx.WrapGormError(err)
	}

	userBO := u.WithContext(ctx)
	isDirty := false
	if req.Email != nil {
		user.Email = *req.Email
		userBO = userBO.Select(u.Email)
		isDirty = true
	}
	if req.NickName != nil {
		user.NickName = *req.NickName
		userBO = userBO.Select(u.NickName)
		isDirty = true
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
		userBO = userBO.Select(u.Phone)
		isDirty = true
	}
	if req.Wechat != nil {
		user.Wechat = *req.Wechat
		userBO = userBO.Select(u.Wechat)
		isDirty = true
	}

	if !isDirty {
		return errorx.ErrBadRequest.WithMsg(locales.User.Str("Nothing to update"))
	}

	_, err = userBO.Updates(user)

	return err
}
