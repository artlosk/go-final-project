package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const tokenTTL = 8 * time.Hour

type Claims struct {
	PassHash string `json:"pass_hash"`
	jwt.RegisteredClaims
}

func hashPassword(pass string) string {
	sum := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(sum[:])
}

func MakeToken(pass string) (string, error) {
	if pass == "" {
		return "", errors.New("пароль не введен")
	}
	claims := Claims{
		PassHash: hashPassword(pass),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(pass))
}

func ValidateToken(tokenString, pass string) bool {
	if tokenString == "" || pass == "" {
		return false
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("ошибка подписи токена")
		}
		return []byte(pass), nil
	})
	if err != nil || !token.Valid {
		return false
	}
	return claims.PassHash == hashPassword(pass)
}
