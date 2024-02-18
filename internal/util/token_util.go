package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ostheperson/go-auth-service/internal/domain"
	"github.com/ostheperson/go-auth-service/internal/helper"
)

type JwtCustomClaims struct {
	Username string      `json:"username"`
	ID       uint        `json:"id"`
	Role     domain.Role `json:"role"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func CreateAccessToken(
	user *domain.Users,
	secret string,
	expiry int,
) (accessToken string, err error) {
	exp := jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * time.Hour))
	claims := &JwtCustomClaims{
		Username: user.Username,
		ID:       user.ID,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func CreateRefreshToken(
	user *domain.Users,
	secret string,
	expiry int,
) (refreshToken string, err error) {
	claimsRefresh := &JwtCustomRefreshClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return rt, err
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func VerifyAndExtract(requestToken string, secret string) (*JwtCustomClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}
	token, err := jwt.ParseWithClaims(
		requestToken,
		&JwtCustomClaims{},
		keyFunc,
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtCustomClaims)

	if !ok && !token.Valid {
		return claims, fmt.Errorf(helper.ErrInvalidToken)
	}

	if token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("internal error")
}
