BINARY      = poly-api
MAIN_API    = ./cmd/api
MAIN_WORKER = ./cmd/worker
MIGRATE_PATH = ./migrations
DB_URL      ?= $(DATABASE_URL)

.PHONY: run build test sqlc migrate-up migrate-down migrate-create lint tidy

run:
	go run $(MAIN_API)

run-worker:
	go run $(MAIN_WORKER)

build:
	go build -o bin/$(BINARY) $(MAIN_API)

test:
	go test ./internal/domain/... -v -count=1

test-all:
	go test ./... -v -count=1

sqlc:
	sqlc generate

migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $$name

tidy:
	go mod tidy

lint:
	golangci-lint run ./...
