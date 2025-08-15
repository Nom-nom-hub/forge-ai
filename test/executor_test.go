package test

import (
	"testing"
	"time"

	"forgeai/pkg/executor"
)

func TestLocalExecutor(t *testing.T) {
	// Test that we can create a basic executor
	exec := executor.NewLocalExecutor()
	if exec == nil {
		t.Error("Failed to create LocalExecutor instance")
	}

	languages := exec.SupportedLanguages()
	if len(languages) == 0 {
		t.Error("No supported languages found")
	}

	// Test timeout configuration
	exec.Timeout = 5 * time.Second
	if exec.Timeout != 5*time.Second {
		t.Error("Failed to set timeout")
	}

	// Test memory limit configuration
	exec.MemoryLimit = 64
	if exec.MemoryLimit != 64 {
		t.Error("Failed to set memory limit")
	}

	t.Logf("Supported languages: %v", languages)
}