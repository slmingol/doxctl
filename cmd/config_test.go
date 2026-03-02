/*
Package cmd - ...

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
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

// TestConfigCmd_Initialization tests the config command structure
func TestConfigCmd_Initialization(t *testing.T) {
	if configCmd == nil {
		t.Fatal("configCmd should not be nil")
	}

	if configCmd.Use != "config" {
		t.Errorf("Expected Use to be 'config', got: %s", configCmd.Use)
	}

	if !configCmd.HasSubCommands() {
		t.Error("configCmd should have subcommands")
	}

	// Check for expected subcommands
	expectedCmds := []string{"validate", "show", "init"}
	for _, expectedCmd := range expectedCmds {
		found := false
		for _, cmd := range configCmd.Commands() {
			if cmd.Name() == expectedCmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' not found", expectedCmd)
		}
	}
}

// TestConfigValidateCmd_ValidConfig tests validating a valid config file
func TestConfigValidateCmd_ValidConfig(t *testing.T) {
	// Create a temporary valid config file
	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	validConfig := `
domainName: "example.local"
sites:
  - lab1
  - rdu1
wellKnownSvcs:
  - svc: idm
    svrs:
      - idm-01a.lab1.example.local
      - idm-01b.lab1.example.local
minVpnRoutes: 5
pingTimeout: 250
dnsLookupTimeout: 100
failThreshold: 5
`
	if _, err := tmpFile.WriteString(validConfig); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	// Set up viper to use this config
	viper.Reset()
	viper.SetConfigFile(tmpFile.Name())

	// Test the validation logic directly
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	c := &config{}
	err = viper.Unmarshal(c)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	c.setDefaults()
	err = c.Validate()
	if err != nil {
		t.Errorf("Valid config should pass validation, got error: %v", err)
	}
}

// TestConfigValidateCmd_InvalidConfig tests validating an invalid config
func TestConfigValidateCmd_InvalidConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		expectError string
	}{
		{
			name: "missing domainName",
			config: `
sites:
  - lab1
wellKnownSvcs:
  - svc: idm
    svrs:
      - idm-01a.lab1.example.local
`,
			expectError: "domainName is required",
		},
		{
			name: "missing sites",
			config: `
domainName: "example.local"
wellKnownSvcs:
  - svc: idm
    svrs:
      - idm-01a.lab1.example.local
`,
			expectError: "at least one site must be defined",
		},
		{
			name: "missing wellKnownSvcs",
			config: `
domainName: "example.local"
sites:
  - lab1
`,
			expectError: "at least one service must be defined",
		},
		{
			name: "wellKnownSvcs with empty svc name",
			config: `
domainName: "example.local"
sites:
  - lab1
wellKnownSvcs:
  - svc: ""
    svrs:
      - idm-01a.lab1.example.local
`,
			expectError: "'svc' field is required",
		},
		{
			name: "wellKnownSvcs with empty svrs",
			config: `
domainName: "example.local"
sites:
  - lab1
wellKnownSvcs:
  - svc: idm
    svrs: []
`,
			expectError: "at least one server must be defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tt.config); err != nil {
				t.Fatalf("Failed to write config: %v", err)
			}
			tmpFile.Close()

			// Set up viper
			viper.Reset()
			viper.SetConfigFile(tmpFile.Name())

			err = viper.ReadInConfig()
			if err != nil {
				t.Fatalf("Failed to read config: %v", err)
			}

			c := &config{}
			err = viper.Unmarshal(c)
			if err != nil {
				t.Fatalf("Failed to unmarshal config: %v", err)
			}

			c.setDefaults()
			err = c.Validate()
			if err == nil {
				t.Errorf("Expected validation error for %s, got nil", tt.name)
			} else if !strings.Contains(err.Error(), tt.expectError) {
				t.Errorf("Expected error containing '%s', got: %v", tt.expectError, err)
			}
		})
	}
}

// TestConfigInitCmd_CreatesFile tests that config init creates a file
func TestConfigInitCmd_CreatesFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "doxctl-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test the getExampleConfig function
	exampleConfig := getExampleConfig()
	if exampleConfig == "" {
		t.Error("getExampleConfig() should return non-empty string")
	}

	// Check that example config contains expected keys
	expectedKeys := []string{
		"minVpnRoutes:",
		"domainName:",
		"sites:",
		"wellKnownSvcs:",
		"pingTimeout:",
		"dnsLookupTimeout:",
		"failThreshold:",
	}

	for _, key := range expectedKeys {
		if !strings.Contains(exampleConfig, key) {
			t.Errorf("Example config should contain '%s'", key)
		}
	}

	// Write the example config to temp dir
	configPath := filepath.Join(tmpDir, ".doxctl.yaml")
	err = os.WriteFile(configPath, []byte(exampleConfig), 0600)
	if err != nil {
		t.Fatalf("Failed to write example config: %v", err)
	}

	// Verify the file was created and can be read
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read created config: %v", err)
	}

	if string(content) != exampleConfig {
		t.Error("Written config does not match example config")
	}

	// Try to parse it with viper to ensure it's valid YAML
	viper.Reset()
	viper.SetConfigFile(configPath)
	err = viper.ReadInConfig()
	if err != nil {
		t.Errorf("Example config should be valid YAML: %v", err)
	}

	cfg := &config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		t.Errorf("Example config should unmarshal successfully: %v", err)
	}

	cfg.setDefaults()
	err = cfg.Validate()
	if err != nil {
		t.Errorf("Example config should be valid: %v", err)
	}
}

// TestConfigShowCmd_DisplaysConfig tests that config show displays configuration
func TestConfigShowCmd_DisplaysConfig(t *testing.T) {
	// Create a temporary valid config file
	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	validConfig := `
domainName: "test.local"
sites:
  - site1
  - site2
wellKnownSvcs:
  - svc: testSvc
    svrs:
      - server1.test.local
      - server2.test.local
minVpnRoutes: 10
pingTimeout: 300
dnsLookupTimeout: 150
failThreshold: 3
domNameChk: "test.local"
domSearchChk: "search.test"
domAddrChk: "10.0.0.0"
digProbeServerA: "dns1"
digProbeServerB: "dns2"
`
	if _, err := tmpFile.WriteString(validConfig); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	// Set up viper
	viper.Reset()
	viper.SetConfigFile(tmpFile.Name())

	// Read and unmarshal config
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	c := &config{}
	err = viper.Unmarshal(c)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	c.setDefaults()

	// Verify the config values
	if c.DomainName != "test.local" {
		t.Errorf("Expected domainName 'test.local', got: %s", c.DomainName)
	}

	if len(c.Sites) != 2 {
		t.Errorf("Expected 2 sites, got: %d", len(c.Sites))
	}

	if len(c.Svcs) != 1 {
		t.Errorf("Expected 1 service, got: %d", len(c.Svcs))
	}

	if c.Svcs[0].Svc != "testSvc" {
		t.Errorf("Expected service name 'testSvc', got: %s", c.Svcs[0].Svc)
	}

	if len(c.Svcs[0].Svrs) != 2 {
		t.Errorf("Expected 2 servers, got: %d", len(c.Svcs[0].Svrs))
	}

	if c.MinVpnRoutes != 10 {
		t.Errorf("Expected MinVpnRoutes 10, got: %d", c.MinVpnRoutes)
	}
}

// TestConfig_SetDefaults tests that setDefaults works correctly
func TestConfig_SetDefaults(t *testing.T) {
	tests := []struct {
		name     string
		config   *config
		expected *config
	}{
		{
			name:   "all zero values",
			config: &config{},
			expected: &config{
				PingTimeout:      250 * time.Millisecond,
				DNSLookupTimeout: 100 * time.Millisecond,
				FailThreshold:    5,
				MinVpnRoutes:     5,
			},
		},
		{
			name: "partial zero values",
			config: &config{
				PingTimeout:   500 * time.Millisecond,
				FailThreshold: 10,
			},
			expected: &config{
				PingTimeout:      500 * time.Millisecond,
				DNSLookupTimeout: 100 * time.Millisecond,
				FailThreshold:    10,
				MinVpnRoutes:     5,
			},
		},
		{
			name: "all values set",
			config: &config{
				PingTimeout:      1000 * time.Millisecond,
				DNSLookupTimeout: 200 * time.Millisecond,
				FailThreshold:    15,
				MinVpnRoutes:     20,
			},
			expected: &config{
				PingTimeout:      1000 * time.Millisecond,
				DNSLookupTimeout: 200 * time.Millisecond,
				FailThreshold:    15,
				MinVpnRoutes:     20,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.setDefaults()

			if tt.config.PingTimeout != tt.expected.PingTimeout {
				t.Errorf("PingTimeout: expected %v, got %v", tt.expected.PingTimeout, tt.config.PingTimeout)
			}
			if tt.config.DNSLookupTimeout != tt.expected.DNSLookupTimeout {
				t.Errorf("DNSLookupTimeout: expected %v, got %v", tt.expected.DNSLookupTimeout, tt.config.DNSLookupTimeout)
			}
			if tt.config.FailThreshold != tt.expected.FailThreshold {
				t.Errorf("FailThreshold: expected %v, got %v", tt.expected.FailThreshold, tt.config.FailThreshold)
			}
			if tt.config.MinVpnRoutes != tt.expected.MinVpnRoutes {
				t.Errorf("MinVpnRoutes: expected %v, got %v", tt.expected.MinVpnRoutes, tt.config.MinVpnRoutes)
			}
		})
	}
}

// TestConfig_Validate_EdgeCases tests edge cases in validation
func TestConfig_Validate_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		config      *config
		expectError bool
		errorString string
	}{
		{
			name: "valid config with all fields",
			config: &config{
				DomainName: "example.com",
				Sites:      []string{"site1"},
				Svcs: []svc{
					{Svc: "svc1", Svrs: []string{"server1"}},
				},
			},
			expectError: false,
		},
		{
			name: "empty domain name",
			config: &config{
				Sites: []string{"site1"},
				Svcs: []svc{
					{Svc: "svc1", Svrs: []string{"server1"}},
				},
			},
			expectError: true,
			errorString: "domainName is required",
		},
		{
			name: "empty sites array",
			config: &config{
				DomainName: "example.com",
				Sites:      []string{},
				Svcs: []svc{
					{Svc: "svc1", Svrs: []string{"server1"}},
				},
			},
			expectError: true,
			errorString: "at least one site must be defined",
		},
		{
			name: "nil sites array",
			config: &config{
				DomainName: "example.com",
				Sites:      nil,
				Svcs: []svc{
					{Svc: "svc1", Svrs: []string{"server1"}},
				},
			},
			expectError: true,
			errorString: "at least one site must be defined",
		},
		{
			name: "empty svcs array",
			config: &config{
				DomainName: "example.com",
				Sites:      []string{"site1"},
				Svcs:       []svc{},
			},
			expectError: true,
			errorString: "at least one service must be defined",
		},
		{
			name: "nil svcs array",
			config: &config{
				DomainName: "example.com",
				Sites:      []string{"site1"},
				Svcs:       nil,
			},
			expectError: true,
			errorString: "at least one service must be defined",
		},
		{
			name: "service with empty name",
			config: &config{
				DomainName: "example.com",
				Sites:      []string{"site1"},
				Svcs: []svc{
					{Svc: "", Svrs: []string{"server1"}},
				},
			},
			expectError: true,
			errorString: "'svc' field is required",
		},
		{
			name: "service with empty servers",
			config: &config{
				DomainName: "example.com",
				Sites:      []string{"site1"},
				Svcs: []svc{
					{Svc: "svc1", Svrs: []string{}},
				},
			},
			expectError: true,
			errorString: "at least one server must be defined",
		},
		{
			name: "service with nil servers",
			config: &config{
				DomainName: "example.com",
				Sites:      []string{"site1"},
				Svcs: []svc{
					{Svc: "svc1", Svrs: nil},
				},
			},
			expectError: true,
			errorString: "at least one server must be defined",
		},
		{
			name: "multiple services with one invalid",
			config: &config{
				DomainName: "example.com",
				Sites:      []string{"site1"},
				Svcs: []svc{
					{Svc: "svc1", Svrs: []string{"server1"}},
					{Svc: "", Svrs: []string{"server2"}},
				},
			},
			expectError: true,
			errorString: "wellKnownSvcs[1]: 'svc' field is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.errorString)
				} else if !strings.Contains(err.Error(), tt.errorString) {
					t.Errorf("Expected error containing '%s', got: %v", tt.errorString, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

// TestSvcSlice_Structure tests the svc struct
func TestSvcSlice_Structure(t *testing.T) {
	testSvc := svc{
		Svc:  "testService",
		Svrs: []string{"server1.example.com", "server2.example.com"},
	}

	if testSvc.Svc != "testService" {
		t.Errorf("Expected Svc to be 'testService', got: %s", testSvc.Svc)
	}

	if len(testSvc.Svrs) != 2 {
		t.Errorf("Expected 2 servers, got: %d", len(testSvc.Svrs))
	}

	expectedServers := []string{"server1.example.com", "server2.example.com"}
	for i, srv := range testSvc.Svrs {
		if srv != expectedServers[i] {
			t.Errorf("Expected server %d to be '%s', got: %s", i, expectedServers[i], srv)
		}
	}
}

// TestGetExampleConfig_ValidYAML tests that getExampleConfig returns valid YAML
func TestGetExampleConfig_ValidYAML(t *testing.T) {
	exampleConfig := getExampleConfig()

	// Write to temp file
	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(exampleConfig); err != nil {
		t.Fatalf("Failed to write example config: %v", err)
	}
	tmpFile.Close()

	// Parse with viper
	viper.Reset()
	viper.SetConfigFile(tmpFile.Name())

	if err := viper.ReadInConfig(); err != nil {
		t.Errorf("Example config should be valid YAML: %v", err)
	}

	// Unmarshal and validate
	cfg := &config{}
	if err := viper.Unmarshal(cfg); err != nil {
		t.Errorf("Example config should unmarshal successfully: %v", err)
	}

	cfg.setDefaults()

	if err := cfg.Validate(); err != nil {
		t.Errorf("Example config should pass validation: %v", err)
	}

	// Check specific values from example config
	if cfg.MinVpnRoutes != 5 {
		t.Errorf("Expected MinVpnRoutes to be 5, got: %d", cfg.MinVpnRoutes)
	}

	if len(cfg.Sites) == 0 {
		t.Error("Example config should have sites defined")
	}

	if len(cfg.Svcs) == 0 {
		t.Error("Example config should have services defined")
	}
}
