Master Build Prompt — AI Sandbox Executor in Go
System Role

You are a senior Go systems engineer and CLI developer. You will design and implement a fully functional AI Sandbox Executor in Go that can be easily integrated into other CLI tools. The project must follow professional Go architecture, be modular, secure, and well-documented.

Project Overview

We are building a CLI tool (tentative name: ForgeAI or Runix) that executes AI-generated code in a secure sandboxed environment. The goal is to provide a powerful, extensible execution engine that other AI CLI tools can integrate via a library or binary.

Core Features
1. Secure Sandboxed Code Execution

Run arbitrary code in multiple languages (start with Python, Go, JS).

Isolate processes to prevent host compromise.

Limit CPU, memory, and execution time.

Execute in ephemeral temp directories and auto-delete after completion.

2. Multi-Language Support

Plugin system for interpreters/runtimes.

Default support: Python, Go, JavaScript.

Simple API to add new language executors.

3. CLI Interface

Beautiful UX like Crush CLI (colorized output, clear prompts).

Core commands:

run <language> <code> → Execute code.

exec <file> → Run file in sandbox.

lang list → Show supported languages.

lang add/remove → Manage language plugins.

config → Adjust security limits.

4. Integration-Friendly

Export as Go package (import "forgeai/sandbox").

Provide REST API mode (optional) for remote execution.

JSON-formatted outputs for machine integration.

Technical Requirements
Language

Go 1.21+

Project Structure (Crush CLI Style)
