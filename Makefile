include .env

MIGRATIONS_DIR := ./migrations
DB_DSN := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: run test test-integration build migrate-up migrate-down migrate-status swagger

run:
	go run ./cmd/api

test:
	go test ./...

test-integration:
	go test -tags=integration ./...

build:
	go build ./...

swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.8.1 init -g cmd/api/main.go -o docs

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status
