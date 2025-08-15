// Package api provides the REST API for ForgeAI
//
// This package implements a REST API server for ForgeAI that allows
// programmatic access to code execution capabilities.
//
// Example usage:
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
package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Config holds the API server configuration
type Config struct {
	// Host is the host address to bind to
	Host string
	
	// Port is the port to listen on
	Port int
	
	// TLS configuration
	TLS struct {
		Enabled   bool
		CertFile  string
		KeyFile   string
	}
	
	// Rate limiting
	RateLimit struct {
		RequestsPerHour int
		RequestsPerMinute int
	}
}

// Server represents the API server
type Server struct {
	router *gin.Engine
	config *Config
	httpServer *http.Server
}

// NewServer creates a new API server
func NewServer(config *Config) *Server {
	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)
	
	// Create the router
	router := gin.New()
	
	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// Create the HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: router,
	}
	
	return &Server{
		config:     config,
		router:     router,
		httpServer: httpServer,
	}
}

// Config returns the server configuration
func (s *Server) Config() *Config {
	return s.config
}

// Start starts the API server
func (s *Server) Start() error {
	// Register routes
	s.registerRoutes()
	
	// Start the server
	if s.config.TLS.Enabled {
		if err := s.httpServer.ListenAndServeTLS(s.config.TLS.CertFile, s.config.TLS.KeyFile); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("failed to start TLS server: %w", err)
		}
	} else {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("failed to start server: %w", err)
		}
	}
	
	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}