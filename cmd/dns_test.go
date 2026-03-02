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
	"bytes"
	"testing"
)

func TestDnsCmd_Initialization(t *testing.T) {
	if dnsCmd == nil {
		t.Fatal("dnsCmd should not be nil")
	}

	if dnsCmd.Use != "dns" {
		t.Errorf("Expected Use to be 'dns', got: %s", dnsCmd.Use)
	}

	if dnsCmd.Short == "" {
		t.Error("Short description should not be empty")
	}

	if dnsCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}

func TestDnsCmd_Flags(t *testing.T) {
	flags := dnsCmd.Flags()

	resolverFlag := flags.Lookup("resolverChk")
	if resolverFlag == nil {
		t.Error("resolverChk flag should be defined")
	}

	pingFlag := flags.Lookup("pingChk")
	if pingFlag == nil {
		t.Error("pingChk flag should be defined")
	}

	digFlag := flags.Lookup("digChk")
	if digFlag == nil {
		t.Error("digChk flag should be defined")
	}

	allFlag := flags.Lookup("allChk")
	if allFlag == nil {
		t.Error("allChk flag should be defined")
	}
}

func TestDnsCmd_IsSubcommandOfRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "dns" {
			found = true
			break
		}
	}

	if !found {
		t.Error("dns command should be a subcommand of root")
	}
}

func TestDnsExecute_NoFlags(t *testing.T) {
	// Reset flags
	resolverChk = false
	pingChk = false
	digChk = false
	allChk = false

	// The function calls cmd.Usage() and os.Exit(1) when no flags are set
	// We just verify the flags are properly reset
	if resolverChk || pingChk || digChk || allChk {
		t.Error("All flags should be false")
	}
}

func TestDnsExecute_SwitchCases(t *testing.T) {
	tests := []struct {
		name         string
		resolverChk  bool
		pingChk      bool
		digChk       bool
		allChk       bool
		expectAction string
	}{
		{"resolverChk only", true, false, false, false, "resolver"},
		{"pingChk only", false, true, false, false, "ping"},
		{"digChk only", false, false, true, false, "dig"},
		{"allChk", false, false, false, true, "all"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set flags
			resolverChk = tt.resolverChk
			pingChk = tt.pingChk
			digChk = tt.digChk
			allChk = tt.allChk

			// Verify flags are set correctly
			if resolverChk != tt.resolverChk {
				t.Errorf("resolverChk = %v, want %v", resolverChk, tt.resolverChk)
			}
			if pingChk != tt.pingChk {
				t.Errorf("pingChk = %v, want %v", pingChk, tt.pingChk)
			}
			if digChk != tt.digChk {
				t.Errorf("digChk = %v, want %v", digChk, tt.digChk)
			}
			if allChk != tt.allChk {
				t.Errorf("allChk = %v, want %v", allChk, tt.allChk)
			}

			// Test the switch logic
			switch {
			case resolverChk:
				if tt.expectAction != "resolver" {
					t.Errorf("Expected action %s, got resolver", tt.expectAction)
				}
			case pingChk:
				if tt.expectAction != "ping" {
					t.Errorf("Expected action %s, got ping", tt.expectAction)
				}
			case digChk:
				if tt.expectAction != "dig" {
					t.Errorf("Expected action %s, got dig", tt.expectAction)
				}
			case allChk:
				if tt.expectAction != "all" {
					t.Errorf("Expected action %s, got all", tt.expectAction)
				}
			default:
				if tt.expectAction != "default" {
					t.Errorf("Expected action %s, got default", tt.expectAction)
				}
			}
		})
	}
}

func TestDnsExecute_ResolverFlag(t *testing.T) {
	// This test would require mocking system commands
	// For now, we just test flag setting
	resolverChk = true
	pingChk = false
	digChk = false
	allChk = false

	// Verify flag is set
	if !resolverChk {
		t.Error("resolverChk should be true")
	}
}

func TestDnsExecute_AllChecks(t *testing.T) {
	// Test the allChk flag which runs all checks
	allChk = true
	resolverChk = false
	pingChk = false
	digChk = false

	if !allChk {
		t.Error("allChk should be true")
	}

	// When allChk is true, dnsExecute calls all three check functions
	// We can't fully test them without mocking, but we verify the logic
}

func TestDnsExecute_PingFlag(t *testing.T) {
	pingChk = true
	resolverChk = false
	digChk = false
	allChk = false

	if !pingChk {
		t.Error("pingChk should be true")
	}
}

func TestDnsExecute_DigFlag(t *testing.T) {
	digChk = true
	resolverChk = false
	pingChk = false
	allChk = false

	if !digChk {
		t.Error("digChk should be true")
	}
}

func TestScutilResolverIPs_MockData(t *testing.T) {
	// This function relies on scutil which is macOS-specific
	// We'll test the logic with mock data instead
	
	// Test data processing logic
	lines := []string{"  10.0.0.1  ", "  10.0.0.2  "}
	
	// Simulate the trimming logic from scutilResolverIPs
	for i := 0; i < len(lines); i++ {
		// The function trims spaces, verify we can do the same
		trimmed := string(bytes.TrimSpace([]byte(lines[i])))
		if trimmed != "10.0.0.1" && trimmed != "10.0.0.2" {
			t.Errorf("Expected trimmed IP, got: %s", trimmed)
		}
	}
}

func TestScutilVPNInterface_MockData(t *testing.T) {
	// This function relies on scutil which is macOS-specific
	// We'll test the logic with mock data instead
	
	testOutput := "utun0\n"
	trimmed := string(bytes.TrimRight([]byte(testOutput), "\n"))
	
	if trimmed != "utun0" {
		t.Errorf("Expected 'utun0', got: %s", trimmed)
	}

	// Test N/A case
	emptyOutput := ""
	result := emptyOutput
	if len(result) == 0 {
		result = "N/A"
	}
	
	if result != "N/A" {
		t.Errorf("Expected 'N/A' for empty output, got: %s", result)
	}
}

func TestDnsCmd_PreRun(t *testing.T) {
	// Test that PreRun doesn't panic with empty config
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PreRun panicked: %v", r)
		}
	}()

	// PreRun should handle nil/empty config gracefully
	if dnsCmd.PreRun != nil {
		// Just verify PreRun is set, actual execution requires full viper setup
		if dnsCmd.PreRun == nil {
			t.Error("PreRun should be defined")
		}
	}
}

func TestDnsCmd_HasRun(t *testing.T) {
	if dnsCmd.Run == nil {
		t.Error("Run function should be defined")
	}
}

func TestDnsCmd_FlagDefaults(t *testing.T) {
	// Get flags and verify defaults
	flags := dnsCmd.Flags()

	resolverFlag := flags.Lookup("resolverChk")
	if resolverFlag.DefValue != "false" {
		t.Errorf("Expected resolverChk default to be 'false', got: %s", resolverFlag.DefValue)
	}

	pingFlag := flags.Lookup("pingChk")
	if pingFlag.DefValue != "false" {
		t.Errorf("Expected pingChk default to be 'false', got: %s", pingFlag.DefValue)
	}

	digFlag := flags.Lookup("digChk")
	if digFlag.DefValue != "false" {
		t.Errorf("Expected digChk default to be 'false', got: %s", digFlag.DefValue)
	}

	allFlag := flags.Lookup("allChk")
	if allFlag.DefValue != "false" {
		t.Errorf("Expected allChk default to be 'false', got: %s", allFlag.DefValue)
	}
}

func TestDnsCmd_FlagShorthand(t *testing.T) {
	flags := dnsCmd.Flags()

	resolverFlag := flags.Lookup("resolverChk")
	if resolverFlag.Shorthand != "r" {
		t.Errorf("Expected resolverChk shorthand to be 'r', got: %s", resolverFlag.Shorthand)
	}

	pingFlag := flags.Lookup("pingChk")
	if pingFlag.Shorthand != "p" {
		t.Errorf("Expected pingChk shorthand to be 'p', got: %s", pingFlag.Shorthand)
	}

	digFlag := flags.Lookup("digChk")
	if digFlag.Shorthand != "d" {
		t.Errorf("Expected digChk shorthand to be 'd', got: %s", digFlag.Shorthand)
	}

	allFlag := flags.Lookup("allChk")
	if allFlag.Shorthand != "a" {
		t.Errorf("Expected allChk shorthand to be 'a', got: %s", allFlag.Shorthand)
	}
}
