package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"forgeai/pkg/sandbox"
)

// LocalExecutor is a basic implementation of the Executor interface
// that runs code using the local system's interpreters
type LocalExecutor struct {
	// Timeout for execution
	Timeout time.Duration

	// MemoryLimit in MB
	MemoryLimit int
}

// NewLocalExecutor creates a new LocalExecutor with default settings
func NewLocalExecutor() *LocalExecutor {
	return &LocalExecutor{
		Timeout:     30 * time.Second,
		MemoryLimit: 128, // 128 MB
	}
}

// Execute runs the provided code in a sandboxed environment
func (e *LocalExecutor) Execute(ctx context.Context, language, code string) (*sandbox.ExecutionResult, error) {
	// Create a temporary directory for execution
	tempDir, err := os.MkdirTemp("", "forgeai-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after execution

	// Write code to a temporary file
	filePath, err := e.writeCodeToFile(tempDir, language, code)
	if err != nil {
		return nil, fmt.Errorf("failed to write code to file: %w", err)
	}

	// Execute the file
	return e.ExecuteFile(ctx, filePath)
}

// ExecuteFile runs the provided file in a sandboxed environment
func (e *LocalExecutor) ExecuteFile(ctx context.Context, filePath string) (*sandbox.ExecutionResult, error) {
	// Get the language from the file extension
	language := e.getLanguageFromFile(filePath)

	// Get the command to execute the file
	cmdArgs, err := e.getCommandForLanguage(language, filePath)
	if err != nil {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	// Apply resource limits
	// Note: Full sandboxing would require more sophisticated techniques
	// like containers or system call filtering which are OS-specific

	// Set up context with timeout
	if e.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.Timeout)
		defer cancel()
	}

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

// SupportedLanguages returns a list of supported languages
func (e *LocalExecutor) SupportedLanguages() []string {
	return []string{"python", "go", "javascript"}
}

// writeCodeToFile writes the provided code to a temporary file
func (e *LocalExecutor) writeCodeToFile(tempDir, language, code string) (string, error) {
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
func (e *LocalExecutor) getLanguageFromFile(filePath string) string {
	// Simple implementation based on file extension
	// In a real implementation, this would be more sophisticated
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
func (e *LocalExecutor) getCommandForLanguage(language, filePath string) ([]string, error) {
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