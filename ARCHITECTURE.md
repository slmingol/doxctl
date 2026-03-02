# Architecture Documentation

## Overview

`doxctl` is a diagnostic CLI tool designed to help end users triage connectivity problems stemming from VPN and DNS setups on their laptops. It provides comprehensive diagnostics for VPN, DNS, network routing, and service connectivity.

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        User Interface                        │
│                     (CLI - Cobra Framework)                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                    Command Layer (cmd/)                      │
│  ┌──────┐  ┌─────┐  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐│
│  │ root │  │ vpn │  │ dns  │  │ net  │  │ svrs │  │ svcs ││
│  └──────┘  └─────┘  └──────┘  └──────┘  └──────┘  └──────┘│
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                  Internal Helpers (internal/)                │
│                  - Command Pipelining                        │
│                  - Utility Functions                         │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│              External Dependencies & OS                      │
│  ┌──────────┐  ┌─────────┐  ┌───────────┐  ┌─────────────┐│
│  │ Network  │  │  DNS    │  │ Ping/ICMP │  │   System    ││
│  │ Commands │  │Libraries│  │  Library  │  │  Commands   ││
│  └──────────┘  └─────────┘  └───────────┘  └─────────────┘│
└─────────────────────────────────────────────────────────────┘
```

## Component Design

### 1. Command Layer (`cmd/`)

The command layer implements the CLI interface using the Cobra framework. Each command is responsible for a specific diagnostic area.

#### Root Command (`root.go`)
- Entry point for the CLI
- Manages configuration loading (Viper)
- Handles global flags and configuration
- Initializes sub-commands

#### VPN Command (`vpn.go`)
- Validates VPN connection status
- Checks routing table for VPN routes
- Verifies minimum required routes (`minVpnRoutes` setting)
- Platform-specific route checking (macOS, Linux)

#### DNS Command (`dns.go`)
- Tests DNS resolver configuration
- Validates DNS search paths
- Performs DNS lookups and resolution tests
- Checks resolver reachability via ping
- Validates specific domain resolution

**Sub-checks:**
- `--resolver-chk`: Validates DNS resolver configuration
- `--ping-chk`: Pings DNS resolvers to verify reachability
- `--dig-chk`: Performs DNS queries to validate resolution

#### Network Command (`net.go`)
- General network diagnostics
- Network interface validation
- Basic connectivity checks

#### Servers Command (`svrs.go`)
- Tests connectivity to well-known servers
- Supports server expansion via brace notation
- Performs ping tests to verify reachability
- Configurable via `wellKnownSvcs` in config file

#### Services Command (`svcs.go`)
- Tests connectivity to services
- Service-specific health checks
- Groups servers by service type

### 2. Internal Packages (`internal/`)

#### cmdhelp Package
Contains utility functions for command execution and manipulation.

**Key Functions:**
- `Pipeline()`: Chains multiple commands together using pipes
  - Enables complex command sequences
  - Captures stdout and stderr separately
  - Handles error propagation

### 3. Configuration System

Configuration is managed through Viper and YAML files.

#### Configuration File (`.doxctl.yaml`)

**Location**: `~/.doxctl.yaml` or current directory

**Structure:**
```yaml
# VPN Configuration
minVpnRoutes: 5              # Minimum routes required for VPN validation

# DNS Configuration
domNameChk: "bandwidth.local"     # Domain name to check
domSearchChk: "[0-1].*bandwidth"  # Search path regex pattern
domAddrChk: "[0-1].*10.5"        # Address regex pattern
domainName: "bandwidthclec.local" # Primary domain
digProbeServerA: "idm-01a"        # DNS probe server A
digProbeServerB: "idm-01b"        # DNS probe server B

# Site Configuration
sites:                       # List of sites/locations
  - lab1
  - rdu1
  - atl1

# Service Configuration
wellKnownSvcs:              # Services and their servers
  - svc: openshift
    svrs:
      - ocp-master-01{a,b,c}.{lab1,rdu1}.bandwidthclec.local

# Timeout Configuration
pingTimeout: 250            # Ping timeout in milliseconds
failThreshold: 5            # Failure threshold for checks
dnsLookupTimeout: 100       # DNS lookup timeout in milliseconds
```

**Configuration Precedence:**
1. Command-line flags (highest priority)
2. Environment variables (with `DOXCTL_` prefix)
3. Configuration file
4. Default values (lowest priority)

## Data Flow

### Typical Diagnostic Flow

```
User Command
    │
    ▼
Parse Command & Flags
    │
    ▼
Load Configuration
(Config File + Env Vars + Flags)
    │
    ▼
Initialize Diagnostic Check
    │
    ▼
Execute System Commands
or Network Operations
    │
    ▼
Collect Results
    │
    ▼
Format Output (Tables)
    │
    ▼
Display to User
(Color-coded, Tabular)
```

### Example: DNS Resolver Check

```
doxctl dns --resolver-chk
    │
    ▼
Load DNS Config
(resolvers, timeout, etc.)
    │
    ▼
Read System DNS Config
(/etc/resolv.conf or system API)
    │
    ▼
Validate Each Resolver
    ├─ Check format
    ├─ Perform DNS query
    └─ Measure response time
    │
    ▼
Generate Report Table
    │
    ▼
Display Results
(Green=Pass, Red=Fail)
```

## Key Design Decisions

### 1. CLI Framework: Cobra + Viper

**Why Cobra:**
- Industry-standard for Go CLI applications
- Excellent support for sub-commands
- Built-in help generation
- Flag parsing and validation

**Why Viper:**
- Unified configuration management
- Support for multiple config formats
- Environment variable integration
- Configuration file watching

### 2. Multi-Platform Support

**Approach:**
- Runtime platform detection (`runtime.GOOS`)
- Platform-specific command execution
- Conditional logic for OS differences

**Platforms Supported:**
- macOS (Darwin)
- Linux
- Multi-architecture (amd64, arm64)

### 3. Network Testing Strategy

**Ping Implementation:**
- Uses `github.com/go-ping/ping` library
- Supports both privileged and unprivileged ICMP
- Configurable timeouts and retry logic

**DNS Testing:**
- Uses `github.com/miekg/dns` for low-level DNS operations
- Uses `github.com/lixiangzhong/dnsutil` for high-level DNS utilities
- Direct DNS protocol implementation (not relying on system resolver)

### 4. Output Formatting

**Table-based Output:**
- Uses `github.com/jedib0t/go-pretty/v6/table`
- Consistent, readable format
- Easy to parse visually

**Color Coding:**
- Uses `github.com/gookit/color`
- Green for success
- Red for failures
- Yellow for warnings
- Improves user experience and quick identification of issues

### 5. Server Expansion

**Brace Expansion:**
- Uses `github.com/kujtimiihoxha/go-brace-expansion`
- Allows compact server list notation
- Example: `server-{a,b,c}.{site1,site2}.domain` expands to 6 servers
- Reduces configuration verbosity

## Error Handling

### Strategy

1. **Graceful Degradation**: Continue checks even if one fails
2. **Clear Error Messages**: Provide actionable error information
3. **Exit Codes**: Use standard exit codes for scripting
4. **Logging**: Errors are displayed with context

### Common Error Scenarios

- **Configuration Missing**: Use defaults or prompt user
- **Network Unreachable**: Report and continue with next check
- **Permission Denied**: Inform user about required permissions
- **Timeout**: Report timeout and suggest increasing timeout value

## Security Considerations

### Input Validation

- Configuration values are validated before use
- Command injection prevention through parameterized execution
- No direct shell command interpolation with user input

### Least Privilege

- Runs with user privileges (no sudo required for most operations)
- ICMP ping can fall back to unprivileged mode
- No credential storage or handling

### Network Safety

- Read-only network operations
- No modifications to system configuration
- No data sent to external services

## Build and Release

### Build System

- **Make**: Simple automation via Makefile
- **GoReleaser**: Automated release creation
- **GitHub Actions**: CI/CD pipeline

### Release Process

```
1. Version Bump
   └─ scripts/version-up.sh --patch --apply

2. Build
   └─ goreleaser build --clean

3. Test
   └─ make test

4. Release
   └─ goreleaser release --clean

5. Distribution
   ├─ GitHub Releases
   ├─ Docker Images (ghcr.io, Docker Hub)
   └─ Homebrew Tap
```

### Artifact Types

1. **Binaries**: Platform-specific executables
2. **Docker Images**: Multi-arch containers
3. **Homebrew Formula**: macOS package manager
4. **Checksums**: Integrity verification

## Testing Strategy

### Current Testing

- **Unit Tests**: `go test ./...`
- **Race Detection**: `go test -race`
- **Coverage**: Code coverage tracking

### Test Execution

```bash
make test  # Runs all tests with coverage
```

## Dependencies

### Core Dependencies

- **cobra**: CLI framework
- **viper**: Configuration management
- **go-ping/ping**: ICMP ping implementation
- **miekg/dns**: DNS protocol library
- **go-pretty**: Table formatting
- **gookit/color**: Terminal colors

### Build Dependencies

- **goreleaser**: Release automation
- **docker**: Container builds
- **go 1.25+**: Go compiler

## Future Considerations

### Potential Enhancements

1. **Output Formats**: JSON, XML for machine parsing
2. **Plugins**: Extensible check system
3. **Caching**: Cache DNS results for performance
4. **Parallel Execution**: Run checks concurrently
5. **Web UI**: Optional web-based interface
6. **Continuous Monitoring**: Daemon mode for ongoing checks

### Scalability

- Currently designed for single-user, single-execution scenarios
- Could be extended for continuous monitoring
- Suitable for integration into larger diagnostic systems

## Maintainability

### Code Organization

- **Clear Separation**: Commands, helpers, configuration
- **Consistent Patterns**: Similar structure across commands
- **Documentation**: Inline comments for complex logic
- **Standard Go Practices**: Following Go idioms and conventions

### Development Workflow

1. Local development with `make install`
2. Testing with `make test`
3. Build validation with `make build`
4. Release with `make release`

---

For more information:
- **Contributing**: See CONTRIBUTING.md
- **Security**: See SECURITY.md
- **Changes**: See CHANGELOG.md
