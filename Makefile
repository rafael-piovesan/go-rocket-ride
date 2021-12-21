.SILENT: ; # no need for @
.DEFAULT: help # Running Make will run the help target

help: ## Show Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Install dependencies
	go mod tidy
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/segmentio/golines@latest

migrate: ## Run database migrations
	migrate -path db/migrations -database ${DSN} up

integration:
	go test -v -tags=integration ./...
