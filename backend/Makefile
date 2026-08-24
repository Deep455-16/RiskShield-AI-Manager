.PHONY: build test run migrate

build:
	go build -o server ./cmd/server

test:
	go test ./internal/...

run:
	go run ./cmd/server

migrate:
	MIGRATIONS_DIR=./migrations go run ./cmd/server
