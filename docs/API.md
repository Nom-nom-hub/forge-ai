// Package api provides the REST API for ForgeAI.
//
// This package implements a REST API server for ForgeAI that allows
// programmatic access to code execution capabilities.
//
// # Installation
//
//	go get forgeai/pkg/api
//
// # Usage
//
//	package main
//
//	import (
//	    "forgeai/pkg/api"
//	)
//
//	func main() {
//	    config := &api.Config{
//	        Host: "0.0.0.0",
//	        Port: 8080,
//	    }
//
//	    server := api.NewServer(config)
//	    if err := server.Start(); err != nil {
//	        panic(err)
//	    }
//	}
//
// # Features
//
//   - RESTful API design
//   - JSON request/response format
//   - Asynchronous job execution
//   - Resource limits (CPU, memory, timeout)
//   - Health and readiness checks
//   - Rate limiting
//   - TLS support
//
// # Endpoints
//
//   - GET / - API information
//   - GET /healthz - Health check
//   - GET /readyz - Readiness check
//   - GET /v1/languages - List supported languages
//   - POST /v1/execute - Execute code
//   - POST /v1/execute/file - Execute file
//   - GET /v1/jobs/{id} - Get job status
//   - DELETE /v1/jobs/{id} - Cancel job
//   - GET /v1/jobs - List jobs
//   - GET /v1/status - Server status
package api