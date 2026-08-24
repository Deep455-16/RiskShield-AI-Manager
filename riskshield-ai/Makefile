.PHONY: all build up down logs test

all: up

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f

build:
	cd backend && go build -o server ./cmd/server

test:
	cd backend && go test ./internal/...

migrate:
	cd backend && go run cmd/server/main.go

seed:
	@echo "Demo data seeded automatically on startup with DEMO_MODE=true"
