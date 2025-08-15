package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"forgeai/pkg/container"
	"forgeai/pkg/executor"
	"forgeai/pkg/plugin"
	"forgeai/pkg/sandbox"
)

var (
	jsonOutput   bool
	containerized bool
	pluginDir    string
	timeout      time.Duration
	memoryLimit  int
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

		// Get the appropriate executor
		exec, err := getExecutor()
		if err != nil {
			return fmt.Errorf("failed to get executor: %w", err)
		}

		// Execute code
		result, err := exec.Execute(context.Background(), language, code)
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

		// Get the appropriate executor
		exec, err := getExecutor()
		if err != nil {
			return fmt.Errorf("failed to get executor: %w", err)
		}

		// Execute file
		result, err := exec.ExecuteFile(context.Background(), file)
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
		// Get the appropriate executor
		exec, err := getExecutor()
		if err != nil {
			return fmt.Errorf("failed to get executor: %w", err)
		}

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
	rootCmd.PersistentFlags().StringVar(&pluginDir, "plugin-dir", "", "Directory to load plugins from")
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

// getExecutor returns the appropriate executor based on the flags
func getExecutor() (sandbox.Executor, error) {
	if pluginDir != "" {
		// Use plugin manager
		manager := plugin.NewManager()
		if err := manager.LoadPluginsFromDir(pluginDir); err != nil {
			return nil, fmt.Errorf("failed to load plugins: %w", err)
		}
		
		// If containerized flag is also set, we need to handle this case
		// For now, we'll prioritize plugins over containerized execution
		if containerized {
			fmt.Println("Warning: Both --plugin-dir and --container flags are set. Using plugins.")
		}
		
		// Return a composite executor that can handle both plugins and default executors
		return &CompositeExecutor{
			PluginManager: manager,
			LocalExecutor: executor.NewLocalExecutor(),
			DockerExecutor: container.NewDockerExecutor(),
			UseContainer: containerized,
		}, nil
	} else if containerized {
		// Use containerized executor
		dockerExec := container.NewDockerExecutor()
		dockerExec.Timeout = timeout
		dockerExec.MemoryLimit = memoryLimit
		return dockerExec, nil
	} else {
		// Use local executor
		localExec := executor.NewLocalExecutor()
		localExec.Timeout = timeout
		localExec.MemoryLimit = memoryLimit
		return localExec, nil
	}
}

// CompositeExecutor combines plugin, local, and container executors
type CompositeExecutor struct {
	PluginManager  *plugin.Manager
	LocalExecutor  *executor.LocalExecutor
	DockerExecutor *container.DockerExecutor
	UseContainer   bool
}

func (c *CompositeExecutor) Execute(ctx context.Context, language, code string) (*sandbox.ExecutionResult, error) {
	// Check if we have a plugin for this language
	if executor, ok := c.PluginManager.GetExecutor(language); ok {
		return executor.Execute(ctx, language, code)
	}
	
	// Use the appropriate executor based on the UseContainer flag
	if c.UseContainer {
		c.DockerExecutor.Timeout = c.LocalExecutor.Timeout
		c.DockerExecutor.MemoryLimit = c.LocalExecutor.MemoryLimit
		return c.DockerExecutor.Execute(ctx, language, code)
	}
	
	return c.LocalExecutor.Execute(ctx, language, code)
}

func (c *CompositeExecutor) ExecuteFile(ctx context.Context, filePath string) (*sandbox.ExecutionResult, error) {
	// Get the language from the file extension
	language := getLanguageFromFile(filePath)
	
	// Check if we have a plugin for this language
	if executor, ok := c.PluginManager.GetExecutor(language); ok {
		return executor.ExecuteFile(ctx, filePath)
	}
	
	// Use the appropriate executor based on the UseContainer flag
	if c.UseContainer {
		c.DockerExecutor.Timeout = c.LocalExecutor.Timeout
		c.DockerExecutor.MemoryLimit = c.LocalExecutor.MemoryLimit
		return c.DockerExecutor.ExecuteFile(ctx, filePath)
	}
	
	return c.LocalExecutor.ExecuteFile(ctx, filePath)
}

func (c *CompositeExecutor) SupportedLanguages() []string {
	// Get languages from plugins
	pluginLanguages := c.PluginManager.SupportedLanguages()
	
	// Get languages from the appropriate executor
	var defaultLanguages []string
	if c.UseContainer {
		defaultLanguages = c.DockerExecutor.SupportedLanguages()
	} else {
		defaultLanguages = c.LocalExecutor.SupportedLanguages()
	}
	
	// Combine the lists
	languages := make([]string, 0, len(pluginLanguages)+len(defaultLanguages))
	languages = append(languages, pluginLanguages...)
	languages = append(languages, defaultLanguages...)
	
	return languages
}

// getLanguageFromFile determines the language from the file extension
func getLanguageFromFile(filePath string) string {
	switch {
	case filepath.Ext(filePath) == ".py":
		return "python"
	case filepath.Ext(filePath) == ".go":
		return "go"
	case filepath.Ext(filePath) == ".js":
		return "javascript"
	case filepath.Ext(filePath) == ".rs":
		return "rust"
	default:
		return "unknown"
	}
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