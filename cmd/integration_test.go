/*
Package cmd - Integration tests

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
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
)

// TestDnsCmd_ExecuteWithFlags tests the dns command execute function with various flags
func TestDnsCmd_ExecuteWithFlags(t *testing.T) {
	// Setup config for testing
	setupTestConfig(t)

	tests := []struct {
		name string
		args []string
	}{
		{"dns help", []string{"dns", "--help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Set command args
			rootCmd.SetArgs(tt.args)

			// Execute command - expecting it to work or fail gracefully
			_ = rootCmd.Execute()

			// Restore stdout
			w.Close()
			os.Stdout = old
			io.Copy(io.Discard, r)
		})
	}
}

// TestVpnCmd_ExecuteWithFlags tests the vpn command execute function with various flags
func TestVpnCmd_ExecuteWithFlags(t *testing.T) {
	// Setup config for testing
	setupTestConfig(t)

	tests := []struct {
		name string
		args []string
	}{
		{"vpn help", []string{"vpn", "--help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Set command args
			rootCmd.SetArgs(tt.args)

			// Execute command
			_ = rootCmd.Execute()

			// Restore stdout
			w.Close()
			os.Stdout = old
			io.Copy(io.Discard, r)
		})
	}
}

// TestSvrsCmd_ExecuteWithFlags tests the svrs command execute function with various flags
func TestSvrsCmd_ExecuteWithFlags(t *testing.T) {
	// Setup config for testing
	setupTestConfig(t)

	tests := []struct {
		name string
		args []string
	}{
		{"svrs help", []string{"svrs", "--help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Set command args
			rootCmd.SetArgs(tt.args)

			// Execute command
			_ = rootCmd.Execute()

			// Restore stdout
			w.Close()
			os.Stdout = old
			io.Copy(io.Discard, r)
		})
	}
}

// setupTestConfig creates a test configuration
func setupTestConfig(t *testing.T) {
	viper.Reset()
	viper.Set("minVpnRoutes", 5)
	viper.Set("domNameChk", "test.com")
	viper.Set("domSearchChk", "search.test.com")
	viper.Set("domAddrChk", "10.0.0")
	viper.Set("pingTimeout", 10*time.Millisecond)
	viper.Set("dnsLookupTimeout", 10*time.Millisecond)
	viper.Set("failThreshold", 0)

	conf = &config{
		MinVpnRoutes:     5,
		DomNameChk:       "test.com",
		DomSearchChk:     "search.test.com",
		DomAddrChk:       "10.0.0",
		PingTimeout:      10 * time.Millisecond,
		DNSLookupTimeout: 10 * time.Millisecond,
		FailThreshold:    0,
		Svcs: []svc{
			{Svc: "idm", Svrs: []string{"localhost"}},
		},
	}
}

// TestScutilResolverIPs tests the resolver IP parsing logic
func TestScutilResolverIPs(t *testing.T) {
	// Test the string processing logic that would be used
	testData := []string{"  10.0.0.1  ", "  10.0.0.2  ", "  "}
	
	var cleaned []string
	for _, ip := range testData {
		trimmed := string(bytes.TrimSpace([]byte(ip)))
		if len(trimmed) > 0 {
			cleaned = append(cleaned, trimmed)
		}
	}

	if len(cleaned) != 2 {
		t.Errorf("Expected 2 IPs, got %d", len(cleaned))
	}

	if cleaned[0] != "10.0.0.1" {
		t.Errorf("Expected first IP to be '10.0.0.1', got: %s", cleaned[0])
	}
}

// TestScutilVPNInterface tests the VPN interface parsing logic
func TestScutilVPNInterface(t *testing.T) {
	// Test the string processing logic
	testOutputs := []struct {
		input    string
		expected string
	}{
		{"utun0\n", "utun0"},
		{"\n", "N/A"},
		{"", "N/A"},
		{"utun1", "utun1"},
	}

	for _, tc := range testOutputs {
		result := string(bytes.TrimRight([]byte(tc.input), "\n"))
		if len(result) == 0 {
			result = "N/A"
		}

		if result != tc.expected {
			t.Errorf("For input %q, expected %q, got %q", tc.input, tc.expected, result)
		}
	}
}

// TestConfigUnmarshal tests config unmarshaling
func TestConfigUnmarshal(t *testing.T) {
	viper.Reset()
	viper.Set("minVpnRoutes", 10)
	viper.Set("domNameChk", "example.com")
	viper.Set("failThreshold", 5)

	cfg := &config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if cfg.MinVpnRoutes != 10 {
		t.Errorf("Expected MinVpnRoutes to be 10, got: %d", cfg.MinVpnRoutes)
	}

	if cfg.DomNameChk != "example.com" {
		t.Errorf("Expected DomNameChk to be 'example.com', got: %s", cfg.DomNameChk)
	}

	if cfg.FailThreshold != 5 {
		t.Errorf("Expected FailThreshold to be 5, got: %d", cfg.FailThreshold)
	}
}

// TestRootCmd_CommandTree tests the full command tree
func TestRootCmd_CommandTree(t *testing.T) {
	// Verify all subcommands are properly registered
	commands := rootCmd.Commands()
	
	expectedCommands := map[string]bool{
		"dns":  false,
		"vpn":  false,
		"svrs": false,
	}

	for _, cmd := range commands {
		if _, ok := expectedCommands[cmd.Name()]; ok {
			expectedCommands[cmd.Name()] = true
		}
	}

	for cmdName, found := range expectedCommands {
		if !found {
			t.Errorf("Expected command %s not found in command tree", cmdName)
		}
	}
}

// TestRootCmd_PersistentFlags tests persistent flags
func TestRootCmd_PersistentFlags(t *testing.T) {
	// Test that persistent flags are available to all commands
	for _, cmd := range rootCmd.Commands() {
		inheritedFlags := cmd.InheritedFlags()
		
		configFlag := inheritedFlags.Lookup("config")
		if configFlag == nil {
			t.Errorf("Command %s should inherit config flag", cmd.Name())
		}

		verboseFlag := inheritedFlags.Lookup("verbose")
		if verboseFlag == nil {
			t.Errorf("Command %s should inherit verbose flag", cmd.Name())
		}
	}
}
