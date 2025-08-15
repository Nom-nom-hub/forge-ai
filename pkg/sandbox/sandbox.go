package sandbox

import (
	"context"
	"time"
)

// ExecutionResult represents the result of a sandboxed execution
type ExecutionResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
}

// Executor defines the interface for executing code in a sandbox
type Executor interface {
	// Execute runs the provided code in a sandboxed environment
	Execute(ctx context.Context, language, code string) (*ExecutionResult, error)

	// ExecuteFile runs the provided file in a sandboxed environment
	ExecuteFile(ctx context.Context, filePath string) (*ExecutionResult, error)

	// SupportedLanguages returns a list of supported languages
	SupportedLanguages() []string
}