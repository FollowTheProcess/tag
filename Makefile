.PHONY: help tidy fmt test lint cover clean check sloc gen build install uninstall
.DEFAULT_GOAL := help

export GOEXPERIMENT := loopvar

help: ## Show the list of available tasks
	@echo "Available Tasks:\n"
	@grep -E '^[a-zA-Z_0-9%-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-10s %s\n", $$1, $$2}'

tidy: ## Tidy dependencies in go.mod
	go mod tidy

fmt: ## Run go fmt on all source files
	go fmt ./...

test: ## Run the test suite
	go test -race ./...

lint: ## Run the linters and auto-fix if possible
	golangci-lint run --fix

clean: ## Remove build artifacts and other clutter
	go clean ./...
	rm -rf ./bin ./dist

check: test lint ## Run tests and linting in one go

sloc: ## Print lines of code (for fun)
	find . -name "*.go" | xargs wc -l | sort -nr | head

gen: ## Run go generate
	go generate ./...

build: gen ## Compile the project binary
	mkdir -p ./bin
	goreleaser build --single-target --skip-before --snapshot --clean --output ./bin/tag

install: uninstall build ## Install the project on your machine
	cp ./bin/tag ${GOBIN}

uninstall: ## Uninstall the project from your machine
	rm -rf ${GOBIN}/tag
