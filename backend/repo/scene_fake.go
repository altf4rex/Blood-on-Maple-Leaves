package repo

import (
	"errors"

	"blood-on-maple-leaves/backend/domain"
)

// FakeSceneRepo — фейковая реализация SceneRepo для тестов.
type FakeSceneRepo struct {
	Scenes map[string]domain.Scene
}

// Load возвращает сцену по ID или ошибку, если нет в карте.
func (f *FakeSceneRepo) Load(id string) (domain.Scene, error) {
	scene, ok := f.Scenes[id]
	if !ok {
		return domain.Scene{}, errors.New("scene not found")
	}
	return scene, nil
}
