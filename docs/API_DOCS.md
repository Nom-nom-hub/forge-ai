# API Documentation

## Overview

The ForgeAI REST API provides programmatic access to code execution capabilities. It allows you to execute code in multiple languages with resource limits and security controls.

## Base URL

```
http://localhost:8080/v1
```

## Authentication

Currently, the API does not require authentication. In production, authentication would be implemented using API keys or JWT tokens.

## Rate Limiting

The API implements rate limiting to prevent abuse. Limits are:
- 1000 requests per hour per IP
- 100 requests per minute per IP

## Error Handling

All errors are returned as JSON with the following structure:

```json
{
  "error": "error message"
}
```

## Endpoints

### Get API Information
```
GET /
```

Returns information about the API server.

**Response:**
```json
{
  "message": "ForgeAI API Server",
  "version": "1.0.0",
  "docs": "/v1/docs"
}
```

### Health Check
```
GET /healthz
```

Returns the health status of the server.

**Response:**
```json
{
  "status": "healthy",
  "time": "2023-01-01T00:00:00Z"
}
```

### Readiness Check
```
GET /readyz
```

Returns the readiness status of the server.

**Response:**
```json
{
  "status": "ready",
  "time": "2023-01-01T00:00:00Z"
}
```

### List Supported Languages
```
GET /v1/languages
```

Returns a list of supported programming languages.

**Response:**
```json
{
  "languages": ["python", "go", "javascript"],
  "timestamp": "2023-01-01T00:00:00Z"
}
```

### Execute Code
```
POST /v1/execute
```

Executes code in the specified language.

**Request:**
```json
{
  "language": "python",
  "code": "print('Hello, World!')",
  "timeout": 30,
  "memory_limit": 128,
  "network_access": false
}
```

**Response:**
```json
{
  "job_id": "job-1234567890",
  "status": "pending"
}
```

### Execute File
```
POST /v1/execute/file
```

Executes code from a file.

**Request:**
```json
{
  "file_path": "/path/to/file.py",
  "timeout": 30,
  "memory_limit": 128,
  "network_access": false
}
```

**Response:**
```json
{
  "job_id": "job-1234567890",
  "status": "pending"
}
```

### Get Job Status
```
GET /v1/jobs/{job_id}
```

Returns the status and results of a job.

**Response (running):**
```json
{
  "job_id": "job-1234567890",
  "status": "running",
  "language": "python",
  "timeout": 30,
  "memory_limit": 128,
  "network_access": false,
  "created_at": "2023-01-01T00:00:00Z",
  "started_at": "2023-01-01T00:00:01Z"
}
```

**Response (completed):**
```json
{
  "job_id": "job-1234567890",
  "status": "completed",
  "language": "python",
  "timeout": 30,
  "memory_limit": 128,
  "network_access": false,
  "created_at": "2023-01-01T00:00:00Z",
  "started_at": "2023-01-01T00:00:01Z",
  "completed_at": "2023-01-01T00:00:02Z",
  "stdout": "Hello, World!\n",
  "stderr": "",
  "exit_code": 0,
  "duration": "100ms"
}
```

### Cancel Job
```
DELETE /v1/jobs/{job_id}
```

Cancels a running job.

**Response:**
```json
{
  "job_id": "job-1234567890",
  "status": "cancelled",
  "message": "Job cancelled successfully"
}
```

### List Jobs
```
GET /v1/jobs
```

Lists all jobs with optional filtering.

**Query Parameters:**
- `status`: Filter by status (pending, running, completed, failed, cancelled)
- `language`: Filter by language

**Response:**
```json
{
  "jobs": [
    {
      "job_id": "job-1234567890",
      "status": "completed",
      "language": "python",
      "created_at": "2023-01-01T00:00:00Z",
      "started_at": "2023-01-01T00:00:01Z",
      "completed_at": "2023-01-01T00:00:02Z"
    }
  ],
  "count": 1
}
```

### Get Server Status
```
GET /v1/status
```

Returns the current status of the server.

**Response:**
```json
{
  "version": "1.0.0",
  "uptime": "2h30m",
  "jobs_running": 5,
  "jobs_queued": 2,
  "cpu_usage": 45.2,
  "memory_usage": 1024,
  "disk_usage": 5120,
  "timestamp": "2023-01-01T00:00:00Z"
}
```

## Job Statuses

- `pending`: Job is waiting to be executed
- `running`: Job is currently executing
- `completed`: Job completed successfully
- `failed`: Job failed to execute
- `cancelled`: Job was cancelled

## Resource Limits

- **Timeout**: Maximum execution time in seconds (default: 30, max: 300)
- **Memory Limit**: Memory limit in MB (default: 128, max: 1024)
- **Network Access**: Allow network connections (default: false)

## Security

- Code execution occurs in isolated environments
- Resource limits prevent DoS attacks
- Network access can be restricted
- File system access is limited