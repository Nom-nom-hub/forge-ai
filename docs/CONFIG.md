# Configuration Guide

## Overview

ForgeAI can be configured through command-line flags, environment variables, and configuration files. This guide explains all configuration options.

## Configuration Methods

### 1. Command-Line Flags
```bash
forgeai --timeout=60s --memory-limit=256 run python "print('Hello, World!')"
```

### 2. Environment Variables
```bash
export FORGEAI_TIMEOUT=60s
export FORGEAI_MEMORY_LIMIT=256
forgeai run python "print('Hello, World!')"
```

### 3. Configuration File
Create `.forgeai.yaml` in your home directory:
```yaml
timeout: 60s
memory_limit: 256
plugin_dir: ./plugins
container: false
debug: false
```

## Configuration Precedence

Configuration values are applied in the following order (later values override earlier ones):

1. Default values
2. Configuration file
3. Environment variables
4. Command-line flags

## Core Configuration Options

### Timeout
Maximum execution time for code.

**Flag:** `--timeout`
**Env Var:** `FORGEAI_TIMEOUT`
**Config:** `timeout`
**Default:** `30s`
**Format:** Duration (e.g., `10s`, `1m`, `5m30s`)

### Memory Limit
Maximum memory usage in MB.

**Flag:** `--memory-limit`
**Env Var:** `FORGEAI_MEMORY_LIMIT`
**Config:** `memory_limit`
**Default:** `128`
**Range:** 1-1024 MB

### CPU Shares
CPU shares for resource allocation (Linux only).

**Flag:** `--cpu-shares`
**Env Var:** `FORGEAI_CPU_SHARES`
**Config:** `cpu_shares`
**Default:** `100`
**Range:** 2-1024

### Network Access
Allow network connections during execution.

**Flag:** `--network-access`
**Env Var:** `FORGEAI_NETWORK_ACCESS`
**Config:** `network_access`
**Default:** `false`

### Debug Mode
Enable debug output for troubleshooting.

**Flag:** `--debug`
**Env Var:** `FORGEAI_DEBUG`
**Config:** `debug`
**Default:** `false`

## Execution Configuration

### Container Mode
Use containerized execution for stronger isolation.

**Flag:** `--container`
**Env Var:** `FORGEAI_CONTAINER`
**Config:** `container`
**Default:** `false`

### Plugin Directory
Directory containing language plugins.

**Flag:** `--plugin-dir`
**Env Var:** `FORGEAI_PLUGIN_DIR`
**Config:** `plugin_dir`
**Default:** `./plugins`

## API Configuration

### API Host
Host address for the API server.

**Env Var:** `FORGEAI_API_HOST`
**Config:** `api.host`
**Default:** `0.0.0.0`

### API Port
Port for the API server.

**Env Var:** `FORGEAI_API_PORT`
**Config:** `api.port`
**Default:** `8080`

### TLS Enabled
Enable TLS for the API server.

**Env Var:** `FORGEAI_API_TLS_ENABLED`
**Config:** `api.tls.enabled`
**Default:** `false`

### TLS Certificate File
TLS certificate file path.

**Env Var:** `FORGEAI_API_TLS_CERT_FILE`
**Config:** `api.tls.cert_file`
**Default:** (empty)

### TLS Key File
TLS private key file path.

**Env Var:** `FORGEAI_API_TLS_KEY_FILE`
**Config:** `api.tls.key_file`
**Default:** (empty)

## Security Configuration

### Read-Only Root
Make the root filesystem read-only (container mode).

**Flag:** `--read-only-root`
**Env Var:** `FORGEAI_READ_ONLY_ROOT`
**Config:** `security.read_only_root`
**Default:** `true`

### Seccomp Profile
Seccomp profile for system call filtering (Linux only).

**Env Var:** `FORGEAI_SECCOMP_PROFILE`
**Config:** `security.seccomp_profile`
**Default:** (empty)

### AppArmor Profile
AppArmor profile for access control (Linux only).

**Env Var:** `FORGEAI_APPARMOR_PROFILE`
**Config:** `security.apparmor_profile`
**Default:** (empty)

## Resource Limits

### Default Values
```yaml
timeout: 30s
memory_limit: 128
cpu_shares: 100
network_access: false
```

### Maximum Values
```yaml
timeout: 300s
memory_limit: 1024
cpu_shares: 1024
```

## Example Configuration Files

### Basic Configuration
```yaml
# .forgeai.yaml
timeout: 60s
memory_limit: 256
debug: false
```

### Advanced Configuration
```yaml
# .forgeai.yaml
timeout: 30s
memory_limit: 128
cpu_shares: 100
network_access: false
container: false
plugin_dir: ./plugins
debug: false

api:
  host: 0.0.0.0
  port: 8080
  tls:
    enabled: false
    cert_file: ""
    key_file: ""

security:
  read_only_root: true
  seccomp_profile: ""
  apparmor_profile: ""
```

## Environment Variable Examples

### Linux/macOS
```bash
export FORGEAI_TIMEOUT=60s
export FORGEAI_MEMORY_LIMIT=256
export FORGEAI_CONTAINER=false
export FORGEAI_PLUGIN_DIR=./plugins
export FORGEAI_DEBUG=false
```

### Windows (PowerShell)
```powershell
$env:FORGEAI_TIMEOUT="60s"
$env:FORGEAI_MEMORY_LIMIT="256"
$env:FORGEAI_CONTAINER="false"
$env:FORGEAI_PLUGIN_DIR="./plugins"
$env:FORGEAI_DEBUG="false"
```

### Windows (Command Prompt)
```cmd
set FORGEAI_TIMEOUT=60s
set FORGEAI_MEMORY_LIMIT=256
set FORGEAI_CONTAINER=false
set FORGEAI_PLUGIN_DIR=./plugins
set FORGEAI_DEBUG=false
```

## Configuration Validation

ForgeAI validates configuration values at startup:

- **Timeout**: Must be a valid duration
- **Memory Limit**: Must be between 1 and 1024 MB
- **CPU Shares**: Must be between 2 and 1024
- **File Paths**: Must be valid paths that exist

Invalid configuration values will cause ForgeAI to exit with an error message.

## Troubleshooting

### Common Issues

#### Invalid Duration Format
```
Error: invalid timeout value
```
**Solution**: Use valid duration format (e.g., `30s`, `1m`, `5m30s`)

#### Memory Limit Too High
```
Error: memory limit exceeds maximum
```
**Solution**: Reduce memory limit to 1024 MB or less

#### Plugin Directory Not Found
```
Error: plugin directory not found
```
**Solution**: Ensure plugin directory exists and is accessible

#### Configuration File Parse Error
```
Error: failed to parse config file
```
**Solution**: Check YAML syntax in configuration file

### Debugging Configuration
Enable debug mode to see configuration values:
```bash
forgeai --debug run python "print('Hello, World!')"
```