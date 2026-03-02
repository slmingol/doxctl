//go:build darwin || linux
// +build darwin linux

package cmd

import (
	"testing"
)

// TestGetResolverIPs tests that getResolverIPs returns a slice (could be empty)
func TestGetResolverIPs(t *testing.T) {
	// This test just ensures the function can be called without panic
	// It may return an empty slice if no resolvers are configured
	ips := getResolverIPs()
	if ips == nil {
		t.Error("getResolverIPs should return a non-nil slice")
	}
}

// TestGetVPNInterface tests that getVPNInterface returns a string
func TestGetVPNInterface(t *testing.T) {
	// This test just ensures the function can be called without panic
	// It may return "N/A" if no VPN interface is detected
	iface := getVPNInterface()
	if iface == "" {
		t.Error("getVPNInterface should return a non-empty string")
	}
}

// TestGetDNSConfig tests that getDNSConfig returns three strings
func TestGetDNSConfig(t *testing.T) {
	// This test just ensures the function can be called without panic
	domain, search, servers := getDNSConfig()

	if domain == "" {
		t.Error("getDNSConfig should return a non-empty domain string")
	}
	if search == "" {
		t.Error("getDNSConfig should return a non-empty search string")
	}
	if servers == "" {
		t.Error("getDNSConfig should return a non-empty servers string")
	}
}
