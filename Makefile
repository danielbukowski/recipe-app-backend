include .env

dev:
	air

build:
	go build -o ./tmp/api ./cmd/api/main.go

lint-code:
	golangci-lint run ./... 

generate-sql:
	DATABASE_URL=${DATABASE_URL} sqlc generate

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

generate-docs:
	swag init -dir ./cmd/api/

.PHONY: dev, build, lint-code, generate-sql, test, db-up, db-reset, db-check-migration-files, check-build, lint-queries, generate-swagger-docs