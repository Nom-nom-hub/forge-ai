# Integration Guide

## Overview

This guide explains how to integrate ForgeAI into your applications using the CLI, REST API, or Go SDK.

## Integration Methods

### 1. Command-Line Interface (CLI)
Best for simple scripts and automation.

### 2. REST API
Best for web applications and services.

### 3. Go SDK
Best for Go applications requiring deep integration.

## CLI Integration

### Basic Usage
```bash
# Execute Python code
forgeai run python "print('Hello, World!')"

# Execute a file
forgeai exec script.py

# List supported languages
forgeai lang list
```

### Advanced CLI Usage
```bash
# Execute with custom timeout
forgeai --timeout=60s run python "import time; time.sleep(30)"

# Execute with custom memory limit
forgeai --memory-limit=256 run python "print('Hello, World!')"

# Execute in container mode
forgeai --container run python "print('Hello, World!')"

# Execute with plugins
forgeai --plugin-dir=./plugins run rust "fn main() { println!(\"Hello, World!\"); }"
```

### Scripting Examples

#### Bash Script
```bash
#!/bin/bash

# Execute code and capture output
result=$(forgeai run python "print('Hello, World!')")
echo "Result: $result"

# Execute file and check exit code
forgeai exec script.py
if [ $? -eq 0 ]; then
    echo "Script executed successfully"
else
    echo "Script failed"
fi
```

#### PowerShell Script
```powershell
# Execute code and capture output
$result = forgeai run python "print('Hello, World!')"
Write-Host "Result: $result"

# Execute file and check exit code
forgeai exec script.py
if ($LASTEXITCODE -eq 0) {
    Write-Host "Script executed successfully"
} else {
    Write-Host "Script failed"
}
```

## REST API Integration

### API Endpoints
```bash
# Base URL
http://localhost:8080/v1

# Execute code
POST /v1/execute

# Execute file
POST /v1/execute/file

# Get job status
GET /v1/jobs/{job_id}

# Cancel job
DELETE /v1/jobs/{job_id}
```

### API Integration Examples

#### Python (requests)
```python
import requests
import time

# Execute code
response = requests.post('http://localhost:8080/v1/execute', json={
    'language': 'python',
    'code': 'print("Hello, World!")'
})

job_id = response.json()['job_id']

# Poll for result
while True:
    response = requests.get(f'http://localhost:8080/v1/jobs/{job_id}')
    job = response.json()
    
    if job['status'] == 'completed':
        print(job['stdout'])
        break
    elif job['status'] == 'failed':
        print(f"Job failed: {job['stderr']}")
        break
    
    time.sleep(1)
```

#### JavaScript (fetch)
```javascript
// Execute code
fetch('http://localhost:8080/v1/execute', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({
        language: 'javascript',
        code: 'console.log("Hello, World!");'
    })
})
.then(response => response.json())
.then(data => {
    const jobId = data.job_id;
    
    // Poll for result
    const poll = () => {
        fetch(`http://localhost:8080/v1/jobs/${jobId}`)
        .then(response => response.json())
        .then(job => {
            if (job.status === 'completed') {
                console.log(job.stdout);
            } else if (job.status === 'failed') {
                console.error(`Job failed: ${job.stderr}`);
            } else {
                setTimeout(poll, 1000);
            }
        });
    };
    
    poll();
});
```

#### cURL
```bash
# Execute code
curl -X POST http://localhost:8080/v1/execute \
  -H "Content-Type: application/json" \
  -d '{
    "language": "python",
    "code": "print(\"Hello, World!\")"
  }'

# Get job status
curl http://localhost:8080/v1/jobs/job-1234567890
```

## Go SDK Integration

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

### Advanced Go SDK Usage
```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "forgeai/pkg/config"
    "forgeai/pkg/executor"
)

func main() {
    // Create executor with custom config
    exec := executor.NewLocalExecutor()
    exec.Timeout = 60 * time.Second
    exec.MemoryLimit = 256 // 256 MB
    
    // Execute code
    ctx := context.Background()
    result, err := exec.Execute(ctx, "python", "print('Hello, World!')")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Output: %s\n", result.Stdout)
    fmt.Printf("Exit Code: %d\n", result.ExitCode)
    fmt.Printf("Duration: %v\n", result.Duration)
}
```

### Container Executor
```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "forgeai/pkg/container"
)

func main() {
    // Create container executor
    exec := container.NewDockerExecutor()
    exec.Timeout = 30 * time.Second
    exec.MemoryLimit = 128 // 128 MB
    
    // Execute code
    ctx := context.Background()
    result, err := exec.Execute(ctx, "python", "print('Hello, World!')")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Output: %s\n", result.Stdout)
}
```

### Plugin Executor
```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "forgeai/pkg/plugin"
)

func main() {
    // Create plugin manager
    manager := plugin.NewManager()
    manager.LoadPluginsFromDir("./plugins")
    
    // Get executor for language
    exec, ok := manager.GetExecutor("rust")
    if !ok {
        log.Fatal("Rust executor not found")
    }
    
    // Execute code
    ctx := context.Background()
    result, err := exec.Execute(ctx, "rust", "fn main() { println!(\"Hello, World!\"); }")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Output: %s\n", result.Stdout)
}
```

## Error Handling

### CLI Error Handling
```bash
# Check exit code
forgeai run python "print('Hello, World!')"
exit_code=$?
if [ $exit_code -ne 0 ]; then
    echo "Command failed with exit code $exit_code"
fi
```

### API Error Handling
```python
import requests

response = requests.post('http://localhost:8080/v1/execute', json={
    'language': 'python',
    'code': 'print("Hello, World!")'
})

if response.status_code != 200:
    print(f"API error: {response.status_code}")
    print(response.text)
else:
    print("Success!")
```

### Go SDK Error Handling
```go
result, err := exec.Execute(context.Background(), "python", "print('Hello, World!')")
if err != nil {
    // Handle execution error
    log.Printf("Execution error: %v", err)
    return
}

if result.ExitCode != 0 {
    // Handle non-zero exit code
    log.Printf("Non-zero exit code: %d", result.ExitCode)
    log.Printf("Stderr: %s", result.Stderr)
}
```

## Best Practices

### Security
1. Always validate and sanitize input code
2. Use appropriate resource limits
3. Disable network access when not needed
4. Use containerized execution for untrusted code
5. Monitor resource usage

### Performance
1. Reuse executors when possible
2. Use appropriate timeout values
3. Batch multiple executions when possible
4. Monitor execution performance
5. Implement caching for repeated executions

### Reliability
1. Implement retry logic for transient failures
2. Handle timeouts gracefully
3. Log execution results for debugging
4. Monitor system resources
5. Implement health checks

## Troubleshooting

### Common Issues

#### Network Connectivity
```
Error: connection refused
```
**Solution**: Ensure ForgeAI API server is running

#### Invalid Language
```
Error: unsupported language
```
**Solution**: Check supported languages or install required plugin

#### Resource Limits Exceeded
```
Error: memory limit exceeded
```
**Solution**: Increase memory limit or optimize code

#### Timeout
```
Error: execution timed out
```
**Solution**: Increase timeout or optimize code

### Debugging Tips

#### Enable Debug Mode
```bash
forgeai --debug run python "print('Hello, World!')"
```

#### Check API Server Logs
```bash
# Start API server with verbose logging
forgeai-api --debug
```

#### Monitor System Resources
```bash
# Check system resources during execution
htop
iotop
```

## Production Considerations

### Deployment
1. Use containerized deployment for consistency
2. Implement load balancing for high availability
3. Use persistent storage for job history
4. Implement monitoring and alerting
5. Configure appropriate resource limits

### Security
1. Implement authentication for API access
2. Use TLS for encrypted communication
3. Implement rate limiting to prevent abuse
4. Regularly update container images
5. Scan dependencies for vulnerabilities

### Monitoring
1. Monitor API server health
2. Track execution performance metrics
3. Monitor system resource usage
4. Log execution results for auditing
5. Implement alerting for failures

## Support

For integration questions and support, please:
1. Check the documentation
2. Review the examples
3. Search existing issues
4. Contact support@forgeai.dev