package main

import (
	"fmt"
	"os"

	"forgeai/pkg/security"
)

func main() {
	fmt.Println("ForgeAI Security Testing Framework")
	fmt.Println("==================================")
	
	// Create test framework
	framework := security.NewTestFramework()
	
	// Run tests
	fmt.Println("Running security tests...")
	reports := framework.RunTests()
	
	// Generate and print report
	report := framework.GenerateReport(reports)
	fmt.Println(report)
	
	// Count passed and failed tests
	passed := 0
	failed := 0
	for _, r := range reports {
		if r.Passed {
			passed++
		} else {
			failed++
		}
	}
	
	// Exit with appropriate code
	if failed > 0 {
		fmt.Printf("Security testing completed with %d failures.\n", failed)
		os.Exit(1)
	} else {
		fmt.Println("All security tests passed!")
		os.Exit(0)
	}
}
