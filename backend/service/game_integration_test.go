package service

import (
	"context"
	"testing"
	"time"

	"blood-on-maple-leaves/backend/domain"
	"blood-on-maple-leaves/backend/repo"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
)

func setupService(t *testing.T) (*GameService, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{ /* same as above */ }
	pgC, _ := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
	host, _ := pgC.Host(ctx)
	port, _ := pgC.MappedPort(ctx, "5432")
	dsn := "postgres://test:test@" + host + ":" + port.Port() + "/test?sslmode=disable"
	pool, _ := pgxpool.New(ctx, dsn)
	// run migrations...
	saveRepo := repo.NewSaveRepoPG(pool)
	// fake scene repo with single scene+choice
	scene := domain.Scene{
		ID:      "intro",
		Text:    "…",
		Choices: []domain.Choice{{ID: "attack", Next: "hallway", Effects: map[string]int{"rage": 1}}},
	}
	sceneRepo := &repo.FakeSceneRepo{Scenes: map[string]domain.Scene{"intro": scene}}
	svc := NewGameService(sceneRepo, saveRepo)
	return svc, func() { pool.Close(); pgC.Terminate(ctx) }
}

func TestChooseForPlayer(t *testing.T) {
	svc, teardown := setupService(t)
	defer teardown()

	playerID := uuid.New()
	// First time: no save → GetLatestByPlayer returns error; handle by creating initial save manually
	initial := domain.Save{ID: uuid.New(), PlayerID: playerID, SceneID: "intro", Honor: 0, Rage: 0, Karma: 0, CreatedAt: time.Now()}
	svc.SaveRepo.Create(context.Background(), initial)

	next, newSave, err := svc.ChooseForPlayer(context.Background(), playerID, "intro", "attack")
	if err != nil {
		t.Fatalf("ChooseForPlayer failed: %v", err)
	}
	if next != "hallway" {
		t.Errorf("expected next 'hallway', got '%s'", next)
	}
	if newSave.Rage != 1 {
		t.Errorf("expected rage=1, got %d", newSave.Rage)
	}
}
