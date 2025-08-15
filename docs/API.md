# API Reference

This document provides a reference for the ForgeAI Go SDK API.

## Packages

### cli
Package cli implements the command-line interface for ForgeAI.

### sandbox
Package sandbox defines the interfaces for sandboxed code execution.

#### Types
- `ExecutionResult` - Represents the result of a sandboxed execution
- `Executor` - Interface for executing code in a sandbox

### executor
Package executor implements the code execution logic.

#### Types
- `LocalExecutor` - Basic implementation of the Executor interface

#### Functions
- `NewLocalExecutor()` - Creates a new LocalExecutor with default settings

### config
Package config manages configuration for the sandbox executor.

#### Types
- `Config` - Configuration for the sandbox executor

#### Functions
- `DefaultConfig()` - Returns a Config with default values

### output
Package output handles formatting and printing of execution results.

#### Types
- `Printer` - Handles formatting and printing of execution results

#### Functions
- `NewPrinter()` - Creates a new Printer