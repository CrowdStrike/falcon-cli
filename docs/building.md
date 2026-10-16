# Building and Testing

## Prerequisites

- Go 1.26+
- Make

## Build Targets

```bash
make build       # Build the falcon binary
make test        # Run unit tests with race detector and coverage
make test-e2e    # Run end-to-end tests (requires live Falcon tenant)
make lint        # Run golangci-lint
make fmt         # Format code
make vet         # Run go vet
make snapshot    # Build cross-platform binaries via goreleaser
make clean       # Remove build artifacts
make help        # Show all available targets
```

## Testing

[![CI](https://github.com/CrowdStrike/falcon-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/CrowdStrike/falcon-cli/actions/workflows/ci.yml)

Unit tests run without external dependencies:

```bash
make test
```

### End-to-end tests

E2E tests run against a live Falcon tenant:

```bash
export FALCON_CLIENT_ID=<id>
export FALCON_CLIENT_SECRET=<secret>
make test-e2e
```

Filter by label:

```bash
GINKGO_LABEL_FILTER="risks" make test-e2e
```

## Project Structure

```
cmd/falcon/          Entry point (main.go)
pkg/
  cli/               CLI bootstrap, config init, auth check
  cmd/               Command implementations (cobra)
    auth/            Authentication commands
    fcs/             Falcon Cloud Security command tree
    root/            Root command registration
    version/         Version command
  cmdutil/           Shared command utilities
  config/            Configuration management (viper)
  factory/           Dependency injection
  iostreams/         IO streams abstraction
  output/            Output formatting
  utils/             Shared helpers
  version/           Version metadata
test/e2e/            End-to-end test suite (Ginkgo)
```
