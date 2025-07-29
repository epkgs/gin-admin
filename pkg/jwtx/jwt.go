package jwtx

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Auther interface {
	// Generate a JWT (JSON Web Token) with the provided subject.
	GenerateToken(ctx context.Context, subject string) (TokenInfo, error)
	// Invalidate a token by removing it from the token store.
	DestroyToken(ctx context.Context, accessToken string) error
	// Parse from a given access token.
	ParseToken(ctx context.Context, accessToken string) (TokenClaims, error)
	// Parse from a given refresh token.
	ParseRefreshToken(ctx context.Context, refreshToken string) (TokenClaims, error)
	// Release any resources held by the JWTAuth instance.
	Release(ctx context.Context) error
}

const defaultSigningKey = "CG24SDVP8OHPK395GB5G"
const defaultRefreshKey = "IOW3846N73946NLS0"

var ErrInvalidToken = errors.New("invalid token")

type options struct {
	signingMethod jwt.SigningMethod
	signingKey    []byte
	signingKey2   []byte
	refreshKey    []byte
	keyFuncs      []func(*jwt.Token) (interface{}, error)
	expired       int // second
	tokenType     string
}

type Option func(*options)

func SetSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

func SetSigningKey(key, oldKey string) Option {
	return func(o *options) {
		o.signingKey = []byte(key)
		if oldKey != "" && key != oldKey {
			o.signingKey2 = []byte(oldKey)
		}
	}
}

func SetRefreshKey(key string) Option {
	return func(o *options) {
		o.refreshKey = []byte(key)
	}
}

func SetExpired(expired int) Option {
	return func(o *options) {
		o.expired = expired
	}
}

func New(store Storer, opts ...Option) Auther {
	o := options{
		tokenType:     "Bearer",
		expired:       7200,
		signingMethod: jwt.SigningMethodHS512,
		signingKey:    []byte(defaultSigningKey),
		refreshKey:    []byte(defaultRefreshKey),
	}

	for _, opt := range opts {
		opt(&o)
	}

	o.keyFuncs = append(o.keyFuncs, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return o.signingKey, nil
	})

	if o.signingKey2 != nil {
		o.keyFuncs = append(o.keyFuncs, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return o.signingKey2, nil
		})
	}

	return &JWTAuth{
		opts:  &o,
		store: store,
	}
}

type JWTAuth struct {
	opts  *options
	store Storer
}

func (a *JWTAuth) GenerateToken(ctx context.Context, subject string) (TokenInfo, error) {
	now := time.Now()
	r, _ := uuid.NewRandom()
	claimID := r.String()
	expiresAt := now.Add(time.Duration(a.opts.expired) * time.Second)
	refreshExpiresAt := now.Add(time.Duration(a.opts.expired*30) * time.Second)

	accessClaims := jwt.RegisteredClaims{}
	accessClaims.ID = claimID
	accessClaims.IssuedAt = &jwt.NumericDate{Time: now}
	accessClaims.ExpiresAt = &jwt.NumericDate{Time: expiresAt}
	accessClaims.NotBefore = &jwt.NumericDate{Time: now}
	accessClaims.Subject = subject
	accessToken := jwt.NewWithClaims(a.opts.signingMethod, &accessClaims)
	accessTokenStr, err := accessToken.SignedString(a.opts.signingKey)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.RegisteredClaims{}
	refreshClaims.ID = claimID
	refreshClaims.IssuedAt = &jwt.NumericDate{Time: now}
	refreshClaims.ExpiresAt = &jwt.NumericDate{Time: refreshExpiresAt}
	refreshClaims.NotBefore = &jwt.NumericDate{Time: now}
	refreshClaims.Subject = subject
	refreshToken := jwt.NewWithClaims(a.opts.signingMethod, &refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString(a.opts.refreshKey)
	if err != nil {
		return nil, err
	}

	tokenInfo := &tokenInfo{
		Expires:      expiresAt.Unix(),
		TokenType:    a.opts.tokenType,
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}
	return tokenInfo, nil
}

func (a *JWTAuth) parseToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	for _, keyFunc := range a.opts.keyFuncs {
		token, err = jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, keyFunc)
		if err != nil || token == nil || !token.Valid {
			continue
		}
		break
	}

	if err != nil || token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return token.Claims.(*jwt.RegisteredClaims), nil
}

func (a *JWTAuth) parseRefreshToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	token, err = jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return a.opts.refreshKey, nil
	})

	if err != nil || token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return token.Claims.(*jwt.RegisteredClaims), nil
}

func (a *JWTAuth) callStore(fn func(Storer) error) error {
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

	return a.callStore(func(store Storer) error {
		expired := time.Until(expiresAt.Time)
		return store.Set(ctx, claims.ID, expired)
	})
}

func (a *JWTAuth) ParseToken(ctx context.Context, tokenStr string) (TokenClaims, error) {
	if tokenStr == "" {
		return nil, ErrInvalidToken
	}

	claims, err := a.parseToken(tokenStr)
	if err != nil {
		return nil, err
	}

	err = a.callStore(func(store Storer) error {
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

func (a *JWTAuth) ParseRefreshToken(ctx context.Context, tokenStr string) (TokenClaims, error) {
	if tokenStr == "" {
		return nil, ErrInvalidToken
	}

	claims, err := a.parseRefreshToken(tokenStr)
	if err != nil {
		return nil, err
	}

	err = a.callStore(func(store Storer) error {
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
	return a.callStore(func(store Storer) error {
		return store.Close(ctx)
	})
}
