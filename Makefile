include .env

dev:
	air

build:
	go build -o ./tmp/api ./cmd/api/main.go

lint-code:
	golangci-lint run ./... 

db-up:
	goose -dir $(GOOSE_MIGRATION_DIR) postgres $(DATABASE_URL) up

db-reset:
	goose -dir $(GOOSE_MIGRATION_DIR) postgres $(DATABASE_URL) reset

db-check-migration-files:
	goose -dir $(GOOSE_MIGRATION_DIR) validate

check-build:
	go build -v ./...

lint-queries:
	DATABASE_URL=${DATABASE_URL} sqlc vet

test:
	go test -v -race ./internal/...

.PHONY: dev, build, lint-code, test, db-up, db-reset, db-check-migration-files, lint-queries