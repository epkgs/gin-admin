package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
)

var userI18n = i18n.NewCatalog("user")

func init() {
	userI18n.LoadTranslations()
}

type userErrors struct {
	NotFound                *Definition                         // 用户不存在
	NotLogin                *Definition                         // 用户未登录
	Disabled                *DefinitionF[struct{ Name string }] // 用户 {{.Name}} 已被禁用
	Exists                  *DefinitionF[struct{ Name string }] // 用户 {{.Name}} 已存在
	PasswordExpired         *Definition                         // 密码已过期
	UsernamePassword        *Definition                         // 用户名或密码错误
	UserTokenError          *Definition                         // 用户令牌错误
	GenVisitToken           *Definition                         // 生成访问令牌失败
	GenRefreshToken         *Definition                         // 生成刷新令牌失败
	ParseToken              *Definition                         // 解析令牌失败
	InvalidToken            *Definition                         // 无效的令牌
	PasswordEncrypt         *Definition                         // 密码加密失败
	PasswordDecrypt         *Definition                         // 密码解密失败
	UserNameOrPasswordEmpty *Definition                         // 用户名或密码不能为空
	Password                *Definition                         // 密码错误
	ModifySuperUser         *Definition                         // 超级用户不能修改
	RoleCodeExists          *Definition                         // 角色编码已存在
	RoleNotFount            *Definition                         // 角色不存在
	Incorrect               *Definition                         // 用户信息错误
	OldPassword             *Definition                         // 旧密码错误
	CaptchaIDNotFound       *Definition                         // 验证码ID不存在
}

var User = &userErrors{
	NotFound:                Define(userI18n, 2000, "user not found", http.StatusNotFound),                                        // 用户不存在
	NotLogin:                Define(userI18n, 2001, "user is not logged in", http.StatusUnauthorized),                             // 用户未登录
	Disabled:                Definef[struct{ Name string }](userI18n, 2002, "user {{.Name}} is disabled", http.StatusForbidden),   // 用户 {{.Name}} 已被禁用
	Exists:                  Definef[struct{ Name string }](userI18n, 2003, "user {{.Name}} already exists", http.StatusConflict), // 用户 {{.Name}} 已存在
	PasswordExpired:         Define(userI18n, 2004, "password expired", http.StatusForbidden),                                     // 密码已过期
	UsernamePassword:        Define(userI18n, 2005, "incorrect username or password", http.StatusUnauthorized),                    // 用户名或密码错误
	UserTokenError:          Define(userI18n, 2006, "wrong user token", http.StatusUnauthorized),                                  // 用户令牌错误
	GenVisitToken:           Define(userI18n, 2007, "generate visit token failed", http.StatusInternalServerError),                // 生成访问令牌失败
	GenRefreshToken:         Define(userI18n, 2008, "generate refresh token failed", http.StatusInternalServerError),              // 生成刷新令牌失败
	ParseToken:              Define(userI18n, 2009, "parse token failed", http.StatusUnauthorized),                                // 解析令牌失败
	InvalidToken:            Define(userI18n, 2010, "invalid token", http.StatusUnauthorized),                                     // 无效的令牌
	PasswordEncrypt:         Define(userI18n, 2011, "password encrypt failed", http.StatusInternalServerError),                    // 密码加密失败
	PasswordDecrypt:         Define(userI18n, 2012, "password decrypt failed", http.StatusInternalServerError),                    // 密码解密失败
	UserNameOrPasswordEmpty: Define(userI18n, 2013, "username or password empty", http.StatusBadRequest),                          // 用户名或密码不能为空
	Password:                Define(userI18n, 2014, "password error", http.StatusUnauthorized),                                    // 密码错误
	ModifySuperUser:         Define(userI18n, 2015, "super user can not modify", http.StatusForbidden),                            // 超级用户不能修改
	RoleCodeExists:          Define(userI18n, 2016, "role code already exists", http.StatusBadRequest),                            // 角色编码已存在
	RoleNotFount:            Define(userI18n, 2017, "role not found", http.StatusNotFound),                                        // 角色不存在
	Incorrect:               Define(userI18n, 2018, "incorrect user", http.StatusBadRequest),                                      // 用户信息错误
	OldPassword:             Define(userI18n, 2019, "old password incorrect", http.StatusBadRequest),                              // 旧密码错误
	CaptchaIDNotFound:       Define(userI18n, 2020, "captcha id not found", http.StatusBadRequest),                                // 验证码ID不存在
}
