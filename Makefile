.PHONY: init test test-migration test-pgx test-gorm test-service run tidy generate

init:
	@go install github.com/google/wire/cmd/wire@latest

test:
	@go test ./...

test-migration:
	@go test ./pkg/usecases/migrations/goose

test-pgx:
	@go test ./internal/usecases/database/postgres

test-gorm:
	@go test ./internal/usecases/database/gorm

test-service:
	@go test ./internal/usecases/service

run:
	@go run ./cmd/app/

tidy:
	@go mod tidy

generate: init
	@go generate ./...