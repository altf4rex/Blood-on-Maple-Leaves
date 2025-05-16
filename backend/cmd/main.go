package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"blood-on-maple-leaves/backend/handlers"
	"blood-on-maple-leaves/backend/middleware"
	"blood-on-maple-leaves/backend/repo"
	"blood-on-maple-leaves/backend/service"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// initPostgres создаёт и возвращает пул соединений к Postgres.
func initPostgres() *pgxpool.Pool {
	dsn := os.Getenv("DB_DSN")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	return pool
}

// initRedis создаёт и возвращает клиент Redis.
func initRedis() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis connect error: %v", err)
	}
	return rdb
}

// runMigrations пытается до 10 раз выполнить m.Up(), чтобы подождать Postgres.
func runMigrations(dsn, dir string) {
	var (
		m   *migrate.Migrate
		err error
	)
	for i := 0; i < 10; i++ {
		m, err = migrate.New("file://"+dir, dsn)
		if err == nil {
			if err = m.Up(); err == nil || err == migrate.ErrNoChange {
				log.Println("migrations applied")
				return
			}
		}
		log.Printf("migrations retry (%d/10): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("migrations failed: %v", err)
}

func main() {
	// 1) Прогон миграций до открытия пула Postgres
	dsn := os.Getenv("DB_DSN")
	runMigrations(dsn, "./migrations")

	// 2) Инициализация БД и кеша
	db := initPostgres()
	defer db.Close()
	rdb := initRedis()
	defer rdb.Close()

	// 3) Репозитории
	playerRepo := repo.NewPlayerRepo(db)
	tokenRepo := repo.NewTokenRepo(rdb)
	saveRepo := repo.NewSaveRepoPG(db)
	sceneRepo := repo.NewSceneRepoFS("./scenes")

	// 4) Сервисы
	authSvc := service.NewAuthService(playerRepo, tokenRepo)
	gameSvc := service.NewGameService(sceneRepo, saveRepo)

	// 5) HTTP-обработчики
	sceneH := handlers.NewSceneHandler(gameSvc)

	r := chi.NewRouter()
	r.Post("/signup", handlers.SignupHandler(authSvc))
	r.Post("/login", handlers.LoginHandler(authSvc))
	r.With(middleware.AuthMiddleware).Get("/me", handlers.MeHandler(authSvc))

	r.With(middleware.AuthMiddleware).Get("/scenes/{id}", sceneH.GetScene)
	r.With(middleware.AuthMiddleware).Post("/scenes/{id}/choose", sceneH.Choose)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Println("🚀 API started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
