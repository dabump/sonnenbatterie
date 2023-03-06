GO_FLAGS	?=
NAME	 	:= sonnenbatterie
PACKAGE		:= github.com/dabump/$(NAME)
BIN_PREFIX	:= bin

default: help

.PHONY: test
test: ## Run all tests
	@go clean --testcache && go test ./...

generate: ## Generate mocks
	@go generate ./...

cover: ## Run test coverage suite
	@go test ./... --coverprofile=cov.out
	@go tool cover --html=cov.out

build: ## Builds the token bucket
	@go build ${GO_FLAGS} -o ${BIN_PREFIX}/status ./cmd/status/main.go

clean: ## Cleans the build binaries
	@rm -rf ${BIN_PREFIX}

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[38;5;69m%-30s\033[38;5;38m %s\033[0m\n", $$1, $$2}'