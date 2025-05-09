package middleware

import (
	"context"
	"net/http"
	"strings"

	"blood-on-maple-leaves/backend/internal/token"
)

// Ключ для userID в контексте
type contextKey string

const ContextUserID contextKey = "userID"

// AuthMiddleware — middleware для проверки access-токена
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Извлекаем заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		// 2. Проверяем формат "Bearer <token>"
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "invalid Authorization format", http.StatusUnauthorized)
			return
		}

		// 3. Получаем сам токен (без "Bearer ")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 4. Проверяем токен через JWT-модуль
		claims, err := token.VerifyAccessToken(tokenStr)
		if err != nil {
			http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		// 5. Получаем userID из токена
		userID, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		// 6. Добавляем userID в контекст запроса
		ctx := context.WithValue(r.Context(), ContextUserID, userID)

		// 7. Передаём запрос дальше, уже с userID в контексте
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
