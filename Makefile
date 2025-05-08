# запускает всё окружение
up:
	docker-compose up -d

# останавливает и удаляет
down:
	docker-compose down -v

# показывает статус контейнеров
ps:
	docker-compose ps

# «хвост» логов backend
logs-api:
	docker-compose logs -f api

# билд только backend-образа (кэш Dockerfile)
build:
	docker-compose build api
