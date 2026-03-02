<p align="center">
  <img src="https://github.com/slmingol/doxctl/blob/main/imgs/whats_up_dox__banner.png?raw=true" alt="dox logo">
</p>

[![Release Version](https://img.shields.io/github/v/release/slmingol/doxctl?label=version&color=blue)](https://github.com/slmingol/doxctl/releases/latest)
[![Build & Release](https://github.com/slmingol/doxctl/actions/workflows/build-release.yml/badge.svg?event=push)](https://github.com/slmingol/doxctl/actions/workflows/build-release.yml)
[![CodeQL](https://github.com/slmingol/doxctl/actions/workflows/codeql.yml/badge.svg)](https://github.com/slmingol/doxctl/actions/workflows/codeql.yml)
[![Lint](https://github.com/slmingol/doxctl/actions/workflows/lint.yml/badge.svg)](https://github.com/slmingol/doxctl/actions/workflows/lint.yml)
[![Code Coverage](https://github.com/slmingol/doxctl/actions/workflows/codecoverage.yml/badge.svg)](https://github.com/slmingol/doxctl/actions/workflows/codecoverage.yml)
[![Security Scan](https://github.com/slmingol/doxctl/actions/workflows/security-scan.yml/badge.svg)](https://github.com/slmingol/doxctl/actions/workflows/security-scan.yml)

# TLDR
`doxctl` is a diagnostic CLI tool that endusers can use to triage connectivity problems stemming from their VPN & DNS setups on their laptops. It can help with the following areas:

| Command | Description |
| ------- | ----------- |
| **dns** | Verify DNS resolver configuration, connectivity, and query responses |
| **vpn** | Check VPN connection status, network interfaces, and route configuration |
| **svrs** | Test connectivity to well-known servers through VPN |
| **svcs** | Multi-datacenter service-level health checks for critical endpoints |
| **net** | Network performance testing and SLO validation (latency, jitter, packet loss) |

## Features

- **Human-readable table output** - Easy to read diagnostic results
- **Machine-readable formats** - JSON and YAML output for automation, CI/CD pipelines, and monitoring systems
- **VPN diagnostics** - Check VPN connectivity, routes, and interface status
- **DNS diagnostics** - Verify DNS resolver configuration and connectivity
- **Server reachability** - Test connectivity to well-known servers

For details on using JSON/YAML output for automation, see [docs/OUTPUT_FORMATS.md](docs/OUTPUT_FORMATS.md).

# Requirements

For building from source or development:
- **Go 1.25.0+** is required (as specified in go.mod)

# Installation
## MacOS
<details><summary>Tree - CLICK ME</summary>
<p>

```
$ brew install slmingol/tap/doxctl
```

### Examples
#### Install
```
$ brew install slmingol/tap/doxctl
==> Installing doxctl from slmingol/tap
==> Downloading https://ghcr.io/v2/homebrew/core/go/manifests/1.25.0
######################################################################## 100.0%
==> Downloading https://ghcr.io/v2/homebrew/core/go/blobs/sha256:...
######################################################################## 100.0%
==> Downloading https://github.com/slmingol/doxctl/releases/download/1.0.1/doxctl_1.0.1_Darwin_x86_64.tar.gz
######################################################################## 100.0%
==> Installing dependencies for slmingol/tap/doxctl: go
==> Installing slmingol/tap/doxctl dependency: go
==> Pouring go--1.25.0.arm64_sequoia.bottle.tar.gz
🍺  /opt/homebrew/Cellar/go/1.25.0: 9,956 files, 503.6MB
==> Installing slmingol/tap/doxctl
🍺  /opt/homebrew/Cellar/doxctl/1.0.1: 5 files, 9.1MB, built in 4 seconds
```

#### search
```
$ brew search doxctl
==> Formulae
slmingol/tap/doxctl ✔
```

#### Upgrade
```
$ brew update
Updated 2 taps (homebrew/cask and slmingol/tap).
==> Updated Formulae
slmingol/tap/doxctl ✔

You have 96 outdated formulae and 10 outdated casks installed.
You can upgrade them with brew upgrade
or list them with brew outdated.


$ brew upgrade doxctl
==> Upgrading 1 outdated package:
slmingol/tap/doxctl 1.0.0 -> 1.0.1
==> Upgrading slmingol/tap/doxctl 1.0.0 -> 1.0.1
==> Downloading https://github.com/slmingol/doxctl/releases/download/1.0.1/doxctl_1.0.1_Darwin_x86_64.tar.gz
######################################################################## 100.0%
🍺  /opt/homebrew/Cellar/doxctl/1.0.1: 5 files, 9.1MB, built in 4 seconds
Removing: /opt/homebrew/Cellar/doxctl/1.0.0... (5 files, 9.1MB)
Removing: /Users/smingolelli/Library/Caches/Homebrew/doxctl--1.0.0.tar.gz... (3.4MB)
```

#### Uninstall
```
$ brew uninstall doxctl
Uninstalling /opt/homebrew/Cellar/doxctl/1.0.1... (5 files, 9.1MB)
```
</p>
</details>

## Docker

Docker images are available for both `amd64` and `arm64` architectures, making `doxctl` compatible with:
- x86_64 systems (traditional servers, Intel-based laptops)
- ARM-based systems (Apple M1/M2, Raspberry Pi, AWS Graviton)

### Available Images

Images are published to both GitHub Container Registry (GHCR) and Docker Hub:
- `ghcr.io/slmingol/doxctl:latest` (multi-arch)
- `slmingol/doxctl:latest` (multi-arch)
- Tagged versions: `ghcr.io/slmingol/doxctl:1.0.1` or `slmingol/doxctl:1.0.1`

### Quick Start

```bash
# Pull and run the latest version
docker run --rm ghcr.io/slmingol/doxctl:latest --help

# Run DNS diagnostics
docker run --rm ghcr.io/slmingol/doxctl:latest dns --help

# Run VPN diagnostics
docker run --rm ghcr.io/slmingol/doxctl:latest vpn --help

# Run server connectivity checks
docker run --rm ghcr.io/slmingol/doxctl:latest svrs --help
```

### Advanced Usage

#### With Configuration File

```bash
# Create a config file
cat > doxctl.yaml <<EOF
# Your configuration here
EOF

# Mount and use the config file
docker run --rm \
  -v $(pwd)/doxctl.yaml:/root/.doxctl.yaml \
  ghcr.io/slmingol/doxctl:latest dns -a
```

#### Interactive Shell

The Docker image is based on Alpine Linux, which includes a shell for debugging:

```bash
# Start an interactive shell
docker run --rm -it --entrypoint /bin/sh ghcr.io/slmingol/doxctl:latest

# Inside the container, you can run doxctl commands
/usr/bin/doxctl --help
```

#### Network Troubleshooting

For network diagnostics, you may need host network mode:

```bash
# Run with host network access
docker run --rm --network host ghcr.io/slmingol/doxctl:latest dns -p
```

### Architecture-Specific Images

While the multi-arch manifest automatically selects the correct image for your platform, you can explicitly pull architecture-specific images:

```bash
# AMD64 (x86_64)
docker pull ghcr.io/slmingol/doxctl:latest-amd64

# ARM64 (Apple Silicon, Graviton, etc.)
docker pull ghcr.io/slmingol/doxctl:latest-arm64
```

# Usage
## General
```
$ doxctl -h

'doxctl' is a collection of tools which can be used to diagnose & triage problems
stemming from the following areas with a laptop or desktop system:

  - DNS, specifically with the configuration of resolvers
  - VPN configuration and network connectivity over it
  - General access to well-known servers
  - ... or general network connectivity issues

Usage:
  doxctl [command]

Available Commands:
  dns         Run diagnostics related to DNS servers (aka. resolvers) configurations
  help        Help about any command
  net         Network performance testing and SLO validation
  svcs        Service-level health checks for multi-datacenter endpoints
  svrs        Run diagnostics verifying connectivity to well known servers thru a VPN connection
  vpn         Run diagnostics related to VPN connections, network interfaces & configurations

Flags:
  -c, --config string   config file (default is $HOME/.doxctl.yaml)
  -h, --help            help for doxctl
  -o, --output string   Output format: table, json, yaml (default "table")
  -v, --verbose         Enable verbose output of commands

Use "doxctl [command] --help" for more information about a command.
```

------------------------------------------------------------------------------

## NET
```
$ doxctl net -h

Test network connectivity and performance against defined SLO thresholds.

This command measures:
  - Average, minimum, and maximum latency
  - Jitter (latency variance)
  - Packet loss percentage
  - SLO compliance (latency threshold)

Examples:
  # Test network performance to configured targets
  doxctl net --perf

  # Set custom SLO threshold (default: 50ms)
  doxctl net --perf --slo 100

  # Specify number of packets to send (default: 10)
  doxctl net --perf --packets 20

Usage:
  doxctl net [flags]

Flags:
  -h, --help          help for net
  -n, --packets int   Number of packets to send (default 10)
  -p, --perf          Run network performance tests
  -s, --slo float     SLO threshold in milliseconds (default 50)

Global Flags:
  -c, --config string   config file (default is $HOME/.doxctl.yaml)
  -o, --output string   Output format: table, json, yaml (default "table")
  -v, --verbose         Enable verbose output of commands
```

### Example Output

#### Network Performance Check

<details><summary>Tree - CLICK ME</summary>
<p>

##### With Good Performance
```
$ doxctl net --perf

┌─────────────────────────────────────────────────────────────────────────────────┐
│ Network Performance & SLO Validation                                            │
├───────────────────────┬──────────┬──────────┬──────────┬─────────────┬────────┬─────────┬────────┤
│ TARGET                │ AVG (MS) │ MIN (MS) │ MAX (MS) │ JITTER (MS) │ LOSS % │ SLO     │ STATUS │
├───────────────────────┼──────────┼──────────┼──────────┼─────────────┼────────┼─────────┼────────┤
│ lab1.example.com      │ 12.50    │ 10.20    │ 15.80    │ 1.20        │ 0.0    │ 50 ms   │ ✓ PASS │
│ rdu1.example.com      │ 18.30    │ 15.10    │ 22.50    │ 2.10        │ 0.0    │ 50 ms   │ ✓ PASS │
│ dfw1.example.com      │ 45.20    │ 42.00    │ 48.90    │ 1.80        │ 0.0    │ 50 ms   │ ✓ PASS │
└───────────────────────┴──────────┴──────────┴──────────┴─────────────┴────────┴─────────┴────────┘

Summary: 3/3 targets meeting SLO (100.0% success rate)
```

##### With Performance Issues
```
$ doxctl net --perf

┌─────────────────────────────────────────────────────────────────────────────────┐
│ Network Performance & SLO Validation                                            │
├───────────────────────┬──────────┬──────────┬──────────┬─────────────┬────────┬─────────┬────────┤
│ TARGET                │ AVG (MS) │ MIN (MS) │ MAX (MS) │ JITTER (MS) │ LOSS % │ SLO     │ STATUS │
├───────────────────────┼──────────┼──────────┼──────────┼─────────────┼────────┼─────────┼────────┤
│ lab1.example.com      │ 12.50    │ 10.20    │ 15.80    │ 1.20        │ 0.0    │ 50 ms   │ ✓ PASS │
│ rdu1.example.com      │ 85.40    │ 78.10    │ 95.20    │ 5.30        │ 2.5    │ 50 ms   │ ✗ FAIL │
│ dfw1.example.com      │ 45.20    │ 42.00    │ 48.90    │ 1.80        │ 8.0    │ 50 ms   │ ✗ FAIL │
└───────────────────────┴──────────┴──────────┴──────────┴─────────────┴────────┴─────────┴────────┘

Summary: 1/3 targets meeting SLO (33.3% success rate)
```

##### JSON Output
```
$ doxctl net --perf -o json

{
  "timestamp": "2026-03-02T10:30:00Z",
  "results": [
    {
      "timestamp": "2026-03-02T10:30:00Z",
      "target": "lab1.example.com",
      "avgLatencyMs": 12.5,
      "minLatencyMs": 10.2,
      "maxLatencyMs": 15.8,
      "jitterMs": 1.2,
      "packetLoss": 0.0,
      "meetsSLO": true,
      "sloThreshold": 50.0
    }
  ],
  "summary": {
    "totalTargets": 1,
    "passing": 1,
    "failing": 0
  }
}
```
</p>
</details>

------------------------------------------------------------------------------

## SVCS
```
$ doxctl svcs -h

Check the health and availability of services across multiple datacenters.

This command performs HTTP/HTTPS health checks on service endpoints and measures:
  - Response time
  - HTTP status codes
  - Service availability
  - Multi-datacenter service health

Examples:
  # Check health of all configured services
  doxctl svcs --health

  # Set custom timeout (default: 5 seconds)
  doxctl svcs --health --timeout 10

  # Skip TLS verification for self-signed certificates
  doxctl svcs --health --insecure

Usage:
  doxctl svcs [flags]

Flags:
  -H, --health        Run service health checks
  -h, --help          help for svcs
  -k, --insecure      Skip TLS certificate verification
  -t, --timeout int   HTTP request timeout in seconds (default 5)

Global Flags:
  -c, --config string   config file (default is $HOME/.doxctl.yaml)
  -o, --output string   Output format: table, json, yaml (default "table")
  -v, --verbose         Enable verbose output of commands
```

### Example Output

#### Service Health Checks

<details><summary>Tree - CLICK ME</summary>
<p>

##### All Services Healthy
```
$ doxctl svcs --health

┌──────────────────────────────────────────────────────────────────────────────────────────┐
│ Service Health Checks                                                                    │
├───────────┬───────────────────────────────────────────────┬─────────────┬──────────────┬───────────┬───────┤
│ SERVICE   │ ENDPOINT                                      │ STATUS CODE │ RESPONSE (MS)│ STATUS    │ ERROR │
├───────────┼───────────────────────────────────────────────┼─────────────┼──────────────┼───────────┼───────┤
│ openshift │ https://ocp-master-01a.lab1.example.com:6443  │ 200         │ 45.23        │ ✓ Healthy │       │
│ openshift │ https://ocp-master-01b.rdu1.example.com:6443  │ 200         │ 32.18        │ ✓ Healthy │       │
│ elastic   │ https://es-master-01a.lab1.example.com:6443   │ 200         │ 28.45        │ ✓ Healthy │       │
│ elastic   │ https://es-master-01b.rdu1.example.com:6443   │ 200         │ 22.91        │ ✓ Healthy │       │
└───────────┴───────────────────────────────────────────────┴─────────────┴──────────────┴───────────┴───────┘

Summary: 4/4 services healthy (100.0% availability)
```

##### With Service Failures
```
$ doxctl svcs --health

┌──────────────────────────────────────────────────────────────────────────────────────────────────┐
│ Service Health Checks                                                                            │
├───────────┬───────────────────────────────────────────────┬─────────────┬──────────────┬──────────┬──────────────────────┤
│ SERVICE   │ ENDPOINT                                      │ STATUS CODE │ RESPONSE (MS)│ STATUS   │ ERROR                │
├───────────┼───────────────────────────────────────────────┼─────────────┼──────────────┼──────────┼──────────────────────┤
│ openshift │ https://ocp-master-01a.lab1.example.com:6443  │ 200         │ 45.23        │ ✓ Healthy│                      │
│ openshift │ https://ocp-master-01b.rdu1.example.com:6443  │ -           │ -            │ ✗ Failed │ connection refused   │
│ elastic   │ https://es-master-01a.lab1.example.com:6443   │ 503         │ 125.40       │ ✗ Failed │                      │
│ elastic   │ https://es-master-01b.rdu1.example.com:6443   │ 200         │ 22.91        │ ✓ Healthy│                      │
└───────────┴───────────────────────────────────────────────┴─────────────┴──────────────┴──────────┴──────────────────────┘

Summary: 2/4 services healthy (50.0% availability)
```

##### YAML Output
```
$ doxctl svcs --health -o yaml

timestamp: 2026-03-02T10:35:00Z
results:
  - timestamp: 2026-03-02T10:35:00Z
    service: openshift
    endpoint: https://ocp-master-01a.lab1.example.com:6443/healthz
    responseTimeMs: 45.23
    statusCode: 200
    healthy: true
    error: ""
summary:
  total: 1
  healthy: 1
  failed: 0
```
</p>
</details>

------------------------------------------------------------------------------

## DNS
```
$ doxctl dns -h

doxctl's 'dns' subcommand can help triage DNS resovler configuration issues,
general access to DNS resolvers and name resolution against DNS resolvers.

Usage:
  doxctl dns [flags]

Flags:
  -a, --allChk        Run all the checks in this subcommand module
  -d, --digChk        Check if VPN defined resolvers respond with well-known servers in DCs
  -h, --help          help for dns
  -p, --pingChk       Check if VPN defined resolvers are pingable & reachable
  -r, --resolverChk   Check if VPN designated DNS resolvers are configured

Global Flags:
  -c, --config string   config file (default is $HOME/.doxctl.yaml)
  -v, --verbose         Enable verbose output of commands
```

### Example Output

#### resolverChk

<details><summary>Tree - CLICK ME</summary>
<p>

##### Off VPN
```
$ doxctl dns -r

**NOTE:** Using config file: /Users/smingolelli/.doxctl.yaml


┌───────────────────────────────────────────────────────────────────────────┐
│ VPN defined DNS Resolver Checks                                           │
├──────────────────────────────────────────┬────────────────────────────────┤
│ PROPERTY DESCRIPTION                     │ VALUE                          │
├──────────────────────────────────────────┼────────────────────────────────┤
│ DomainName defined?                      │ unset                          │
│ SearchDomains defined?                   │ unset                          │
│ ServerAddresses defined?                 │ unset                          │
└──────────────────────────────────────────┴────────────────────────────────┘

INFO: Any values of unset indicate that the VPN client is not defining DNS resolver(s) properly!


```

##### On VPN
```
$ doxctl dns -r

**NOTE:** Using config file: /Users/smingolelli/.doxctl.yaml


┌───────────────────────────────────────────────────────────────────────────┐
│ VPN defined DNS Resolver Checks                                           │
├──────────────────────────────────────────┬────────────────────────────────┤
│ PROPERTY DESCRIPTION                     │ VALUE                          │
├──────────────────────────────────────────┼────────────────────────────────┤
│ DomainName defined?                      │ set                            │
│ SearchDomains defined?                   │ set                            │
│ ServerAddresses defined?                 │ set                            │
└──────────────────────────────────────────┴────────────────────────────────┘

INFO: Any values of unset indicate that the VPN client is not defining DNS resolver(s) properly!


```
</p>
</details>

#### pingChk

<details><summary>Tree - CLICK ME</summary>
<p>

##### Off VPN
```
$ doxctl dns -p

**NOTE:** Using config file: /Users/smingolelli/.doxctl.yaml

┌──────────────────────────────────────────────────────────────────────────────────┐
│ VPN defined DNS Resolver Connectivity Checks                                     │
├──────────────────────────────────────────┬───────────────┬───────────────┬───────┤
│                     PROPERTY DESCRIPTION │            IP │       NET I/F │ VALUE │
├──────────────────────────────────────────┼───────────────┼───────────────┼───────┤
└──────────────────────────────────────────┴───────────────┴───────────────┴───────┘

WARNING:

   Your VPN client does not appear to be defining any DNS resolver(s) properly,
   you're either not connected via VPN or it's misconfigured!



```

##### On VPN
```
$ doxctl dns -p

**NOTE:** Using config file: /Users/smingolelli/.doxctl.yaml

┌──────────────────────────────────────────────────────────────────────────────────┐
│ VPN defined DNS Resolver Connectivity Checks                                     │
├──────────────────────────────────────────┬───────────────┬───────────────┬───────┤
│ PROPERTY DESCRIPTION                     │ IP            │ NET I/F       │ VALUE │
├──────────────────────────────────────────┼───────────────┼───────────────┼───────┤
│ Resovler is pingable?                    │ 10.5.0.18     │ utun2         │ true  │
│ Reachable via TCP?                       │ 10.5.0.18     │ utun2         │ true  │
│ Reachable via UDP?                       │ 10.5.0.18     │ utun2         │ true  │
├──────────────────────────────────────────┼───────────────┼───────────────┼───────┤
│ Resovler is pingable?                    │ 10.5.0.19     │ utun2         │ true  │
│ Reachable via TCP?                       │ 10.5.0.19     │ utun2         │ true  │
│ Reachable via UDP?                       │ 10.5.0.19     │ utun2         │ true  │
└──────────────────────────────────────────┴───────────────┴───────────────┴───────┘



```
</p>
</details>

#### digChk

<details><summary>Tree - CLICK ME</summary>
<p>

##### Off VPN
```
$ doxctl dns -d

**NOTE:** Using config file: /Users/smingolelli/.doxctl.yaml

┌──────────────────────────────────────────────────────────────────────────────┐
│ Dig Check against VPN defined DNS Resolvers                                  │
├──────────────────────────────────────────┬─────────────────┬─────────────────┤
│ HOSTNAME TO 'DIG'                        │ RESOLVER IP     │ IS RESOLVABLE?  │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.lab1.somedom.local               │                 │ false           │
│ idm-01b.lab1.somedom.local               │                 │ false           │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.rdu1.somedom.local               │                 │ false           │
│ idm-01b.rdu1.somedom.local               │                 │ false           │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.atl1.somedom.local               │                 │ false           │
│ idm-01b.atl1.somedom.local               │                 │ false           │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.dfw1.somedom.local               │                 │ false           │
│ idm-01b.dfw1.somedom.local               │                 │ false           │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.lax2.somedom.local               │                 │ false           │
│ idm-01b.lax2.somedom.local               │                 │ false           │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.jfk1.somedom.local               │                 │ false           │
│ idm-01b.jfk1.somedom.local               │                 │ false           │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ SUCCESSESFUL QUERIES                     │ RESOLVER #1: 0  │                 │
│                                          │ RESOLVER #2: 0  │                 │
└──────────────────────────────────────────┴─────────────────┴─────────────────┘

WARNING:

   Your VPN client does not appear to be defining any DNS resolver(s) properly,
   you're either not connected via VPN or it's misconfigured!



```

##### On VPN
```
$ doxctl dns -d

**NOTE:** Using config file: /Users/smingolelli/.doxctl.yaml

┌──────────────────────────────────────────────────────────────────────────────┐
│ Dig Check against VPN defined DNS Resolvers                                  │
├──────────────────────────────────────────┬─────────────────┬─────────────────┤
│ HOSTNAME TO 'DIG'                        │ RESOLVER IP     │ IS RESOLVABLE?  │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.lab1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01b.lab1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01a.lab1.somedom.local               │ 10.5.0.19       │ true            │
│ idm-01b.lab1.somedom.local               │ 10.5.0.19       │ true            │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.rdu1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01b.rdu1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01a.rdu1.somedom.local               │ 10.5.0.19       │ true            │
│ idm-01b.rdu1.somedom.local               │ 10.5.0.19       │ true            │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.atl1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01b.atl1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01a.atl1.somedom.local               │ 10.5.0.19       │ true            │
│ idm-01b.atl1.somedom.local               │ 10.5.0.19       │ true            │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.dfw1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01b.dfw1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01a.dfw1.somedom.local               │ 10.5.0.19       │ true            │
│ idm-01b.dfw1.somedom.local               │ 10.5.0.19       │ true            │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.lax2.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01b.lax2.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01a.lax2.somedom.local               │ 10.5.0.19       │ true            │
│ idm-01b.lax2.somedom.local               │ 10.5.0.19       │ true            │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ idm-01a.jfk1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01b.jfk1.somedom.local               │ 10.5.0.18       │ true            │
│ idm-01a.jfk1.somedom.local               │ 10.5.0.19       │ true            │
│ idm-01b.jfk1.somedom.local               │ 10.5.0.19       │ true            │
├──────────────────────────────────────────┼─────────────────┼─────────────────┤
│ SUCCESSESFUL QUERIES                     │ RESOLVER #1: 12 │                 │
│                                          │ RESOLVER #2: 12 │                 │
└──────────────────────────────────────────┴─────────────────┴─────────────────┘



```
</p>
</details>

------------------------------------------------------------------------------

## VPN
```
$ doxctl vpn -h

doxctl's 'vpn' subcommand can help triage VPN related configuration issues,
& routes related to a VPN connection.

Usage:
  doxctl vpn [flags]

Flags:
  -a, --allChk           Run all the checks in this subcommand module
  -h, --help             help for vpn
  -i, --ifReachableChk   Check if network interfaces are reachable
  -r, --vpnRoutesChk     Check if >5 VPN routes are defined
  -s, --vpnStatusChk     Check if VPN client's status reports as 'Connected'

Global Flags:
  -c, --config string   config file (default is $HOME/.doxctl.yaml)
  -v, --verbose         Enable verbose output of commands
```

### Example Output

#### ifReachableChk
<details><summary>Tree - CLICK ME</summary>
<p>

##### Off VPN
```
$ doxctl vpn -i

NOTE: Using config file: /Users/smingolelli/.doxctl.yaml

┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│ Interfaces Reachable Checks                                                                 │
├────────────────────────────────────────────────────┬────────────────────────────────┬───────┤
│ PROPERTY DESCRIPTION                               │ VALUE                          │ NOTES │
├────────────────────────────────────────────────────┼────────────────────────────────┼───────┤
│ How many network interfaces found?                 │ 1                              │ [en0] │
│ At least 1 interface's a utun device?              │ false                          │ []    │
│ All active interfaces are reporting as reachable?  │ true                           │       │
└────────────────────────────────────────────────────┴────────────────────────────────┴───────┘

WARNING:

   Your VPN client does not appear to be defining a TUN interface properly,
   your VPN is either not connected or it's misconfigured!



```

##### On VPN
```
$ doxctl vpn -i

NOTE: Using config file: /Users/smingolelli/.doxctl.yaml

┌───────────────────────────────────────────────────────────────────────────────────────────────────┐
│ Interfaces Reachable Checks                                                                       │
├────────────────────────────────────────────────────┬────────────────────────────────┬─────────────┤
│ PROPERTY DESCRIPTION                               │ VALUE                          │ NOTES       │
├────────────────────────────────────────────────────┼────────────────────────────────┼─────────────┤
│ How many network interfaces found?                 │ 2                              │ [en0 utun2] │
│ At least 1 interface's a utun device?              │ true                           │ [utun2]     │
│ All active interfaces are reporting as reachable?  │ true                           │             │
└────────────────────────────────────────────────────┴────────────────────────────────┴─────────────┘



```
</p>
</details>

#### vpnRoutesChk
<details><summary>Tree - CLICK ME</summary>
<p>

##### Off VPN
```
$ doxctl vpn -r

NOTE: Using config file: /Users/smingolelli/.doxctl.yaml

┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│ VPN Interface Route Checks                                                                  │
├────────────────────────────────────────────────────┬────────────────────────────────┬───────┤
│ PROPERTY DESCRIPTION                               │ VALUE                          │ NOTES │
├────────────────────────────────────────────────────┼────────────────────────────────┼───────┤
│ At least [5] routes using interface [NIL]?         │ false                          │     0 │
└────────────────────────────────────────────────────┴────────────────────────────────┴───────┘

WARNING:

   Your VPN client does not appear to be defining a TUN interface properly,
   it's either not connected or it's misconfigured!



```

##### On VPN
```
$ doxctl vpn -r

NOTE: Using config file: /Users/smingolelli/.doxctl.yaml

┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│ VPN Interface Route Checks                                                                  │
├────────────────────────────────────────────────────┬────────────────────────────────┬───────┤
│ PROPERTY DESCRIPTION                               │ VALUE                          │ NOTES │
├────────────────────────────────────────────────────┼────────────────────────────────┼───────┤
│ At least [5] routes using interface [utun2]?       │ true                           │   148 │
└────────────────────────────────────────────────────┴────────────────────────────────┴───────┘



```
</p>
</details>

#### vpnStatusChk
<details><summary>Tree - CLICK ME</summary>
<p>

##### Off VPN
```
$ doxctl vpn -s

NOTE: Using config file: /Users/smingolelli/.doxctl.yaml

┌───────────────────────────────────────────────────────────────────────────────────────────────┐
│ VPN Connection Status Checks                                                                  │
├──────────────────────────────────────────────────────┬────────────────────────────────┬───────┤
│ PROPERTY DESCRIPTION                                 │ VALUE                          │ NOTES │
├──────────────────────────────────────────────────────┼────────────────────────────────┼───────┤
│ VPN Client reports connection status as 'Connected'? │ false                          │       │
└──────────────────────────────────────────────────────┴────────────────────────────────┴───────┘

WARNING:

   Your VPN client's does not appear to be a state of 'connected',
   it's either down or misconfigured!"



```

##### On VPN
```
$ doxctl vpn -s

NOTE: Using config file: /Users/smingolelli/.doxctl.yaml

┌───────────────────────────────────────────────────────────────────────────────────────────────┐
│ VPN Connection Status Checks                                                                  │
├──────────────────────────────────────────────────────┬────────────────────────────────┬───────┤
│ PROPERTY DESCRIPTION                                 │ VALUE                          │ NOTES │
├──────────────────────────────────────────────────────┼────────────────────────────────┼───────┤
│ VPN Client reports connection status as 'Connected'? │ true                           │       │
└──────────────────────────────────────────────────────┴────────────────────────────────┴───────┘



```
</p>
</details>

------------------------------------------------------------------------------

## SVRS
```
$ doxctl svrs -h

doxctl's 'svrs' subcommand can help triage & test connectivity to 'well known servers'
thru a VPN connection to servers which have been defined in your '.doxctl.yaml'
configuration file.

Usage:
  doxctl svrs [flags]

Flags:
  -a, --allChk             Run all the checks in this subcommand module
  -h, --help               help for svrs
  -s, --svrsReachableChk   Check if well known servers are reachable

Global Flags:
  -c, --config string   config file (default is $HOME/.doxctl.yaml)
  -v, --verbose         Enable verbose output of commands
```

### Example Output

#### svrsReachableChk
<details><summary>Tree - CLICK ME</summary>
<p>

##### Off VPN
```
$ doxctl svrs -s

NOTE: Using config file: /Users/smingolelli/.doxctl.yaml

INFO: Attempting to ping all well known servers, this may take a few...

   --- Working through svc: openshift
   --- Working through svc: elastic
   --- Working through svc: idm


   ...one sec, preparing `ping` results...


WARNING: More than 5 hosts appear to be unreachable, aborting remainder....


┌─────────────────────────────────────────────────────────────────────────────────────────────────┐
│ Well known Servers Reachable Checks                                                             │
├──────────────────────────────────────────┬──────────────────────┬────────────┬──────────────────┤
│ HOST                                     │ SERVICE              │ REACHABLE? │ PING PERFORMANCE │
├──────────────────────────────────────────┼──────────────────────┼────────────┼──────────────────┤
│ ocp-master-01a.lab1.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01a.rdu1.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01a.dfw1.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01a.lax2.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01a.jfk1.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01b.lab1.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01b.rdu1.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01b.dfw1.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01b.lax2.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01b.jfk1.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01c.lab1.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01c.rdu1.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01c.dfw1.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01c.lax2.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01c.jfk1.somedom.local        │ openshift            │ false      │ N/A              │
│ ocp-master-01a.lhr1.somedom.us           │ openshift            │ false      │ N/A              │
│ ocp-master-01a.fra1.somedom.us           │ openshift            │ false      │ N/A              │
│ ocp-master-01b.lhr1.somedom.us           │ openshift            │ false      │ N/A              │
│ ocp-master-01b.fra1.somedom.us           │ openshift            │ false      │ N/A              │
│ ocp-master-01c.lhr1.somedom.us           │ openshift            │ false      │ N/A              │
│ ocp-master-01c.fra1.somedom.us           │ openshift            │ false      │ N/A              │
├──────────────────────────────────────────┼──────────────────────┼────────────┼──────────────────┤
│ es-master-01a.lab1.somedom.local         │ elastic              │ false      │ N/A              │
│ es-master-01a.rdu1.somedom.local         │ elastic              │ false      │ N/A              │
│ es-master-01b.lab1.somedom.local         │ elastic              │ false      │ N/A              │
│ es-master-01b.rdu1.somedom.local         │ elastic              │ false      │ N/A              │
│ es-master-01c.lab1.somedom.local         │ elastic              │ false      │ N/A              │
│ es-master-01c.rdu1.somedom.local         │ elastic              │ false      │ N/A              │
├──────────────────────────────────────────┼──────────────────────┼────────────┼──────────────────┤
│ idm-01a.lab1.somedom.local               │ idm                  │ false      │ N/A              │
│ idm-01a.rdu1.somedom.local               │ idm                  │ false      │ N/A              │
│ idm-01a.dfw1.somedom.local               │ idm                  │ false      │ N/A              │
│ idm-01a.lax2.somedom.local               │ idm                  │ false      │ N/A              │
│ idm-01a.jfk1.somedom.local               │ idm                  │ false      │ N/A              │
│ idm-01b.lab1.somedom.local               │ idm                  │ false      │ N/A              │
│ idm-01b.rdu1.somedom.local               │ idm                  │ false      │ N/A              │
│ idm-01b.dfw1.somedom.local               │ idm                  │ false      │ N/A              │
│ idm-01b.lax2.somedom.local               │ idm                  │ false      │ N/A              │
│ idm-01b.jfk1.somedom.local               │ idm                  │ false      │ N/A              │
│ idm-01a.lhr1.somedom.us                  │ idm                  │ false      │ N/A              │
│ idm-01a.fra1.somedom.us                  │ idm                  │ false      │ N/A              │
│ idm-01b.lhr1.somedom.us                  │ idm                  │ false      │ N/A              │
│ idm-01b.fra1.somedom.us                  │ idm                  │ false      │ N/A              │
└──────────────────────────────────────────┴──────────────────────┴────────────┴──────────────────┘

WARNING:

   Your VPN client does not appear to be functioning properly, it's likely one or more of the following:

      - Well known servers are unreachable via ping   --- try running 'doxctl vpn -h'
      - Servers are unresovlable in DNS               --- try running 'doxctl dns -h'
      - VPN client is otherwise misconfigured!




```

##### On VPN
```
doxctl svrs -s

NOTE: Using config file: /Users/smingolelli/.doxctl.yaml

INFO: Attempting to ping all well known servers, this may take a few...

   --- Working through svc: openshift
   --- Working through svc: elastic
   --- Working through svc: idm


   ...one sec, preparing `ping` results...

┌────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│ Well known Servers Reachable Checks                                                                    │
├──────────────────────────────────────────┬──────────────────────┬────────────┬─────────────────────────┤
│ HOST                                     │ SERVICE              │ REACHABLE? │ PING PERFORMANCE        │
├──────────────────────────────────────────┼──────────────────────┼────────────┼─────────────────────────┤
│ ocp-master-01a.lab1.somedom.local        │ openshift            │ true       │ rnd-trp avg = 44.525ms  │
│ ocp-master-01a.rdu1.somedom.local        │ openshift            │ true       │ rnd-trp avg = 24.337ms  │
│ ocp-master-01a.dfw1.somedom.local        │ openshift            │ true       │ rnd-trp avg = 55.118ms  │
│ ocp-master-01a.lax2.somedom.local        │ openshift            │ true       │ rnd-trp avg = 97.183ms  │
│ ocp-master-01a.jfk1.somedom.local        │ openshift            │ true       │ rnd-trp avg = 36.187ms  │
│ ocp-master-01b.lab1.somedom.local        │ openshift            │ true       │ rnd-trp avg = 44.237ms  │
│ ocp-master-01b.rdu1.somedom.local        │ openshift            │ true       │ rnd-trp avg = 17.678ms  │
│ ocp-master-01b.dfw1.somedom.local        │ openshift            │ true       │ rnd-trp avg = 56.559ms  │
│ ocp-master-01b.lax2.somedom.local        │ openshift            │ true       │ rnd-trp avg = 96.493ms  │
│ ocp-master-01b.jfk1.somedom.local        │ openshift            │ true       │ rnd-trp avg = 43.273ms  │
│ ocp-master-01c.lab1.somedom.local        │ openshift            │ true       │ rnd-trp avg = 41.358ms  │
│ ocp-master-01c.rdu1.somedom.local        │ openshift            │ true       │ rnd-trp avg = 31.427ms  │
│ ocp-master-01c.dfw1.somedom.local        │ openshift            │ true       │ rnd-trp avg = 55.095ms  │
│ ocp-master-01c.lax2.somedom.local        │ openshift            │ true       │ rnd-trp avg = 103.423ms │
│ ocp-master-01c.jfk1.somedom.local        │ openshift            │ true       │ rnd-trp avg = 37.22ms   │
│ ocp-master-01a.lhr1.somedom.us           │ openshift            │ true       │ rnd-trp avg = 133.023ms │
│ ocp-master-01a.fra1.somedom.us           │ openshift            │ true       │ rnd-trp avg = 136.647ms │
│ ocp-master-01b.lhr1.somedom.us           │ openshift            │ true       │ rnd-trp avg = 127.451ms │
│ ocp-master-01b.fra1.somedom.us           │ openshift            │ true       │ rnd-trp avg = 139.85ms  │
│ ocp-master-01c.lhr1.somedom.us           │ openshift            │ true       │ rnd-trp avg = 132.362ms │
│ ocp-master-01c.fra1.somedom.us           │ openshift            │ true       │ rnd-trp avg = 137.558ms │
├──────────────────────────────────────────┼──────────────────────┼────────────┼─────────────────────────┤
│ es-master-01a.lab1.somedom.local         │ elastic              │ true       │ rnd-trp avg = 44.029ms  │
│ es-master-01a.rdu1.somedom.local         │ elastic              │ true       │ rnd-trp avg = 32.187ms  │
│ es-master-01b.lab1.somedom.local         │ elastic              │ true       │ rnd-trp avg = 48.833ms  │
│ es-master-01b.rdu1.somedom.local         │ elastic              │ true       │ rnd-trp avg = 22.477ms  │
│ es-master-01c.lab1.somedom.local         │ elastic              │ true       │ rnd-trp avg = 55.587ms  │
│ es-master-01c.rdu1.somedom.local         │ elastic              │ true       │ rnd-trp avg = 25.79ms   │
├──────────────────────────────────────────┼──────────────────────┼────────────┼─────────────────────────┤
│ idm-01a.lab1.somedom.local               │ idm                  │ true       │ rnd-trp avg = 47.484ms  │
│ idm-01a.rdu1.somedom.local               │ idm                  │ true       │ rnd-trp avg = 22.766ms  │
│ idm-01a.dfw1.somedom.local               │ idm                  │ true       │ rnd-trp avg = 54.07ms   │
│ idm-01a.lax2.somedom.local               │ idm                  │ true       │ rnd-trp avg = 94.755ms  │
│ idm-01a.jfk1.somedom.local               │ idm                  │ true       │ rnd-trp avg = 36.26ms   │
│ idm-01b.lab1.somedom.local               │ idm                  │ true       │ rnd-trp avg = 41.171ms  │
│ idm-01b.rdu1.somedom.local               │ idm                  │ true       │ rnd-trp avg = 27.097ms  │
│ idm-01b.dfw1.somedom.local               │ idm                  │ true       │ rnd-trp avg = 51.547ms  │
│ idm-01b.lax2.somedom.local               │ idm                  │ true       │ rnd-trp avg = 94.203ms  │
│ idm-01b.jfk1.somedom.local               │ idm                  │ true       │ rnd-trp avg = 36.956ms  │
│ idm-01a.lhr1.somedom.us                  │ idm                  │ true       │ rnd-trp avg = 145.853ms │
│ idm-01a.fra1.somedom.us                  │ idm                  │ true       │ rnd-trp avg = 146.425ms │
│ idm-01b.lhr1.somedom.us                  │ idm                  │ true       │ rnd-trp avg = 127.987ms │
│ idm-01b.fra1.somedom.us                  │ idm                  │ true       │ rnd-trp avg = 135.593ms │
└──────────────────────────────────────────┴──────────────────────┴────────────┴─────────────────────────┘



```
</p>
</details>

------------------------------------------------------------------------------

## Debugging - TBD/WIP
All the CLI subcommands can make use of either the `-v` or the `--verbose` switch to gather further diagnostics which can be helpful when triaging connectivity issues.

<details><summary>Tree - CLICK ME</summary>
<p>

For example:
```
$ doxctl dns -r -v
+ printf '\n\nDNS Resolver Checks\n===================\n\n\n'


DNS Resolver Checks
===================


++ printf 'get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n'
++ scutil
+ vpn_resolvers='<dictionary> {
  DomainName : somedom.local
  SearchDomains : <array> {
    0 : somedom.local
  }
  SearchOrder : 1
  ServerAddresses : <array> {
    0 : 10.5.0.18
    1 : 10.5.0.19
    2 : 192.168.7.85
  }
  SupplementalMatchDomains : <array> {
    0 :
    1 : somedom.local
  }
}'
+ column -t
+ echo '<dictionary> {
  DomainName : somedom.local
  SearchDomains : <array> {
    0 : somedom.local
  }
  SearchOrder : 1
  ServerAddresses : <array> {
    0 : 10.5.0.18
    1 : 10.5.0.19
    2 : 192.168.7.85
  }
  SupplementalMatchDomains : <array> {
    0 :
    1 : somedom.local
  }
}'
+ grep -q 'DomainName.*somedom.local'
+ echo 'DomainName set'
+ echo '<dictionary> {
  DomainName : somedom.local
  SearchDomains : <array> {
    0 : somedom.local
  }
  SearchOrder : 1
  ServerAddresses : <array> {
    0 : 10.5.0.18
    1 : 10.5.0.19
    2 : 192.168.7.85
  }
  SupplementalMatchDomains : <array> {
    0 :
    1 : somedom.local
  }
}'
+ grep -A1 SearchDomains
+ grep -qE '[0-1].*somedom'
+ echo 'SearchDomains set'
+ echo '<dictionary> {
  DomainName : somedom.local
  SearchDomains : <array> {
    0 : somedom.local
  }
  SearchOrder : 1
  ServerAddresses : <array> {
    0 : 10.5.0.18
    1 : 10.5.0.19
    2 : 192.168.7.85
  }
  SupplementalMatchDomains : <array> {
    0 :
    1 : somedom.local
  }
}'
+ grep -A3 ServerAddresses
+ grep -qE '[0-1].*10.5'
+ echo 'ServerAddresses set'
DomainName       set
SearchDomains    set
ServerAddresses  set
+ [[ 1 -eq 1 ]]
+ set +x




```
</p>
</details>

## Development

### Prerequisites

- **Go 1.25.0+** (as specified in go.mod)
- **Git** for version control

### Building from Source

```bash
# Clone the repository
git clone https://github.com/slmingol/doxctl.git
cd doxctl

# Build the binary
go build -o doxctl .

# Run the binary
./doxctl --help
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Quality

The project uses [golangci-lint](https://golangci-lint.run/) for code quality enforcement:

```bash
# Install golangci-lint (if not already installed)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Run with auto-fix for some issues
golangci-lint run --fix
```

### Development Workflow

1. **Make changes** to the code
2. **Run tests**: `go test ./...`
3. **Run linter**: `golangci-lint run`
4. **Build**: `go build -o doxctl .`
5. **Test locally**: `./doxctl [command]`

## Security Features

This project implements multiple layers of security and code quality checks:

### Automated Security Scanning

- **CodeQL Analysis**: Runs on every push and pull request to identify security vulnerabilities
- **Security Scan Workflow**: Dedicated security scanning workflow runs on all commits
- **Dependabot**: Automatically monitors dependencies and creates PRs for security updates
  - Checks Go modules weekly
  - Checks GitHub Actions weekly

### Code Quality Enforcement

- **golangci-lint**: Comprehensive linting with multiple analyzers:
  - `errcheck`: Ensures error handling
  - `gosec`: Security-focused static analysis
  - `staticcheck`: Advanced Go static analysis
  - `gofmt`: Code formatting
  - `misspell`: Spelling checks
  - And many more (see `.golangci.yml`)

### Reporting Security Issues

Please see [SECURITY.md](SECURITY.md) for our security policy and instructions on reporting vulnerabilities.

## Platform Support

### Current Status

**macOS**: Fully supported
- Uses native `scutil` command for DNS configuration
- Tested on Intel and Apple Silicon (M1/M2)

**Linux**: Partial support
- Multi-arch builds available (amd64, arm64)
- Some DNS features rely on macOS-specific `scutil` command
- VPN and server connectivity checks work across platforms

### Docker Deployment

Docker images support both architectures:
- **amd64**: Traditional x86_64 systems
- **arm64**: Apple M1/M2, AWS Graviton, Raspberry Pi

See the [Docker](#docker) section for usage examples.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for:
- Code of conduct
- Development guidelines
- Pull request process
- Coding standards

For security-related contributions, please review [SECURITY.md](SECURITY.md) first.

## Resources

### References
- [Exploring Go Packages: Cobra](https://levelup.gitconnected.com/exploring-go-packages-cobra-fce6c4e331d6)
- [CLI Command Line SDK - Cobra](https://github.com/spf13/cobra)
- [Building a multipurpose CLI tool with Cobra and Go](https://dev.to/lumexralph/building-a-multipurpose-cli-tool-with-cobra-and-go-2492)
- [How to create a CLI in golang with cobra](https://towardsdatascience.com/how-to-create-a-cli-in-golang-with-cobra-d729641c7177)
- [How to pipe several commands in Go?](https://stackoverflow.com/questions/10781516/how-to-pipe-several-commands-in-go)
- [Executing System Commands With Golang](https://tutorialedge.net/golang/executing-system-commands-with-golang/)
- [gookit/color](https://github.com/gookit/color)
- [go-ping/ping](https://github.com/go-ping/ping)
- [go-ping example `ping` CLI](https://github.com/go-ping/ping/blob/master/cmd/ping/ping.go)
- [go-pretty demo-table](https://github.com/jedib0t/go-pretty/tree/main/cmd/demo-table)

### Example CLI tools written in Go
- [docker/hub-tool](https://github.com/docker/hub-tool/tree/main/internal/commands)

### MacOS
- [SCNetworkReachability](https://developer.apple.com/documentation/systemconfiguration/scnetworkreachability-g7d)
- [scutil generalized interface to "dynamic Store" and Network Services](https://www.real-world-systems.com/docs/scutil.1.html)
- [Homebrew formulas - slmingol/homebrew-tap](https://github.com/slmingol/homebrew-tap)
