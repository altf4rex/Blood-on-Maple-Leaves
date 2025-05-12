package service

import (
	"blood-on-maple-leaves/backend/domain"
	"blood-on-maple-leaves/backend/repo"
	"errors"
)

type GameService struct {
	SceneRepo repo.SceneRepo
}

// NewGameService создаёт сервис с доступом к сценам.
func NewGameService(sceneRepo repo.SceneRepo) *GameService {
	return &GameService{
		SceneRepo: sceneRepo,
	}
}

// todo: ApplyChoice
// Это метод, который получает сцену и выбор игрока (choiceID).
// Он ищет нужный выбор по ID внутри сцены и возвращает его.
// Если выбор не найден — возвращаем ошибку.
// Далее мы будем применять "эффекты" из этого выбора (rage, honor, и т.д.).

func (g *GameService) ApplyChoice(scene domain.Scene, choiceID string) (domain.Choice, error) {
	for _, choice := range scene.Choices {
		if choice.ID == choiceID {
			return choice, nil
		}
	}

	return domain.Choice{}, errors.New("invalid choice")
}

// Choose загружает сцену, применяет выбор, и возвращает следующую сцену.
func (g *GameService) Choose(sceneID, choiceID string) (string, error) {
	// 1. Загружаем YAML-сцену
	scene, err := g.SceneRepo.Load(sceneID)
	if err != nil {
		return "", err
	}

	// 2. Применяем выбор
	choice, err := g.ApplyChoice(scene, choiceID)
	if err != nil {
		return "", err
	}

	// 3. Возвращаем ID следующей сцены
	return choice.Next, nil
}

func (g *GameService) GetScene(sceneID string) (domain.Scene, error) {
	return g.SceneRepo.Load(sceneID)
}
