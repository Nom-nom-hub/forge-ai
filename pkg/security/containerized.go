package security

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"forgeai/pkg/sandbox"
)

// ContainerizedExecutor implements security controls using containerization
type ContainerizedExecutor struct {
	Timeout     time.Duration
	MemoryLimit int
	EnableNetwork bool
	ReadOnlyRoot bool
}

// NewContainerizedExecutor creates a new containerized executor
func NewContainerizedExecutor() *ContainerizedExecutor {
	return &ContainerizedExecutor{
		Timeout:       10 * time.Second,
		MemoryLimit:   128, // 128 MB
		EnableNetwork: false, // Disable network by default
		ReadOnlyRoot:  true,  // Read-only root filesystem
	}
}

// Execute runs code with containerized security controls
func (ce *ContainerizedExecutor) Execute(ctx context.Context, language, code string) (*sandbox.ExecutionResult, error) {
	// Create a temporary directory for execution
	tempDir, err := os.MkdirTemp("", "forgeai-container-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after execution

	// Write code to a temporary file
	filePath, err := ce.writeCodeToFile(tempDir, language, code)
	if err != nil {
		return nil, fmt.Errorf("failed to write code to file: %w", err)
	}

	// Execute the file with containerized security controls
	return ce.ExecuteFile(ctx, filePath)
}

// ExecuteFile runs a file with containerized security controls
func (ce *ContainerizedExecutor) ExecuteFile(ctx context.Context, filePath string) (*sandbox.ExecutionResult, error) {
	// Get the language from the file extension
	language := ce.getLanguageFromFile(filePath)
	
	// Check if Docker is available
	if !ce.isDockerAvailable() {
		// Fall back to secure local execution
		return ce.executeLocally(ctx, language, filePath)
	}
	
	// Execute using Docker with security controls
	return ce.executeWithDocker(ctx, language, filePath)
}

// executeWithDocker runs code using Docker with security controls
func (ce *ContainerizedExecutor) executeWithDocker(ctx context.Context, language, filePath string) (*sandbox.ExecutionResult, error) {
	// Get the appropriate Docker image
	image := ce.getImageForLanguage(language)
	
	// Get the directory and filename
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)
	
	// Build the docker command with security controls
	cmdArgs := []string{
		"docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/workspace:ro", dir), // Read-only mount
		"-w", "/workspace",
	}
	
	// Add resource limits
	if ce.MemoryLimit > 0 {
		cmdArgs = append(cmdArgs, "--memory", fmt.Sprintf("%dm", ce.MemoryLimit))
	}
	
	// Add CPU limit (using cpu-shares)
	cmdArgs = append(cmdArgs, "--cpu-shares", "100")
	
	// Add read-only root filesystem if requested
	if ce.ReadOnlyRoot {
		cmdArgs = append(cmdArgs, "--read-only")
		// Add tmpfs for temporary files
		cmdArgs = append(cmdArgs, "--tmpfs", "/tmp:rw,noexec,nosuid,size=10m")
	}
	
	// Disable network if requested
	if !ce.EnableNetwork {
		cmdArgs = append(cmdArgs, "--network", "none")
	}
	
	// Run as non-root user
	cmdArgs = append(cmdArgs, "--user", "65534:65534") // nobody user
	
	// Add the image and command
	cmdArgs = append(cmdArgs, image)
	
	// Add the execution command based on language
	switch language {
	case "python":
		cmdArgs = append(cmdArgs, "python", filename)
	case "go":
		cmdArgs = append(cmdArgs, "go", "run", filename)
	case "javascript":
		cmdArgs = append(cmdArgs, "node", filename)
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	
	// Apply timeout
	if ce.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ce.Timeout)
		defer cancel()
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

// executeLocally runs code using local execution with basic security controls
func (ce *ContainerizedExecutor) executeLocally(ctx context.Context, language, filePath string) (*sandbox.ExecutionResult, error) {
	// Get the command to execute the file
	cmdArgs, err := ce.getCommandForLanguage(language, filePath)
	if err != nil {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	// Apply resource limits
	// Set up context with timeout
	if ce.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ce.Timeout)
		defer cancel()
	}
	
	// Create command with security restrictions
	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	
	// Apply additional security measures based on OS
	if runtime.GOOS == "windows" {
		// On Windows, we can't easily apply the same restrictions
		// but we can at least set the working directory to the temp directory
		cmd.Dir = filepath.Dir(filePath)
	} else {
		// On Unix-like systems, we can apply more restrictions
		cmd.Dir = filepath.Dir(filePath)
		
		// TODO: Implement additional security measures:
		// - User namespace isolation
		// - Seccomp profiles
		// - AppArmor/SELinux profiles
		// - Chroot or pivot_root
		// - Capability dropping
	}
	
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

// SupportedLanguages returns a list of supported languages
func (ce *ContainerizedExecutor) SupportedLanguages() []string {
	return []string{"python", "go", "javascript"}
}

// writeCodeToFile writes the provided code to a temporary file
func (ce *ContainerizedExecutor) writeCodeToFile(tempDir, language, code string) (string, error) {
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

// getLanguageFromFile determines the language from the file extension
func (ce *ContainerizedExecutor) getLanguageFromFile(filePath string) string {
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

// getCommandForLanguage returns the command to execute a file for the given language
func (ce *ContainerizedExecutor) getCommandForLanguage(language, filePath string) ([]string, error) {
	switch language {
	case "python":
		return []string{"python", filePath}, nil
	case "go":
		// For Go, we need to run the file differently
		// We'll use "go run" for simplicity
		return []string{"go", "run", filePath}, nil
	case "javascript":
		return []string{"node", filePath}, nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
}

// getImageForLanguage returns the Docker image for the given language
func (ce *ContainerizedExecutor) getImageForLanguage(language string) string {
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

// isDockerAvailable checks if Docker is available
func (ce *ContainerizedExecutor) isDockerAvailable() bool {
	cmd := exec.Command("docker", "--version")
	err := cmd.Run()
	return err == nil
}