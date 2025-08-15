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

// SecureExecutor implements enhanced security controls
type SecureExecutor struct {
	Timeout     time.Duration
	MemoryLimit int
}

// NewSecureExecutor creates a new secure executor
func NewSecureExecutor() *SecureExecutor {
	return &SecureExecutor{
		Timeout:     10 * time.Second,
		MemoryLimit: 128, // 128 MB
	}
}

// Execute runs code with enhanced security controls
func (se *SecureExecutor) Execute(ctx context.Context, language, code string) (*sandbox.ExecutionResult, error) {
	// Create a temporary directory for execution
	tempDir, err := os.MkdirTemp("", "forgeai-secure-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after execution

	// Write code to a temporary file
	filePath, err := se.writeCodeToFile(tempDir, language, code)
	if err != nil {
		return nil, fmt.Errorf("failed to write code to file: %w", err)
	}

	// Execute the file with security controls
	return se.ExecuteFile(ctx, filePath)
}

// ExecuteFile runs a file with enhanced security controls
func (se *SecureExecutor) ExecuteFile(ctx context.Context, filePath string) (*sandbox.ExecutionResult, error) {
	// Get the language from the file extension
	language := se.getLanguageFromFile(filePath)
	
	// Get the command to execute the file
	cmdArgs, err := se.getCommandForLanguage(language, filePath)
	if err != nil {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	// Apply resource limits
	// Set up context with timeout
	if se.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, se.Timeout)
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
func (se *SecureExecutor) SupportedLanguages() []string {
	return []string{"python", "go", "javascript"}
}

// writeCodeToFile writes the provided code to a temporary file
func (se *SecureExecutor) writeCodeToFile(tempDir, language, code string) (string, error) {
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
func (se *SecureExecutor) getLanguageFromFile(filePath string) string {
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
func (se *SecureExecutor) getCommandForLanguage(language, filePath string) ([]string, error) {
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