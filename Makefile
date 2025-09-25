
.PHONY: help all build install test cover lint tidy docs clean debug-mcp login-ghcr install-tools pin-actions fix-struct-tags doc-website

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: tidy build test lint docs install ## Build, test, lint, and generate docs

build: ## Build the binary
	mkdir -p bin
	go build -o bin/backlog main.go

install: ## Install the binary
	go install .

test: ## Run tests
	go test -race -v ./...

cover: ## Run tests with coverage
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

lint: ## Linting
	go vet ./...

tidy: ## Run go mod tidy on all modules
	go mod tidy
	
docs: ## Generate documentation
	rm -rf ./docs/cli
	go generate -x ./...

clean: ## Clean up build artifacts
	rm -rf bin/
	rm -rf coverage.out coverage.html

release: ## push a new release
	git push origin main
	git push origin main --tags
	goreleaser release --clean

release-test: ## test a release locally
	goreleaser release --snapshot --clean

debug-mcp: ## Debug MCP issues
	npx @modelcontextprotocol/inspector go run . mcp

login-ghcr:
	@echo $GITHUB_TOKEN | docker login ghcr.io -u veggiemonk --password-stdin

install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/stacklok/frizbee@latest
	go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest

pin-actions: ## pin github actions
	frizbee actions .github/workflows

fix-struct-tags: ## format struct tags
	golangci-lint run -E tagalign --fix

doc-website: ## documentation website
	uv run mkdocs serve
