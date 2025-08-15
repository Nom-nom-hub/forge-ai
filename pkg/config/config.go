package config

import (
	"time"
)

// Config represents the configuration for the sandbox executor
type Config struct {
	// Timeout for execution
	Timeout time.Duration

	// MemoryLimit in MB
	MemoryLimit int

	// CPU shares (Linux only)
	CPUShares int

	// Network access allowed
	NetworkAccess bool

	// Allowed directories for file access
	AllowedDirs []string
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		Timeout:       30 * time.Second,
		MemoryLimit:   128, // 128 MB
		CPUShares:     100, // 10% of CPU (Linux only)
		NetworkAccess: false,
		AllowedDirs:   []string{},
	}
}