# Makefile for ForgeAI

# Variables
BINARY=forgeai
API_BINARY=forgeai-api
MAIN_FILE=cmd/forgeai/main.go
API_MAIN_FILE=cmd/api/main.go

# Default target
all: build

# Build the binary
build:
	go build -o ${BINARY} ${MAIN_FILE}

# Build the API server
build-api:
	go build -o ${API_BINARY} ${API_MAIN_FILE}

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
	rm -f ${BINARY} ${API_BINARY}

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

# Help
help:
	@echo "ForgeAI Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make              Build the CLI binary"
	@echo "  make build        Build the CLI binary"
	@echo "  make build-api    Build the API server"
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
	@echo "  make help         Show this help"

.PHONY: all build build-api deps test test-coverage test-verbose test-integration clean install fmt vet lint docs release release-api help