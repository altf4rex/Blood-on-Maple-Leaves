package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"blood-on-maple-leaves/backend/handlers"
	"blood-on-maple-leaves/backend/middleware"
	"blood-on-maple-leaves/backend/repo"
	"blood-on-maple-leaves/backend/service"
)

func initPostgres() *pgxpool.Pool {
	dsn := os.Getenv("DB_DSN")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	return pool
}

func initRedis() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis connect error: %v", err)
	}
	return rdb
}

func main() {
	r := chi.NewRouter()

	// Инициализируем инфраструктуру
	db := initPostgres()
	defer db.Close()
	rdb := initRedis()
	defer rdb.Close()

	// Репозитории (конкретные реализации)
	playerRepo := repo.NewPlayerRepo(db)
	tokenRepo := repo.NewTokenRepo(rdb)

	// Сервис авторизации
	authSvc := service.NewAuthService(playerRepo, tokenRepo)

	// Маршруты
	r.Post("/signup", handlers.SignupHandler(authSvc))
	r.Post("/login", handlers.LoginHandler(authSvc))
	r.With(middleware.AuthMiddleware).Get("/me", handlers.MeHandler(authSvc))
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Println("🚀 API started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
