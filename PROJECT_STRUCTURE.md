# ForgeAI Project Structure

## Root Directory
```
forgeai/
├── .gitignore                 # Git ignore patterns
├── CHANGELOG.md               # Version history
├── CONTRIBUTING.md            # Contribution guidelines
├── go.mod                     # Go module definition
├── go.sum                     # Go dependency checksums
├── LICENSE                    # MIT license
├── Makefile                   # Build automation
├── PRD.md                     # Product requirements document
├── README.md                  # Project overview
├── forgeai.exe                # CLI application
├── forgeai-api.exe            # REST API server
├── forgeai-plugin.exe         # Plugin manager
├── forgeai-security.exe       # Security testing tool
├── forgeai-perf.exe           # Performance testing tool
└── forgeai-performance.exe    # Performance testing tool
```

## Source Code Structure
```
forgeai/
├── cmd/                       # Command-line applications
│   ├── forgeai/               # CLI application
│   │   └── main.go            # CLI entry point
│   ├── api/                   # REST API server
│   │   └── main.go            # API server entry point
│   ├── plugin/                # Plugin manager
│   │   └── main.go            # Plugin manager entry point
│   ├── security/              # Security testing tool
│   │   └── main.go            # Security testing entry point
│   └── performance/           # Performance testing tool
│       └── main.go            # Performance testing entry point
├── pkg/                       # Go packages
│   ├── api/                   # REST API implementation
│   │   ├── server.go          # API server
│   │   └── jobs.go            # Job management
│   ├── cli/                   # CLI interface
│   │   └── cli.go             # Command-line parsing
│   ├── config/                # Configuration management
│   │   └── config.go          # Configuration handling
│   ├── container/              # Containerized execution
│   │   ├── container.go        # Container executor interface
│   │   └── docker.go          # Docker executor implementation
│   ├── executor/              # Local execution
│   │   └── executor.go        # Local executor implementation
│   ├── output/                # Output formatting
│   │   └── output.go          # Output handling
│   ├── plugin/                # Plugin system
│   │   ├── manager.go         # Plugin manager
│   │   └── plugin.go          # Plugin interface
│   ├── registry/              # Plugin registry
│   │   └── client.go          # Registry client
│   ├── sandbox/              # Sandbox interface
│   │   └── sandbox.go         # Sandbox interface definition
│   ├── security/              # Security testing
│   │   ├── containerized.go   # Containerized executor
│   │   ├── executor.go        # Secure executor
│   │   └── testing.go         # Security testing framework
│   └── performance/           # Performance testing
│       └── testing.go          # Performance testing framework
├── docs/                      # Documentation
│   ├── API_DOCS.md            # REST API documentation
│   ├── CONFIG.md              # Configuration guide
│   ├── INTEGRATION.md         # Integration guide
│   └── API.md                 # Go SDK API reference
├── examples/                 # Example files
│   ├── hello_world.py         # Python example
│   ├── hello_world.go         # Go example
│   └── hello_world.js        # JavaScript example
├── plugins/                   # Plugin directory
│   ├── hello-plugin/         # Hello plugin
│   │   ├── manifest.json      # Plugin manifest
│   │   └── hello-plugin.exe  # Plugin executable
│   └── rust-plugin/          # Rust plugin
│       ├── manifest.json      # Plugin manifest
│       └── rust-plugin.exe    # Plugin executable
└── test/                      # Tests
    └── executor_test.go       # Unit tests
```

## Documentation
```
forgeai/
├── README.md                  # Project overview
├── CHANGELOG.md               # Version history
├── CONTRIBUTING.md            # Contribution guidelines
├── LICENSE                    # MIT license
├── docs/                      # Comprehensive documentation
│   ├── API_DOCS.md            # REST API documentation
│   ├── CONFIG.md              # Configuration guide
│   ├── INTEGRATION.md         # Integration guide
│   └── API.md                 # Go SDK API reference


## Testing
```
forgeai/
├── test/                      # Unit and integration tests
│   └── executor_test.go       # Executor unit tests
├── cmd/security/main.go       # Security testing tool
├── cmd/performance/main.go    # Performance testing tool
├── forgeai-security.exe       # Security testing executable
├── forgeai-perf.exe           # Performance testing executable
└── forgeai-performance.exe    # Performance testing executable
```

## Key Features

### 1. CLI Tool
- `forgeai run <language> <code>` - Execute code directly
- `forgeai exec <file>` - Execute file
- `forgeai lang list` - List supported languages
- `forgeai --container` - Containerized execution
- `forgeai --plugin-dir` - Plugin execution

### 2. REST API
- `GET /` - API information
- `GET /healthz` - Health check
- `GET /readyz` - Readiness check
- `GET /v1/languages` - Supported languages
- `POST /v1/execute` - Code execution
- `POST /v1/execute/file` - File execution
- `GET /v1/jobs/{id}` - Job status
- `DELETE /v1/jobs/{id}` - Cancel job
- `GET /v1/jobs` - List jobs
- `GET /v1/status` - Server status

### 3. Go SDK
- `executor.NewLocalExecutor()` - Local execution
- `container.NewDockerExecutor()` - Containerized execution
- `plugin.NewManager()` - Plugin management
- `plugin.NewExternalExecutor()` - Plugin execution

### 4. Plugin System
- Cross-platform plugin support (Windows, Linux, macOS)
- External executable plugins
- Plugin manifest system
- Automatic plugin discovery
- Plugin registry integration

### 5. Security Features
- Process isolation
- Container isolation
- Plugin isolation
- Resource limits (CPU, memory, timeout)
- Network access controls
- File system restrictions

### 6. Performance Features
- Asynchronous job execution
- Job queuing
- Resource pooling
- Concurrent execution
- Performance monitoring

## Supported Languages

### Built-in Support
- Python (3.9)
- Go (1.19)
- JavaScript (Node.js 16)

### Plugin Support
- Hello (example plugin)
- Rust (simulated plugin)
- Extensible for any language

### Container Support
- Python (python:3.9-alpine)
- Go (golang:1.19-alpine)
- JavaScript (node:16-alpine)

## Enterprise Features

### Scalability
- Horizontal scaling
- Load balancing
- Job queuing
- Resource management

### Monitoring
- Health checks
- Performance metrics
- Execution logging
- Error tracking

### Operations
- Graceful shutdown
- Configuration management
- Deployment flexibility
- Upgrade support

This structure provides a comprehensive, well-organized foundation for a production-ready code execution platform with enterprise-grade security, extensibility, and performance features.