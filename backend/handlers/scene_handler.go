package handlers

import (
	"encoding/json"
	"net/http"

	"blood-on-maple-leaves/backend/service"

	"github.com/go-chi/chi/v5"
)

// SceneHandler — HTTP-адаптер для работы со сценами.
// Инжектит GameService, чтобы handler не работал напрямую с репозиториями.
type SceneHandler struct {
	GameSvc *service.GameService
}

// NewSceneHandler — конструктор для SceneHandler.
// Принимает указатель на GameService и возвращает новый экземпляр handler’а.
func NewSceneHandler(gs *service.GameService) *SceneHandler {
	return &SceneHandler{
		GameSvc: gs,
	}
}

// GetScene обрабатывает GET /scenes/{id}.
// 1. Считывает параметр sceneID из URL.
// 2. Вызывает GameService.GetScene для загрузки сцены.
// 3. При ошибке возвращает 404 Not Found.
// 4. Иначе код 200 и JSON-ответ со всей структурой domain.Scene.
func (h *SceneHandler) GetScene(w http.ResponseWriter, r *http.Request) {
	// 1. Извлекаем sceneID из URL: "{id}"
	sceneID := chi.URLParam(r, "id")

	// 2. Загружаем сцену из репозитория через сервис
	scene, err := h.GameSvc.GetScene(sceneID)
	if err != nil {
		// 3. Если сцена не найдена — 404
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// 4. Успешный ответ: устанавливаем заголовок и код, кодируем в JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(scene)
}

// ChooseRequest описывает тело POST-запроса /scenes/{id}/choose
// В JSON должно быть: { "choice_id": "<ID варианта>" }
type ChooseRequest struct {
	ChoiceID string `json:"choice_id"`
}

// ChooseResponse описывает тело успешного ответа на выбор.
// Возвращаем ID следующей сцены.
type ChooseResponse struct {
	NextSceneID string `json:"next_scene_id"`
}

// Choose обрабатывает POST /scenes/{id}/choose.
// 1. Извлекает sceneID из URL.
// 2. Декодирует JSON-тело запроса в ChooseRequest.
// 3. Вызывает GameService.Choose для применения выбора.
// 4. При ошибке возвращает 400 Bad Request или 404 Not Found.
// 5. Иначе код 200 и JSON-ответ с NextSceneID.
func (h *SceneHandler) Choose(w http.ResponseWriter, r *http.Request) {
	// 1. Извлекаем sceneID из URL
	sceneID := chi.URLParam(r, "id")

	// 2. Декодируем тело запроса
	var req ChooseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Некорректный JSON → 400 Bad Request
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 3. Применяем выбор через GameService
	nextID, err := h.GameSvc.Choose(sceneID, req.ChoiceID)
	if err != nil {
		// Если сцена не найдена или выбор неверен → 400/404
		// Здесь можем дифференцировать, но для простоты — 400 Bad Request
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 4. Формируем ответ
	resp := ChooseResponse{NextSceneID: nextID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
