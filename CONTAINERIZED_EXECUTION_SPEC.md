# Containerized Execution Technical Specification

## Overview

This document specifies the technical implementation of containerized execution for ForgeAI, providing stronger isolation and resource controls than the current local execution model.

## Architecture

### Component Diagram
```
+-------------------+
|     CLI/API       |
+-------------------+
         |
         v
+-------------------+
|  Container Engine |
| (Docker/gVisor/   |
|  Firecracker)     |
+-------------------+
         |
         v
+-------------------+
|  Container Images |
| (Python, Go, JS,  |
|  etc.)            |
+-------------------+
```

### Key Components

1. **Container Executor**: Implements the sandbox.Executor interface using containerization
2. **Image Manager**: Handles container image lifecycle (pull, cache, update)
3. **Resource Controller**: Enforces CPU, memory, and I/O limits
4. **Network Manager**: Controls network access and isolation
5. **Volume Manager**: Manages ephemeral storage and file system access

## Container Executor Implementation

### Interface Implementation
```go
type ContainerExecutor struct {
    engine     ContainerEngine
    config     *config.Config
    imageMgr   *ImageManager
    resourceCtl *ResourceController
}

func (c *ContainerExecutor) Execute(ctx context.Context, language, code string) (*sandbox.ExecutionResult, error) {
    // Implementation
}

func (c *ContainerExecutor) ExecuteFile(ctx context.Context, filePath string) (*sandbox.ExecutionResult, error) {
    // Implementation
}

func (c *ContainerExecutor) SupportedLanguages() []string {
    // Implementation
}
```

### Execution Flow
1. Receive execution request
2. Select appropriate container image
3. Create ephemeral container with resource limits
4. Copy code into container
5. Execute code with timeout
6. Capture output and exit code
7. Clean up container and resources

## Container Images

### Base Images
- **Python**: python:3.9-alpine with security restrictions
- **Go**: golang:1.19-alpine with module restrictions
- **JavaScript**: node:16-alpine with fs restrictions
- **Rust**: rust:1.64-alpine with cargo restrictions

### Security Hardening
- Non-root user execution
- Read-only root filesystem
- Restricted capabilities (drop all, add only necessary)
- Seccomp profiles for system call filtering
- AppArmor/SELinux profiles where applicable

### Example Dockerfile (Python)
```dockerfile
FROM python:3.9-alpine

# Create non-root user
RUN adduser -D forgeai

# Install security tools
RUN apk add --no-cache su-exec

# Set working directory
WORKDIR /workspace

# Copy security profiles
COPY seccomp.json /etc/seccomp.json
COPY apparmor-profile /etc/apparmor.d/forgeai

# Switch to non-root user
USER forgeai

# Default command
CMD ["python", "main.py"]
```

## Resource Management

### CPU Limits
- CPU shares for relative allocation
- CPU quota for hard limits
- CPU period for time-based allocation

### Memory Limits
- Memory limit in bytes
- Memory swap limit
- Kernel memory limit
- OOM killer configuration

### I/O Limits
- Block I/O weight
- Block I/O bandwidth limits
- Block I/O IOPS limits

### Network Limits
- Bandwidth limits
- Connection limits
- Protocol restrictions

## Security Controls

### File System Isolation
- Ephemeral volumes for code execution
- Read-only root filesystem
- Temporary directory for outputs
- No host directory mounting by default

### Network Isolation
- User-defined networks by default
- No external network access by default
- Port exposure controls
- DNS restrictions

### Process Isolation
- PID namespace isolation
- User namespace mapping
- Capability dropping
- Seccomp system call filtering

### Image Security
- Image signature verification
- Vulnerability scanning
- Base image update policies
- Content trust enforcement

## Configuration

### Container Config Structure
```go
type ContainerConfig struct {
    // Engine selection
    Engine string // docker, gvisor, firecracker
    
    // Resource limits
    CPUShares    int
    MemoryLimit  int64 // bytes
    DiskQuota    int64 // bytes
    
    // Network settings
    NetworkMode  string // bridge, none, host
    DNS          []string
    
    // Security settings
    ReadOnlyRoot bool
    Capabilities []string
    SeccompProfile string
    
    // Image settings
    ImageRegistry string
    ImagePullPolicy string // always, if-not-present, never
}
```

### Default Configuration
```yaml
container:
  engine: docker
  cpu_shares: 100
  memory_limit: 134217728 # 128MB
  disk_quota: 1073741824 # 1GB
  network_mode: none
  read_only_root: true
  capabilities: []
  seccomp_profile: "forgeai-seccomp.json"
  image_registry: "forgeai"
  image_pull_policy: "if-not-present"
```

## Error Handling

### Container Startup Failures
- Image not found
- Resource limits exceeded
- Security policy violations
- Engine not available

### Execution Failures
- Timeout exceeded
- Memory limit exceeded
- Disk quota exceeded
- Process killed by OOM

### Cleanup Failures
- Container deletion failures
- Volume cleanup failures
- Network cleanup failures

## Monitoring and Metrics

### Key Metrics
- Container startup time
- Execution duration
- Resource usage (CPU, memory, disk, network)
- Error rates
- Success rates

### Logging
- Container lifecycle events
- Resource limit violations
- Security policy violations
- Error details

## Performance Optimization

### Startup Optimization
- Image layer caching
- Container snapshotting
- Pre-warmed containers
- Parallel container creation

### Runtime Optimization
- Resource pooling
- Connection pooling
- Memory pre-allocation
- CPU affinity

### Cleanup Optimization
- Asynchronous cleanup
- Batch cleanup operations
- Resource reclaiming
- Garbage collection

## Testing

### Unit Tests
- Container executor functionality
- Resource limit enforcement
- Security policy enforcement
- Error handling

### Integration Tests
- End-to-end execution flow
- Resource limit testing
- Security testing
- Performance testing

### Security Tests
- Host escape attempts
- Resource exhaustion attacks
- Network access attempts
- File system access attempts

## Deployment

### Requirements
- Docker Engine 20.10+ or equivalent
- Linux kernel 4.0+ (for advanced features)
- Sufficient storage for container images
- Appropriate permissions for container operations

### Installation
- Container engine installation
- Image pre-pulling
- Security profile deployment
- Configuration setup

### Maintenance
- Image updates
- Security profile updates
- Log rotation
- Resource monitoring

This technical specification provides a comprehensive blueprint for implementing containerized execution in ForgeAI, ensuring strong isolation, resource control, and security while maintaining performance and usability.