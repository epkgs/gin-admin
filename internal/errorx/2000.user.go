package errorx

import (
	"gin-admin/locales"
	"net/http"

	"github.com/epkgs/i18n/errors"
)

type userErrors struct {
	NotFound                errors.I18nError                         // 用户不存在
	NotLogin                errors.I18nError                         // 用户未登录
	Disabled                errors.Definition[struct{ Name string }] // 用户 {{.Name}} 已被禁用
	Exists                  errors.Definition[struct{ Name string }] // 用户 {{.Name}} 已存在
	PasswordExpired         errors.I18nError                         // 密码已过期
	UsernamePassword        errors.I18nError                         // 用户名或密码错误
	UserTokenError          errors.I18nError                         // 用户令牌错误
	GenVisitToken           errors.I18nError                         // 生成访问令牌失败
	GenRefreshToken         errors.I18nError                         // 生成刷新令牌失败
	ParseToken              errors.I18nError                         // 解析令牌失败
	InvalidToken            errors.I18nError                         // 无效的令牌
	PasswordEncrypt         errors.I18nError                         // 密码加密失败
	PasswordDecrypt         errors.I18nError                         // 密码解密失败
	UserNameOrPasswordEmpty errors.I18nError                         // 用户名或密码不能为空
	Password                errors.I18nError                         // 密码错误
	ModifySuperUser         errors.I18nError                         // 超级用户不能修改
	RoleCodeExists          errors.I18nError                         // 角色编码已存在
	RoleNotFount            errors.I18nError                         // 角色不存在
	Incorrect               errors.I18nError                         // 用户信息错误
	OldPassword             errors.I18nError                         // 旧密码错误
	CaptchaIDNotFound       errors.I18nError                         // 验证码ID不存在
}

var User = &userErrors{
	NotFound:                New(locales.User, 2000, "user not found", http.StatusNotFound),                                          // 用户不存在
	NotLogin:                New(locales.User, 2001, "user is not logged in", http.StatusUnauthorized),                               // 用户未登录
	Disabled:                Define[struct{ Name string }](locales.User, 2002, "user {{.Name}} is disabled", http.StatusForbidden),   // 用户 {{.Name}} 已被禁用
	Exists:                  Define[struct{ Name string }](locales.User, 2003, "user {{.Name}} already exists", http.StatusConflict), // 用户 {{.Name}} 已存在
	PasswordExpired:         New(locales.User, 2004, "password expired", http.StatusForbidden),                                       // 密码已过期
	UsernamePassword:        New(locales.User, 2005, "incorrect username or password", http.StatusUnauthorized),                      // 用户名或密码错误
	UserTokenError:          New(locales.User, 2006, "wrong user token", http.StatusUnauthorized),                                    // 用户令牌错误
	GenVisitToken:           New(locales.User, 2007, "generate visit token failed", http.StatusInternalServerError),                  // 生成访问令牌失败
	GenRefreshToken:         New(locales.User, 2008, "generate refresh token failed", http.StatusInternalServerError),                // 生成刷新令牌失败
	ParseToken:              New(locales.User, 2009, "parse token failed", http.StatusUnauthorized),                                  // 解析令牌失败
	InvalidToken:            New(locales.User, 2010, "invalid token", http.StatusUnauthorized),                                       // 无效的令牌
	PasswordEncrypt:         New(locales.User, 2011, "password encrypt failed", http.StatusInternalServerError),                      // 密码加密失败
	PasswordDecrypt:         New(locales.User, 2012, "password decrypt failed", http.StatusInternalServerError),                      // 密码解密失败
	UserNameOrPasswordEmpty: New(locales.User, 2013, "username or password empty", http.StatusBadRequest),                            // 用户名或密码不能为空
	Password:                New(locales.User, 2014, "password error", http.StatusUnauthorized),                                      // 密码错误
	ModifySuperUser:         New(locales.User, 2015, "super user can not modify", http.StatusForbidden),                              // 超级用户不能修改
	RoleCodeExists:          New(locales.User, 2016, "role code already exists", http.StatusBadRequest),                              // 角色编码已存在
	RoleNotFount:            New(locales.User, 2017, "role not found", http.StatusNotFound),                                          // 角色不存在
	Incorrect:               New(locales.User, 2018, "incorrect user", http.StatusBadRequest),                                        // 用户信息错误
	OldPassword:             New(locales.User, 2019, "old password incorrect", http.StatusBadRequest),                                // 旧密码错误
	CaptchaIDNotFound:       New(locales.User, 2020, "captcha id not found", http.StatusBadRequest),                                  // 验证码ID不存在
}
