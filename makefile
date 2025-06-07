export GO111MODULE=on
# update app name. this is the name of binary
APP=dns-stress
APP_EXECUTABLE="./out/$(APP)"
ALL_PACKAGES=$(shell go list ./...)
SHELL := /bin/sh # Use bash syntax
TAG="latest"

# Optional colors to beautify output
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

## Quality
check-quality: ## runs code quality checks
	make lint
	make fmt
	make vet

# Append || true below if blocking local developement
lint: ## go linting. Update and use specific lint tool and options
	golangci-lint run ./...

vet: ## go vet
	go vet ./...

fmt: ## runs go formatter
	go fmt ./...

tidy: ## runs tidy to fix go.mod dependencies
	go mod tidy

## Test
test: ## runs tests and create generates coverage report
	make tidy
	make vendor
	go test -v -timeout 10m ./... -coverprofile=coverage.out -json > report.json

coverage: ## displays test coverage report in html mode
	make test
	go tool cover -html=coverage.out

## Build
build: ## build the go application
	mkdir -p out/
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -extldflags '-static'" -o $(APP_EXECUTABLE) 
	@echo "Build passed"

dockerize: ## Create docker image
	docker build --no-cache -t kunaev/dns-stress:$(TAG) .

run: ## runs the go binary. use additional options if required.
	make dockerize
	docker run --rm kunaev/dns-stress:$(TAG)

clean: ## cleans binary and other generated files
	go clean
	rm -rf out/
	rm -f coverage*.out
	rm -f report.json
	rm -rf vendor/

vendor: ## all packages required to support builds and tests in the /vendor directory
	go mod vendor

.PHONY: all test build vendor
## All
all: ## runs setup, quality checks and builds
	make check-quality
	make test
	make build

.PHONY: help
## Help
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)