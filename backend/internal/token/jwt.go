package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid" // импорт для GenerateRefreshToken
)

var jwtSecret = []byte("supersecret") // в продакшене — os.Getenv("JWT_SECRET")

// GenerateAccessToken создаёт JWT для указанного userID и срока жизни ttl
func GenerateAccessToken(userID string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,                     // subject — ID пользователя
		"exp": time.Now().Add(ttl).Unix(), // expiry — срок действия
		"iat": time.Now().Unix(),          // issued at — время создания
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret)
}

// VerifyAccessToken проверяет подпись и срок жизни JWT, возвращает claims
func VerifyAccessToken(tokenStr string) (jwt.MapClaims, error) {
	tkn, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// убеждаемся, что метод подписи — HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !tkn.Valid {
		return nil, errors.New("invalid or expired token")
	}
	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

// GenerateRefreshToken создаёт новый UUIDv4, который используем как refresh-токен
func GenerateRefreshToken() (string, error) {
	return uuid.NewString(), nil
}
