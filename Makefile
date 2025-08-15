# Makefile for ForgeAI

# Variables
BINARY=forgeai
MAIN_FILE=cmd/forgeai/main.go

# Default target
all: build

# Build the binary
build:
	go build -o ${BINARY} ${MAIN_FILE}

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
	rm -f ${BINARY}

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

# Help
help:
	@echo "ForgeAI Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make              Build the binary"
	@echo "  make build        Build the binary"
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
	@echo "  make release      Release build"
	@echo "  make help         Show this help"

.PHONY: all build deps test test-coverage test-verbose test-integration clean install fmt vet lint docs release help