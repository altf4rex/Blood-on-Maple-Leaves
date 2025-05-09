package handlers

import (
	"encoding/json"
	"net/http"

	"blood-on-maple-leaves/backend/middleware"
	"blood-on-maple-leaves/backend/service"
)

// MeHandler возвращает информацию о текущем игроке.
// Захватывает authSvc и читает userID из контекста, установленного в middleware.
func MeHandler(authSvc *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Извлечь userID из контекста
		uidVal := r.Context().Value(middleware.ContextUserID)
		userID, ok := uidVal.(string)
		if !ok || userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// 2. Получить данные игрока по ID
		player, err := authSvc.PlayerRepo.GetByID(r.Context(), userID)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		// 3. Ответить JSON-ом с данными игрока
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(player)
	}
}
