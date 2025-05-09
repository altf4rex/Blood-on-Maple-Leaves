package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenRepo — структура для работы с Redis
type TokenRepo struct {
	RDB *redis.Client // клиент Redis
}

// NewTokenRepo — конструктор, создающий TokenRepo
func NewTokenRepo(rdb *redis.Client) *TokenRepo {
	return &TokenRepo{RDB: rdb}
}

// SaveRefreshToken сохраняет refresh-токен в Redis с TTL 30 дней
func (r *TokenRepo) SaveRefreshToken(ctx context.Context, token, userID string) error {
	key := fmt.Sprintf("refresh:%s", token) // ключ в Redis
	return r.RDB.Set(ctx, key, userID, 30*24*time.Hour).Err()
}

// GetUserIDByRefresh возвращает userID по refresh-токену
func (r *TokenRepo) GetUserIDByRefresh(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("refresh:%s", token)
	return r.RDB.Get(ctx, key).Result()
}

// DeleteRefreshToken — удаляет refresh-токен (logout)
func (r *TokenRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("refresh:%s", token)
	return r.RDB.Del(ctx, key).Err()
}
