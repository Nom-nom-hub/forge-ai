# ForgeAI - Secure Sandboxed Code Executor

[![Go Report Card](https://goreportcard.com/badge/github.com/forgeai/forgeai)](https://goreportcard.com/report/github.com/forgeai/forgeai)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

ForgeAI is a CLI tool that executes AI-generated code in a secure sandboxed environment. It supports multiple languages and provides isolation to prevent host compromise.

## Features

- 🔒 **Secure Sandboxed Execution**: Run untrusted code in isolation
- 🚀 **Multi-language Support**: Python, Go, and JavaScript out of the box
- ⚙️ **Configurable Resource Limits**: CPU, memory, and execution time limits
- 🎨 **Beautiful CLI Interface**: Colorized output and clear prompts
- 📦 **Go SDK**: Integrate ForgeAI into your own applications
- 📡 **JSON Output**: Machine-readable output for integration
- 🧩 **Extensible Architecture**: Easy to add new language executors

## Installation

### Prerequisites
- Go 1.21+

### Building from Source
```bash
git clone https://github.com/forgeai/forgeai.git
cd forgeai
go build -o forgeai cmd/forgeai/main.go
```

### Using Go Install
```bash
go install github.com/forgeai/forgeai/cmd/forgeai@latest
```

## Usage

### Execute Code Directly
```bash
# Run Python code
forgeai run python "print('Hello, World!')"

# Run JavaScript code
forgeai run javascript "console.log('Hello, World!')"

# Run Go code
forgeai run go "package main; import \"fmt\"; func main() { fmt.Println(\"Hello, World!\") }"
```

### Execute a File
```bash
# Execute a Python file
forgeai exec examples/hello_world.py

# Execute a JavaScript file
forgeai exec examples/hello_world.js

# Execute a Go file
forgeai exec examples/hello_world.go
```

### List Supported Languages
```bash
# List supported languages
forgeai lang list

# Get JSON output
forgeai lang list --json
```

### Configuration
```bash
# Run with custom timeout
forgeai --timeout=10s run python "print('Hello, World!')"

# Run with custom memory limit
forgeai --memory-limit=64 run python "print('Hello, World!')"

# Get JSON output
forgeai --json run python "print('Hello, World!')"
```

## Go SDK

ForgeAI provides a Go SDK for integrating code execution capabilities into your applications:

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "forgeai/pkg/executor"
)

func main() {
    // Create executor
    exec := executor.NewLocalExecutor()
    
    // Execute code
    result, err := exec.Execute(context.Background(), "python", "print('Hello, World!')")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Output: %s\n", result.Stdout)
}
```

## Security

ForgeAI implements several security measures to protect the host system:

1. **Isolation**: Code execution occurs in ephemeral temporary directories
2. **Resource Limits**: CPU, memory, and execution time limits prevent resource exhaustion
3. **File System Access**: Limited file system access to prevent data leakage
4. **Automatic Cleanup**: Temporary directories are automatically deleted after execution

## Architecture

```
+-------------------+
|     CLI/API       |
+-------------------+
|   Sandbox Layer   |
+-------------------+
|  Executor Plugins |
+-------------------+
|     Security      |
+-------------------+
```

## Documentation

- [API Reference](docs/API.md)
- [Configuration](docs/CONFIG.md)
- [Integration Guide](docs/INTEGRATION.md)

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) for the CLI framework
```