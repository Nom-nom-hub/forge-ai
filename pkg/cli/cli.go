package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"forgeai/pkg/container"
	"forgeai/pkg/executor"
	"forgeai/pkg/sandbox"
)

var (
	jsonOutput  bool
	containerized bool
	timeout     time.Duration
	memoryLimit int
)

var rootCmd = &cobra.Command{
	Use:   "forgeai",
	Short: "ForgeAI is a secure sandboxed code executor",
	Long: `ForgeAI is a CLI tool that executes AI-generated code in a secure sandboxed environment.
It supports multiple languages and provides isolation to prevent host compromise.`,
}

var runCmd = &cobra.Command{
	Use:   "run [language] [code]",
	Short: "Execute code in a sandbox",
	Long:  `Execute the provided code in the specified language within a secure sandbox.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		language := args[0]
		code := args[1]

		var executor sandbox.Executor

		if containerized {
			// Use containerized executor
			dockerExec := container.NewDockerExecutor()
			dockerExec.Timeout = timeout
			dockerExec.MemoryLimit = memoryLimit
			executor = dockerExec
		} else {
			// Use local executor
			localExec := executor.NewLocalExecutor()
			localExec.Timeout = timeout
			localExec.MemoryLimit = memoryLimit
			executor = localExec
		}

		// Execute code
		result, err := executor.Execute(context.Background(), language, code)
		if err != nil {
			return fmt.Errorf("failed to execute code: %w", err)
		}

		return printResult(result)
	},
}

var execCmd = &cobra.Command{
	Use:   "exec [file]",
	Short: "Execute a file in a sandbox",
	Long:  `Execute the provided file within a secure sandbox.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		var executor sandbox.Executor

		if containerized {
			// Use containerized executor
			dockerExec := container.NewDockerExecutor()
			dockerExec.Timeout = timeout
			dockerExec.MemoryLimit = memoryLimit
			executor = dockerExec
		} else {
			// Use local executor
			localExec := executor.NewLocalExecutor()
			localExec.Timeout = timeout
			localExec.MemoryLimit = memoryLimit
			executor = localExec
		}

		// Execute file
		result, err := executor.ExecuteFile(context.Background(), file)
		if err != nil {
			return fmt.Errorf("failed to execute file: %w", err)
		}

		return printResult(result)
	},
}

var langCmd = &cobra.Command{
	Use:   "lang",
	Short: "Manage language support",
	Long:  `List, add, or remove language executors.`,
}

var langListCmd = &cobra.Command{
	Use:   "list",
	Short: "List supported languages",
	RunE: func(cmd *cobra.Command, args []string) error {
		exec := executor.NewLocalExecutor()
		languages := exec.SupportedLanguages()

		if jsonOutput {
			return json.NewEncoder(os.Stdout).Encode(languages)
		}

		fmt.Println("Supported languages:")
		for _, lang := range languages {
			fmt.Printf("  - %s\n", lang)
		}
		return nil
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Adjust security limits",
	Long:  `Configure security limits for the sandbox execution.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Configuring security limits...")
		// TODO: Implement configuration management
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	rootCmd.PersistentFlags().BoolVar(&containerized, "container", false, "Use containerized execution")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 30*time.Second, "Execution timeout")
	rootCmd.PersistentFlags().IntVar(&memoryLimit, "memory-limit", 128, "Memory limit in MB")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(execCmd)

	langCmd.AddCommand(langListCmd)
	rootCmd.AddCommand(langCmd)

	rootCmd.AddCommand(configCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func printResult(result *sandbox.ExecutionResult) error {
	if jsonOutput {
		return json.NewEncoder(os.Stdout).Encode(result)
	}

	fmt.Printf("Execution completed in %v\n", result.Duration)
	fmt.Printf("Exit code: %d\n", result.ExitCode)

	if result.Stdout != "" {
		fmt.Printf("Stdout:\n%s\n", result.Stdout)
	}

	if result.Stderr != "" {
		fmt.Printf("Stderr:\n%s\n", result.Stderr)
	}

	return nil
}