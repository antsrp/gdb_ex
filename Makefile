.PHONY: init test run tidy generate

init:
	@go install github.com/google/wire/cmd/wire@latest

test:
	@go test ./...

run:
	@go run ./cmd/app/

tidy:
	@go mod tidy

generate: init
	@go generate ./...