/*
Package cmd - Comprehensive tests for svrsReachChkWithDeps

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
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-ping/ping"
)

// ========== Mock Implementations ==========

// mockDNSResolver allows controlling DNS lookup behavior in tests
type mockDNSResolver struct {
	hosts map[string][]string // Map of hostname to IP addresses
	err   error               // Error to return
}

func (m *mockDNSResolver) LookupHost(ctx context.Context, host string) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	if ips, ok := m.hosts[host]; ok {
		return ips, nil
	}
	return nil, errors.New("host not found")
}

// mockPinger allows controlling ping behavior in tests
type mockPinger struct {
	timeout    time.Duration
	runErr     error
	packetLoss float64
	packetsRcv int
	avgRtt     time.Duration
}

func (m *mockPinger) SetTimeout(duration time.Duration) {
	m.timeout = duration
}

func (m *mockPinger) Run() error {
	return m.runErr
}

func (m *mockPinger) Statistics() *ping.Statistics {
	return &ping.Statistics{
		PacketLoss:  m.packetLoss,
		PacketsRecv: m.packetsRcv,
		AvgRtt:      m.avgRtt,
	}
}

// mockBraceExpander allows controlling brace expansion in tests
type mockBraceExpander struct {
	expansions map[string][]string // Map of pattern to expanded results
}

func (m *mockBraceExpander) Expand(pattern string) []string {
	if expansions, ok := m.expansions[pattern]; ok {
		return expansions
	}
	// Default: return pattern as-is
	return []string{pattern}
}

// ========== Test Helpers ==========

func setupTestConfigWithServers() {
	conf = &config{
		DNSLookupTimeout: 1000,
		PingTimeout:      2000,
		FailThreshold:    3,
		Svcs: []svc{
			{
				Svc:  "web",
				Svrs: []string{"server1.example.com", "server2.example.com"},
			},
			{
				Svc:  "db",
				Svrs: []string{"db{1..3}.example.com"},
			},
		},
	}
}

// ========== Tests ==========

func TestSvrsReachChkWithDeps_DNSLookupFailure(t *testing.T) {
	setupTestConfigWithServers()
	outputFormat = "json" // Suppress table output

	mockResolver := &mockDNSResolver{
		err: errors.New("DNS lookup failed"),
	}

	mockExpander := &mockBraceExpander{
		expansions: map[string][]string{
			"server1.example.com":  {"server1.example.com"},
			"server2.example.com":  {"server2.example.com"},
			"db{1..3}.example.com": {"db1.example.com", "db2.example.com", "db3.example.com"},
		},
	}

	// Create a mock pinger factory that should never be called
	originalNewPinger := NewPinger
	defer func() { NewPinger = originalNewPinger }()

	NewPinger = func(host string) (Pinger, error) {
		t.Error("NewPinger should not be called when DNS lookup fails")
		return nil, errors.New("should not be called")
	}

	// Should not panic
	svrsReachChkWithDeps(mockResolver, mockExpander)
}

func TestSvrsReachChkWithDeps_EmptyDNSResult(t *testing.T) {
	setupTestConfigWithServers()
	outputFormat = "json"

	mockResolver := &mockDNSResolver{
		hosts: map[string][]string{
			"server1.example.com": {}, // Empty result
		},
	}

	mockExpander := &mockBraceExpander{
		expansions: map[string][]string{
			"server1.example.com":  {"server1.example.com"},
			"server2.example.com":  {"server2.example.com"},
			"db{1..3}.example.com": {"db1.example.com"},
		},
	}

	originalNewPinger := NewPinger
	defer func() { NewPinger = originalNewPinger }()

	NewPinger = func(host string) (Pinger, error) {
		t.Error("NewPinger should not be called when DNS returns empty result")
		return nil, errors.New("should not be called")
	}

	svrsReachChkWithDeps(mockResolver, mockExpander)
}

func TestSvrsReachChkWithDeps_PingerCreationFailure(t *testing.T) {
	setupTestConfigWithServers()
	outputFormat = "json"

	mockResolver := &mockDNSResolver{
		hosts: map[string][]string{
			"server1.example.com": {"192.168.1.1"},
		},
	}

	mockExpander := &mockBraceExpander{
		expansions: map[string][]string{
			"server1.example.com":  {"server1.example.com"},
			"server2.example.com":  {"server2.example.com"},
			"db{1..3}.example.com": {"db1.example.com"},
		},
	}

	originalNewPinger := NewPinger
	defer func() { NewPinger = originalNewPinger }()

	NewPinger = func(host string) (Pinger, error) {
		return nil, errors.New("failed to create pinger")
	}

	// Should not panic
	svrsReachChkWithDeps(mockResolver, mockExpander)
}

func TestSvrsReachChkWithDeps_PingTimeout(t *testing.T) {
	setupTestConfigWithServers()
	outputFormat = "json"

	mockResolver := &mockDNSResolver{
		hosts: map[string][]string{
			"server1.example.com": {"192.168.1.1"},
		},
	}

	mockExpander := &mockBraceExpander{
		expansions: map[string][]string{
			"server1.example.com":  {"server1.example.com"},
			"server2.example.com":  {"server2.example.com"},
			"db{1..3}.example.com": {"db1.example.com"},
		},
	}

	originalNewPinger := NewPinger
	defer func() { NewPinger = originalNewPinger }()

	NewPinger = func(host string) (Pinger, error) {
		return &mockPinger{
			runErr:     nil,
			packetLoss: 100.0, // Total packet loss
			packetsRcv: 0,     // No packets received
			avgRtt:     0,
		}, nil
	}

	// Should not panic
	svrsReachChkWithDeps(mockResolver, mockExpander)
}

func TestSvrsReachChkWithDeps_SuccessfulPing(t *testing.T) {
	setupTestConfigWithServers()
	outputFormat = "json"

	mockResolver := &mockDNSResolver{
		hosts: map[string][]string{
			"server1.example.com": {"192.168.1.1"},
			"server2.example.com": {"192.168.1.2"},
			"db1.example.com":     {"192.168.1.3"},
			"db2.example.com":     {"192.168.1.4"},
			"db3.example.com":     {"192.168.1.5"},
		},
	}

	mockExpander := &mockBraceExpander{
		expansions: map[string][]string{
			"server1.example.com":  {"server1.example.com"},
			"server2.example.com":  {"server2.example.com"},
			"db{1..3}.example.com": {"db1.example.com", "db2.example.com", "db3.example.com"},
		},
	}

	originalNewPinger := NewPinger
	defer func() { NewPinger = originalNewPinger }()

	NewPinger = func(host string) (Pinger, error) {
		return &mockPinger{
			runErr:     nil,
			packetLoss: 0.0,
			packetsRcv: 4,
			avgRtt:     10 * time.Millisecond,
		}, nil
	}

	// Should not panic
	svrsReachChkWithDeps(mockResolver, mockExpander)
}

func TestSvrsReachChkWithDeps_PartialPacketLoss(t *testing.T) {
	setupTestConfigWithServers()
	outputFormat = "json"

	mockResolver := &mockDNSResolver{
		hosts: map[string][]string{
			"server1.example.com": {"192.168.1.1"},
		},
	}

	mockExpander := &mockBraceExpander{
		expansions: map[string][]string{
			"server1.example.com":  {"server1.example.com"},
			"server2.example.com":  {"server2.example.com"},
			"db{1..3}.example.com": {"db1.example.com"},
		},
	}

	originalNewPinger := NewPinger
	defer func() { NewPinger = originalNewPinger }()

	NewPinger = func(host string) (Pinger, error) {
		return &mockPinger{
			runErr:     nil,
			packetLoss: 50.0, // 50% packet loss
			packetsRcv: 2,    // Some packets still received
			avgRtt:     20 * time.Millisecond,
		}, nil
	}

	// Should not panic
	svrsReachChkWithDeps(mockResolver, mockExpander)
}

func TestSvrsReachChkWithDeps_BraceExpansion(t *testing.T) {
	setupTestConfigWithServers()
	outputFormat = "json"

	mockResolver := &mockDNSResolver{
		hosts: map[string][]string{
			"db1.example.com": {"192.168.1.1"},
			"db2.example.com": {"192.168.1.2"},
			"db3.example.com": {"192.168.1.3"},
		},
	}

	mockExpander := &mockBraceExpander{
		expansions: map[string][]string{
			"server1.example.com":  {"server1.example.com"},
			"server2.example.com":  {"server2.example.com"},
			"db{1..3}.example.com": {"db1.example.com", "db2.example.com", "db3.example.com"},
		},
	}

	originalNewPinger := NewPinger
	defer func() { NewPinger = originalNewPinger }()

	pingedHosts := make(map[string]bool)
	NewPinger = func(host string) (Pinger, error) {
		pingedHosts[host] = true
		return &mockPinger{
			runErr:     nil,
			packetLoss: 0.0,
			packetsRcv: 4,
			avgRtt:     10 * time.Millisecond,
		}, nil
	}

	svrsReachChkWithDeps(mockResolver, mockExpander)

	// Verify all expanded hosts were pinged
	expectedHosts := []string{"db1.example.com", "db2.example.com", "db3.example.com"}
	for _, host := range expectedHosts {
		if !pingedHosts[host] {
			t.Errorf("Expected host %s to be pinged", host)
		}
	}
}

func TestSvrsReachChkWithDeps_MixedResults(t *testing.T) {
	setupTestConfigWithServers()
	outputFormat = "json"

	mockResolver := &mockDNSResolver{
		hosts: map[string][]string{
			"server1.example.com": {"192.168.1.1"},
			"server2.example.com": {"192.168.1.2"},
			// db1 will fail DNS
			"db2.example.com": {"192.168.1.4"},
			// db3 will fail DNS
		},
	}

	mockExpander := &mockBraceExpander{
		expansions: map[string][]string{
			"server1.example.com":  {"server1.example.com"},
			"server2.example.com":  {"server2.example.com"},
			"db{1..3}.example.com": {"db1.example.com", "db2.example.com", "db3.example.com"},
		},
	}

	originalNewPinger := NewPinger
	defer func() { NewPinger = originalNewPinger }()

	NewPinger = func(host string) (Pinger, error) {
		// server1 succeeds, server2 has packet loss
		if host == "server1.example.com" || host == "db2.example.com" {
			return &mockPinger{
				runErr:     nil,
				packetLoss: 0.0,
				packetsRcv: 4,
				avgRtt:     10 * time.Millisecond,
			}, nil
		}
		return &mockPinger{
			runErr:     nil,
			packetLoss: 100.0,
			packetsRcv: 0,
			avgRtt:     0,
		}, nil
	}

	// Should not panic with mixed results
	svrsReachChkWithDeps(mockResolver, mockExpander)
}

func TestSvrsReachChkWithDeps_EmptyConfig(t *testing.T) {
	conf = &config{
		DNSLookupTimeout: 1000,
		PingTimeout:      2000,
		FailThreshold:    3,
		Svcs:             []svc{}, // Empty services
	}
	outputFormat = "json"

	mockResolver := &mockDNSResolver{
		hosts: map[string][]string{},
	}

	mockExpander := &mockBraceExpander{
		expansions: map[string][]string{},
	}

	// Should handle empty config gracefully
	svrsReachChkWithDeps(mockResolver, mockExpander)
}

func TestSvrsReachChkWithDeps_TableOutput(t *testing.T) {
	setupTestConfigWithServers()
	outputFormat = "table" // Test table output path

	mockResolver := &mockDNSResolver{
		hosts: map[string][]string{
			"server1.example.com": {"192.168.1.1"},
		},
	}

	mockExpander := &mockBraceExpander{
		expansions: map[string][]string{
			"server1.example.com":  {"server1.example.com"},
			"server2.example.com":  {"server2.example.com"},
			"db{1..3}.example.com": {"db1.example.com"},
		},
	}

	originalNewPinger := NewPinger
	defer func() { NewPinger = originalNewPinger }()

	NewPinger = func(host string) (Pinger, error) {
		return &mockPinger{
			runErr:     nil,
			packetLoss: 0.0,
			packetsRcv: 4,
			avgRtt:     10 * time.Millisecond,
		}, nil
	}

	// Should not panic with table output
	svrsReachChkWithDeps(mockResolver, mockExpander)
}

func TestSvrsReachChk_CallsWithDeps(t *testing.T) {
	setupTestConfigWithServers()
	outputFormat = "json"

	// Override factory functions
	originalNewDNSResolver := NewDNSResolver
	originalNewBraceExpander := NewBraceExpander
	originalNewPinger := NewPinger
	defer func() {
		NewDNSResolver = originalNewDNSResolver
		NewBraceExpander = originalNewBraceExpander
		NewPinger = originalNewPinger
	}()

	resolverCalled := false
	expanderCalled := false

	NewDNSResolver = func() DNSResolver {
		resolverCalled = true
		return &mockDNSResolver{
			hosts: map[string][]string{
				"server1.example.com": {"192.168.1.1"},
			},
		}
	}

	NewBraceExpander = func() BraceExpander {
		expanderCalled = true
		return &mockBraceExpander{
			expansions: map[string][]string{
				"server1.example.com":  {"server1.example.com"},
				"server2.example.com":  {"server2.example.com"},
				"db{1..3}.example.com": {"db1.example.com"},
			},
		}
	}

	NewPinger = func(host string) (Pinger, error) {
		return &mockPinger{
			runErr:     nil,
			packetLoss: 0.0,
			packetsRcv: 4,
			avgRtt:     10 * time.Millisecond,
		}, nil
	}

	svrsReachChk()

	if !resolverCalled {
		t.Error("NewDNSResolver should have been called")
	}
	if !expanderCalled {
		t.Error("NewBraceExpander should have been called")
	}
}
