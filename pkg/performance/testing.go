package performance

import (
	"context"
	"fmt"
	"sync"
	"time"

	"forgeai/pkg/security"
)

// PerformanceMetrics holds performance metrics for execution
type PerformanceMetrics struct {
	ExecutorName     string
	TestCaseName     string
	TotalExecutions  int
	SuccessfulExecutions int
	FailedExecutions int
	TotalTime        time.Duration
	AverageTime      time.Duration
	MinTime          time.Duration
	MaxTime          time.Duration
	SuccessRate      float64
}

// PerformanceTest holds a performance test case
type PerformanceTest struct {
	Name string
	Code string
	Language string
}

// PerformanceFramework manages performance testing
type PerformanceFramework struct {
	executors map[string]security.Executor
	tests     []PerformanceTest
}

// NewPerformanceFramework creates a new performance testing framework
func NewPerformanceFramework() *PerformanceFramework {
	executors := map[string]security.Executor{
		"Local Executor":      createLocalExecutor(),
		"Secure Executor":     security.NewSecureExecutor(),
		"Containerized Executor": security.NewContainerizedExecutor(),
	}
	
	tests := []PerformanceTest{
		{
			Name:     "Simple Print",
			Code:     "print('Hello, World!')",
			Language: "python",
		},
		{
			Name:     "Loop Calculation",
			Code:     "total = 0\nfor i in range(1000):\n    total += i\nprint(total)",
			Language: "python",
		},
		{
			Name:     "Function Definition",
			Code:     "def fibonacci(n):\n    if n <= 1:\n        return n\n    return fibonacci(n-1) + fibonacci(n-2)\n\nprint(fibonacci(10))",
			Language: "python",
		},
		{
			Name:     "File I/O",
			Code:     "with open('test.txt', 'w') as f:\n    f.write('Hello, World!')\nwith open('test.txt', 'r') as f:\n    print(f.read())",
			Language: "python",
		},
	}
	
	return &PerformanceFramework{
		executors: executors,
		tests:     tests,
	}
}

// RunPerformanceTests runs all performance tests
func (pf *PerformanceFramework) RunPerformanceTests() []PerformanceMetrics {
	var metrics []PerformanceMetrics
	
	for executorName, executor := range pf.executors {
		for _, test := range pf.tests {
			testMetrics := pf.runPerformanceTest(executorName, executor, test)
			metrics = append(metrics, testMetrics)
		}
	}
	
	return metrics
}

// runPerformanceTest runs a single performance test
func (pf *PerformanceFramework) runPerformanceTest(executorName string, executor security.Executor, test PerformanceTest) PerformanceMetrics {
	metrics := PerformanceMetrics{
		ExecutorName: executorName,
		TestCaseName: test.Name,
		MinTime:      time.Hour, // Initialize to a large value
		MaxTime:      0,
	}
	
	// Warm up
	for i := 0; i < 3; i++ {
		_, err := executor.Execute(context.Background(), test.Language, test.Code)
		if err != nil {
			metrics.FailedExecutions++
		} else {
			metrics.SuccessfulExecutions++
		}
	}
	
	// Reset counters for actual test
	metrics.SuccessfulExecutions = 0
	metrics.FailedExecutions = 0
	metrics.TotalTime = 0
	
	// Run timed test
	const numTests = 10
	start := time.Now()
	
	var times []time.Duration
	var mu sync.Mutex
	
	// Run tests concurrently
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit concurrent executions to 5
	
	for i := 0; i < numTests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			execStart := time.Now()
			_, err := executor.Execute(context.Background(), test.Language, test.Code)
			execDuration := time.Since(execStart)
			
			mu.Lock()
			defer mu.Unlock()
			
			if err != nil {
				metrics.FailedExecutions++
			} else {
				metrics.SuccessfulExecutions++
			}
			
			metrics.TotalTime += execDuration
			times = append(times, execDuration)
			
			// Update min/max times
			if execDuration < metrics.MinTime {
				metrics.MinTime = execDuration
			}
			if execDuration > metrics.MaxTime {
				metrics.MaxTime = execDuration
			}
		}()
	}
	
	wg.Wait()
	
	metrics.TotalExecutions = numTests
	metrics.AverageTime = metrics.TotalTime / time.Duration(numTests)
	metrics.SuccessRate = float64(metrics.SuccessfulExecutions) / float64(numTests) * 100
	
	totalDuration := time.Since(start)
	fmt.Printf("Test %s with %s completed in %v\n", test.Name, executorName, totalDuration)
	
	return metrics
}

// GenerateReport generates a performance report
func (pf *PerformanceFramework) GenerateReport(metrics []PerformanceMetrics) string {
	report := fmt.Sprintf("Performance Testing Report\n")
	report += fmt.Sprintf("========================\n\n")
	
	// Group metrics by executor
	executorMetrics := make(map[string][]PerformanceMetrics)
	for _, metric := range metrics {
		executorMetrics[metric.ExecutorName] = append(executorMetrics[metric.ExecutorName], metric)
	}
	
	// Generate report for each executor
	for executorName, metrics := range executorMetrics {
		report += fmt.Sprintf("Executor: %s\n", executorName)
		report += fmt.Sprintf("--------------------\n")
		
		totalTime := time.Duration(0)
		totalExecutions := 0
		totalSuccessful := 0
		
		for _, metric := range metrics {
			report += fmt.Sprintf("  Test: %s\n", metric.TestCaseName)
			report += fmt.Sprintf("    Executions: %d\n", metric.TotalExecutions)
			report += fmt.Sprintf("    Successful: %d\n", metric.SuccessfulExecutions)
			report += fmt.Sprintf("    Failed: %d\n", metric.FailedExecutions)
			report += fmt.Sprintf("    Success Rate: %.2f%%\n", metric.SuccessRate)
			report += fmt.Sprintf("    Total Time: %v\n", metric.TotalTime)
			report += fmt.Sprintf("    Average Time: %v\n", metric.AverageTime)
			report += fmt.Sprintf("    Min Time: %v\n", metric.MinTime)
			report += fmt.Sprintf("    Max Time: %v\n", metric.MaxTime)
			report += fmt.Sprintf("\n")
			
			totalTime += metric.TotalTime
			totalExecutions += metric.TotalExecutions
			totalSuccessful += metric.SuccessfulExecutions
		}
		
		overallSuccessRate := float64(totalSuccessful) / float64(totalExecutions) * 100
		report += fmt.Sprintf("  Overall Success Rate: %.2f%%\n", overallSuccessRate)
		report += fmt.Sprintf("  Overall Total Time: %v\n", totalTime)
		report += fmt.Sprintf("\n")
	}
	
	return report
}

// createLocalExecutor creates a local executor with optimized settings
func createLocalExecutor() security.Executor {
	// For this example, we'll use the existing local executor
	// In a real implementation, we might create a specialized version
	return security.NewSecureExecutor()
}