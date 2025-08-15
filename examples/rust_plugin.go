// Package main implements a simple plugin for executing Rust code
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"forgeai/pkg/plugin"
	"forgeai/pkg/sandbox"
)

// RustExecutor is a simple executor that runs Rust code
type RustExecutor struct {
	// Timeout for execution
	Timeout time.Duration

	// MemoryLimit in MB
	MemoryLimit int
}

// New creates a new RustExecutor
func New() (plugin.Executor, error) {
	return &RustExecutor{
		Timeout:     30 * time.Second,
		MemoryLimit: 128, // 128 MB
	}, nil
}

// Execute runs the provided Rust code
func (r *RustExecutor) Execute(ctx context.Context, language, code string) (*sandbox.ExecutionResult, error) {
	// Only support "rust" language
	if language != "rust" {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	// Create a temporary directory for execution
	tempDir, err := os.MkdirTemp("", "forgeai-rust-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after execution

	// Write code to a temporary file
	filePath := filepath.Join(tempDir, "main.rs")
	err = os.WriteFile(filePath, []byte(code), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write code to file: %w", err)
	}

	// Execute the file
	return r.ExecuteFile(ctx, filePath)
}

// ExecuteFile runs the provided Rust file
func (r *RustExecutor) ExecuteFile(ctx context.Context, filePath string) (*sandbox.ExecutionResult, error) {
	// Set up context with timeout
	if r.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.Timeout)
		defer cancel()
	}

	// Get the directory containing the file
	dir := filepath.Dir(filePath)

	// Create a simple Cargo.toml file
	cargoToml := `[package]
name = "forgeai-exec"
version = "0.1.0"
edition = "2021"

[[bin]]
name = "main"
path = "main.rs"
`
	cargoPath := filepath.Join(dir, "Cargo.toml")
	err := os.WriteFile(cargoPath, []byte(cargoToml), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write Cargo.toml: %w", err)
	}

	// Build the Rust code
	buildCmd := exec.CommandContext(ctx, "cargo", "build", "--release")
	buildCmd.Dir = dir

	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		return &sandbox.ExecutionResult{
			Stdout:   "",
			Stderr:   string(buildOutput),
			ExitCode: -1,
			Duration: 0,
		}, nil
	}

	// Execute the built binary
	binaryPath := filepath.Join(dir, "target", "release", "main")
	cmd := exec.CommandContext(ctx, binaryPath)

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
func (r *RustExecutor) SupportedLanguages() []string {
	return []string{"rust"}
}