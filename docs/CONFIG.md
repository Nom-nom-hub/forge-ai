# Configuration

This document describes the configuration options for ForgeAI.

## Configuration File

ForgeAI can be configured using a configuration file in YAML format.

### Example Configuration

```yaml
timeout: 30s
memory_limit: 128
cpu_shares: 100
network_access: false
allowed_dirs: []
```

## Command Line Options

### Global Options
- `--json` - Output results in JSON format
- `--timeout` - Execution timeout (default: 30s)
- `--memory-limit` - Memory limit in MB (default: 128)

### Command Specific Options
- `run` - Execute code directly
- `exec` - Execute a file
- `lang list` - List supported languages
- `config` - Adjust security limits

## Configuration Options

### timeout
- Type: duration
- Default: 30s
- Description: Maximum execution time for code

### memory_limit
- Type: integer
- Default: 128
- Description: Memory limit in MB

### cpu_shares
- Type: integer
- Default: 100
- Description: CPU shares (Linux only)

### network_access
- Type: boolean
- Default: false
- Description: Allow network access

### allowed_dirs
- Type: array of strings
- Default: []
- Description: Allowed directories for file access