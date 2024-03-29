GO_FLAGS	?=
NAME		:= sonnenbatterie
PACKAGE		:= github.com/dabump/$(NAME)
BIN_PREFIX	:= bin
DOCKER_BIN	:= podman

default: help

.PHONY: test
test: ## Run all tests
	@go clean --testcache && go test ./...

docker: ## Build docker container & start
	@${DOCKER_BIN} build -t sonnenbatterie/daemon . && \
	${DOCKER_BIN} run -d --restart=always -p 8881:8881 --name sonnenbatterie-daemon sonnenbatterie/daemon

generate: ## Generate mocks
	@go generate ./...

cover: ## Run test coverage suite
	@go test ./... --coverprofile=cov.out
	@go tool cover --html=cov.out

build: ## Builds the token bucket
	@go build ${GO_FLAGS} -o ${BIN_PREFIX}/status ./cmd/status/main.go
	@go build ${GO_FLAGS} -o ${BIN_PREFIX}/daemon ./cmd/daemon/main.go

format: ## Format and organise imports
	@go install mvdan.cc/gofumpt@latest
	@gofumpt -l -w .

clean: ## Cleans the build binaries
	@rm -rf ${BIN_PREFIX}

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[38;5;69m%-30s\033[38;5;38m %s\033[0m\n", $$1, $$2}'