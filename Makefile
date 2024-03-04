
GO_ENV=CGO_ENABLED=0 GO111MODULE=on
Revision=$(shell git rev-parse --short HEAD 2>/dev/null || echo "")
GO_FLAGS=-ldflags="-X 'github.com/htquangg/a-wasm/cmd.Revision=$(Revision)' -X 'github.com/apache/incubator-answer/cmd.Time=`date +%s`' -extldflags -static"
GO=$(GO_ENV) $(shell which go)

ifndef GOPATH
	GOPATH := $(shell go env GOPATH)
endif
ifndef GOBIN # derive value from gopath (default to first entry, similar to 'go get')
	GOBIN := $(shell go env GOPATH | sed 's/:.*//')/bin
endif

BIN=awasm

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
# https://copyprogramming.com/howto/how-to-pass-argument-to-makefile-from-command-line
%:      # thanks to chakrit
    @:    # thanks to William Pursell

.PHONY: build
build: ## build the service binary
	@echo "Building $(BIN) server"
	@$(GO) build $(GO_FLAGS) -o $(BIN) main.go

.PHONY: build
clean: ## clean all build result
	@$(GO) clean ./...
	@rm -f $(BIN)

.PHONY: dev
dev: ## run the dev application
	@go run main.go $(filter-out $@,$(MAKECMDGOALS))

.PHONY: dev-run
dev-run: ## run the dev application and serve
	@go run main.go run

.PHONY: air
air: ## live reloading the application
	@air -c .air.toml -- $(filter-out $@,$(MAKECMDGOALS))

.PHONY: air-run
air-run: ## live reloading the application and serve
	@air -c .air.toml -- run

.PHONY: test
test: ## run the go tests
	@echo "Running tests"
	go test ./... -v --cover

test-report: ## run the go tests and report
	@echo "Running tests"
	go test ./... -v --cover -coverprofile=coverage.out
	@echo "Reporting tests"
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
DB_DRIVER =  postgres
DB_NAME = dev-local-awasm-001
DB_HOST = 127.0.0.1
DB_PORT = 5432
DB_USER = postgres
DB_PASS = localdb
DB_PARSE_TIME = true

# Go migrate postgres https://github.com/pressly/goose
.PHONY: migrate-create
migrate-create: ## create new migration file
	@./scripts/goose-migrate.sh -p ./migrations/schemas -c create ${name}

.PHONY: migrate-up
migrate-up: ## migrate the DB to the most recent version available
	@./scripts/goose-migrate.sh -p ./migrations/schemas -c up -o "${DB_DRIVER}://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

.PHONY: migrate-down
migrate-down: ## roll back the version by 1
	@./scripts/goose-migrate.sh -p ./migrations/schemas -c down -o "${DB_DRIVER}://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

.PHONY: migrate-status
migrate-status: ## check the migration status for the current DB
	@./scripts/goose-migrate.sh -p ./migrations/schemas -c status -o "${DB_DRIVER}://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

.PHONY: help
help: ## print help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
