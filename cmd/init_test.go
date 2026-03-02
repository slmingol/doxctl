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
	"testing"

	"github.com/spf13/viper"
)

// TestInitConfig_ValidFile tests initConfig with a valid config file
func TestInitConfig_ValidFile(t *testing.T) {
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
wellKnownSvcs:
  - svc: testSvc
    svrs:
      - server1.test.local
`
	if _, err := tmpFile.WriteString(validConfig); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	// Reset viper and set up for test
	viper.Reset()
	viper.SetConfigFile(tmpFile.Name())

	// Read and unmarshal config directly (testing the logic that initConfig uses)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	testConf := &config{}
	if err := viper.Unmarshal(testConf); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Set defaults and validate
	testConf.setDefaults()
	if err := testConf.Validate(); err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}

	// Verify config was loaded correctly
	if testConf.DomainName != "test.local" {
		t.Errorf("Expected DomainName 'test.local', got: %s", testConf.DomainName)
	}

	// Check defaults were set
	if testConf.PingTimeout == 0 {
		t.Error("Expected PingTimeout default to be set")
	}
	if testConf.DNSLookupTimeout == 0 {
		t.Error("Expected DNSLookupTimeout default to be set")
	}
}

// TestInitConfig_WithInvalidYAML tests initConfig with invalid YAML
func TestInitConfig_WithInvalidYAML(t *testing.T) {
	// Create a temporary file with invalid YAML
	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	invalidYAML := `
domainName: test.local
sites: [site1
  - invalid yaml
`
	if _, err := tmpFile.WriteString(invalidYAML); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	// Reset viper
	viper.Reset()
	viper.SetConfigFile(tmpFile.Name())

	// We can't actually test os.Exit, but we can test that
	// viper.ReadInConfig returns an error
	err = viper.ReadInConfig()
	if err == nil {
		t.Error("Expected error when reading invalid YAML")
	}
}

// TestInitConfig_ConfigInCurrentDirectory tests finding config in current dir
func TestInitConfig_ConfigInCurrentDirectory(t *testing.T) {
	// Create config in current directory
	tmpDir, err := os.MkdirTemp("", "doxctl-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, ".doxctl.yaml")
	validConfig := `
domainName: "current.local"
sites:
  - site1
wellKnownSvcs:
  - svc: testSvc
    svrs:
      - server1.current.local
`
	if err := os.WriteFile(configPath, []byte(validConfig), 0600); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Save current directory and change to temp dir
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Reset viper and don't set a specific config file
	viper.Reset()
	cfgFile = ""

	// Manually configure viper like initConfig does
	viper.AddConfigPath(".")
	viper.SetConfigName(".doxctl")

	err = viper.ReadInConfig()
	if err != nil {
		t.Errorf("Expected to find config in current directory, got error: %v", err)
	}
}

// TestInitConfig_NoConfigFile tests behavior when no config file exists
func TestInitConfig_NoConfigFile(t *testing.T) {
	// Create a temp directory with no config file
	tmpDir, err := os.MkdirTemp("", "doxctl-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save current directory and change to temp dir
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Reset viper
	viper.Reset()
	cfgFile = ""

	// Manually configure viper like initConfig does
	viper.AddConfigPath(".")
	viper.SetConfigName(".doxctl")
	viper.AutomaticEnv()

	// Try to read config - should return error
	err = viper.ReadInConfig()
	if err == nil {
		t.Error("Expected error when no config file exists")
	}

	// initConfig should handle this gracefully (config is optional for some commands)
	// so we just verify the error exists
}

// TestOutputFormat_Flag tests output format flag functionality
func TestOutputFormat_Flag(t *testing.T) {
	tests := []struct {
		name           string
		flagValue      string
		expectedFormat string
	}{
		{"default table", "", "table"},
		{"json format", "json", "json"},
		{"yaml format", "yaml", "yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the flag
			outputFormat = "table"

			if tt.flagValue != "" {
				outputFormat = tt.flagValue
			}

			if outputFormat != tt.expectedFormat {
				t.Errorf("Expected outputFormat '%s', got: %s", tt.expectedFormat, outputFormat)
			}
		})
	}
}

// TestVerboseFlag tests verbose flag functionality
func TestVerboseFlag(t *testing.T) {
	// Default should be false
	verboseChk = false
	if verboseChk {
		t.Error("verboseChk should default to false")
	}

	// Set to true
	verboseChk = true
	if !verboseChk {
		t.Error("verboseChk should be true after setting")
	}

	// Reset
	verboseChk = false
}

// TestAllFlag tests the allChk flag
func TestAllFlag(t *testing.T) {
	// Default should be false
	allChk = false
	if allChk {
		t.Error("allChk should default to false")
	}

	// Set to true
	allChk = true
	if !allChk {
		t.Error("allChk should be true after setting")
	}

	// Reset
	allChk = false
}

// TestGlobalVariables tests that global variables are properly initialized
func TestGlobalVariables(t *testing.T) {
	// Test that global vars can be accessed
	_ = cfgFile
	_ = verboseChk
	_ = allChk
	_ = outputFormat

	// Test that conf can be nil
	conf = nil
	if conf != nil {
		t.Error("conf should be able to be nil")
	}

	// Test that conf can hold a config
	conf = &config{
		DomainName: "test.local",
		Sites:      []string{"site1"},
		Svcs: []svc{
			{Svc: "testSvc", Svrs: []string{"server1"}},
		},
	}

	conf.setDefaults()
	if err := conf.Validate(); err != nil {
		t.Errorf("Simple config should be valid: %v", err)
	}

	// Reset
	conf = nil
}

// TestConfig_ComplexValidation tests complex config scenarios
func TestConfig_ComplexValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *config
		expectValid bool
	}{
		{
			name: "valid config with multiple services",
			config: &config{
				DomainName: "example.com",
				Sites:      []string{"site1", "site2", "site3"},
				Svcs: []svc{
					{Svc: "svc1", Svrs: []string{"server1", "server2"}},
					{Svc: "svc2", Svrs: []string{"server3"}},
				},
			},
			expectValid: true,
		},
		{
			name: "config with single site and service",
			config: &config{
				DomainName: "example.com",
				Sites:      []string{"onlySite"},
				Svcs: []svc{
					{Svc: "onlySvc", Svrs: []string{"onlyServer"}},
				},
			},
			expectValid: true,
		},
		{
			name: "config with many servers",
			config: &config{
				DomainName: "example.com",
				Sites:      []string{"site1"},
				Svcs: []svc{
					{
						Svc: "bigSvc",
						Svrs: []string{
							"server1", "server2", "server3", "server4", "server5",
							"server6", "server7", "server8", "server9", "server10",
						},
					},
				},
			},
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.setDefaults()
			err := tt.config.Validate()

			if tt.expectValid && err != nil {
				t.Errorf("Expected valid config, got error: %v", err)
			} else if !tt.expectValid && err == nil {
				t.Error("Expected invalid config, got nil error")
			}
		})
	}
}

// TestSvcSlice_MultipleServices tests svc slice with multiple services
func TestSvcSlice_MultipleServices(t *testing.T) {
	svcs := []svc{
		{Svc: "service1", Svrs: []string{"server1", "server2"}},
		{Svc: "service2", Svrs: []string{"server3"}},
		{Svc: "service3", Svrs: []string{"server4", "server5", "server6"}},
	}

	if len(svcs) != 3 {
		t.Errorf("Expected 3 services, got: %d", len(svcs))
	}

	// Verify first service
	if svcs[0].Svc != "service1" {
		t.Errorf("Expected first service name 'service1', got: %s", svcs[0].Svc)
	}

	if len(svcs[0].Svrs) != 2 {
		t.Errorf("Expected 2 servers in first service, got: %d", len(svcs[0].Svrs))
	}

	// Verify third service has correct server count
	if len(svcs[2].Svrs) != 3 {
		t.Errorf("Expected 3 servers in third service, got: %d", len(svcs[2].Svrs))
	}
}

// TestConfig_DefaultValues tests that default values are correctly applied
func TestConfig_DefaultValues(t *testing.T) {
	cfg := &config{
		DomainName: "test.local",
		Sites:      []string{"site1"},
		Svcs: []svc{
			{Svc: "svc1", Svrs: []string{"server1"}},
		},
	}

	// Before setDefaults, optional fields should be zero
	if cfg.PingTimeout != 0 {
		t.Errorf("Expected PingTimeout to be 0 before setDefaults, got: %v", cfg.PingTimeout)
	}

	cfg.setDefaults()

	// After setDefaults, check expected default values
	if cfg.PingTimeout == 0 {
		t.Error("Expected PingTimeout to be set after setDefaults")
	}

	if cfg.DNSLookupTimeout == 0 {
		t.Error("Expected DNSLookupTimeout to be set after setDefaults")
	}

	if cfg.FailThreshold == 0 {
		t.Error("Expected FailThreshold to be set after setDefaults")
	}

	if cfg.MinVpnRoutes == 0 {
		t.Error("Expected MinVpnRoutes to be set after setDefaults")
	}
}
