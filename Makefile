dev:
	air

build:
	go build -o ./tmp/api ./cmd/api/main.go

lint:
	golangci-lint run ./... 

check:
	go build -v ./...

test:
	go test -v -race ./internal/...

.PHONY: dev, build, lint, test