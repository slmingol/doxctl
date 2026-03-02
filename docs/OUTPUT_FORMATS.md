# Output Format Examples

This document demonstrates the new JSON and YAML output formats for automation.

## Usage

All commands now support the `--output` (or `-o`) flag with three options:
- `table` (default) - Human-readable table format
- `json` - Machine-readable JSON format
- `yaml` - Machine-readable YAML format

## Examples

### DNS Resolver Check

#### Table Output (default)
```bash
doxctl dns -r
```

#### JSON Output
```bash
doxctl dns -r -o json
```
Example output:
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "domainNameSet": true,
  "searchDomainsSet": true,
  "serverAddressesSet": true
}
```

#### YAML Output
```bash
doxctl dns -r -o yaml
```
Example output:
```yaml
timestamp: 2024-01-01T12:00:00Z
domainNameSet: true
searchDomainsSet: true
serverAddressesSet: true
```

### DNS Resolver Ping Check

#### JSON Output
```bash
doxctl dns -p -o json
```
Example output:
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "resolvers": [
    {
      "resolverIP": "10.5.0.1",
      "netInterface": "utun0",
      "pingReachable": true,
      "tcpReachable": true,
      "udpReachable": true
    }
  ]
}
```

### DNS Dig Check

#### JSON Output
```bash
doxctl dns -d -o json
```
Example output:
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "results": [
    {
      "hostname": "idm-01a.lab1.example.com",
      "resolverIP": "10.5.0.1",
      "isResolvable": true
    }
  ],
  "summary": {
    "10.5.0.1": 10,
    "10.5.0.2": 10
  }
}
```

### VPN Interface Check

#### JSON Output
```bash
doxctl vpn -i -o json
```
Example output:
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "interfaceCount": 3,
  "interfaces": ["eth0", "utun0", "utun1"],
  "hasTunInterface": true,
  "tunInterfaces": ["utun0", "utun1"],
  "allInterfacesReachable": true
}
```

### VPN Routes Check

#### JSON Output
```bash
doxctl vpn -r -o json
```
Example output:
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "vpnInterface": "utun0",
  "routeCount": 25,
  "minRoutesRequired": 5,
  "hasSufficientRoutes": true
}
```

### VPN Connection Status Check

#### JSON Output
```bash
doxctl vpn -s -o json
```
Example output:
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "isConnected": true
}
```

### Server Reachability Check

#### JSON Output
```bash
doxctl svrs -s -o json
```
Example output:
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "servers": [
    {
      "host": "server1.example.com",
      "service": "web",
      "reachable": true,
      "performance": "rnd-trp avg = 12ms"
    },
    {
      "host": "server2.example.com",
      "service": "db",
      "reachable": false,
      "performance": "N/A"
    }
  ],
  "pingFailures": 1,
  "reachFailures": 1
}
```

## Automation Use Cases

### CI/CD Pipeline Example

Check VPN connectivity before deployment:
```bash
#!/bin/bash
result=$(doxctl vpn -s -o json)
is_connected=$(echo "$result" | jq -r '.isConnected')

if [ "$is_connected" != "true" ]; then
  echo "VPN is not connected. Aborting deployment."
  exit 1
fi

echo "VPN connected. Proceeding with deployment."
```

### Monitoring Integration

Scrape DNS health metrics:
```bash
#!/bin/bash
doxctl dns -p -o json | jq '{
  timestamp: .timestamp,
  total_resolvers: (.resolvers | length),
  reachable_resolvers: [.resolvers[] | select(.pingReachable == true)] | length
}'
```

### Automated Alerting

Parse server reachability and alert on failures:
```bash
#!/bin/bash
result=$(doxctl svrs -s -o json)
failures=$(echo "$result" | jq '.pingFailures + .reachFailures')

if [ "$failures" -gt 5 ]; then
  # Send alert
  echo "High failure count: $failures"
  # curl -X POST https://alerting-service/alert -d "failures=$failures"
fi
```

### YAML Processing with Python

```python
import subprocess
import yaml

# Get VPN status
result = subprocess.run(
    ['doxctl', 'vpn', '-s', '-o', 'yaml'],
    capture_output=True,
    text=True
)

data = yaml.safe_load(result.stdout)
if not data['isConnected']:
    print("VPN is disconnected!")
```

## Benefits

1. **Machine-readable**: Easy to parse with standard tools (jq, yq, etc.)
2. **Scriptable**: Can be used in automation scripts
3. **CI/CD Integration**: Parse results in pipelines
4. **Monitoring**: Feed data to monitoring systems
5. **Alerting**: Automated alerting based on results
6. **Cross-platform**: Standard formats work everywhere
