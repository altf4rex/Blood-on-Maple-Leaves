package service

import (
	"errors"
	"testing"

	"blood-on-maple-leaves/backend/domain"
)

// fakeSceneRepo — фейковая реализация SceneRepo для unit-тестов.
type fakeSceneRepo struct {
	scenes map[string]domain.Scene
}

func (f *fakeSceneRepo) Load(id string) (domain.Scene, error) {
	scene, ok := f.scenes[id]
	if !ok {
		return domain.Scene{}, errors.New("scene not found")
	}
	return scene, nil
}

func TestApplyChoice(t *testing.T) {
	scene := domain.Scene{
		ID:   "intro",
		Text: "Text",
		Choices: []domain.Choice{
			{ID: "attack", Text: "A", Next: "hallway", Effects: map[string]int{"rage": 1}},
			{ID: "sneak", Text: "S", Next: "backdoor", Effects: map[string]int{"honor": 1}},
		},
	}

	cases := []struct {
		name     string
		scene    domain.Scene
		choiceID string
		wantErr  bool
	}{
		{"valid attack", scene, "attack", false},
		{"valid sneak", scene, "sneak", false},
		{"invalid", scene, "run", true},
	}

	svc := NewGameService(nil, nil) // SaveRepo не нужен тут

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			choice, err := svc.ApplyChoice(tc.scene, tc.choiceID)
			if tc.wantErr {
				if err == nil {
					t.Errorf("ApplyChoice(%s) expected error", tc.choiceID)
				}
			} else {
				if err != nil {
					t.Errorf("ApplyChoice(%s) unexpected error: %v", tc.choiceID, err)
				}
				if choice.ID != tc.choiceID {
					t.Errorf("got ID=%s; want %s", choice.ID, tc.choiceID)
				}
			}
		})
	}
}

func TestChoose(t *testing.T) {
	scene := domain.Scene{
		ID:   "intro",
		Text: "Text",
		Choices: []domain.Choice{
			{ID: "attack", Text: "A", Next: "hallway", Effects: map[string]int{"rage": 1}},
		},
	}

	fakeRepo := &fakeSceneRepo{scenes: map[string]domain.Scene{"intro": scene}}
	svc := NewGameService(fakeRepo, nil)

	cases := []struct {
		name     string
		sceneID  string
		choiceID string
		wantNext string
		wantErr  bool
	}{
		{"ok", "intro", "attack", "hallway", false},
		{"badScene", "bad", "attack", "", true},
		{"badChoice", "intro", "run", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			next, err := svc.Choose(tc.sceneID, tc.choiceID)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Choose(%s,%s) expected error", tc.sceneID, tc.choiceID)
				}
			} else {
				if err != nil {
					t.Errorf("Choose(%s,%s) unexpected error: %v", tc.sceneID, tc.choiceID, err)
				}
				if next != tc.wantNext {
					t.Errorf("got next=%s; want %s", next, tc.wantNext)
				}
			}
		})
	}
}
