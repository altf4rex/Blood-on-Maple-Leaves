package repo

import (
	"context"
	"testing"
	"time"

	"blood-on-maple-leaves/backend/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgres(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		Env:          map[string]string{"POSTGRES_USER": "test", "POSTGRES_PASSWORD": "test", "POSTGRES_DB": "test"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}
	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req, Started: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	host, _ := pgC.Host(ctx)
	port, _ := pgC.MappedPort(ctx, "5432")
	dsn := "postgres://test:test@" + host + ":" + port.Port() + "/test?sslmode=disable"
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	// run migrations here (or exec SQL directly)
	// e.g. pool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; ... create tables`)
	return pool, func() {
		pool.Close()
		pgC.Terminate(ctx)
	}
}

func TestSaveRepoPG_CreateAndGetLatest(t *testing.T) {
	pool, teardown := setupPostgres(t)
	defer teardown()

	repo := NewSaveRepoPG(pool)
	playerID := uuid.New()
	save1 := domain.Save{ID: uuid.New(), PlayerID: playerID, SceneID: "intro", Honor: 0, Rage: 0, Karma: 0, CreatedAt: time.Now()}
	if err := repo.Create(context.Background(), save1); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	// wait a bit and insert second
	time.Sleep(10 * time.Millisecond)
	save2 := domain.Save{ID: uuid.New(), PlayerID: playerID, SceneID: "hallway", Honor: 0, Rage: 1, Karma: 0, CreatedAt: time.Now()}
	if err := repo.Create(context.Background(), save2); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	latest, err := repo.GetLatestByPlayer(context.Background(), playerID)
	if err != nil {
		t.Fatalf("GetLatestByPlayer failed: %v", err)
	}
	if latest.SceneID != "hallway" {
		t.Errorf("expected latest scene 'hallway', got '%s'", latest.SceneID)
	}
	if latest.Rage != 1 {
		t.Errorf("expected rage=1, got=%d", latest.Rage)
	}
}
