package repo

import (
	"blood-on-maple-leaves/backend/domain"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SaveRepo — контракт для работы с saves
type SaveRepo interface {
	// Create сохраняет новую запись в таблицу saves.
	Create(ctx context.Context, s domain.Save) error
	// GetLatestByPlayer возвращает последнее сохранение для данного игрока.
	GetLatestByPlayer(ctx context.Context, playerID uuid.UUID) (domain.Save, error)
}

// SaveRepoPG — конкретная реализация SaveRepo через pgxpool.Pool
type SaveRepoPG struct {
	DB *pgxpool.Pool
}

// NewSaveRepoPG — конструктор, принимает пул Postgres.
func NewSaveRepoPG(db *pgxpool.Pool) *SaveRepoPG {
	return &SaveRepoPG{
		DB: db,
	}
}

// todo: (SaveRepoPG) Create
func (r *SaveRepoPG) Create(ctx context.Context, s domain.Save) error {
	_, err := r.DB.Exec(
		ctx,
		`INSERT INTO saves (id, player_id, scene_id, honor, rage, karma, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		s.ID, s.PlayerID, s.SceneID, s.Honor, s.Rage, s.Karma, s.CreatedAt,
	)
	return err
}

// todo: (SaveRepoPG) GetLatestByPlayer
func (r *SaveRepoPG) GetLatestByPlayer(ctx context.Context, playerID uuid.UUID) (domain.Save, error) {
	var s domain.Save
	row := r.DB.QueryRow(
		ctx,
		`
		SELECT id, player_id, scene_id, honor, rage, karma, created_at
		FROM saves
		WHERE player_id = $1
		ORDER BY created_at DESC
		LIMIT 1
		`,
		playerID,
	)
	err := row.Scan(&s.ID, &s.PlayerID, &s.SceneID,
		&s.Honor, &s.Rage, &s.Karma, &s.CreatedAt)
	return s, err
}
