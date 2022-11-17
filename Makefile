SHELL = /bin/bash

# Build-time variables to inject into binaries
export VERSION = $(shell (test "$(shell git describe --tags)" = "$(shell git describe --tags --abbrev=0)" && echo $(shell git describe --tags)) || echo $(shell git describe --tags --abbrev=0)+git)
export GIT_VERSION = $(shell git describe --dirty --tags --always)
export GIT_COMMIT = $(shell git rev-parse HEAD)

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

# Build settings
REPO = $(shell go list -m)
BUILD_DIR = build
GO_ASMFLAGS = -asmflags "all=-trimpath=$(shell dirname $(PWD))"
GO_GCFLAGS = -gcflags "all=-trimpath=$(shell dirname $(PWD))"
GO_BUILD_ARGS = \
  $(GO_GCFLAGS) $(GO_ASMFLAGS) \
  -ldflags " \
    -X '$(REPO)/pkg/version.Version=$(SIMPLE_VERSION)' \
    -X '$(REPO)/pkg/version.GitVersion=$(GIT_VERSION)' \
    -X '$(REPO)/pkg/version.GitCommit=$(GIT_COMMIT)' \
  " \

export GO111MODULE = on
export CGO_ENABLED = 0
export PATH := $(PWD)/$(BUILD_DIR):$(PWD)/$(LOCALBIN):$(PATH)

GOLANGCI_VERSION ?= latest
GOLANGCI_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh"

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: fix
fix: ## Fixup files in the repo.
	go mod tidy
	go fmt ./...
	make setup-lint
	$(LOCALBIN)/golangci-lint run --fix

.PHONY: setup-lint
setup-lint: ## Setup the lint
	test -s $(LOCALBIN)/golangci-lint || { curl -sSfL $(GOLANGCI_INSTALL_SCRIPT) | bash -s -- -b $(LOCALBIN) $(subst v,,$(GOLANGCI_VERSION)); }

.PHONY: license
license: ## Add license header to files
	test -s $(LOCALBIN)/addlicense || GOBIN=$(LOCALBIN) go install github.com/google/addlicense@latest
	$(LOCALBIN)/addlicense -f LICENSE -l mit pkg/ cmd/ internal/

.PHONY: lint
lint: setup-lint ## Run the lint check
	$(LOCALBIN)/golangci-lint run

.PHONY: clean
clean: ## Cleanup build artifacts and tool binaries.
	rm -rf $(BUILD_DIR) dist $(LOCALBIN)

##@ Build

.PHONY: install
install: ## Install falcon
	go install $(GO_BUILD_ARGS) ./cmd/falcon

.PHONY: build
build: ## Build falcon.
	@echo $(VERSION)
	@mkdir -p $(BUILD_DIR)
	go build $(GO_BUILD_ARGS) -o $(BUILD_DIR) ./cmd/falcon

##@ Test

.PHONY: test
test: test-static #test-e2e ## Run all tests

.PHONY: test-static
test-static: test-sanity test-unit ## Run all golang tests

.PHONY: test-sanity
test-sanity: fix ## Test repo formatting, linting, etc.
	git diff --exit-code # fast-fail if fix produced changes
	go vet ./...
	make setup-lint
	make lint
	git diff --exit-code # diff again to ensure other checks don't change repo

.PHONY: test-unit
TEST_PKGS = $(shell go list ./... | grep -v -E 'github.com/crowdstrike/falcon-cli/test/')
test-unit: ## Run unit tests
	go test -coverprofile=coverage.out -covermode=count -short $(TEST_PKGS)

.DEFAULT_GOAL := help
.PHONY: help
help: ## Show this help screen.
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
