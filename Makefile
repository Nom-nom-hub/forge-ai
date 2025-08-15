# Makefile for ForgeAI

# Variables
BINARY=forgeai
API_BINARY=forgeai-api
PLUGIN_BINARY=forgeai-plugin
SECURITY_BINARY=forgeai-security
PERF_BINARY=forgeai-perf
MAIN_FILE=cmd/forgeai/main.go
API_MAIN_FILE=cmd/api/main.go
PLUGIN_MAIN_FILE=cmd/plugin/main.go
SECURITY_MAIN_FILE=cmd/security/main.go
PERF_MAIN_FILE=cmd/performance/main.go

# Default target
all: build build-api build-plugin build-security build-perf

# Build the CLI binary
build:
	go build -o ${BINARY} ${MAIN_FILE}

# Build the API server
build-api:
	go build -o ${API_BINARY} ${API_MAIN_FILE}

# Build the plugin manager
build-plugin:
	go build -o ${PLUGIN_BINARY} ${PLUGIN_MAIN_FILE}

# Build the security testing tool
build-security:
	go build -o ${SECURITY_BINARY} ${SECURITY_MAIN_FILE}

# Build the performance testing tool
build-perf:
	go build -o ${PERF_BINARY} ${PERF_MAIN_FILE}

# Install dependencies
deps:
	go mod tidy

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run integration tests
test-integration:
	go test ./test/integration

# Clean build artifacts
clean:
	rm -f ${BINARY} ${API_BINARY} ${PLUGIN_BINARY} ${SECURITY_BINARY} ${PERF_BINARY}

# Install the binary
install:
	go install ${MAIN_FILE}

# Format the code
fmt:
	go fmt ./...

# Vet the code
vet:
	go vet ./...

# Lint the code (requires golangci-lint)
lint:
	golangci-lint run

# Generate documentation
docs:
	go doc -all ./...

# Release build
release:
	go build -ldflags "-s -w" -o ${BINARY} ${MAIN_FILE}

# Release build for API
release-api:
	go build -ldflags "-s -w" -o ${API_BINARY} ${API_MAIN_FILE}

# Release build for plugin manager
release-plugin:
	go build -ldflags "-s -w" -o ${PLUGIN_BINARY} ${PLUGIN_MAIN_FILE}

# Help
help:
	@echo "ForgeAI Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make              Build all binaries"
	@echo "  make build        Build the CLI binary"
	@echo "  make build-api    Build the API server"
	@echo "  make build-plugin Build the plugin manager"
	@echo "  make build-security Build the security testing tool"
	@echo "  make build-perf   Build the performance testing tool"
	@echo "  make deps         Install dependencies"
	@echo "  make test         Run tests"
	@echo "  make test-coverage Run tests with coverage"
	@echo "  make test-verbose Run tests with verbose output"
	@echo "  make test-integration Run integration tests"
	@echo "  make clean        Clean build artifacts"
	@echo "  make install      Install the binary"
	@echo "  make fmt          Format the code"
	@echo "  make vet          Vet the code"
	@echo "  make lint         Lint the code"
	@echo "  make docs         Generate documentation"
	@echo "  make release      Release build for CLI"
	@echo "  make release-api  Release build for API"
	@echo "  make release-plugin Release build for plugin manager"
	@echo "  make help         Show this help"

.PHONY: all build build-api build-plugin build-security build-perf deps test test-coverage test-verbose test-integration clean install fmt vet lint docs release release-api release-plugin help