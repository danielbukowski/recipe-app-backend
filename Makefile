include .env

dev:
	air

build:
	go build -o ./tmp/api ./cmd/api/main.go

lint:
	golangci-lint run ./... 

db-up:
	goose -dir $(GOOSE_MIGRATION_DIR) postgres $(GOOSE_DBSTRING) up

db-reset:
	goose -dir $(GOOSE_MIGRATION_DIR) postgres $(GOOSE_DBSTRING) reset

db-check-migration-files:
	goose -dir $(GOOSE_MIGRATION_DIR) validate

check-build:
	go build -v ./...

test:
	go test -v -race ./internal/...

.PHONY: dev, build, lint, test, db-up, db-reset, db-check-migration-files