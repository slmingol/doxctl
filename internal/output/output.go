/*
Package output - utilities for formatting output in different formats

Copyright © 2021 Sam Mingolelli <github@lamolabs.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Format represents the output format type
type Format string

const (
	// TableFormat specifies table formatted output (default)
	TableFormat Format = "table"
	// JSONFormat specifies JSON formatted output
	JSONFormat Format = "json"
	// YAMLFormat specifies YAML formatted output
	YAMLFormat Format = "yaml"
)

// Result is a generic interface for check results
type Result interface{}

// DNSResolverCheckResult represents DNS resolver configuration check results
type DNSResolverCheckResult struct {
	Timestamp          time.Time `json:"timestamp" yaml:"timestamp"`
	DomainNameSet      bool      `json:"domainNameSet" yaml:"domainNameSet"`
	SearchDomainsSet   bool      `json:"searchDomainsSet" yaml:"searchDomainsSet"`
	ServerAddressesSet bool      `json:"serverAddressesSet" yaml:"serverAddressesSet"`
}

// DNSResolverPingCheckResult represents DNS resolver connectivity check results
type DNSResolverPingCheckResult struct {
	Timestamp time.Time                    `json:"timestamp" yaml:"timestamp"`
	Resolvers []ResolverConnectivityResult `json:"resolvers" yaml:"resolvers"`
}

// ResolverConnectivityResult represents connectivity check for a single resolver
type ResolverConnectivityResult struct {
	ResolverIP    string `json:"resolverIP" yaml:"resolverIP"`
	NetInterface  string `json:"netInterface" yaml:"netInterface"`
	PingReachable bool   `json:"pingReachable" yaml:"pingReachable"`
	TCPReachable  bool   `json:"tcpReachable" yaml:"tcpReachable"`
	UDPReachable  bool   `json:"udpReachable" yaml:"udpReachable"`
}

// DNSResolverDigCheckResult represents DNS dig check results
type DNSResolverDigCheckResult struct {
	Timestamp time.Time        `json:"timestamp" yaml:"timestamp"`
	Results   []DigCheckResult `json:"results" yaml:"results"`
	Summary   map[string]int   `json:"summary" yaml:"summary"`
}

// DigCheckResult represents a single dig check result
type DigCheckResult struct {
	Hostname     string `json:"hostname" yaml:"hostname"`
	ResolverIP   string `json:"resolverIP" yaml:"resolverIP"`
	IsResolvable bool   `json:"isResolvable" yaml:"isResolvable"`
}

// VPNInterfaceCheckResult represents VPN interface reachability check results
type VPNInterfaceCheckResult struct {
	Timestamp              time.Time `json:"timestamp" yaml:"timestamp"`
	InterfaceCount         int       `json:"interfaceCount" yaml:"interfaceCount"`
	Interfaces             []string  `json:"interfaces" yaml:"interfaces"`
	HasTunInterface        bool      `json:"hasTunInterface" yaml:"hasTunInterface"`
	TunInterfaces          []string  `json:"tunInterfaces" yaml:"tunInterfaces"`
	AllInterfacesReachable bool      `json:"allInterfacesReachable" yaml:"allInterfacesReachable"`
}

// VPNRoutesCheckResult represents VPN routes check results
type VPNRoutesCheckResult struct {
	Timestamp           time.Time `json:"timestamp" yaml:"timestamp"`
	VPNInterface        string    `json:"vpnInterface" yaml:"vpnInterface"`
	RouteCount          int       `json:"routeCount" yaml:"routeCount"`
	MinRoutesRequired   int       `json:"minRoutesRequired" yaml:"minRoutesRequired"`
	HasSufficientRoutes bool      `json:"hasSufficientRoutes" yaml:"hasSufficientRoutes"`
}

// VPNConnectionStatusResult represents VPN connection status check results
type VPNConnectionStatusResult struct {
	Timestamp   time.Time `json:"timestamp" yaml:"timestamp"`
	IsConnected bool      `json:"isConnected" yaml:"isConnected"`
}

// ServerReachabilityCheckResult represents server reachability check results
type ServerReachabilityCheckResult struct {
	Timestamp     time.Time           `json:"timestamp" yaml:"timestamp"`
	Servers       []ServerCheckResult `json:"servers" yaml:"servers"`
	PingFailures  int                 `json:"pingFailures" yaml:"pingFailures"`
	ReachFailures int                 `json:"reachFailures" yaml:"reachFailures"`
}

// ServerCheckResult represents a single server check result
type ServerCheckResult struct {
	Host        string `json:"host" yaml:"host"`
	Service     string `json:"service" yaml:"service"`
	Reachable   bool   `json:"reachable" yaml:"reachable"`
	Performance string `json:"performance" yaml:"performance"`
}

// Print outputs the result in the specified format
func Print(format string, result interface{}) error {
	switch Format(format) {
	case JSONFormat:
		return PrintJSON(result)
	case YAMLFormat:
		return PrintYAML(result)
	case TableFormat:
		// Table format is handled by the calling function
		return nil
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// PrintJSON outputs the result as JSON
func PrintJSON(result interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// PrintYAML outputs the result as YAML
func PrintYAML(result interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	return encoder.Encode(result)
}
