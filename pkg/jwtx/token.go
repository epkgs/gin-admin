package jwtx

import (
	jsoniter "github.com/json-iterator/go"
)

type TokenInfo interface {
	GetRefreshToken() string
	GetAccessToken() string
	GetTokenType() string
	GetExpires() int64
	EncodeToJSON() ([]byte, error)
}

type tokenInfo struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
	TokenType    string `json:"tokenType"`
	Expires      int64  `json:"expires"`
}

func (t *tokenInfo) GetRefreshToken() string {
	return t.RefreshToken
}

func (t *tokenInfo) GetAccessToken() string {
	return t.AccessToken
}

func (t *tokenInfo) GetTokenType() string {
	return t.TokenType
}

func (t *tokenInfo) GetExpires() int64 {
	return t.Expires
}

func (t *tokenInfo) EncodeToJSON() ([]byte, error) {
	return jsoniter.Marshal(t)
}
