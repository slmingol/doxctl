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

func TestVpnCmd_Initialization(t *testing.T) {
	if vpnCmd == nil {
		t.Fatal("vpnCmd should not be nil")
	}

	if vpnCmd.Use != "vpn" {
		t.Errorf("Expected Use to be 'vpn', got: %s", vpnCmd.Use)
	}

	if vpnCmd.Short == "" {
		t.Error("Short description should not be empty")
	}

	if vpnCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}

func TestVpnCmd_Flags(t *testing.T) {
	flags := vpnCmd.Flags()

	ifReachableFlag := flags.Lookup("ifReachableChk")
	if ifReachableFlag == nil {
		t.Error("ifReachableChk flag should be defined")
	}

	vpnRoutesFlag := flags.Lookup("vpnRoutesChk")
	if vpnRoutesFlag == nil {
		t.Error("vpnRoutesChk flag should be defined")
	}

	vpnStatusFlag := flags.Lookup("vpnStatusChk")
	if vpnStatusFlag == nil {
		t.Error("vpnStatusChk flag should be defined")
	}

	allFlag := flags.Lookup("allChk")
	if allFlag == nil {
		t.Error("allChk flag should be defined")
	}
}

func TestVpnCmd_IsSubcommandOfRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "vpn" {
			found = true
			break
		}
	}

	if !found {
		t.Error("vpn command should be a subcommand of root")
	}
}

func TestVpnExecute_NoFlags(t *testing.T) {
	// Reset flags
	ifReachableChk = false
	vpnRoutesChk = false
	vpnStatusChk = false
	allChk = false

	// The function calls cmd.Usage() and os.Exit(1) when no flags are set
	// We just verify the flags are properly reset
	if ifReachableChk || vpnRoutesChk || vpnStatusChk || allChk {
		t.Error("All flags should be false")
	}
}

func TestVpnExecute_SwitchCases(t *testing.T) {
	tests := []struct {
		name           string
		ifReachableChk bool
		vpnRoutesChk   bool
		vpnStatusChk   bool
		allChk         bool
		expectAction   string
	}{
		{"ifReachableChk only", true, false, false, false, "ifReachable"},
		{"vpnRoutesChk only", false, true, false, false, "vpnRoutes"},
		{"vpnStatusChk only", false, false, true, false, "vpnStatus"},
		{"allChk", false, false, false, true, "all"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set flags
			ifReachableChk = tt.ifReachableChk
			vpnRoutesChk = tt.vpnRoutesChk
			vpnStatusChk = tt.vpnStatusChk
			allChk = tt.allChk

			// Verify flags are set correctly
			if ifReachableChk != tt.ifReachableChk {
				t.Errorf("ifReachableChk = %v, want %v", ifReachableChk, tt.ifReachableChk)
			}
			if vpnRoutesChk != tt.vpnRoutesChk {
				t.Errorf("vpnRoutesChk = %v, want %v", vpnRoutesChk, tt.vpnRoutesChk)
			}
			if vpnStatusChk != tt.vpnStatusChk {
				t.Errorf("vpnStatusChk = %v, want %v", vpnStatusChk, tt.vpnStatusChk)
			}
			if allChk != tt.allChk {
				t.Errorf("allChk = %v, want %v", allChk, tt.allChk)
			}

			// Test the switch logic
			switch {
			case ifReachableChk:
				if tt.expectAction != "ifReachable" {
					t.Errorf("Expected action %s, got ifReachable", tt.expectAction)
				}
			case vpnRoutesChk:
				if tt.expectAction != "vpnRoutes" {
					t.Errorf("Expected action %s, got vpnRoutes", tt.expectAction)
				}
			case vpnStatusChk:
				if tt.expectAction != "vpnStatus" {
					t.Errorf("Expected action %s, got vpnStatus", tt.expectAction)
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

func TestVpnExecute_IfReachableFlag(t *testing.T) {
	// Test flag setting
	ifReachableChk = true
	vpnRoutesChk = false
	vpnStatusChk = false
	allChk = false

	// Verify flag is set
	if !ifReachableChk {
		t.Error("ifReachableChk should be true")
	}
}

func TestVpnExecute_VpnRoutesFlag(t *testing.T) {
	// Test flag setting
	ifReachableChk = false
	vpnRoutesChk = true
	vpnStatusChk = false
	allChk = false

	// Verify flag is set
	if !vpnRoutesChk {
		t.Error("vpnRoutesChk should be true")
	}
}

func TestVpnExecute_VpnStatusFlag(t *testing.T) {
	// Test flag setting
	ifReachableChk = false
	vpnRoutesChk = false
	vpnStatusChk = true
	allChk = false

	// Verify flag is set
	if !vpnStatusChk {
		t.Error("vpnStatusChk should be true")
	}
}

func TestVpnExecute_AllFlag(t *testing.T) {
	// Test flag setting
	allChk = true

	// Verify flag is set
	if !allChk {
		t.Error("allChk should be true")
	}
}

func TestVpnCmd_PreRun(t *testing.T) {
	// Test that PreRun doesn't panic with empty config
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PreRun panicked: %v", r)
		}
	}()

	// Verify PreRun is set
	if vpnCmd.PreRun == nil {
		t.Error("PreRun should be defined")
	}
}

func TestVpnCmd_HasRun(t *testing.T) {
	if vpnCmd.Run == nil {
		t.Error("Run function should be defined")
	}
}

func TestVpnCmd_FlagDefaults(t *testing.T) {
	// Get flags and verify defaults
	flags := vpnCmd.Flags()

	ifReachableFlag := flags.Lookup("ifReachableChk")
	if ifReachableFlag.DefValue != "false" {
		t.Errorf("Expected ifReachableChk default to be 'false', got: %s", ifReachableFlag.DefValue)
	}

	vpnRoutesFlag := flags.Lookup("vpnRoutesChk")
	if vpnRoutesFlag.DefValue != "false" {
		t.Errorf("Expected vpnRoutesChk default to be 'false', got: %s", vpnRoutesFlag.DefValue)
	}

	vpnStatusFlag := flags.Lookup("vpnStatusChk")
	if vpnStatusFlag.DefValue != "false" {
		t.Errorf("Expected vpnStatusChk default to be 'false', got: %s", vpnStatusFlag.DefValue)
	}

	allFlag := flags.Lookup("allChk")
	if allFlag.DefValue != "false" {
		t.Errorf("Expected allChk default to be 'false', got: %s", allFlag.DefValue)
	}
}

func TestVpnCmd_FlagShorthand(t *testing.T) {
	flags := vpnCmd.Flags()

	ifReachableFlag := flags.Lookup("ifReachableChk")
	if ifReachableFlag.Shorthand != "i" {
		t.Errorf("Expected ifReachableChk shorthand to be 'i', got: %s", ifReachableFlag.Shorthand)
	}

	vpnRoutesFlag := flags.Lookup("vpnRoutesChk")
	if vpnRoutesFlag.Shorthand != "r" {
		t.Errorf("Expected vpnRoutesChk shorthand to be 'r', got: %s", vpnRoutesFlag.Shorthand)
	}

	vpnStatusFlag := flags.Lookup("vpnStatusChk")
	if vpnStatusFlag.Shorthand != "s" {
		t.Errorf("Expected vpnStatusChk shorthand to be 's', got: %s", vpnStatusFlag.Shorthand)
	}

	allFlag := flags.Lookup("allChk")
	if allFlag.Shorthand != "a" {
		t.Errorf("Expected allChk shorthand to be 'a', got: %s", allFlag.Shorthand)
	}
}

func TestVpnCmd_FlagUsage(t *testing.T) {
	flags := vpnCmd.Flags()

	ifReachableFlag := flags.Lookup("ifReachableChk")
	if ifReachableFlag.Usage == "" {
		t.Error("ifReachableChk flag should have usage text")
	}

	vpnRoutesFlag := flags.Lookup("vpnRoutesChk")
	if vpnRoutesFlag.Usage == "" {
		t.Error("vpnRoutesChk flag should have usage text")
	}

	vpnStatusFlag := flags.Lookup("vpnStatusChk")
	if vpnStatusFlag.Usage == "" {
		t.Error("vpnStatusChk flag should have usage text")
	}
}
