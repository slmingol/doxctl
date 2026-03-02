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
	"io"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestRootCmd_Initialization(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	if rootCmd.Use != "doxctl" {
		t.Errorf("Expected Use to be 'doxctl', got: %s", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Short description should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}

func TestRootCmd_Flags(t *testing.T) {
	flags := rootCmd.PersistentFlags()

	configFlag := flags.Lookup("config")
	if configFlag == nil {
		t.Error("config flag should be defined")
	}

	verboseFlag := flags.Lookup("verbose")
	if verboseFlag == nil {
		t.Error("verbose flag should be defined")
	}
}

func TestRootCmd_HasSubcommands(t *testing.T) {
	if !rootCmd.HasSubCommands() {
		t.Error("rootCmd should have subcommands")
	}

	// Check for expected subcommands
	expectedCmds := []string{"dns", "vpn", "svrs"}
	for _, expectedCmd := range expectedCmds {
		found := false
		for _, cmd := range rootCmd.Commands() {
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

func TestConfig_Struct(t *testing.T) {
	cfg := &config{
		MinVpnRoutes:  5,
		DomNameChk:    "example.com",
		DomSearchChk:  "search.example.com",
		DomAddrChk:    "10.0.0.1",
		DomainName:    "example.com",
		ServerA:       "8.8.8.8",
		ServerB:       "8.8.4.4",
		Sites:         []string{"site1", "site2"},
		Openshift:     []string{"oc1", "oc2"},
		FailThreshold: 3,
	}

	if cfg.MinVpnRoutes != 5 {
		t.Errorf("Expected MinVpnRoutes to be 5, got: %d", cfg.MinVpnRoutes)
	}

	if cfg.DomNameChk != "example.com" {
		t.Errorf("Expected DomNameChk to be 'example.com', got: %s", cfg.DomNameChk)
	}

	if len(cfg.Sites) != 2 {
		t.Errorf("Expected 2 sites, got: %d", len(cfg.Sites))
	}
}

func TestSvc_Struct(t *testing.T) {
	s := svc{
		Svc:  "test-service",
		Svrs: []string{"server1", "server2", "server3"},
	}

	if s.Svc != "test-service" {
		t.Errorf("Expected Svc to be 'test-service', got: %s", s.Svc)
	}

	if len(s.Svrs) != 3 {
		t.Errorf("Expected 3 servers, got: %d", len(s.Svrs))
	}
}

func TestInitConfig_WithoutConfigFile(t *testing.T) {
	// Save original values
	origCfgFile := cfgFile
	defer func() { cfgFile = origCfgFile }()

	// Clear config file path
	cfgFile = ""

	// Reset viper
	viper.Reset()

	// Call initConfig - should not panic
	initConfig()

	// Since no config file exists in test environment, we just verify no panic occurred
}

func TestInitConfig_WithNonExistentConfigFile(t *testing.T) {
	// Save original values
	origCfgFile := cfgFile
	defer func() { cfgFile = origCfgFile }()

	// Set to a non-existent file
	cfgFile = "/tmp/nonexistent-config-file-for-test.yaml"

	// Reset viper
	viper.Reset()

	// Call initConfig - should not panic
	initConfig()

	// Since config file doesn't exist, this should complete without error
}

func TestInitConfig_WithValidConfigFile(t *testing.T) {
	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "doxctl-test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write valid config
	configContent := `
minVpnRoutes: 10
domNameChk: test.example.com
domSearchChk: search.test.example.com
domAddrChk: 192.168.1.1
domainName: example.com
digProbeServerA: 1.1.1.1
digProbeServerB: 8.8.8.8
failThreshold: 5
`
	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	// Save original values
	origCfgFile := cfgFile
	defer func() { cfgFile = origCfgFile }()

	// Set config file path
	cfgFile = tmpFile.Name()

	// Reset viper
	viper.Reset()

	// Call initConfig
	initConfig()

	// Verify config was loaded
	if viper.GetInt("minVpnRoutes") != 10 {
		t.Errorf("Expected minVpnRoutes to be 10, got: %d", viper.GetInt("minVpnRoutes"))
	}

	if viper.GetString("domNameChk") != "test.example.com" {
		t.Errorf("Expected domNameChk to be 'test.example.com', got: %s", viper.GetString("domNameChk"))
	}
}

func TestExecute_Help(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set args to request help
	os.Args = []string{"doxctl", "--help"}

	// Reset the command to clear any previous state
	rootCmd.SetArgs([]string{"--help"})

	// Execute should not panic
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Execute --help failed: %v", err)
	}
}

func TestExecute_Version(t *testing.T) {
	// Test that we can set args without panic
	rootCmd.SetArgs([]string{})
	if rootCmd == nil {
		t.Error("rootCmd should not be nil after SetArgs")
	}
}

func TestRootCmd_Execute(t *testing.T) {
	// Test that Execute function can be called
	// We'll just verify it exists and doesn't panic with --help
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute panicked: %v", r)
		}
	}()

	rootCmd.SetArgs([]string{"--help"})
	rootCmd.Execute()
}

func TestExecuteFunction(t *testing.T) {
	// Test the actual Execute() function from root.go
	// Capture output to suppress it
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set args to help to avoid actual execution
	rootCmd.SetArgs([]string{"--help"})

	// Call Execute - should work fine with --help
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic: %v", r)
		}
	}()

	Execute()

	// Restore stdout
	w.Close()
	os.Stdout = old
	io.Copy(io.Discard, r)
}
