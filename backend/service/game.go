package service

import (
	"context"
	"errors"
	"time"

	"blood-on-maple-leaves/backend/domain"
	"blood-on-maple-leaves/backend/repo"

	"github.com/google/uuid"
)

// GameService управляет игровой логикой: загрузкой сцен, применением выбора и сохранением прогресса.
type GameService struct {
	SceneRepo repo.SceneRepo // для загрузки YAML-сцен
	SaveRepo  repo.SaveRepo  // для чтения/записи прогресса из Postgres
}

// NewGameService создаёт сервис с необходимыми репозиториями.
func NewGameService(sceneRepo repo.SceneRepo, saveRepo repo.SaveRepo) *GameService {
	return &GameService{
		SceneRepo: sceneRepo,
		SaveRepo:  saveRepo,
	}
}

// ApplyChoice находит выбор по его идентификатору в сцене.
// Возвращает объект Choice или ошибку, если выбор не найден.
func (g *GameService) ApplyChoice(scene domain.Scene, choiceID string) (domain.Choice, error) {
	for _, choice := range scene.Choices {
		if choice.ID == choiceID {
			return choice, nil
		}
	}
	return domain.Choice{}, errors.New("invalid choice ID: " + choiceID)
}

// Choose загружает сцену и возвращает идентификатор следующей сцены после применения выбора.
func (g *GameService) Choose(sceneID, choiceID string) (string, error) {
	scene, err := g.SceneRepo.Load(sceneID)
	if err != nil {
		return "", err
	}
	choice, err := g.ApplyChoice(scene, choiceID)
	if err != nil {
		return "", err
	}
	return choice.Next, nil
}

// GetScene возвращает структуру сцены по её идентификатору.
func (g *GameService) GetScene(sceneID string) (domain.Scene, error) {
	return g.SceneRepo.Load(sceneID)
}

// ChooseForPlayer обрабатывает выбор игрока с учётом сохранённого прогресса.
// Он загружает последнюю запись Save, применяет выбранный вариант и сохраняет новое состояние.
func (g *GameService) ChooseForPlayer(
	ctx context.Context,
	playerID uuid.UUID,
	sceneID, choiceID string,
) (string, domain.Save, error) {
	current, err := g.SaveRepo.GetLatestByPlayer(ctx, playerID)
	if err != nil {
		return "", domain.Save{}, err
	}

	scene, err := g.SceneRepo.Load(sceneID)
	if err != nil {
		return "", domain.Save{}, err
	}

	choice, err := g.ApplyChoice(scene, choiceID)
	if err != nil {
		return "", domain.Save{}, err
	}

	newSave := domain.Save{
		ID:        uuid.New(),
		PlayerID:  playerID,
		SceneID:   choice.Next,
		Honor:     current.Honor + choice.Effects["honor"],
		Rage:      current.Rage + choice.Effects["rage"],
		Karma:     current.Karma + choice.Effects["karma"],
		CreatedAt: time.Now(),
	}

	if err := g.SaveRepo.Create(ctx, newSave); err != nil {
		return "", domain.Save{}, err
	}

	return choice.Next, newSave, nil
}

// GetLatestSave возвращает последнее сохранение игрока.
func (g *GameService) GetLatestSave(ctx context.Context, playerID uuid.UUID) (domain.Save, error) {
	return g.SaveRepo.GetLatestByPlayer(ctx, playerID)
}
