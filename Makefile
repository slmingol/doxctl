.DEFAULT_GOAL := help

export GITHUB_TOKEN = ${GO_RELEASER_GITHUB_TOKEN}

##@ General

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: list
list: ## List all available targets (simple)
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null \
		| awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' \
		| sort \
		| egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

##@ Development

.PHONY: fmt
fmt: ## Format Go code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run --timeout=5m

.PHONY: test
test: ## Run tests with coverage
	go test -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: test-race
test-race: ## Run tests with race detector
	go test -v -race ./...

.PHONY: coverage-html
coverage-html: test ## Generate and open HTML coverage report
	go tool cover -html=coverage.txt

.PHONY: security-scan
security-scan: ## Run govulncheck for security vulnerabilities
	@which govulncheck > /dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

.PHONY: clean
clean: ## Remove build artifacts and coverage files
	rm -rf dist/
	rm -f doxctl
	rm -f coverage.txt coverage.out
	rm -f cmd_cov*.txt cmd_coverage.txt

.PHONY: all
all: clean fmt vet lint test build ## Run all checks and build

##@ Build & Release

.PHONY: build
build: ## Build binaries with goreleaser
	goreleaser build --clean --debug

.PHONY: dryrun
dryrun: ## Run goreleaser in snapshot mode (no publish)
	goreleaser --snapshot --skip-publish --clean --debug

.PHONY: install
install: ## Install binary locally
	go install doxctl

.PHONY: docker-build
docker-build: ## Build Docker image locally for testing
	docker build -t doxctl:local .

.PHONY: tag
tag: ## Bump version and create git tag
	scripts/version-up.sh --patch --apply

.PHONY: release
release: ## Build and release with goreleaser
	goreleaser release --clean

##@ Git Operations

.PHONY: add_commit_push
add_commit_push: ## Quick commit and push (use with caution)
	git add .
	git commit -m "Makefile commit"
	git push

.PHONY: commit
commit: ## Bump version, commit, and push
	make tag ; git add . ; git commit -m "Makefile commit" ; git push

##########################################################
# REFERENCES:
#   - https://towardsdatascience.com/how-to-create-a-cli-in-golang-with-cobra-d729641c7177
