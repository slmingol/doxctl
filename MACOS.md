# Running doxctl in Containers on macOS

This guide shows how to run doxctl in a container on macOS with VPN DNS support.

## Problem

When running doxctl in containers on macOS, the container cannot access macOS-specific VPN DNS settings because:
- Containers run in a Linux VM on macOS
- The VM doesn't have access to macOS's `scutil` or VPN network interfaces
- VPN DNS settings aren't automatically passed through to containers

## Solution

Use the `doxctl-container` wrapper script which:
1. Checks if VPN is connected
2. Extracts VPN DNS settings from macOS
3. Creates a custom `resolv.conf`
4. Launches the container with proper DNS configuration

## Prerequisites

1. **VPN Connection**: You must be connected to your VPN
2. **Container Image**: Build the doxctl container image:
   ```bash
   cd doxctl
   GOOS=linux GOARCH=amd64 go build -o doxctl .
   podman build -t doxctl:latest .
   ```
3. **Configuration**: Create a `doxctl.yaml` config file (optional)

## Usage

### Basic Usage

```bash
./doxctl-container dns -a
```

### Example Output

```
=== doxctl Container Launcher (macOS) ===

✓ VPN Connected
  DNS Servers: 10.168.112.10, 10.168.112.11
  Search Domain: bandwidth.local
  Domain: local

✓ Created DNS config
✓ Launching container...
```

### Without VPN

If you run without VPN connected:

```bash
./doxctl-container dns -a
```

```
=== doxctl Container Launcher (macOS) ===

ERROR: No DNS servers detected!
Please connect to your VPN first.
```

## How It Works

The wrapper script:

1. **Detects VPN DNS** using `scutil --dns`
2. **Creates `/resolv.conf.vpn`** with:
   - Domain name
   - Search domains
   - VPN nameservers
3. **Launches container** with:
   - `--dns` flags for each VPN DNS server
   - `--dns-search` for search domains
   - Mounted `resolv.conf.vpn` as `/etc/resolv.conf`
   - Mounted `doxctl.yaml` config

## Configuration

Place a `doxctl.yaml` in the same directory as the wrapper script. Example:

```yaml
# DNS checks
domNameChk: "bandwidth.local"
domSearchChk: "bandwidth"
domAddrChk: "10.168"

# Sites and services
sites:
  - lab1
  - rdu1
  - dfw1

# Timeouts
pingTimeout: 250
dnsLookupTimeout: 100
```

## Troubleshooting

### DNS servers not detected

**Problem**: `ERROR: No DNS servers detected!`

**Solution**: 
- Ensure VPN is connected
- Check DNS settings: `scutil --dns`
- Verify nameservers appear in output

### Container can't resolve internal hosts

**Problem**: Container launches but can't resolve VPN hosts

**Solutions**:
1. Verify VPN DNS is reachable: `ping 10.168.112.10`
2. Check VPN tunnel is up: `ifconfig | grep utun`
3. Rebuild container if DNS config changed

### Config file not found

**Problem**: `WARNING: doxctl.yaml not found`

**Solution**: Either:
- Create `doxctl.yaml` in the same directory as wrapper
- Or let it use container defaults

## Native macOS Alternative

For better performance and integration, you can also run doxctl natively on macOS:

```bash
# Build native macOS binary
go build -o doxctl-mac .

# Run directly
./doxctl-mac dns -a
```

The native macOS version has been enhanced to support any VPN (not just Cisco AnyConnect).

## See Also

- Main README: [../README.md](../README.md)
- Configuration examples: [examples/](examples/)
