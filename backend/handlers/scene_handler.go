package handlers

import (
	"encoding/json"
	"net/http"

	"blood-on-maple-leaves/backend/middleware"
	"blood-on-maple-leaves/backend/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// SceneHandler отвечает за HTTP-эндпоинты работы со сценами и прогрессом.
type SceneHandler struct {
	GameSvc *service.GameService
}

// NewSceneHandler создаёт SceneHandler с внедрённым GameService.
func NewSceneHandler(gs *service.GameService) *SceneHandler {
	return &SceneHandler{GameSvc: gs}
}

// GetScene обрабатывает GET /scenes/{id}.
// Возвращает JSON вида:
//
//	{
//	  "scene": { ...domain.Scene... },
//	  "stats": { "honor": X, "rage": Y, "karma": Z }
//	}
func (h *SceneHandler) GetScene(w http.ResponseWriter, r *http.Request) {
	sceneID := chi.URLParam(r, "id")

	// Загружаем сцену
	scene, err := h.GameSvc.GetScene(sceneID)
	if err != nil {
		http.Error(w, "scene not found", http.StatusNotFound)
		return
	}

	// Получаем playerID из контекста (AuthMiddleware)
	playerIDstr, ok := r.Context().Value(middleware.ContextUserID).(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	playerID, err := uuid.Parse(playerIDstr)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	// Получаем последнее сохранение
	save, err := h.GameSvc.GetLatestSave(r.Context(), playerID)
	// Если сохранения нет или ошибка, просто не добавляем stats

	// Формируем ответ
	resp := map[string]interface{}{"scene": scene}
	if err == nil {
		resp["stats"] = map[string]int{
			"honor": save.Honor,
			"rage":  save.Rage,
			"karma": save.Karma,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ChooseRequest описывает входной JSON для POST /scenes/{id}/choose.
type ChooseRequest struct {
	ChoiceID string `json:"choice_id"`
}

// ChooseResponse описывает выходной JSON после применения выбора.
type ChooseResponse struct {
	NextSceneID string         `json:"next_scene_id"`
	Stats       map[string]int `json:"stats"`
}

// Choose обрабатывает POST /scenes/{id}/choose.
// Принимает выбор игрока, сохраняет новое состояние и возвращает:
//
//	{
//	  "next_scene_id": "...",
//	  "stats": { "honor": X, "rage": Y, "karma": Z }
//	}
func (h *SceneHandler) Choose(w http.ResponseWriter, r *http.Request) {
	sceneID := chi.URLParam(r, "id")

	// Парсим тело запроса
	var req ChooseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Получаем playerID из контекста
	playerIDstr, ok := r.Context().Value(middleware.ContextUserID).(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	playerID, err := uuid.Parse(playerIDstr)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	// Применяем выбор и сохраняем новое состояние
	nextID, save, err := h.GameSvc.ChooseForPlayer(r.Context(), playerID, sceneID, req.ChoiceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Формируем ответ
	resp := ChooseResponse{
		NextSceneID: nextID,
		Stats: map[string]int{
			"honor": save.Honor,
			"rage":  save.Rage,
			"karma": save.Karma,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
