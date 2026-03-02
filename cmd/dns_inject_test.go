/*
Package cmd - Comprehensive tests for DNS functions

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
package cmd

import (
	"errors"
	"testing"
)

// ========== Mock Implementations ==========

// mockCommandExecutor allows controlling command execution in tests
type mockCommandExecutor struct {
	commands map[string][]byte                       // Map of command to output
	err      error                                   // Error to return
	execFunc func(string, ...string) ([]byte, error) // Optional custom execution function
}

func (m *mockCommandExecutor) Execute(name string, args ...string) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}

	// If a custom function is provided, use it
	if m.execFunc != nil {
		return m.execFunc(name, args...)
	}

	// Create a key from the command and args
	key := name
	for _, arg := range args {
		key += " " + arg
	}

	if output, ok := m.commands[key]; ok {
		return output, nil
	}

	return []byte(""), nil
}

// mockFileReader allows controlling file reading in tests
type mockFileReader struct {
	files map[string][]byte // Map of filename to content
	err   error             // Error to return
}

func (m *mockFileReader) ReadFile(filename string) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}

	if content, ok := m.files[filename]; ok {
		return content, nil
	}

	return nil, errors.New("file not found")
}

// ========== Test Helpers ==========

func setupDNSTestConfig() {
	conf = &config{
		DNSLookupTimeout: 1000,
		PingTimeout:      2000,
		FailThreshold:    3,
		DomNameChk:       "example.com",
		DomSearchChk:     "example",
		DomAddrChk:       "10\\.20\\.",
		Svcs: []svc{
			{
				Svc:  "idm",
				Svrs: []string{"server1.example.com", "server2.example.com"},
			},
		},
	}
}

// ========== High-level Function Tests ==========
// Platform-specific tests are in dns_darwin_test.go and dns_linux_test.go

func TestDNSResolverChkWithDeps_Success(t *testing.T) {
	setupDNSTestConfig()
	outputFormat = "json"

	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{},
	}

	mockFile := &mockFileReader{
		files: map[string][]byte{
			"/etc/resolv.conf": []byte(`nameserver 10.20.1.1
search example.com
domain example.com
`),
		},
	}

	// Should not panic
	dnsResolverChkWithDeps(mockExec, mockFile)
}

func TestDNSResolverChkWithDeps_TableOutput(t *testing.T) {
	setupDNSTestConfig()
	outputFormat = "table"

	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{},
	}

	mockFile := &mockFileReader{
		files: map[string][]byte{
			"/etc/resolv.conf": []byte(`nameserver 10.20.1.1
search example.com
domain example.com
`),
		},
	}

	// Should not panic with table output
	dnsResolverChkWithDeps(mockExec, mockFile)
}

func TestDNSResolverPingChkWithDeps_Success(t *testing.T) {
	setupDNSTestConfig()
	outputFormat = "json"

	mockExec := &mockCommandExecutor{
		execFunc: func(name string, args ...string) ([]byte, error) {
			if name == "bash" && len(args) >= 2 && args[0] == "-c" {
				cmd := args[1]
				// Scutil call (macOS)
				if contains(cmd, "scutil") && !contains(cmd, "echo") {
					return []byte(`<dictionary> {
  ServerAddresses : <array> {
    0 : 10.20.1.1
  }
}`), nil
				}
				// Grep for server addresses (macOS)
				if contains(cmd, "grep") && contains(cmd, "ServerAddresses") {
					return []byte("       10.20.1.1\n"), nil
				}
			}
			// Linux: ip route get command
			if name == "ip" && len(args) >= 3 && args[0] == "route" && args[1] == "get" {
				// Mock output: "10.20.1.1 via 192.168.1.1 dev eth0 src 192.168.1.100"
				return []byte("10.20.1.1 via 192.168.1.1 dev eth0 src 192.168.1.100\n"), nil
			}
			return []byte(""), nil
		},
	}

	mockFile := &mockFileReader{
		files: map[string][]byte{
			"/etc/resolv.conf": []byte("nameserver 10.20.1.1\n"),
		},
	}

	pingedHosts := make(map[string]bool)
	mockPingerFactory := func(host string) (Pinger, error) {
		pingedHosts[host] = true
		return &mockPinger{
			runErr:     nil,
			packetLoss: 0.0,
			packetsRcv: 4,
			avgRtt:     10,
		}, nil
	}

	// Should not panic
	dnsResolverPingChkWithDeps(mockExec, mockFile, mockPingerFactory)

	if !pingedHosts["10.20.1.1"] {
		t.Error("Expected to ping 10.20.1.1")
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestDNSResolverPingChkWithDeps_PingFailure(t *testing.T) {
	setupDNSTestConfig()
	outputFormat = "json"

	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{},
	}

	mockFile := &mockFileReader{
		files: map[string][]byte{
			"/etc/resolv.conf": []byte(`nameserver 10.20.1.1
`),
		},
	}

	mockPingerFactory := func(host string) (Pinger, error) {
		return nil, errors.New("ping failed")
	}

	// Should not panic even when ping fails
	dnsResolverPingChkWithDeps(mockExec, mockFile, mockPingerFactory)
}

func TestDNSResolverPingChkWithDeps_NoResolvers(t *testing.T) {
	setupDNSTestConfig()
	outputFormat = "json"

	mockExec := &mockCommandExecutor{
		execFunc: func(name string, args ...string) ([]byte, error) {
			// Return empty results for all commands
			return []byte(""), nil
		},
	}

	mockFile := &mockFileReader{
		files: map[string][]byte{},
	}

	pingerCalled := false
	mockPingerFactory := func(host string) (Pinger, error) {
		// Empty string might be passed due to strings.Split("", "\n") returning [""]
		if host != "" {
			pingerCalled = true
			t.Errorf("Should not attempt to ping non-empty host when no resolvers: %s", host)
		}
		return nil, errors.New("no pinger for empty host")
	}

	// Should not panic even with no resolvers
	dnsResolverPingChkWithDeps(mockExec, mockFile, mockPingerFactory)

	if pingerCalled {
		t.Error("Pinger was called for a real host when no resolvers were available")
	}
}

func TestDNSResolverDigChkWithDeps_Success(t *testing.T) {
	setupDNSTestConfig()
	outputFormat = "json"

	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{},
	}

	mockFile := &mockFileReader{
		files: map[string][]byte{},
	}

	mockExpander := &mockBraceExpander{
		expansions: map[string][]string{
			"server1.example.com": {"server1.example.com"},
			"server2.example.com": {"server2.example.com"},
		},
	}

	// Should not panic
	dnsResolverDigChkWithDeps(mockExec, mockFile, mockExpander)
}

func TestDNSResolverDigChkWithDeps_TableOutput(t *testing.T) {
	setupDNSTestConfig()
	outputFormat = "table"

	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{},
	}

	mockFile := &mockFileReader{
		files: map[string][]byte{
			"/etc/resolv.conf": []byte(`nameserver 10.20.1.1
`),
		},
	}

	mockExpander := &mockBraceExpander{
		expansions: map[string][]string{
			"server1.example.com": {"server1.example.com"},
			"server2.example.com": {"server2.example.com"},
		},
	}

	// Should not panic with table output
	dnsResolverDigChkWithDeps(mockExec, mockFile, mockExpander)
}

func TestDNSResolverDigChkWithDeps_NoResolvers(t *testing.T) {
	setupDNSTestConfig()
	outputFormat = "json"

	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{},
	}

	mockFile := &mockFileReader{
		files: map[string][]byte{
			"/etc/resolv.conf": []byte(`# No nameservers
`),
		},
	}

	mockExpander := &mockBraceExpander{
		expansions: map[string][]string{
			"server1.example.com": {"server1.example.com"},
		},
	}

	// Should handle empty resolver list gracefully
	dnsResolverDigChkWithDeps(mockExec, mockFile, mockExpander)
}
