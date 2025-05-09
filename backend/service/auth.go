package service

import (
	"context"
	"errors"
	"time"

	"blood-on-maple-leaves/backend/domain"
	"blood-on-maple-leaves/backend/internal/token"
	"blood-on-maple-leaves/backend/repo"
)

// Tokens — структура с двумя токенами, которую вернём клиенту
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AuthService — слой бизнес-логики авторизации
// Содержит указатели на конкретные реализации репозиториев
type AuthService struct {
	PlayerRepo *repo.PlayerRepo
	TokenRepo  *repo.TokenRepo
}

// NewAuthService — конструктор AuthService
// Принимает указатели на репозитории
func NewAuthService(playerRepo *repo.PlayerRepo, tokenRepo *repo.TokenRepo) *AuthService {
	return &AuthService{
		PlayerRepo: playerRepo,
		TokenRepo:  tokenRepo,
	}
}

// Signup — логика регистрации нового пользователя
func (s *AuthService) Signup(ctx context.Context, username, password string) (*Tokens, error) {
	// 1. Проверка уникальности username
	exists, err := s.PlayerRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// 2. Создание игрока (ID + хеш пароля)
	player, err := domain.NewPlayer(username, password)
	if err != nil {
		return nil, err
	}

	// 3. Сохранение в базу
	if err := s.PlayerRepo.Create(ctx, &player); err != nil {
		return nil, err
	}

	// 4. Генерация access-токена (TTL = 15 минут)
	accessToken, err := token.GenerateAccessToken(player.ID.String(), 15*time.Minute)
	if err != nil {
		return nil, err
	}

	// 5. Генерация refresh-токена (UUID)
	refreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// 6. Сохранение refresh-токена в Redis
	if err := s.TokenRepo.SaveRefreshToken(ctx, refreshToken, player.ID.String()); err != nil {
		return nil, err
	}

	// 7. Возврат токенов
	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login — логика входа существующего пользователя
func (s *AuthService) Login(ctx context.Context, username, password string) (*Tokens, error) {
	// 1. Получаем игрока по username
	player, err := s.PlayerRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 2. Проверяем пароль
	if !player.CheckPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	// 3. Генерация токенов
	accessToken, err := token.GenerateAccessToken(player.ID.String(), 15*time.Minute)
	if err != nil {
		return nil, err
	}
	refreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// 4. Сохранение refresh-токена
	if err := s.TokenRepo.SaveRefreshToken(ctx, refreshToken, player.ID.String()); err != nil {
		return nil, err
	}

	// 5. Возврат токенов
	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
