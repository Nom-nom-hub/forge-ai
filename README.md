# ForgeAI - Secure Code Execution Platform

## What is ForgeAI?

ForgeAI is a secure, sandboxed code execution platform that allows you to run untrusted code safely. It provides multiple isolation layers, resource limits, and extensibility through a plugin system.

## Key Features

- **Multi-layered Security**: Process, container, and plugin isolation
- **Resource Controls**: CPU, memory, and time limits
- **Language Support**: Python, Go, JavaScript, and extensible via plugins
- **Multiple Interfaces**: CLI, REST API, and Go SDK
- **Cross-platform**: Works on Windows, Linux, and macOS

## Installation

### Prerequisites
- Go 1.19 or later
- Docker (for containerized execution)
- Git

### Quick Install
```bash
# Clone the repository
git clone https://github.com/forgeai/forgeai.git
cd forgeai

# Build binaries
make all

# Install (optional)
make install
```

### Manual Installation
```bash
# Download the latest release
wget https://github.com/forgeai/forgeai/releases/latest/download/forgeai-linux-amd64.tar.gz
tar -xzf forgeai-linux-amd64.tar.gz
sudo mv forgeai* /usr/local/bin/
```

## Getting Started

### CLI Usage
```bash
# Run Python code
forgeai run python "print('Hello, World!')"

# Execute a file
forgeai exec examples/hello_world.py

# List supported languages
forgeai lang list

# Use containerized execution
forgeai --container run python "print('Hello, World!')"

# Use plugins
forgeai --plugin-dir=./plugins run rust "fn main() { println!(\"Hello, World!\"); }"
```

### REST API
```bash
# Start the API server
forgeai-api

# Execute code via API
curl -X POST http://localhost:8080/v1/execute \
  -H "Content-Type: application/json" \
  -d '{
    "language": "python",
    "code": "print(\"Hello, World!\")"
  }'
```

### Go SDK
```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "forgeai/pkg/executor"
)

func main() {
    exec := executor.NewLocalExecutor()
    result, err := exec.Execute(context.Background(), "python", "print('Hello, World!')")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Output: %s\n", result.Stdout)
}
```

## Security Model

ForgeAI implements multiple layers of security:

### 1. Process Isolation
- Code executes in temporary directories
- Automatic cleanup after execution
- Resource limits (timeout, memory)

### 2. Container Isolation
- Docker-based sandboxing
- Network and file system restrictions
- CPU and memory quotas

### 3. Plugin Isolation
- Plugins run as separate processes
- JSON-based communication
- Language-specific security controls

## Resource Limits

### Default Limits
- **Timeout**: 30 seconds
- **Memory**: 128 MB
- **CPU**: 10% of total CPU (Linux only)

### Custom Limits
```bash
# Set custom timeout
forgeai --timeout=60s run python "import time; time.sleep(30)"

# Set custom memory limit
forgeai --memory-limit=256 run python "print('Hello, World!')"

# Combine limits
forgeai --timeout=10s --memory-limit=64 run python "print('Hello, World!')"
```

## Plugin System

### Installing Plugins
```bash
# List available plugins
forgeai-plugin list

# Install a plugin
forgeai-plugin install rust-plugin

# Remove a plugin
forgeai-plugin remove rust-plugin

# Update a plugin
forgeai-plugin update rust-plugin
```

### Creating Plugins
Plugins are external executables that communicate via JSON:

```go
// plugin.go
package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type ExecutionResult struct {
    Stdout   string `json:"stdout"`
    Stderr   string `json:"stderr"`
    ExitCode int    `json:"exit_code"`
    Duration int64  `json:"duration"`
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: plugin <command> [args...]")
        os.Exit(1)
    }
    
    command := os.Args[1]
    
    switch command {
    case "execute":
        // Handle code execution
        result := ExecutionResult{
            Stdout:   "Hello from plugin!",
            ExitCode: 0,
            Duration: 100,
        }
        output, _ := json.Marshal(result)
        fmt.Print(string(output))
    }
}
```

Plugin manifest (`manifest.json`):
```json
{
  "name": "rust-plugin",
  "languages": ["rust"]
}
```

## Configuration

### Environment Variables
```bash
FORGEAI_TIMEOUT=60s          # Execution timeout
FORGEAI_MEMORY_LIMIT=256     # Memory limit in MB
FORGEAI_PLUGIN_DIR=./plugins # Plugin directory
FORGEAI_CONTAINER=true       # Enable containerized execution
```

### Configuration File
Create `.forgeai.yaml` in your home directory:
```yaml
timeout: 30s
memory_limit: 128
plugin_dir: ./plugins
container: false
```

## API Endpoints

### Core Endpoints
- `GET /` - API information
- `GET /healthz` - Health check
- `GET /readyz` - Readiness check
- `GET /v1/languages` - Supported languages
- `POST /v1/execute` - Code execution
- `POST /v1/execute/file` - File execution
- `GET /v1/jobs/{id}` - Job status
- `DELETE /v1/jobs/{id}` - Cancel job

### Example Request
```json
POST /v1/execute
{
  "language": "python",
  "code": "print('Hello, World!')",
  "timeout": 30,
  "memory_limit": 128
}
```

### Example Response
```json
{
  "job_id": "job-12345",
  "status": "completed",
  "stdout": "Hello, World!\n",
  "stderr": "",
  "exit_code": 0,
  "duration": "150ms"
}
```

## Performance

### Execution Models
1. **Local**: Fastest, lowest isolation
2. **Container**: Strong isolation, moderate overhead
3. **Plugin**: Flexible, process overhead

### Resource Usage
Typical resource usage per execution:
- **Memory**: 10-50 MB
- **CPU**: Minimal when idle
- **Disk**: Temporary files cleaned up automatically

## Troubleshooting

### Common Issues

#### Docker Not Found
```bash
Error: docker is not available
```
**Solution**: Install Docker or use local execution mode.

#### Plugin Not Found
```bash
Error: unsupported language: rust
```
**Solution**: Install the required plugin or check plugin directory.

#### Timeout Errors
```bash
Error: Execution timed out
```
**Solution**: Increase timeout with `--timeout` flag.

### Debugging
Enable debug output:
```bash
# CLI debugging
forgeai --debug run python "print('Hello, World!')"

# API debugging
FORGEAI_DEBUG=true forgeai-api
```

## Contributing

### Development Setup
```bash
# Clone repository
git clone https://github.com/forgeai/forgeai.git
cd forgeai

# Install dependencies
make deps

# Run tests
make test

# Build binaries
make build
```

### Code Structure
```
forgeai/
├── cmd/          # Command-line applications
├── pkg/          # Go packages
├── plugins/      # Plugin directory
├── docs/         # Documentation
├── examples/      # Example files
└── test/         # Tests
```

### Pull Request Process
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Update documentation
6. Submit pull request

## License

ForgeAI is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## Support

### Community
- GitHub Issues: https://github.com/forgeai/forgeai/issues
- Discussions: https://github.com/forgeai/forgeai/discussions

### Commercial Support
For enterprise support options, contact support@forgeai.dev

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history.

## Roadmap

### Upcoming Features
- Java, C#, Ruby plugin support
- Advanced security profiles
- Kubernetes integration
- Machine learning optimization

### Long-term Goals
- Cloud-native architecture
- AI-powered code analysis
- Enterprise security features
- Global scale deployment