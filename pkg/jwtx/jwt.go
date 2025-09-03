package jwtx

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Auther interface {
	// Generate a JWT (JSON Web Token) with the provided userID.
	GenerateToken(ctx context.Context, userID string) (*TokenInfo, error)
	// Invalidate a token by removing it from the token store.
	DestroyToken(ctx context.Context, accessToken string) error
	// Parse from a given access token.
	ParseToken(ctx context.Context, accessToken string) (*TokenClaims, error)
	// Release any resources held by the JWTAuth instance.
	Release(ctx context.Context) error
}

// TokenClaims 自定义声明
type TokenClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

type TokenInfo struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
	TokenType    string `json:"tokenType"`
	Expires      int64  `json:"expires"`
}

// TokenType 令牌类型
type TokenType string

const (
	// AccessToken 访问令牌
	AccessToken TokenType = "access"
	// RefreshToken 刷新令牌
	RefreshToken TokenType = "refresh"
)

var DefaultSigningKey = "CG24SDVP8OHPK395GB5G"
var ErrInvalidToken = errors.New("invalid token")

type options struct {
	signingMethod jwt.SigningMethod
	signingKey    []byte
	expires       time.Duration // second
	tokenType     string
}

type Option func(*options)

func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

func WithSigningKey(key string) Option {
	return func(o *options) {
		o.signingKey = []byte(key)
	}
}

func WithExpired(expires int) Option {
	return func(o *options) {
		o.expires = time.Duration(expires) * time.Second
	}
}

func New(cacher Cacher, opts ...Option) Auther {
	o := options{
		tokenType:     "Bearer",
		expires:       2 * 24 * time.Hour,
		signingMethod: jwt.SigningMethodHS512,
		signingKey:    []byte(DefaultSigningKey),
	}

	for _, opt := range opts {
		opt(&o)
	}

	var store storer = nil
	if cacher != nil {
		store = newStore(cacher)
	}

	return &JWTAuth{
		opts:  &o,
		store: store,
	}
}

type JWTAuth struct {
	opts  *options
	store storer
}

func (a *JWTAuth) GenerateToken(ctx context.Context, userID string) (*TokenInfo, error) {
	r, _ := uuid.NewRandom()
	claimID := r.String()

	accessToken, expiredAt, err := a.generateToken(AccessToken, claimID, userID, a.opts.expires)
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := a.generateToken(RefreshToken, claimID, userID, a.opts.expires*10)
	if err != nil {
		return nil, err
	}

	tokenInfo := &TokenInfo{
		Expires:      expiredAt.Unix(),
		TokenType:    a.opts.tokenType,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return tokenInfo, nil
}

func (a *JWTAuth) generateToken(tokenType TokenType, claimID string, userID string, expiration time.Duration) (token string, expiredAt time.Time, err error) {

	now := time.Now()

	// 添加随机抖动，确保每次生成的token不同（最多10分钟）
	jitter := time.Duration(now.UnixNano()%600) * time.Second
	expiredAt = now.Add(expiration).Add(jitter)

	claims := TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        claimID,
			Subject:   string(tokenType),
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	jwtToken := jwt.NewWithClaims(a.opts.signingMethod, claims)
	token, err = jwtToken.SignedString(a.opts.signingKey)

	return
}

func (a *JWTAuth) parseToken(tokenStr string) (*TokenClaims, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (any, error) { return a.opts.signingKey, nil })
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (a *JWTAuth) callStore(fn func(storer) error) error {
	if store := a.store; store != nil {
		return fn(store)
	}
	return nil
}

func (a *JWTAuth) DestroyToken(ctx context.Context, tokenStr string) error {
	claims, err := a.parseToken(tokenStr)
	if err != nil {
		return err
	}

	expiresAt, _ := claims.GetExpirationTime()

	return a.callStore(func(store storer) error {
		expired := time.Until(expiresAt.Time)
		return store.Set(ctx, claims.ID, expired)
	})
}

func (a *JWTAuth) ParseToken(ctx context.Context, tokenStr string) (*TokenClaims, error) {
	if tokenStr == "" {
		return nil, ErrInvalidToken
	}

	claims, err := a.parseToken(tokenStr)
	if err != nil {
		return nil, err
	}

	err = a.callStore(func(store storer) error {
		if exists, err := store.Check(ctx, claims.ID); err != nil {
			return err
		} else if exists {
			return ErrInvalidToken
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (a *JWTAuth) Release(ctx context.Context) error {
	return a.callStore(func(store storer) error {
		return store.Close()
	})
}
