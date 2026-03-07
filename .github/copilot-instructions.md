# Copilot Agent Instructions for doxctl

## Repository Overview

**doxctl** is a Go-based diagnostic CLI tool for troubleshooting VPN and DNS connectivity issues on laptops. The tool performs comprehensive checks on VPN connections, DNS resolvers, network performance, and service health across multiple data centers.

- **Language**: Go 1.25.0+
- **Type**: CLI application using Cobra framework
- **Size**: ~50-60 files, well-structured monorepo
- **Primary Platforms**: macOS (fully supported), Linux (partial)
- **Architectures**: amd64, arm64
- **Distribution**: Homebrew, Docker (multi-arch), GitHub Releases

## Build & Test Commands

### Bootstrap/Setup
```bash
# Install dependencies
go mod download

# Verify Go version (must be 1.25.0+)
go version
```

### Build
```bash
# Standard build (always works)
go build -o doxctl .

# Using Makefile
make build

# Using goreleaser (for release builds)
goreleaser build --clean
```

### Test
**CRITICAL**: Always run tests before committing:
```bash
# Standard test command (with coverage)
go test -coverprofile=coverage.txt -covermode=atomic ./...

# Or use Makefile (preferred)
make test

# Tests with race detection
make test-race
```

**Important**: Tests may skip on non-macOS systems since some checks are platform-specific (use `runtime.GOOS` checks). This is expected behavior.

### Lint
**ALWAYS run linting before submitting changes**:
```bash
# Primary linting command
golangci-lint run --timeout=5m

# Or via Makefile
make lint

# Auto-fix some issues
golangci-lint run --fix
```

If `golangci-lint` is not installed:
```bash
# macOS
brew install golangci-lint

# Other platforms
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Run
```bash
# After build
./doxctl --help

# Specific commands require configuration file
./doxctl config init    # Create example config
./doxctl dns -a        # Run all DNS checks
./doxctl vpn -a        # Run all VPN checks
./doxctl svrs -s       # Check server reachability
```

### Other Important Commands
```bash
# Format code (MUST do before commit)
go fmt ./...

# Vet code
go vet ./...

# Clean build artifacts
make clean

# Security scan
make security-scan  # Requires govulncheck

# Full pre-commit check
make all  # Runs: clean fmt vet lint test build
```

## Configuration Requirements

**Critical**: Most commands require a `.doxctl.yaml` configuration file.

**Location**: `$HOME/.doxctl.yaml` or current directory

**Generate example**:
```bash
doxctl config init
```

**Validate configuration**:
```bash
doxctl config validate
```

**Key configuration fields**:
- `minVpnRoutes`: Minimum VPN routes (default: 5)
- `domainName`: Primary domain for DNS checks
- `sites`: List of datacenter sites (e.g., lab1, rdu1, dfw1)
- `wellKnownSvcs`: Services and their endpoints for health checks
- `pingTimeout`: Timeout for ping operations (default: 250ms)
- `failThreshold`: Failure threshold (default: 5)

## Project Structure

### Root Files
- `main.go`: Entry point, calls `cmd.Execute()`
- `Makefile`: Build automation (use this for common tasks)
- `go.mod`, `go.sum`: Go dependencies
- `.golangci.yml`: Linter configuration (strict rules)
- `README.md` → redirects to `docs/README.md`

### Key Directories

#### `/cmd/` - Command implementations
- `root.go`: Root command, configuration loading (Viper), global flags
- `dns.go`: DNS resolver checks (resolverChk, pingChk, digChk)
- `vpn.go`: VPN connectivity checks (ifReachableChk, vpnRoutesChk, vpnStatusChk)  
- `svrs.go`: Server reachability checks
- `svcs.go`: Service health checks (HTTP/HTTPS endpoints)
- `net.go`: Network performance & SLO validation
- `config.go`: Configuration management commands (init, validate, show)
- `version.go`: Version information
- `*_test.go`: Unit tests (platform-specific tests use build tags)
- `*_inject_test.go`: Dependency injection tests
- `interfaces.go`: Abstractions for testing (CommandExecutor, DNSResolver, Pinger)

#### `/internal/` - Internal packages
- `output/`: Output formatting (table, JSON, YAML)
- `cmdhelp/`: Command helper utilities (pipeline execution)

#### `/docs/` - Documentation
- `README.md`: Primary documentation (comprehensive)
- `CONTRIBUTING.md`: Development guidelines
- `ARCHITECTURE.md`: System design documentation
- `SECURITY.md`: Security policy
- `CHANGELOG.md`: Version history

#### `/docker/` - Container support
- `Dockerfile`: Multi-stage Alpine-based image
- `README.md`: Docker usage instructions

#### `/scripts/` - Utility scripts
- `version-up.sh`: Version bumping script

## Architecture & Code Patterns

### Command Structure
All commands follow this pattern:
```go
var someCmd = &cobra.Command{
    Use:   "command",
    Short: "Brief description",
    Long:  "Detailed description",
    PreRun: func(cmd *cobra.Command, args []string) {
        // Load and validate configuration
        // Exit on config errors
    },
    Run: executeFunction,
}
```

### Dependency Injection for Testing
Functions that interact with system (exec, network, file I/O) have `*WithDeps` variants:
```go
func someCheck() {
    someCheckWithDeps(NewCommandExecutor())
}

func someCheckWithDeps(executor CommandExecutor) {
    // Testable implementation using interface
}
```

### Output Format Support
All diagnostic commands support three output formats via `--output` flag:
- `table` (default): Human-readable tables
- `json`: Machine-readable JSON
- `yaml`: Machine-readable YAML

**Pattern**:
```go
if outputFormat != "table" {
    result := SomeCheckResult{...}
    output.Print(outputFormat, result)
    return
}
// Table output code
```

### Configuration Loading
Configuration is loaded via Viper in each command's `PreRun`:
- Searches: current directory, `$HOME/.doxctl.yaml`
- Environment variables: `DOXCTL_*` prefix
- Validates using `config.Validate()`
- Sets defaults via `config.setDefaults()`

### Platform-Specific Code
- macOS: Uses `scutil` for network/DNS info
- Linux: Uses `/etc/resolv.conf` and `ip route`
- Build tags: `//go:build darwin` or `//go:build linux`

## Common Pitfalls & Gotchas

### 1. Configuration File Required
**Error**: "Configuration file not found"
**Fix**: Run `doxctl config init` to create example configuration, then edit it.

### 2. macOS vs Linux Differences
Some commands are macOS-specific (use `scutil`). Tests skip on other platforms:
```go
if runtime.GOOS != "darwin" {
    t.Skip("Skipping macOS-specific test")
}
```

### 3. Test Coverage Requirements
Project aims for 70%+ coverage. Add tests for new features. Use injection pattern for testability.

### 4. Linting Failures
Common issues:
- Unchecked errors: Always check `err != nil`
- Unused variables: Remove or prefix with `_`
- Code formatting: Run `go fmt ./...`
- Complexity: golangci-lint enforces cyclomatic complexity limits

### 5. Docker Container Networking
macOS VPN settings don't transfer to Docker. Use the `doxctl-container` wrapper script that extracts host VPN config.

### 6. Timeout Settings
Network operations have configurable timeouts:
- `pingTimeout`: 250ms (default)
- `dnsLookupTimeout`: 100ms (default)
- Service health checks: 5s (default)

### 7. Build Tag Requirements
When adding platform-specific code, use proper build tags:
```go
//go:build darwin
// +build darwin
```

## CI/CD & Workflows

### GitHub Actions Workflows
- `build-release.yml`: Build and release on tags
- `codeql.yml`: Security scanning
- `lint.yml`: Linting checks
- `codecoverage.yml`: Coverage reporting
- `security-scan.yml`: Vulnerability scanning

### Pre-commit Checklist
1. Run `go fmt ./...`
2. Run `go vet ./...`
3. Run `make lint`
4. Run `make test` (ensure all tests pass)
5. Update documentation if adding features
6. Commit using conventional commits format

### Commit Message Format
```
<type>: <subject>

<body>

<footer>
```

**Types**: feat, fix, docs, style, refactor, perf, test, chore

**Example**:
```
feat: add custom DNS timeout flag

Add --dns-timeout flag to allow configuring DNS query timeout.
Useful for slow network environments.

Closes #123
```

## Testing Strategy

### Test Organization
- `*_test.go`: Standard unit tests
- `*_inject_test.go`: Dependency injection tests
- `integration_test.go`: Integration tests
- `coverage_test.go`: Additional coverage tests
- Platform-specific tests use build tags

### Mock Patterns
```go
type mockCommandExecutor struct {
    commands map[string][]byte
}

func (m *mockCommandExecutor) Execute(name string, args ...string) ([]byte, error) {
    // Return mocked data
}
```

### Running Specific Tests
```bash
# Run tests for specific package
go test ./cmd/

# Run specific test
go test -run TestDnsCmd_Initialization ./cmd/

# Verbose output
go test -v ./...
```

## Key Dependencies

### Core Libraries
- `github.com/spf13/cobra`: CLI framework
- `github.com/spf13/viper`: Configuration management
- `github.com/go-ping/ping`: ICMP ping
- `github.com/miekg/dns`: DNS operations
- `github.com/jedib0t/go-pretty/v6`: Table formatting
- `github.com/gookit/color`: Terminal colors

### Build Tools
- `goreleaser`: Release automation
- `golangci-lint`: Code quality enforcement

## Validation Checklist

Before making changes, verify:

1. **Go version**: `go version` shows 1.25.0+
2. **Dependencies**: `go mod download` successful
3. **Build**: `make build` or `go build -o doxctl .` successful
4. **Tests**: `make test` passes (some skips on non-macOS OK)
5. **Lint**: `make lint` passes with no errors
6. **Format**: `go fmt ./...` shows no changes
7. **Config**: Create test config with `doxctl config init`

After making changes:

1. **Tests updated**: Added tests for new functionality
2. **Documentation**: Updated README.md or relevant docs
3. **Errors handled**: All errors checked and handled
4. **Platform compatibility**: Build tags used if platform-specific
5. **Output formats**: JSON/YAML support added if needed
6. **Configuration**: Updated example config if new fields added

## Trust These Instructions

These instructions are generated from comprehensive analysis of the actual repository state. When information conflicts with assumptions:

1. **Trust the build commands** documented here - they are verified to work
2. **Trust the test commands** - the test suite is comprehensive
3. **Trust the linting setup** - `.golangci.yml` is properly configured
4. **Trust the project structure** - paths and organization are accurate

Only perform additional exploration if:
- Instructions are incomplete for your specific task
- Repository has been significantly updated since generation
- Encountering errors not covered by common pitfalls section

---

**Last Updated**: Based on repository state as of March 2026
**Maintainer**: Sam Mingolelli (github@lamolabs.org)
**License**: MIT
