package main

import (
	"fmt"
	"os"

	"forgeai/pkg/performance"
)

func main() {
	fmt.Println("ForgeAI Performance Testing Framework")
	fmt.Println("====================================")
	
	// Create performance framework
	framework := performance.NewPerformanceFramework()
	
	// Run performance tests
	fmt.Println("Running performance tests...")
	metrics := framework.RunPerformanceTests()
	
	// Generate and print report
	report := framework.GenerateReport(metrics)
	fmt.Println(report)
	
	// Exit successfully
	os.Exit(0)
}