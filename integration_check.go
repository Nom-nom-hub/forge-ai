package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// Test 1: Help command
	fmt.Println("Test 1: Help command")
	cmd := exec.Command("forgeai.exe", "--help")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: Help command works\n")
	}

	// Test 2: Language list
	fmt.Println("\nTest 2: Language list")
	cmd = exec.Command("forgeai.exe", "lang", "list")
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: Language list works\n")
	}

	// Test 3: JSON output
	fmt.Println("\nTest 3: JSON output")
	cmd = exec.Command("forgeai.exe", "lang", "list", "--json")
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: JSON output works\n")
	}

	// Test 4: Run Python code
	fmt.Println("\nTest 4: Run Python code")
	cmd = exec.Command("forgeai.exe", "run", "python", "print('Hello, World!')")
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: Python execution works\n")
	}

	// Test 5: Run JavaScript code
	fmt.Println("\nTest 5: Run JavaScript code")
	cmd = exec.Command("forgeai.exe", "run", "javascript", "console.log('Hello, World!')")
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: JavaScript execution works\n")
	}

	// Test 6: Run Go code
	fmt.Println("\nTest 6: Run Go code")
	cmd = exec.Command("forgeai.exe", "run", "go", "package main; import \"fmt\"; func main() { fmt.Println(\"Hello, World!\") }")
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: Go execution works\n")
	}

	// Test 7: Execute Python file
	fmt.Println("\nTest 7: Execute Python file")
	cmd = exec.Command("forgeai.exe", "exec", "examples/hello_world.py")
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: Python file execution works\n")
	}

	fmt.Println("\nAll integration tests completed!")
}