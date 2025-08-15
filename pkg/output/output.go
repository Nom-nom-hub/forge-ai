package output

import (
	"time"
)

// ExecutionResult represents the result of a sandboxed execution
type ExecutionResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
}

// Printer handles formatting and printing of execution results
type Printer struct {
	JSONOutput bool
}

// NewPrinter creates a new Printer
func NewPrinter(jsonOutput bool) *Printer {
	return &Printer{
		JSONOutput: jsonOutput,
	}
}