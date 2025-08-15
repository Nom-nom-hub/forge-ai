# REST API Mode Technical Specification

## Overview

This document specifies the technical implementation of the REST API mode for ForgeAI, enabling remote code execution through HTTP endpoints.

## Architecture

### Component Diagram
```
+-------------------+
|   HTTP Client     |
+-------------------+
         |
         v
+-------------------+
|   API Gateway     |
+-------------------+
         |
         v
+-------------------+
|   Auth Manager    |
+-------------------+
         |
         v
+-------------------+
|   Rate Limiter    |
+-------------------+
         |
         v
+-------------------+
|   Job Manager     |
+-------------------+
         |
         v
+-------------------+
|   ForgeAI Core    |
| (CLI/SDK Engine)  |
+-------------------+
         |
         v
+-------------------+
|  Container Engine |
+-------------------+
```

### Key Components

1. **API Gateway**: Handles HTTP requests and responses
2. **Auth Manager**: Manages authentication and authorization
3. **Rate Limiter**: Controls request rate and quotas
4. **Job Manager**: Manages execution jobs and queuing
5. **ForgeAI Core**: Core execution engine (reused from CLI/SDK)
6. **Container Engine**: Containerized execution backend

## API Design

### REST Principles
- Resource-based URLs
- Standard HTTP methods
- Statelessness
- JSON request/response format
- Proper HTTP status codes
- HATEOAS for API navigation

### Versioning
- URL versioning: `/v1/execute`
- Media type versioning: `Accept: application/vnd.forgeai.v1+json`

### Error Handling
- Standard error response format
- Detailed error messages
- Error codes and categories
- Request ID for tracing

## Endpoints

### Core Endpoints

#### Execute Code
```bash
POST /v1/execute
```

**Request:**
```json
{
  "language": "python",
  "code": "print('Hello, World!')",
  "timeout": 30,
  "memory_limit": 128,
  "network_access": false,
  "environment": {
    "ENV_VAR": "value"
  }
}
```

**Response:**
```json
{
  "job_id": "abc123",
  "status": "completed",
  "stdout": "Hello, World!
",
  "stderr": "",
  "exit_code": 0,
  "duration": "150ms",
  "created_at": "2023-01-01T12:00:00Z",
  "completed_at": "2023-01-01T12:00:01Z"
}
```

#### Execute File
```bash
POST /v1/execute/file
```

**Request (multipart/form-data):**
```json
{
  "language": "python",
  "timeout": 30,
  "memory_limit": 128,
  "network_access": false
}
```
```
File: main.py (content)
```

**Response:**
```json
{
  "job_id": "def456",
  "status": "completed",
  "stdout": "Hello, World!
",
  "stderr": "",
  "exit_code": 0,
  "duration": "150ms",
  "created_at": "2023-01-01T12:00:00Z",
  "completed_at": "2023-01-01T12:00:01Z"
}
```

#### Get Job Status
```bash
GET /v1/jobs/{job_id}
```

**Response:**
```json
{
  "job_id": "abc123",
  "status": "running",
  "language": "python",
  "timeout": 30,
  "memory_limit": 128,
  "created_at": "2023-01-01T12:00:00Z"
}
```

#### List Supported Languages
```bash
GET /v1/languages
```

**Response:**
```json
{
  "languages": ["python", "go", "javascript", "rust"],
  "timestamp": "2023-01-01T12:00:00Z"
}
```

#### Get System Status
```bash
GET /v1/status
```

**Response:**
```json
{
  "version": "1.0.0",
  "uptime": "2h30m",
  "jobs_running": 5,
  "jobs_queued": 2,
  "cpu_usage": 45.2,
  "memory_usage": 1024,
  "disk_usage": 5120
}
```

### Administrative Endpoints

#### List All Jobs
```bash
GET /v1/admin/jobs
```

**Query Parameters:**
- `status`: Filter by status (pending, running, completed, failed)
- `language`: Filter by language
- `limit`: Number of results (default: 100)
- `offset`: Pagination offset

#### Cancel Job
```bash
DELETE /v1/admin/jobs/{job_id}
```

#### Get Server Configuration
```bash
GET /v1/admin/config
```

#### Update Server Configuration
```bash
PUT /v1/admin/config
```

## Authentication & Authorization

### Authentication Methods
1. **API Key**: Simple token-based authentication
2. **JWT**: JSON Web Token for stateless auth
3. **OAuth 2.0**: For third-party integrations
4. **Certificate**: Mutual TLS authentication

### API Key Implementation
```bash
# Header
Authorization: Bearer <api_key>

# Query Parameter
?v=1&api_key=<api_key>
```

### User Roles
- **User**: Can execute code and view own jobs
- **Admin**: Can manage all jobs and server configuration
- **System**: Internal services and integrations

### Rate Limiting
- **Per API Key**: 1000 requests/hour
- **Per IP**: 100 requests/hour
- **Global**: 10000 requests/minute

## Job Management

### Job States
1. **Pending**: Job queued for execution
2. **Running**: Job currently executing
3. **Completed**: Job finished successfully
4. **Failed**: Job failed with error
5. **Cancelled**: Job cancelled by user/admin
6. **Timeout**: Job exceeded time limit

### Job Queue
- **Priority Queue**: High-priority jobs executed first
- **FIFO**: First-in, first-out for same priority
- **Resource-based Scheduling**: Schedule based on resource availability
- **Preemption**: Cancel low-priority jobs for high-priority ones

### Job Persistence
- **Database**: Store job metadata and results
- **Storage**: Store code and output for completed jobs
- **Cache**: In-memory cache for recent jobs
- **Backup**: Regular backup of job data

## Real-time Features

### WebSocket Support
```bash
# Connect to job stream
GET /v1/jobs/{job_id}/stream
Upgrade: websocket
```

**Events:**
- `job_status`: Job status updates
- `job_output`: Real-time output streaming
- `job_completed`: Job completion notification

### Server-Sent Events (SSE)
```bash
# Stream job events
GET /v1/jobs/{job_id}/events
Accept: text/event-stream
```

## Configuration

### API Configuration
```yaml
api:
  # Server settings
  host: "0.0.0.0"
  port: 8080
  tls:
    enabled: false
    cert_file: ""
    key_file: ""
  
  # Authentication
  auth:
    enabled: true
    method: "api_key"
    jwt_secret: "secret"
  
  # Rate limiting
  rate_limit:
    enabled: true
    requests_per_hour: 1000
    requests_per_minute: 100
  
  # Job settings
  job:
    max_timeout: 300
    max_memory: 1024
    max_concurrent: 100
    retention_days: 7
  
  # Security
  security:
    cors_origins: ["*"]
    request_timeout: 30
    max_request_size: "10MB"
```

### Environment Variables
```bash
FORGEAI_API_HOST=0.0.0.0
FORGEAI_API_PORT=8080
FORGEAI_API_TLS_ENABLED=false
FORGEAI_API_AUTH_ENABLED=true
FORGEAI_API_RATE_LIMIT_ENABLED=true
```

## Error Handling

### Standard Error Format
```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "The request is invalid",
    "details": [
      {
        "field": "language",
        "message": "Language is required"
      }
    ],
    "request_id": "abc123",
    "timestamp": "2023-01-01T12:00:00Z"
  }
}
```

### Common Error Codes
- `INVALID_REQUEST`: Malformed request
- `UNAUTHORIZED`: Authentication failed
- `FORBIDDEN`: Insufficient permissions
- `NOT_FOUND`: Resource not found
- `TIMEOUT`: Request timeout
- `INTERNAL_ERROR`: Server error
- `RATE_LIMITED`: Rate limit exceeded

## Monitoring & Metrics

### Metrics Collection
- **Request Metrics**: Latency, throughput, error rates
- **Job Metrics**: Execution time, success rate, resource usage
- **System Metrics**: CPU, memory, disk, network usage
- **Business Metrics**: User activity, language usage

### Health Checks
```bash
# Liveness probe
GET /healthz

# Readiness probe
GET /readyz

# Detailed health
GET /healthz/detail
```

### Logging
- **Access Logs**: HTTP request/response logging
- **Application Logs**: Job execution logging
- **Audit Logs**: Security-relevant events
- **Debug Logs**: Detailed debugging information

## Performance Optimization

### Caching
- **Job Results**: Cache recent job results
- **Language List**: Cache supported languages
- **Configuration**: Cache server configuration
- **User Data**: Cache user permissions

### Connection Management
- **Connection Pooling**: Reuse database connections
- **HTTP Keep-Alive**: Reuse HTTP connections
- **WebSocket Reuse**: Reuse WebSocket connections
- **Resource Pooling**: Pool expensive resources

### Asynchronous Processing
- **Background Jobs**: Process jobs asynchronously
- **Event-driven Architecture**: Use events for loose coupling
- **Message Queues**: Use queues for job distribution
- **Microservices**: Decompose into smaller services

## Security

### Transport Security
- **HTTPS**: Enforce HTTPS for all requests
- **HSTS**: HTTP Strict Transport Security
- **CORS**: Cross-Origin Resource Sharing controls
- **Content Security**: Content Security Policy headers

### Data Security
- **Encryption**: Encrypt sensitive data at rest
- **Masking**: Mask sensitive data in logs
- **Validation**: Validate all input data
- **Sanitization**: Sanitize output data

### API Security
- **Rate Limiting**: Prevent abuse
- **Input Validation**: Prevent injection attacks
- **Output Encoding**: Prevent XSS attacks
- **Security Headers**: Set appropriate security headers

## Testing

### Unit Tests
- API endpoint handlers
- Authentication middleware
- Job management logic
- Error handling

### Integration Tests
- End-to-end API flows
- Authentication integration
- Job execution integration
- Database integration

### Load Testing
- Concurrent user simulation
- Stress testing scenarios
- Performance benchmarking
- Scalability testing

### Security Testing
- Penetration testing
- Vulnerability scanning
- Dependency security checks
- Security audit trails

## Deployment

### Container Deployment
```dockerfile
FROM golang:1.19-alpine AS builder
# Build steps

FROM alpine:latest
# Runtime setup
COPY --from=builder /app/forgeai /forgeai
ENTRYPOINT ["/forgeai", "api"]
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: forgeai-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: forgeai-api
  template:
    metadata:
      labels:
        app: forgeai-api
    spec:
      containers:
      - name: forgeai-api
        image: forgeai/api:1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: FORGEAI_API_PORT
          value: "8080"
```

### Load Balancing
- **Reverse Proxy**: Nginx, HAProxy, or cloud load balancer
- **SSL Termination**: Handle TLS at the load balancer
- **Health Checks**: Monitor service health
- **Sticky Sessions**: For WebSocket connections

This technical specification provides a comprehensive blueprint for implementing a robust, secure, and scalable REST API mode for ForgeAI, enabling remote code execution with enterprise-grade features.