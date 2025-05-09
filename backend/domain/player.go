package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Player — доменная сущность игрока
type Player struct {
	ID           uuid.UUID // Уникальный идентификатор
	Username     string    // Имя игрока
	PasswordHash string    // Хеш пароля (а не сам пароль)
	CreatedAt    time.Time
}

// NewPlayer — фабричная функция для создания нового игрока
func NewPlayer(username, rawPassword string) (Player, error) {
	// Проверка логина
	if len(username) < 3 {
		return Player{}, errors.New("username must be at least 3 characters")
	}

	// Проверка пароля
	if len(rawPassword) < 6 {
		return Player{}, errors.New("password must be at least 6 characters")
	}

	// Генерация ID
	id := uuid.New()

	// Хеширование пароля с помощью bcrypt (надежный алгоритм)
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return Player{}, err
	}
	passwordHash := string(hashBytes)

	// Возвращаем нового игрока
	return Player{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}, nil
}

// Проверка пароля при входе
func (p *Player) CheckPassword(rawPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.PasswordHash), []byte(rawPassword))
	return err == nil
}
