package service

import (
	"errors"
	"testing"

	"blood-on-maple-leaves/backend/domain"
)

// fakeSceneRepo — фейковая реализация SceneRepo для unit-тестов.
// Она хранит сцены в map по их ID и позволяет имитировать поведение репозитория.
// Это аналог моков или stub'ов в других языках.
type fakeSceneRepo struct {
	scenes map[string]domain.Scene
}

// Load возвращает сцену по ID или ошибку, если сцена не найдена.
func (f *fakeSceneRepo) Load(id string) (domain.Scene, error) {
	scene, ok := f.scenes[id]
	if !ok {
		// возвращаем пустую структуру и ошибку, когда id нет в карте
		return domain.Scene{}, errors.New("scene not found")
	}
	return scene, nil
}

// TestApplyChoice покрывает метод ApplyChoice
// table-driven tests: кейсы описываются в слайсе, а затем прогоняются через t.Run.
func TestApplyChoice(t *testing.T) {
	// Подготавливаем пример сцены с двумя вариантами
	exampleScene := domain.Scene{
		ID:   "intro",
		Text: "Текст сцены",
		Choices: []domain.Choice{
			{ID: "attack", Text: "Атаковать", Next: "hallway", Effects: map[string]int{"rage": 1}},
			{ID: "sneak", Text: "Пробраться", Next: "backdoor", Effects: map[string]int{"honor": 1}},
		},
	}

	// Таблица тестовых случаев
	cases := []struct {
		name     string       // имя кейса, выводится в t.Run
		scene    domain.Scene // входная сцена
		choiceID string       // ID выбора
		wantErr  bool         // ожидаем ошибку или нет
	}{
		{"valid attack", exampleScene, "attack", false},
		{"valid sneak", exampleScene, "sneak", false},
		{"invalid choice", exampleScene, "runaway", true},
	}

	svc := NewGameService(nil) // SceneRepo не нужен для ApplyChoice

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Вызываем ApplyChoice с заданной сценой и ID выбора
			choice, err := svc.ApplyChoice(tc.scene, tc.choiceID)
			// Проверяем наличие или отсутствие ошибки
			if tc.wantErr {
				if err == nil {
					t.Errorf("ApplyChoice(%s) expected error, got nil", tc.choiceID)
				}
				return
			}
			if err != nil {
				t.Errorf("ApplyChoice(%s) unexpected error: %v", tc.choiceID, err)
				return
			}
			// Если ошибки нет, проверяем, что вернулся правильный Choice
			if choice.ID != tc.choiceID {
				t.Errorf("ApplyChoice returned choice.ID=%s, want %s", choice.ID, tc.choiceID)
			}
		})
	}
}

// TestChoose покрывает метод Choose, используя fakeSceneRepo
func TestChoose(t *testing.T) {
	// Готовим фиктивную сцену для репозитория
	sceneIntro := domain.Scene{
		ID:   "intro",
		Text: "Начальная сцена",
		Choices: []domain.Choice{
			{ID: "attack", Text: "Атаковать", Next: "hallway", Effects: map[string]int{"rage": 1}},
		},
	}

	// Инициализируем fakeSceneRepo с нашей сценой
	fakeRepo := &fakeSceneRepo{
		scenes: map[string]domain.Scene{
			"intro": sceneIntro,
		},
	}

	// Создаём GameService, инжектируем fakeRepo
	svc := NewGameService(fakeRepo)

	// Таблица тестов для Choose
	cases := []struct {
		name     string
		sceneID  string
		choiceID string
		wantNext string
		wantErr  bool
	}{
		{"valid choose", "intro", "attack", "hallway", false},
		{"invalid sceneID", "unknown", "attack", "", true},
		{"invalid choiceID", "intro", "runaway", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			next, err := svc.Choose(tc.sceneID, tc.choiceID)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Choose(%s, %s) expected error, got nil", tc.sceneID, tc.choiceID)
				}
				return
			}
			if err != nil {
				t.Errorf("Choose(%s, %s) unexpected error: %v", tc.sceneID, tc.choiceID, err)
				return
			}
			// Проверяем, что вернулся правильный следующий ID сцены
			if next != tc.wantNext {
				t.Errorf("Choose(%s, %s) = %s; want %s", tc.sceneID, tc.choiceID, next, tc.wantNext)
			}
		})
	}
}
