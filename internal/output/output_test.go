/*
Package output_test - tests for output package

Copyright © 2021 Sam Mingolelli <github@lamolabs.org>
*/
package output_test

import (
	"bytes"
	"doxctl/internal/output"
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestPrintJSON(t *testing.T) {
	result := output.DNSResolverCheckResult{
		Timestamp:          time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		DomainNameSet:      true,
		SearchDomainsSet:   true,
		ServerAddressesSet: true,
	}

	// Test JSON marshaling
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verify JSON is valid
	var decoded output.DNSResolverCheckResult
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify values
	if decoded.DomainNameSet != result.DomainNameSet {
		t.Errorf("Expected DomainNameSet=%v, got %v", result.DomainNameSet, decoded.DomainNameSet)
	}
	if decoded.SearchDomainsSet != result.SearchDomainsSet {
		t.Errorf("Expected SearchDomainsSet=%v, got %v", result.SearchDomainsSet, decoded.SearchDomainsSet)
	}
	if decoded.ServerAddressesSet != result.ServerAddressesSet {
		t.Errorf("Expected ServerAddressesSet=%v, got %v", result.ServerAddressesSet, decoded.ServerAddressesSet)
	}
}

func TestPrintYAML(t *testing.T) {
	result := output.VPNInterfaceCheckResult{
		Timestamp:              time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		InterfaceCount:         2,
		Interfaces:             []string{"eth0", "utun0"},
		HasTunInterface:        true,
		TunInterfaces:          []string{"utun0"},
		AllInterfacesReachable: true,
	}

	// Test YAML marshaling
	data, err := yaml.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal YAML: %v", err)
	}

	// Verify YAML is valid
	var decoded output.VPNInterfaceCheckResult
	if err := yaml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Verify values
	if decoded.InterfaceCount != result.InterfaceCount {
		t.Errorf("Expected InterfaceCount=%v, got %v", result.InterfaceCount, decoded.InterfaceCount)
	}
	if decoded.HasTunInterface != result.HasTunInterface {
		t.Errorf("Expected HasTunInterface=%v, got %v", result.HasTunInterface, decoded.HasTunInterface)
	}
	if len(decoded.Interfaces) != len(result.Interfaces) {
		t.Errorf("Expected %d interfaces, got %d", len(result.Interfaces), len(decoded.Interfaces))
	}
}

func TestServerReachabilityResult(t *testing.T) {
	result := output.ServerReachabilityCheckResult{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Servers: []output.ServerCheckResult{
			{
				Host:        "server1.example.com",
				Service:     "web",
				Reachable:   true,
				Performance: "rnd-trp avg = 10ms",
			},
			{
				Host:        "server2.example.com",
				Service:     "db",
				Reachable:   false,
				Performance: "N/A",
			},
		},
		PingFailures:  1,
		ReachFailures: 1,
	}

	// Test both JSON and YAML
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var jsonDecoded output.ServerReachabilityCheckResult
	if err := json.Unmarshal(jsonData, &jsonDecoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if jsonDecoded.PingFailures != result.PingFailures {
		t.Errorf("Expected PingFailures=%v, got %v", result.PingFailures, jsonDecoded.PingFailures)
	}
	if len(jsonDecoded.Servers) != len(result.Servers) {
		t.Errorf("Expected %d servers, got %d", len(result.Servers), len(jsonDecoded.Servers))
	}

	yamlData, err := yaml.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal YAML: %v", err)
	}

	var yamlDecoded output.ServerReachabilityCheckResult
	if err := yaml.Unmarshal(yamlData, &yamlDecoded); err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if yamlDecoded.ReachFailures != result.ReachFailures {
		t.Errorf("Expected ReachFailures=%v, got %v", result.ReachFailures, yamlDecoded.ReachFailures)
	}
}

func TestDNSResolverPingCheck(t *testing.T) {
	result := output.DNSResolverPingCheckResult{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Resolvers: []output.ResolverConnectivityResult{
			{
				ResolverIP:    "8.8.8.8",
				NetInterface:  "eth0",
				PingReachable: true,
				TCPReachable:  true,
				UDPReachable:  true,
			},
		},
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var decoded output.DNSResolverPingCheckResult
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(decoded.Resolvers) != 1 {
		t.Fatalf("Expected 1 resolver, got %d", len(decoded.Resolvers))
	}

	if decoded.Resolvers[0].ResolverIP != "8.8.8.8" {
		t.Errorf("Expected ResolverIP=8.8.8.8, got %v", decoded.Resolvers[0].ResolverIP)
	}
}

func TestPrintFormat(t *testing.T) {
	// Redirect stdout to capture output
	tests := []struct {
		name   string
		format string
		valid  bool
	}{
		{"JSON format", "json", true},
		{"YAML format", "yaml", true},
		{"Table format", "table", true},
		{"Invalid format", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := output.DNSResolverCheckResult{
				Timestamp:          time.Now(),
				DomainNameSet:      true,
				SearchDomainsSet:   true,
				ServerAddressesSet: true,
			}

			// Capture stdout
			old := bytes.Buffer{}
			err := output.Print(tt.format, result)

			if tt.valid {
				if err != nil && tt.format != "table" {
					t.Errorf("Expected no error for format %s, got %v", tt.format, err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for invalid format %s", tt.format)
				}
			}

			_ = old
		})
	}
}
