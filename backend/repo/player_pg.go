package repo

import (
	"context"

	"blood-on-maple-leaves/backend/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PlayerRepo struct {
	DB *pgxpool.Pool
}

func NewPlayerRepo(db *pgxpool.Pool) *PlayerRepo {
	return &PlayerRepo{DB: db}
}

func (r *PlayerRepo) Create(ctx context.Context, p *domain.Player) error {
	_, err := r.DB.Exec(ctx,
		`INSERT INTO players (id, username, password_hash, created_at)
		 VALUES ($1, $2, $3, $4)`,
		p.ID, p.Username, p.PasswordHash, p.CreatedAt,
	)
	return err
}

func (r *PlayerRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM players WHERE username = $1)`,
		username,
	).Scan(&exists)

	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PlayerRepo) GetByUsername(ctx context.Context, username string) (*domain.Player, error) {
	var p domain.Player
	err := r.DB.QueryRow(ctx,
		`SELECT id, username, password_hash, created_at FROM players WHERE username = $1`,
		username,
	).Scan(&p.ID, &p.Username, &p.PasswordHash, &p.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PlayerRepo) GetByID(ctx context.Context, id string) (*domain.Player, error) {
	var p domain.Player
	err := r.DB.QueryRow(ctx,
		`SELECT id, username, password_hash, created_at FROM players WHERE id = $1`,
		id,
	).Scan(&p.ID, &p.Username, &p.PasswordHash, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
