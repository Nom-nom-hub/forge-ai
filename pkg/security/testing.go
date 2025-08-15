package security

import (
	"context"
	"fmt"
	"strings"
	"time"

	"forgeai/pkg/sandbox"
)

// TestCase represents a security test case
type TestCase struct {
	Name        string
	Code        string
	Language    string
	Description string
	Category    string
	ExpectedResult TestResult
}

// TestResult represents the expected result of a test
type TestResult struct {
	ShouldBeContained bool
	ExpectedExitCode  int
	ExpectedOutput    string
}

// TestReport represents the result of a security test
type TestReport struct {
	TestCase    TestCase
	ActualResult *sandbox.ExecutionResult
	ExitCode    int
	Duration    time.Duration
	Passed      bool
	Error       error
}

// TestFramework manages security testing
type TestFramework struct {
	executor Executor
	tests    []TestCase
}

// Executor interface for code execution
type Executor interface {
	Execute(ctx context.Context, language, code string) (*sandbox.ExecutionResult, error)
	ExecuteFile(ctx context.Context, filePath string) (*sandbox.ExecutionResult, error)
	SupportedLanguages() []string
}

// NewTestFramework creates a new security test framework
func NewTestFramework() *TestFramework {
	exec := NewContainerizedExecutor()
	
	return &TestFramework{
		executor: exec,
		tests: []TestCase{
			{
				Name:        "CPU Exhaustion - Infinite Loop",
				Code:        "while True: pass",
				Language:    "python",
				Description: "Tests CPU exhaustion through infinite loop",
				Category:    "Resource Exhaustion",
				ExpectedResult: TestResult{
					ShouldBeContained: true,
					ExpectedExitCode:  -1, // Timeout
				},
			},
			{
				Name:        "Memory Exhaustion - Large List",
				Code:        "a = []; while True: a.append('x' * 1000000)",
				Language:    "python",
				Description: "Tests memory exhaustion through large list allocation",
				Category:    "Resource Exhaustion",
				ExpectedResult: TestResult{
					ShouldBeContained: true,
					ExpectedExitCode:  1, // Error
				},
			},
			{
				Name:        "File System Access - Sensitive File",
				Code:        "try:\n    with open('/etc/passwd', 'r') as f:\n        print(f.read())\nexcept:\n    print('Access denied')",
				Language:    "python",
				Description: "Tests file system access to sensitive files",
				Category:    "File System Attacks",
				ExpectedResult: TestResult{
					ShouldBeContained: true,
					ExpectedOutput:    "Access denied",
				},
			},
			{
				Name:        "Network Access - External Connection",
				Code:        "import socket\ntry:\n    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)\n    s.connect(('google.com', 80))\n    print('Connection successful')\n    s.close()\nexcept:\n    print('Network access denied')",
				Language:    "python",
				Description: "Tests network access to external sites",
				Category:    "Network Attacks",
				ExpectedResult: TestResult{
					ShouldBeContained: true,
					ExpectedOutput:    "Network access denied",
				},
			},
			{
				Name:        "Valid Code Execution",
				Code:        "print('Hello, World!')",
				Language:    "python",
				Description: "Tests valid code execution",
				Category:    "Valid Execution",
				ExpectedResult: TestResult{
					ShouldBeContained: false,
					ExpectedExitCode:  0,
					ExpectedOutput:    "Hello, World!",
				},
			},
		},
	}
}

// RunTests runs all security tests
func (tf *TestFramework) RunTests() []TestReport {
	reports := make([]TestReport, len(tf.tests))
	
	for i, test := range tf.tests {
		reports[i] = tf.RunTest(test)
	}
	
	return reports
}

// RunTest runs a single security test
func (tf *TestFramework) RunTest(test TestCase) TestReport {
	report := TestReport{
		TestCase: test,
	}
	
	start := time.Now()
	result, err := tf.executor.Execute(context.Background(), test.Language, test.Code)
	duration := time.Since(start)
	
	report.Duration = duration
	
	if err != nil {
		report.Error = err
		report.Passed = test.ExpectedResult.ShouldBeContained
		return report
	}
	
	report.ActualResult = result
	report.ExitCode = result.ExitCode
	
	// Validate the test result
	report.Passed = tf.validateTestResult(test, result)
	
	return report
}

// validateTestResult validates if a test passed based on expected results
func (tf *TestFramework) validateTestResult(test TestCase, result *sandbox.ExecutionResult) bool {
	expected := test.ExpectedResult
	
	// If we expect containment
	if expected.ShouldBeContained {
		// Check for timeout (exit code -1)
		if result.ExitCode == -1 {
			return true
		}
		
		// Check for error exit code
		if result.ExitCode != 0 {
			return true
		}
		
		// Check for expected output indicating containment
		if expected.ExpectedOutput != "" && result.Stdout != "" {
			// Normalize line endings for comparison
			stdout := result.Stdout
			if len(stdout) > 0 && stdout[len(stdout)-1] == '\n' {
				stdout = stdout[:len(stdout)-1]
			}
			
			// Check if the output contains the expected containment message
			if strings.Contains(stdout, expected.ExpectedOutput) || stdout == expected.ExpectedOutput {
				return true
			}
		}
		
		// If we expect an error output
		if result.Stderr != "" {
			return true
		}
		
		return false
	}
	
	// For non-contained tests, check for successful execution
	if result.ExitCode == expected.ExpectedExitCode {
		if expected.ExpectedOutput != "" {
			// Check if the output matches the expected output
			if strings.Contains(result.Stdout, expected.ExpectedOutput) || result.Stdout == expected.ExpectedOutput {
				return true
			}
			return false
		}
		return true
	}
	
	return false
}

// GenerateReport generates a summary report of test results
func (tf *TestFramework) GenerateReport(reports []TestReport) string {
	passed := 0
	failed := 0
	
	for _, report := range reports {
		if report.Passed {
			passed++
		} else {
			failed++
		}
	}
	
	result := fmt.Sprintf("Security Testing Report\n")
	result += fmt.Sprintf("======================\n\n")
	result += fmt.Sprintf("Total tests: %d\n", len(reports))
	result += fmt.Sprintf("Passed: %d\n", passed)
	result += fmt.Sprintf("Failed: %d\n", failed)
	result += fmt.Sprintf("Success rate: %.2f%%\n\n", float64(passed)/float64(len(reports))*100)
	
	result += fmt.Sprintf("Test Results:\n")
	result += fmt.Sprintf("-------------\n")
	
	for _, report := range reports {
		status := "❌ FAIL"
		if report.Passed {
			status = "✅ PASS"
		}
		
		reportStr := fmt.Sprintf("%s %s (%s)\n", status, report.TestCase.Name, report.TestCase.Category)
		reportStr += fmt.Sprintf("  Duration: %v\n", report.Duration)
		
		if report.Error != nil {
			reportStr += fmt.Sprintf("  Error: %v\n", report.Error)
		} else {
			reportStr += fmt.Sprintf("  Exit code: %d\n", report.ExitCode)
			if report.ActualResult.Stdout != "" {
				reportStr += fmt.Sprintf("  Stdout: %s\n", report.ActualResult.Stdout)
			}
			if report.ActualResult.Stderr != "" {
				reportStr += fmt.Sprintf("  Stderr: %s\n", report.ActualResult.Stderr)
			}
		}
		
		result += reportStr + "\n"
	}
	
	return result
}