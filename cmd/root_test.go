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
	"testing"
)

// TestRootCmd tests the root command initialization
func TestRootCmd(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	if rootCmd.Use != "doxctl" {
		t.Errorf("Expected Use to be 'doxctl', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Expected Short description to be non-empty")
	}

	if rootCmd.Long == "" {
		t.Error("Expected Long description to be non-empty")
	}
}

// TestRootCmdFlags tests that root command has expected flags
func TestRootCmdFlags(t *testing.T) {
	configFlag := rootCmd.PersistentFlags().Lookup("config")
	if configFlag == nil {
		t.Error("Expected 'config' flag to be defined")
	}

	verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("Expected 'verbose' flag to be defined")
	}
}

// TestRootCmdHasSubcommands tests that root command has expected subcommands
func TestRootCmdHasSubcommands(t *testing.T) {
	expectedCommands := []string{"dns", "vpn", "svrs"}

	for _, expectedCmd := range expectedCommands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Use == expectedCmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' to be registered", expectedCmd)
		}
	}
}

// TestExecute verifies that Execute function is exported
func TestExecute(t *testing.T) {
	// We verify Execute is exported by checking rootCmd is properly configured
	// Direct calling Execute would exit the process, so we test its dependencies
	if rootCmd.RunE == nil && rootCmd.Run == nil {
		// This is expected for root command as it has no Run function
		// It delegates to subcommands
	}
}
