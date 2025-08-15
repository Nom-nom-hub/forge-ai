package container

import (
	"context"
	"fmt"
	"time"

	"forgeai/pkg/sandbox"
)

// ContainerExecutor implements the sandbox.Executor interface using containerization
type ContainerExecutor struct {
	// Engine specifies the container engine to use (docker, gvisor, firecracker)
	Engine string
	
	// Timeout for execution
	Timeout time.Duration
	
	// MemoryLimit in MB
	MemoryLimit int
	
	// CPUShares for CPU allocation
	CPUShares int
	
	// NetworkAccess controls network access
	NetworkAccess bool
	
	// ReadOnlyRoot makes the root filesystem read-only
	ReadOnlyRoot bool
}

// NewContainerExecutor creates a new ContainerExecutor with default settings
func NewContainerExecutor() *ContainerExecutor {
	return &ContainerExecutor{
		Engine:        "docker",
		Timeout:       30 * time.Second,
		MemoryLimit:   128, // 128 MB
		CPUShares:     100, // 10% of CPU (Linux only)
		NetworkAccess: false,
		ReadOnlyRoot:  true,
	}
}

// Execute runs the provided code in a containerized environment
func (c *ContainerExecutor) Execute(ctx context.Context, language, code string) (*sandbox.ExecutionResult, error) {
	// Create a temporary directory for execution
	tempDir, err := c.createTempDir()
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer c.cleanupTempDir(tempDir) // Clean up after execution

	// Write code to a temporary file
	filePath, err := c.writeCodeToFile(tempDir, language, code)
	if err != nil {
		return nil, fmt.Errorf("failed to write code to file: %w", err)
	}

	// Execute the file in a container
	return c.ExecuteFile(ctx, filePath)
}

// ExecuteFile runs the provided file in a containerized environment
func (c *ContainerExecutor) ExecuteFile(ctx context.Context, filePath string) (*sandbox.ExecutionResult, error) {
	// Get the language from the file extension
	language := c.getLanguageFromFile(filePath)
	
	// Validate language support
	if !c.isLanguageSupported(language) {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	
	// Select appropriate container image
	image := c.getImageForLanguage(language)
	
	// Set up context with timeout
	if c.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.Timeout)
		defer cancel()
	}
	
	// Prepare container configuration
	config := &ContainerConfig{
		Image:         image,
		Timeout:       c.Timeout,
		MemoryLimit:   c.MemoryLimit,
		CPUShares:     c.CPUShares,
		NetworkAccess: c.NetworkAccess,
		ReadOnlyRoot:  c.ReadOnlyRoot,
		FilePath:      filePath,
		Language:      language,
	}
	
	// Execute in container
	result, err := c.runContainer(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("container execution failed: %w", err)
	}
	
	return result, nil
}

// SupportedLanguages returns a list of supported languages
func (c *ContainerExecutor) SupportedLanguages() []string {
	// For now, return the same languages as the local executor
	// In a full implementation, this could be dynamic based on available container images
	return []string{"python", "go", "javascript"}
}

// Internal methods would be implemented here in a full implementation
// For now, we'll add placeholder implementations

func (c *ContainerExecutor) createTempDir() (string, error) {
	// In a real implementation, this would create a temporary directory
	return "/tmp/forgeai-container", nil
}

func (c *ContainerExecutor) cleanupTempDir(dir string) {
	// In a real implementation, this would clean up the temporary directory
}

func (c *ContainerExecutor) writeCodeToFile(tempDir, language, code string) (string, error) {
	// In a real implementation, this would write code to a file
	return fmt.Sprintf("%s/main.%s", tempDir, c.getFileExtension(language)), nil
}

func (c *ContainerExecutor) getLanguageFromFile(filePath string) string {
	// In a real implementation, this would determine language from file extension
	return "python"
}

func (c *ContainerExecutor) isLanguageSupported(language string) bool {
	// In a real implementation, this would check language support
	return true
}

func (c *ContainerExecutor) getImageForLanguage(language string) string {
	// In a real implementation, this would return appropriate container images
	switch language {
	case "python":
		return "python:3.9-alpine"
	case "go":
		return "golang:1.19-alpine"
	case "javascript":
		return "node:16-alpine"
	default:
		return "alpine:latest"
	}
}

func (c *ContainerExecutor) getFileExtension(language string) string {
	// In a real implementation, this would return file extensions
	switch language {
	case "python":
		return "py"
	case "go":
		return "go"
	case "javascript":
		return "js"
	default:
		return "txt"
	}
}

func (c *ContainerExecutor) runContainer(ctx context.Context, config *ContainerConfig) (*sandbox.ExecutionResult, error) {
	// In a real implementation, this would:
	// 1. Pull the container image if needed
	// 2. Create and start a container with the specified configuration
	// 3. Copy the code file into the container
	// 4. Execute the code
	// 5. Capture stdout, stderr, and exit code
	// 6. Clean up the container
	
	// For now, return a placeholder result
	result := &sandbox.ExecutionResult{
		Stdout:   fmt.Sprintf("Container execution would run %s code in %s container", config.Language, config.Image),
		Stderr:   "",
		ExitCode: 0,
		Duration: 100 * time.Millisecond,
	}
	
	return result, nil
}

// ContainerConfig holds configuration for container execution
type ContainerConfig struct {
	Image         string
	Timeout       time.Duration
	MemoryLimit   int
	CPUShares     int
	NetworkAccess bool
	ReadOnlyRoot  bool
	FilePath      string
	Language      string
}