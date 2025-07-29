package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
)

var userI18n = i18n.NewCatalog("user")

func init() {
	userI18n.LoadTranslations()
}

var (
	ErrUserNotFound            = Define(userI18n, 2000, "user not found", http.StatusNotFound)                                        // 用户不存在
	ErrUserNotLogin            = Define(userI18n, 2001, "user is not logged in", http.StatusUnauthorized)                             // 用户未登录
	ErrUserDisabled            = Definef[struct{ Name string }](userI18n, 2002, "user {{.Name}} is disabled", http.StatusForbidden)   // 用户 {{.Name}} 已被禁用
	ErrUserExists              = Definef[struct{ Name string }](userI18n, 2003, "user {{.Name}} already exists", http.StatusConflict) // 用户 {{.Name}} 已存在
	ErrPasswordExpired         = Define(userI18n, 2004, "password expired", http.StatusForbidden)                                     // 密码已过期
	ErrUsernamePassword        = Define(userI18n, 2005, "incorrect username or password", http.StatusUnauthorized)                    // 用户名或密码错误
	ErrUserTokenError          = Define(userI18n, 2006, "wrong user token", http.StatusUnauthorized)                                  // 用户令牌错误
	ErrGenVisitToken           = Define(userI18n, 2007, "generate visit token failed", http.StatusInternalServerError)                // 生成访问令牌失败
	ErrGenRefreshToken         = Define(userI18n, 2008, "generate refresh token failed", http.StatusInternalServerError)              // 生成刷新令牌失败
	ErrParseToken              = Define(userI18n, 2009, "parse token failed", http.StatusUnauthorized)                                // 解析令牌失败
	ErrInvalidToken            = Define(userI18n, 2010, "invalid token", http.StatusUnauthorized)                                     // 无效的令牌
	ErrPasswordEncrypt         = Define(userI18n, 2011, "password encrypt failed", http.StatusInternalServerError)                    // 密码加密失败
	ErrPasswordDecrypt         = Define(userI18n, 2012, "password decrypt failed", http.StatusInternalServerError)                    // 密码解密失败
	ErrUserNameOrPasswordEmpty = Define(userI18n, 2013, "username or password empty", http.StatusBadRequest)                          // 用户名或密码不能为空
	ErrPassword                = Define(userI18n, 2014, "password error", http.StatusUnauthorized)                                    // 密码错误
	ErrModifySuperUser         = Define(userI18n, 2015, "super user can not modify", http.StatusForbidden)                            // 超级用户不能修改
	ErrRoleCodeExists          = Define(userI18n, 2016, "role code already exists", http.StatusBadRequest)                            // 角色编码已存在
	ErrRoleNotFount            = Define(userI18n, 2017, "role not found", http.StatusNotFound)                                        // 角色不存在
	ErrUser                    = Define(userI18n, 2018, "incorrect user", http.StatusBadRequest)                                      // 用户信息错误
	ErrOldPassword             = Define(userI18n, 2019, "old password incorrect", http.StatusBadRequest)                              // 旧密码错误
	ErrCaptchaIDNotFound       = Define(userI18n, 2020, "captcha id not found", http.StatusBadRequest)                                // 验证码ID不存在
)
