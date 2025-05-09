package handlers

import (
	"encoding/json"
	"net/http"

	"blood-on-maple-leaves/backend/service"
)

// SignupRequest — форма запроса для /signup
type SignupRequest struct {
	Username string `json:"username"` // {"username": "..."}
	Password string `json:"password"` // {"password": "..."}
}

// LoginRequest — форма запроса для /login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignupHandler возвращает http.HandlerFunc, замыкая authSvc
func SignupHandler(authSvc *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Распарсить тело запроса
		var req SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		// 2. Вызвать сервис регистрации
		tokens, err := authSvc.Signup(r.Context(), req.Username, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 3. Ответить JSON-ом с токенами
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tokens)
	}
}

// LoginHandler возвращает http.HandlerFunc для входа пользователя
func LoginHandler(authSvc *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Распарсить тело запроса
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		// 2. Вызвать сервис логина
		tokens, err := authSvc.Login(r.Context(), req.Username, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// 3. Ответить JSON-ом с токенами
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokens)
	}
}
