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
	Host string
	Port int
}

// Server represents the API server
type Server struct {
	config     *Config
	router     *gin.Engine
	httpServer *http.Server
	jobManager *JobManager
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
		jobManager: NewJobManager(),
	}
}

// Config returns the server configuration
func (s *Server) Config() *Config {
	return s.config
}

// Start starts the API server
func (s *Server) Start(ctx context.Context) error {
	// Register routes
	s.registerRoutes()
	
	// Start the server
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	
	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// registerRoutes sets up the API routes
func (s *Server) registerRoutes() {
	// Root endpoint
	s.router.GET("/", s.handleRoot)
	
	// Health check endpoints
	s.router.GET("/healthz", s.handleHealthCheck)
	s.router.GET("/readyz", s.handleReadinessCheck)
	
	// API v1 routes
	v1 := s.router.Group("/v1")
	{
		v1.GET("/languages", s.handleListLanguages)
		v1.POST("/execute", s.handleExecuteCode)
		v1.POST("/execute/file", s.handleExecuteFile)
		v1.GET("/jobs/:id", s.handleGetJob)
		v1.DELETE("/jobs/:id", s.handleCancelJob)
		v1.GET("/jobs", s.handleListJobs)
		v1.GET("/status", s.handleGetStatus)
	}
}

// handleRoot handles the root endpoint
func (s *Server) handleRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ForgeAI API Server",
		"version": "1.0.0",
		"docs":    "/v1/docs",
	})
}

// handleHealthCheck handles the health check endpoint
func (s *Server) handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().UTC(),
	})
}

// handleReadinessCheck handles the readiness check endpoint
func (s *Server) handleReadinessCheck(c *gin.Context) {
	// In a real implementation, this would check if all dependencies are ready
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"time":   time.Now().UTC(),
	})
}

// handleListLanguages handles listing supported languages
func (s *Server) handleListLanguages(c *gin.Context) {
	// In a real implementation, this would get languages from the executor
	languages := []string{"python", "go", "javascript", "hello"} // Include plugin languages
	
	c.JSON(http.StatusOK, gin.H{
		"languages":  languages,
		"timestamp":  time.Now().UTC(),
	})
}

// handleExecuteCode handles code execution
func (s *Server) handleExecuteCode(c *gin.Context) {
	// Parse the request
	var req struct {
		Language      string `json:"language" binding:"required"`
		Code          string `json:"code" binding:"required"`
		Timeout       int    `json:"timeout"`
		MemoryLimit   int    `json:"memory_limit"`
		NetworkAccess bool   `json:"network_access"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Set default values
	if req.Timeout == 0 {
		req.Timeout = 30
	}
	if req.MemoryLimit == 0 {
		req.MemoryLimit = 128
	}
	
	// Create a job
	job := s.jobManager.CreateJob(req.Language, req.Code)
	job.Timeout = req.Timeout
	job.MemoryLimit = req.MemoryLimit
	job.NetworkAccess = req.NetworkAccess
	
	// Execute the job in a goroutine
	go s.jobManager.ExecuteJob(job)
	
	// Return the job ID
	c.JSON(http.StatusCreated, gin.H{
		"job_id": job.ID,
		"status": job.Status,
	})
}

// handleExecuteFile handles file execution
func (s *Server) handleExecuteFile(c *gin.Context) {
	// Parse the request
	var req struct {
		FilePath      string `json:"file_path" binding:"required"`
		Timeout       int    `json:"timeout"`
		MemoryLimit   int    `json:"memory_limit"`
		NetworkAccess bool   `json:"network_access"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Set default values
	if req.Timeout == 0 {
		req.Timeout = 30
	}
	if req.MemoryLimit == 0 {
		req.MemoryLimit = 128
	}
	
	// Create a job
	job := s.jobManager.CreateFileJob(req.FilePath)
	job.Timeout = req.Timeout
	job.MemoryLimit = req.MemoryLimit
	job.NetworkAccess = req.NetworkAccess
	
	// Execute the job in a goroutine
	go s.jobManager.ExecuteJob(job)
	
	// Return the job ID
	c.JSON(http.StatusCreated, gin.H{
		"job_id": job.ID,
		"status": job.Status,
	})
}

// handleGetJob handles getting job status
func (s *Server) handleGetJob(c *gin.Context) {
	jobID := c.Param("id")
	
	job, ok := s.jobManager.GetJob(jobID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}
	
	// Convert job to response format
	resp := gin.H{
		"job_id":      job.ID,
		"status":      job.Status,
		"language":    job.Language,
		"timeout":     job.Timeout,
		"memory_limit": job.MemoryLimit,
		"network_access": job.NetworkAccess,
		"created_at":  job.CreatedAt,
		"started_at":  job.StartedAt,
		"completed_at": job.CompletedAt,
	}
	
	// Add result if job is completed
	if job.Status == "completed" && job.Result != nil {
		resp["stdout"] = job.Result.Stdout
		resp["stderr"] = job.Result.Stderr
		resp["exit_code"] = job.Result.ExitCode
		resp["duration"] = job.Result.Duration.String()
	}
	
	// Add error if job failed
	if job.Status == "failed" && job.Error != "" {
		resp["error"] = job.Error
	}
	
	c.JSON(http.StatusOK, resp)
}

// handleCancelJob handles canceling a job
func (s *Server) handleCancelJob(c *gin.Context) {
	jobID := c.Param("id")
	
	if s.jobManager.CancelJob(jobID) {
		c.JSON(http.StatusOK, gin.H{
			"job_id":   jobID,
			"status":   "cancelled",
			"message":  "Job cancelled successfully",
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"job_id":   jobID,
			"status":   "error",
			"message":  "Job not found or cannot be cancelled",
		})
	}
}

// handleListJobs handles listing jobs
func (s *Server) handleListJobs(c *gin.Context) {
	status := c.Query("status")
	language := c.Query("language")
	
	jobs := s.jobManager.ListJobs(status, language)
	
	// Convert jobs to response format
	jobList := make([]gin.H, len(jobs))
	for i, job := range jobs {
		jobList[i] = gin.H{
			"job_id":      job.ID,
			"status":      job.Status,
			"language":    job.Language,
			"created_at":  job.CreatedAt,
			"started_at":  job.StartedAt,
			"completed_at": job.CompletedAt,
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"jobs": jobList,
		"count": len(jobList),
	})
}

// handleGetStatus handles getting server status
func (s *Server) handleGetStatus(c *gin.Context) {
	// In a real implementation, this would return actual server metrics
	
	c.JSON(http.StatusOK, gin.H{
		"version":        "1.0.0",
		"uptime":         "2h30m",
		"jobs_running":   5,
		"jobs_queued":    2,
		"cpu_usage":      45.2,
		"memory_usage":   1024,
		"disk_usage":     5120,
		"timestamp":      time.Now().UTC(),
	})
}