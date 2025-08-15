package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"forgeai/pkg/api"
)

func main() {
	// Create a context that listens for interrupt signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal, shutting down...")
		cancel()
	}()

	// Start the API server
	server := api.NewServer(&api.Config{
		Host: "0.0.0.0",
		Port: 8080,
	})

	fmt.Printf("Starting ForgeAI API server on %s:%d\n", server.Config().Host, server.Config().Port)
	
	// Start the server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Start(ctx)
	}()

	// Wait for either the server to exit or context to be cancelled
	select {
	case err := <-errChan:
		if err != nil {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		// Graceful shutdown
		fmt.Println("Shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()
		
		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("Error during shutdown: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Server shutdown complete")
	}
}