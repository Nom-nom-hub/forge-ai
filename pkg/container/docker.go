package container

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"forgeai/pkg/sandbox"
)

// DockerExecutor implements the sandbox.Executor interface using Docker
type DockerExecutor struct {
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

// NewDockerExecutor creates a new DockerExecutor with default settings
func NewDockerExecutor() *DockerExecutor {
	return &DockerExecutor{
		Timeout:       30 * time.Second,
		MemoryLimit:   128, // 128 MB
		CPUShares:     100, // 10% of CPU (Linux only)
		NetworkAccess: false,
		ReadOnlyRoot:  true,
	}
}

// Execute runs the provided code in a Docker container
func (d *DockerExecutor) Execute(ctx context.Context, language, code string) (*sandbox.ExecutionResult, error) {
	// Create a temporary directory for execution
	tempDir, err := os.MkdirTemp("", "forgeai-docker-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after execution

	// Write code to a temporary file
	filePath, err := d.writeCodeToFile(tempDir, language, code)
	if err != nil {
		return nil, fmt.Errorf("failed to write code to file: %w", err)
	}

	// Execute the file in a container
	return d.ExecuteFile(ctx, filePath)
}

// ExecuteFile runs the provided file in a Docker container
func (d *DockerExecutor) ExecuteFile(ctx context.Context, filePath string) (*sandbox.ExecutionResult, error) {
	// Get the language from the file extension
	language := d.getLanguageFromFile(filePath)
	
	// Validate language support
	if !d.isLanguageSupported(language) {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	
	// Select appropriate container image
	image := d.getImageForLanguage(language)
	
	// Set up context with timeout
	if d.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, d.Timeout)
		defer cancel()
	}
	
	// Prepare container configuration
	config := &DockerConfig{
		Image:         image,
		Timeout:       d.Timeout,
		MemoryLimit:   d.MemoryLimit,
		CPUShares:     d.CPUShares,
		NetworkAccess: d.NetworkAccess,
		ReadOnlyRoot:  d.ReadOnlyRoot,
		FilePath:      filePath,
		Language:      language,
	}
	
	// Execute in container
	result, err := d.runContainer(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("container execution failed: %w", err)
	}
	
	return result, nil
}

// SupportedLanguages returns a list of supported languages
func (d *DockerExecutor) SupportedLanguages() []string {
	return []string{"python", "go", "javascript"}
}

// Internal methods

func (d *DockerExecutor) writeCodeToFile(tempDir, language, code string) (string, error) {
	var fileName string
	
	switch language {
	case "python":
		fileName = "main.py"
	case "go":
		fileName = "main.go"
	case "javascript":
		fileName = "main.js"
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}
	
	filePath := filepath.Join(tempDir, fileName)
	
	err := os.WriteFile(filePath, []byte(code), 0644)
	if err != nil {
		return "", err
	}
	
	return filePath, nil
}

func (d *DockerExecutor) getLanguageFromFile(filePath string) string {
	switch {
	case filepath.Ext(filePath) == ".py":
		return "python"
	case filepath.Ext(filePath) == ".go":
		return "go"
	case filepath.Ext(filePath) == ".js":
		return "javascript"
	default:
		return "unknown"
	}
}

func (d *DockerExecutor) isLanguageSupported(language string) bool {
	supported := d.SupportedLanguages()
	for _, lang := range supported {
		if lang == language {
			return true
		}
	}
	return false
}

func (d *DockerExecutor) getImageForLanguage(language string) string {
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

func (d *DockerExecutor) runContainer(ctx context.Context, config *DockerConfig) (*sandbox.ExecutionResult, error) {
	// Check if Docker is available
	if !d.isDockerAvailable() {
		return nil, fmt.Errorf("docker is not available")
	}
	
	// Pull the image if it doesn't exist
	if err := d.pullImage(ctx, config.Image); err != nil {
		return nil, fmt.Errorf("failed to pull image %s: %w", config.Image, err)
	}
	
	// Get the directory and filename
	dir := filepath.Dir(config.FilePath)
	filename := filepath.Base(config.FilePath)
	
	// Build the docker command
	cmdArgs := []string{
		"docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/workspace", dir),
		"-w", "/workspace",
	}
	
	// Add resource limits
	if config.MemoryLimit > 0 {
		cmdArgs = append(cmdArgs, "--memory", fmt.Sprintf("%dm", config.MemoryLimit))
	}
	
	if config.CPUShares > 0 {
		cmdArgs = append(cmdArgs, "--cpu-shares", fmt.Sprintf("%d", config.CPUShares))
	}
	
	// Add read-only root filesystem if requested
	if config.ReadOnlyRoot {
		cmdArgs = append(cmdArgs, "--read-only")
	}
	
	// Disable network if requested
	if !config.NetworkAccess {
		cmdArgs = append(cmdArgs, "--network", "none")
	}
	
	// Add the image and command
	cmdArgs = append(cmdArgs, config.Image)
	
	// Add the execution command based on language
	switch config.Language {
	case "python":
		cmdArgs = append(cmdArgs, "python", filename)
	case "go":
		cmdArgs = append(cmdArgs, "go", "run", filename)
	case "javascript":
		cmdArgs = append(cmdArgs, "node", filename)
	default:
		return nil, fmt.Errorf("unsupported language: %s", config.Language)
	}
	
	// Create the command
	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	
	// Capture output
	result := &sandbox.ExecutionResult{
		Stdout: "",
		Stderr: "",
	}
	
	start := time.Now()
	
	// Run the command
	output, err := cmd.CombinedOutput()
	
	result.Duration = time.Since(start)
	result.Stdout = string(output)
	
	// Check if the context was cancelled (timeout)
	if ctx.Err() == context.DeadlineExceeded {
		result.Stderr = "Execution timed out"
		result.ExitCode = -1
		return result, nil
	}
	
	// Get exit code
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = -1
			result.Stderr = err.Error()
		}
	} else {
		result.ExitCode = 0
	}
	
	return result, nil
}

func (d *DockerExecutor) isDockerAvailable() bool {
	cmd := exec.Command("docker", "--version")
	err := cmd.Run()
	return err == nil
}

func (d *DockerExecutor) pullImage(ctx context.Context, image string) error {
	// Check if image exists locally
	cmd := exec.CommandContext(ctx, "docker", "image", "inspect", image)
	err := cmd.Run()
	if err == nil {
		// Image exists, no need to pull
		return nil
	}
	
	// Image doesn't exist, pull it
	cmd = exec.CommandContext(ctx, "docker", "pull", image)
	return cmd.Run()
}

// DockerConfig holds configuration for Docker execution
type DockerConfig struct {
	Image         string
	Timeout       time.Duration
	MemoryLimit   int
	CPUShares     int
	NetworkAccess bool
	ReadOnlyRoot  bool
	FilePath      string
	Language      string
}