# Integration Guide

This document explains how to integrate ForgeAI into other CLI tools.

## Go SDK

ForgeAI provides a Go SDK that can be used to integrate code execution capabilities into other applications.

### Installation

```bash
go get forgeai
```

### Basic Usage

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

### With Custom Configuration

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "forgeai/pkg/executor"
)

func main() {
    // Create executor with custom config
    exec := executor.NewLocalExecutor()
    exec.Timeout = 10 * time.Second
    exec.MemoryLimit = 64 // 64 MB
    
    // Execute code
    result, err := exec.Execute(context.Background(), "python", "print('Hello, World!')")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Output: %s\n", result.Stdout)
}
```

## REST API Mode

ForgeAI can be run in REST API mode for remote execution.

### Starting the API Server

```bash
forgeai api
```

### API Endpoints

#### POST /execute
Execute code in a sandbox.

Request:
```json
{
  "language": "python",
  "code": "print('Hello, World!')"
}
```

Response:
```json
{
  "stdout": "Hello, World!\n",
  "stderr": "",
  "exit_code": 0,
  "duration": "100ms"
}
```

#### GET /languages
Get supported languages.

Response:
```json
["python", "go", "javascript"]
```

## Plugin System

ForgeAI supports a plugin system for adding new language executors.

### Creating a Plugin

Plugins are Go plugins that implement the Executor interface.

```go
package main

import (
    "context"
    
    "forgeai/pkg/sandbox"
)

type MyExecutor struct{}

func (e *MyExecutor) Execute(ctx context.Context, language, code string) (*sandbox.ExecutionResult, error) {
    // Implementation
}

func (e *MyExecutor) ExecuteFile(ctx context.Context, filePath string) (*sandbox.ExecutionResult, error) {
    // Implementation
}

func (e *MyExecutor) SupportedLanguages() []string {
    return []string{"mylanguage"}
}
```

### Loading Plugins

Plugins can be loaded at startup using the `--plugin-dir` option.

```bash
forgeai --plugin-dir=/path/to/plugins run mylanguage "code"
```