package repo

import (
	"os"
	"path/filepath"

	"blood-on-maple-leaves/backend/domain"

	"gopkg.in/yaml.v3"
)

// SceneRepo — интерфейс, описывающий загрузку сцен
type SceneRepo interface {
	Load(sceneID string) (domain.Scene, error)
}

// SceneRepoFS — файловая реализация, читает YAML-файлы из папки
type SceneRepoFS struct {
	BasePath string // путь к папке со сценами (например, "./scenes")
}

// NewSceneRepoFS — конструктор, принимает путь к папке
func NewSceneRepoFS(basePath string) *SceneRepoFS {
	return &SceneRepoFS{BasePath: basePath}
}

// Load загружает YAML-файл и возвращает сцену
func (r *SceneRepoFS) Load(sceneID string) (domain.Scene, error) {
	var scene domain.Scene

	// 1. Формируем путь к файлу безопасно
	path := filepath.Join(r.BasePath, sceneID+".yaml")

	// 2. Читаем YAML-файл
	data, err := os.ReadFile(path)
	if err != nil {
		return scene, err
	}

	// 3. Парсим YAML в структуру
	err = yaml.Unmarshal(data, &scene)
	if err != nil {
		return scene, err
	}

	return scene, nil
}
