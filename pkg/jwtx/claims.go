package jwtx

import "github.com/golang-jwt/jwt/v5"

type TokenClaims interface {
	jwt.Claims
}
