package main

import (
	"fmt"
	"os"

	"forgeai/pkg/registry"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		listPlugins()
	case "install":
		if len(os.Args) < 3 {
			fmt.Println("Usage: forgeai-plugin install <plugin-name>")
			os.Exit(1)
		}
		installPlugin(os.Args[2])
	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Usage: forgeai-plugin remove <plugin-name>")
			os.Exit(1)
		}
		removePlugin(os.Args[2])
	case "update":
		if len(os.Args) < 3 {
			fmt.Println("Usage: forgeai-plugin update <plugin-name>")
			os.Exit(1)
		}
		updatePlugin(os.Args[2])
	case "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("ForgeAI Plugin Manager")
	fmt.Println("======================")
	fmt.Println("Usage:")
	fmt.Println("  forgeai-plugin list              List installed plugins")
	fmt.Println("  forgeai-plugin install <name>    Install a plugin")
	fmt.Println("  forgeai-plugin remove <name>     Remove a plugin")
	fmt.Println("  forgeai-plugin update <name>     Update a plugin")
	fmt.Println("  forgeai-plugin help              Show this help")
}

func listPlugins() {
	// For now, we'll just list the plugins in the local directory
	// In a real implementation, we would use the PluginManager
	pluginDir := "./plugins"
	
	manager := registry.NewPluginManager(pluginDir, "http://localhost:8080")
	
	plugins, err := manager.ListInstalledPlugins()
	if err != nil {
		fmt.Printf("Error listing plugins: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Installed Plugins:")
	for _, plugin := range plugins {
		fmt.Printf("  - %s\n", plugin)
	}
}

func installPlugin(name string) {
	pluginDir := "./plugins"
	
	manager := registry.NewPluginManager(pluginDir, "http://localhost:8080")
	
	fmt.Printf("Installing plugin: %s\n", name)
	if err := manager.InstallPlugin(name, "latest"); err != nil {
		fmt.Printf("Error installing plugin: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Plugin installed successfully!")
}

func removePlugin(name string) {
	pluginDir := "./plugins"
	
	manager := registry.NewPluginManager(pluginDir, "http://localhost:8080")
	
	fmt.Printf("Removing plugin: %s\n", name)
	if err := manager.RemovePlugin(name); err != nil {
		fmt.Printf("Error removing plugin: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Plugin removed successfully!")
}

func updatePlugin(name string) {
	pluginDir := "./plugins"
	
	manager := registry.NewPluginManager(pluginDir, "http://localhost:8080")
	
	fmt.Printf("Updating plugin: %s\n", name)
	if err := manager.UpdatePlugin(name); err != nil {
		fmt.Printf("Error updating plugin: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Plugin updated successfully!")
}