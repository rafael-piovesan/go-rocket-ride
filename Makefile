.SILENT: ; # no need for @
.DEFAULT: help # Running Make will run the help target

help: ## Show Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Install dependencies
	go mod tidy
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/segmentio/golines@latest
	go install github.com/vektra/mockery/cmd/mockery

migrate: ## Run database migrations
	migrate -path db/migrations -database ${DSN} up

fixtures: ## Load database fixtures for testing locally
	testfixtures -d postgres -c ${DSN} -D db/fixtures/local --dangerous-no-test-database-check

lint: ## Run linter
	golangci-lint  run

format: ## Format source code lines
	golines . -m 120 -w --ignore-generated

mock: ## Generate mocks for interfaces
	mockery -name Datastore
	mockery -name RideUseCase

integration: ## Run integration tests
	go test -v -tags=integration ./...

unit: ## Run unit tests
	go test -v -tags=unit ./...

server: ## Run API server locally
	go run cmd/api/main.go