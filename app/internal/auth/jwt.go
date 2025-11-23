package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"vox/internal/config"
)

func GenerateAccessToken(cfg *config.Config, sub uint) (string, error) {
	secret := cfg.GetAccessTokenSecret()
	if secret == "" {
		return "", errors.New("ACCESS_TOKEN_SECRET not set")
	}
	claims := jwt.MapClaims{
		"sub":	sub,
		"exp":	time.Now().Add(cfg.GetAccessTokenTTL()).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(cfg *config.Config, sub uint) (string, error) {
	secret := cfg.GetRefreshTokenSecret()
	if secret == "" {
		return "", errors.New("REFRESH_TOKEN_SECRET not set")
	}
	claims := jwt.MapClaims{
		"sub":	sub,
		"exp":	time.Now().Add(cfg.GetRefreshTokenTTL()).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateAndExtractSub(tokenString, secret string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected alg %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}
	subf, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("sub not found")
	}
	return uint(subf), nil
}
