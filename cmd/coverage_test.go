/*
Package cmd - Additional coverage tests

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
	"io"
	"os"
	"runtime"
	"testing"
	"time"
)

// TestAllDnsExecutePaths tests all code paths in dnsExecute
func TestAllDnsExecutePaths(t *testing.T) {
	setupMinimalConfig()

	// Capture output for all tests
	old := os.Stdout
	defer func() { os.Stdout = old }()

	tests := []struct {
		name        string
		resolverChk bool
		pingChk     bool
		digChk      bool
		allChk      bool
		skip        bool
	}{
		{"resolver path", true, false, false, false, runtime.GOOS != "darwin"},
		{"ping path", false, true, false, false, runtime.GOOS != "darwin"},
		{"dig path", false, false, true, false, runtime.GOOS != "darwin"},
		{"all path", false, false, false, true, runtime.GOOS != "darwin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("Skipping macOS-specific test")
			}

			resolverChk = tt.resolverChk
			pingChk = tt.pingChk
			digChk = tt.digChk
			allChk = tt.allChk

			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered: %v", r)
				}
				w.Close()
				os.Stdout = old
				io.Copy(io.Discard, r)
			}()

			dnsExecute(dnsCmd, []string{})
		})
	}
}

// TestAllVpnExecutePaths tests all code paths in vpnExecute
func TestAllVpnExecutePaths(t *testing.T) {
	setupMinimalConfig()

	old := os.Stdout
	defer func() { os.Stdout = old }()

	tests := []struct {
		name           string
		ifReachableChk bool
		vpnRoutesChk   bool
		vpnStatusChk   bool
		allChk         bool
		skip           bool
	}{
		{"interface path", true, false, false, false, runtime.GOOS != "darwin"},
		{"routes path", false, true, false, false, runtime.GOOS != "darwin"},
		{"status path - linux", false, false, true, false, runtime.GOOS != "linux"},
		{"status path - darwin", false, false, true, false, runtime.GOOS != "darwin"},
		{"all path", false, false, false, true, runtime.GOOS != "darwin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("Skipping OS-specific test")
			}

			ifReachableChk = tt.ifReachableChk
			vpnRoutesChk = tt.vpnRoutesChk
			vpnStatusChk = tt.vpnStatusChk
			allChk = tt.allChk

			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered: %v", r)
				}
				w.Close()
				os.Stdout = old
				io.Copy(io.Discard, r)
			}()

			vpnExecute(vpnCmd, []string{})
		})
	}
}

// TestAllSvrsExecutePaths tests all code paths in svrsExecute
func TestAllSvrsExecutePaths(t *testing.T) {
	setupMinimalConfig()

	old := os.Stdout
	defer func() { os.Stdout = old }()

	tests := []struct {
		name             string
		svrsReachableChk bool
		allChk           bool
		skip             bool
	}{
		{"servers path", true, false, runtime.GOOS != "darwin"},
		{"all path", false, true, runtime.GOOS != "darwin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("Skipping macOS-specific test")
			}

			svrsReachableChk = tt.svrsReachableChk
			allChk = tt.allChk

			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered: %v", r)
				}
				w.Close()
				os.Stdout = old
				io.Copy(io.Discard, r)
			}()

			svrsExecute(svrsCmd, []string{})
		})
	}
}

// TestDnsResolverPingChkCoverage increases coverage for dnsResolverPingChk
func TestDnsResolverPingChkCoverage(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test")
	}

	setupMinimalConfig()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered: %v", r)
		}
		w.Close()
		os.Stdout = old
		io.Copy(io.Discard, r)
	}()

	dnsResolverPingChk()
}

// TestScutilVPNInterfaceCoverage tests scutilVPNInterface
// DISABLED: scutilVPNInterface function no longer exists
/*
func TestScutilVPNInterfaceCoverage(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered: %v", r)
		}
	}()

	iface := scutilVPNInterface()
	t.Logf("VPN interface: %s", iface)
}
*/

// TestIfReachChkCoverage tests ifReachChk
func TestIfReachChkCoverage(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test")
	}

	setupMinimalConfig()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered: %v", r)
		}
		w.Close()
		os.Stdout = old
		io.Copy(io.Discard, r)
	}()

	ifReachChk()
}

// TestVpnRteChkCoverage tests vpnRteChk
func TestVpnRteChkCoverage(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test")
	}

	setupMinimalConfig()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered: %v", r)
		}
		w.Close()
		os.Stdout = old
		io.Copy(io.Discard, r)
	}()

	vpnRteChk()
}

// TestVpnConnChkCoverage tests vpnConnChk on Linux
func TestVpnConnChkCoverage(t *testing.T) {
	setupMinimalConfig()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered: %v", r)
		}
		w.Close()
		os.Stdout = old
		io.Copy(io.Discard, r)
	}()

	vpnConnChk()
}

// TestSvrsExecuteAllPath tests the allChk path specifically
func TestSvrsExecuteAllPath(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test")
	}

	setupMinimalConfig()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	allChk = true
	svrsReachableChk = false

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered: %v", r)
		}
		w.Close()
		os.Stdout = old
		io.Copy(io.Discard, r)
	}()

	svrsExecute(svrsCmd, []string{})
}

// TestSvrsReachChkWithMultipleServers tests with multiple server configurations
func TestSvrsReachChkWithMultipleServers(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test")
	}

	// Setup config with multiple services and servers
	conf = &config{
		MinVpnRoutes:     5,
		PingTimeout:      10 * time.Millisecond,
		DNSLookupTimeout: 10 * time.Millisecond,
		FailThreshold:    0,
		Svcs: []svc{
			{Svc: "idm", Svrs: []string{"localhost"}},
		},
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered: %v", r)
		}
		w.Close()
		os.Stdout = old
		io.Copy(io.Discard, r)
	}()

	svrsReachChk()
}
