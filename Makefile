dev:
	air

build:
	go build -o ./tmp/api ./cmd/api/main.go

.PHONY: dev, build