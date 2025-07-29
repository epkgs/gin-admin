package dtos

import "strings"

type Login struct {
	Username string `json:"username" binding:"required"` // Login name
	Password string `json:"password" binding:"required"` // Login password (md5 hash)
	// CaptchaID   string `json:"captchaId" binding:"required"`   // Captcha verify id
	// CaptchaCode string `json:"captchaCode" binding:"required"` // Captcha verify code
}

func (a *Login) Trim() *Login {
	a.Username = strings.TrimSpace(a.Username)
	// a.CaptchaCode = strings.TrimSpace(a.CaptchaCode)
	return a
}

type LoginToken struct {
	AccessToken  string `json:"accessToken"`  // Access token (JWT)
	TokenType    string `json:"tokenType"`    // Token type (Usage: Authorization=${token_type} ${access_token})
	Expires      int64  `json:"expires"`      // Expired time (second)s
	RefreshToken string `json:"refreshToken"` // Refresh token (JWT)
}

type AuthUpdatePasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required"` // Old password (md5 hash)
	NewPassword string `json:"newPassword" binding:"required"` // New password (md5 hash)
}

type AuthUpdateUserReq struct {
	Name   *string `json:"name" binding:"omitempty,max=64"`     // Name of user
	Wechat *string `json:"wechat" binding:"omitempty,max=64"`   // Wechat account
	Phone  *string `json:"phone" binding:"omitempty,max=32"`    // Phone number of user
	Email  *string `json:"email" binding:"omitempty,max=128"`   // Email of user
	Remark *string `json:"remark" binding:"omitempty,max=1024"` // Remark of user
}

type Captcha struct {
	CaptchaID string `json:"captchaId"` // Captcha ID
}

type CaptchaImageReq struct {
	ID     string `from:"id"`     // Captcha ID
	Reload bool   `from:"reload"` // Reload captcha image (reload=1)
}
