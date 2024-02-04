
GO_ENV=CGO_ENABLED=0 GO111MODULE=on

ifndef GOPATH
	GOPATH := $(shell go env GOPATH)
endif
ifndef GOBIN # derive value from gopath (default to first entry, similar to 'go get')
	GOBIN := $(shell go env GOPATH | sed 's/:.*//')/bin
endif

SERVICE_NAME=awasm

###############################################################################
#
# Initialization
#
###############################################################################
.PHONY: tidy
tidy: ## add missing and remove unused modules
	@go mod tidy

###############################################################################
#
# Build and testing rules
#
###############################################################################
.PHONY: build
build: ## build the service binary
	@echo "Building ${SERVICE_NAME} server"
	@go build -o ./bin/${SERVICE_NAME} main.go

.PHONY: dev
dev: ## run dev service
	@go run main.go

.PHONY: air
air: ## live reloading dev service
	@air -c .air.toml

.PHONY: test
test: ## run the go tests
	@echo "Running tests"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

###############################################################################
#
# Code formatting and linting
#
###############################################################################
.PHONY: lint
lint: ## lint
	@echo "Linting"
	golangci-lint run ./...

.PHONY: format
format: ## format
	@echo "Formating ..."
	golines -m 120 -w --ignore-generated .
	gofumpt -l -w .
	@echo "Formatting complete"

###############################################################################
#
# Database migration
#
###############################################################################
DB_NAME = dev-local-awasm-001
DB_HOST = localhost
DB_PORT = 3306
DB_USER = root
DB_PASS = toor
DB_PARSE_TIME = true

# Go migrate mysql https://github.com/pressly/goose
.PHONY: migrate-create
migrate-create: ## create new migration file
	@./scripts/goose-migrate.sh -p ./migrations/schemas -c create ${name}

.PHONY: migrate-up
migrate-up: ## migrate the DB to the most recent version available
	@./scripts/goose-migrate.sh -p ./migrations/schemas -c up -o "${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=${DB_PARSE_TIME}"

.PHONY: migrate-down
migrate-down: ## roll back the version by 1
	@./scripts/goose-migrate.sh -p ./migrations/schemas -c down -o "${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=${DB_PARSE_TIME}"

.PHONY: migrate-status
migrate-status: ## check the migration status for the current DB
	@./scripts/goose-migrate.sh -p ./migrations/schemas -c status -o "${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=${DB_PARSE_TIME}"

.PHONY: help
help: ## print help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
