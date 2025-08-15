# ForgeAI - Secure Sandboxed Code Executor

ForgeAI is a CLI tool that executes AI-generated code in a secure sandboxed environment. It supports multiple languages and provides isolation to prevent host compromise.

## Features

- Secure sandboxed execution of code
- Multi-language support (Python, Go, JavaScript)
- Resource limits (CPU, memory, execution time)
- Beautiful CLI interface with colorized output
- JSON output for machine integration
- Extensible architecture for adding new languages

## Installation

To install ForgeAI, you need to have Go 1.21+ installed on your system:

```bash
git clone <repository-url>
cd forgeai
go build -o forgeai cmd/forgeai/main.go
```

## Usage

### Execute code directly

```bash
# Run Python code
forgeai run python "print('Hello, World!')"

# Run JavaScript code
forgeai run javascript "console.log('Hello, World!')"

# Run Go code
forgeai run go "package main; import \"fmt\"; func main() { fmt.Println(\"Hello, World!\") }"
```

### Execute a file

```bash
# Execute a Python file
forgeai exec examples/hello_world.py

# Execute a JavaScript file
forgeai exec examples/hello_world.js

# Execute a Go file
forgeai exec examples/hello_world.go
```

### List supported languages

```bash
# List supported languages
forgeai lang list

# Get JSON output
forgeai lang list --json
```

### Configuration

```bash
# Adjust security limits
forgeai config
```

## Architecture

ForgeAI follows a modular architecture with the following components:

- `cmd/forgeai`: Main entry point for the CLI application
- `pkg/cli`: CLI interface implementation using Cobra
- `pkg/sandbox`: Sandbox execution interfaces
- `pkg/executor`: Implementation of code executors
- `pkg/config`: Configuration management

## Security

ForgeAI implements several security measures to protect the host system:

1. Execution in ephemeral temporary directories
2. Automatic cleanup after execution
3. Resource limits (CPU, memory, execution time)
4. Isolation from the host file system

## Extending Language Support

To add support for a new language:

1. Implement the `Executor` interface in `pkg/executor`
2. Add the language to the supported languages list
3. Register the executor in the CLI

## License

MIT