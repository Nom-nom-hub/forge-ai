package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	"forgeai/pkg/executor"
	"forgeai/pkg/sandbox"
)

// Job represents a code execution job
type Job struct {
	ID          string
	Status      string // pending, running, completed, failed, cancelled
	Language    string
	Code        string
	FilePath    string
	Timeout     int
	MemoryLimit int
	NetworkAccess bool
	Result      *sandbox.ExecutionResult
	Error       string
	CreatedAt   time.Time
	StartedAt   time.Time
	CompletedAt time.Time
}

// JobManager manages execution jobs
type JobManager struct {
	jobs map[string]*Job
	mu   sync.RWMutex
}

// NewJobManager creates a new job manager
func NewJobManager() *JobManager {
	return &JobManager{
		jobs: make(map[string]*Job),
	}
}

// CreateJob creates a new job
func (jm *JobManager) CreateJob(language, code string) *Job {
	job := &Job{
		ID:        generateJobID(),
		Status:    "pending",
		Language:  language,
		Code:      code,
		Timeout:   30,
		MemoryLimit: 128,
		CreatedAt: time.Now(),
	}
	
	jm.mu.Lock()
	jm.jobs[job.ID] = job
	jm.mu.Unlock()
	
	return job
}

// CreateFileJob creates a new file execution job
func (jm *JobManager) CreateFileJob(filePath string) *Job {
	job := &Job{
		ID:        generateJobID(),
		Status:    "pending",
		FilePath:  filePath,
		Timeout:   30,
		MemoryLimit: 128,
		CreatedAt: time.Now(),
	}
	
	jm.mu.Lock()
	jm.jobs[job.ID] = job
	jm.mu.Unlock()
	
	return job
}

// GetJob retrieves a job by ID
func (jm *JobManager) GetJob(id string) (*Job, bool) {
	jm.mu.RLock()
	job, ok := jm.jobs[id]
	jm.mu.RUnlock()
	return job, ok
}

// ListJobs lists all jobs with optional filters
func (jm *JobManager) ListJobs(status, language string) []*Job {
	jm.mu.RLock()
	defer jm.mu.RUnlock()
	
	var jobs []*Job
	for _, job := range jm.jobs {
		if (status == "" || job.Status == status) && 
		   (language == "" || job.Language == language) {
			jobs = append(jobs, job)
		}
	}
	
	return jobs
}

// CancelJob cancels a job
func (jm *JobManager) CancelJob(id string) bool {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	
	job, ok := jm.jobs[id]
	if !ok {
		return false
	}
	
	// Only cancel jobs that are pending or running
	if job.Status == "pending" || job.Status == "running" {
		job.Status = "cancelled"
		job.CompletedAt = time.Now()
		return true
	}
	
	return false
}

// ExecuteJob executes a job
func (jm *JobManager) ExecuteJob(job *Job) {
	jm.mu.Lock()
	job.Status = "running"
	job.StartedAt = time.Now()
	jm.mu.Unlock()
	
	// Create executor
	exec := executor.NewLocalExecutor()
	exec.Timeout = time.Duration(job.Timeout) * time.Second
	exec.MemoryLimit = job.MemoryLimit
	
	var result *sandbox.ExecutionResult
	var err error
	
	// Execute based on job type
	if job.Code != "" {
		result, err = exec.Execute(context.Background(), job.Language, job.Code)
	} else if job.FilePath != "" {
		result, err = exec.ExecuteFile(context.Background(), job.FilePath)
	} else {
		err = fmt.Errorf("invalid job: no code or file path")
	}
	
	// Update job with results
	jm.mu.Lock()
	defer jm.mu.Unlock()
	
	job.CompletedAt = time.Now()
	
	if err != nil {
		job.Status = "failed"
		job.Error = err.Error()
	} else {
		job.Status = "completed"
		job.Result = result
	}
}

// generateJobID generates a unique job ID
func generateJobID() string {
	return fmt.Sprintf("job-%d", time.Now().UnixNano())
}