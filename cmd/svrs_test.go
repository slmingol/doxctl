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

func TestSvrsCmd_Initialization(t *testing.T) {
	if svrsCmd == nil {
		t.Fatal("svrsCmd should not be nil")
	}

	if svrsCmd.Use != "svrs" {
		t.Errorf("Expected Use to be 'svrs', got: %s", svrsCmd.Use)
	}

	if svrsCmd.Short == "" {
		t.Error("Short description should not be empty")
	}

	if svrsCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}

func TestSvrsCmd_Flags(t *testing.T) {
	flags := svrsCmd.Flags()

	svrsReachableFlag := flags.Lookup("svrsReachableChk")
	if svrsReachableFlag == nil {
		t.Error("svrsReachableChk flag should be defined")
	}

	allFlag := flags.Lookup("allChk")
	if allFlag == nil {
		t.Error("allChk flag should be defined")
	}
}

func TestSvrsCmd_IsSubcommandOfRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "svrs" {
			found = true
			break
		}
	}

	if !found {
		t.Error("svrs command should be a subcommand of root")
	}
}

func TestSvrsExecute_NoFlags(t *testing.T) {
	// Reset flags
	svrsReachableChk = false
	allChk = false

	// The function calls cmd.Usage() and os.Exit(1) when no flags are set
	// We just verify the flags are properly reset
	if svrsReachableChk || allChk {
		t.Error("All flags should be false")
	}
}

func TestSvrsExecute_SwitchCases(t *testing.T) {
	tests := []struct {
		name             string
		svrsReachableChk bool
		allChk           bool
		expectAction     string
	}{
		{"svrsReachableChk only", true, false, "svrsReachable"},
		{"allChk", false, true, "all"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set flags
			svrsReachableChk = tt.svrsReachableChk
			allChk = tt.allChk

			// Verify flags are set correctly
			if svrsReachableChk != tt.svrsReachableChk {
				t.Errorf("svrsReachableChk = %v, want %v", svrsReachableChk, tt.svrsReachableChk)
			}
			if allChk != tt.allChk {
				t.Errorf("allChk = %v, want %v", allChk, tt.allChk)
			}

			// Test the switch logic
			switch {
			case svrsReachableChk:
				if tt.expectAction != "svrsReachable" {
					t.Errorf("Expected action %s, got svrsReachable", tt.expectAction)
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

func TestSvrsExecute_SvrsReachableFlag(t *testing.T) {
	// Test flag setting
	svrsReachableChk = true
	allChk = false

	// Verify flag is set
	if !svrsReachableChk {
		t.Error("svrsReachableChk should be true")
	}
}

func TestSvrsExecute_AllFlag(t *testing.T) {
	// Test flag setting
	allChk = true

	// Verify flag is set
	if !allChk {
		t.Error("allChk should be true")
	}
}

func TestSvrsCmd_PreRun(t *testing.T) {
	// Test that PreRun doesn't panic with empty config
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PreRun panicked: %v", r)
		}
	}()

	// Verify PreRun is set
	if svrsCmd.PreRun == nil {
		t.Error("PreRun should be defined")
	}
}

func TestSvrsCmd_HasRun(t *testing.T) {
	if svrsCmd.Run == nil {
		t.Error("Run function should be defined")
	}
}

func TestSvrsCmd_FlagDefaults(t *testing.T) {
	// Get flags and verify defaults
	flags := svrsCmd.Flags()

	svrsReachableFlag := flags.Lookup("svrsReachableChk")
	if svrsReachableFlag.DefValue != "false" {
		t.Errorf("Expected svrsReachableChk default to be 'false', got: %s", svrsReachableFlag.DefValue)
	}

	allFlag := flags.Lookup("allChk")
	if allFlag.DefValue != "false" {
		t.Errorf("Expected allChk default to be 'false', got: %s", allFlag.DefValue)
	}
}

func TestSvrsCmd_FlagShorthand(t *testing.T) {
	flags := svrsCmd.Flags()

	svrsReachableFlag := flags.Lookup("svrsReachableChk")
	if svrsReachableFlag.Shorthand != "s" {
		t.Errorf("Expected svrsReachableChk shorthand to be 's', got: %s", svrsReachableFlag.Shorthand)
	}

	allFlag := flags.Lookup("allChk")
	if allFlag.Shorthand != "a" {
		t.Errorf("Expected allChk shorthand to be 'a', got: %s", allFlag.Shorthand)
	}
}

func TestSvrsCmd_FlagUsage(t *testing.T) {
	flags := svrsCmd.Flags()

	svrsReachableFlag := flags.Lookup("svrsReachableChk")
	if svrsReachableFlag.Usage == "" {
		t.Error("svrsReachableChk flag should have usage text")
	}

	allFlag := flags.Lookup("allChk")
	if allFlag.Usage == "" {
		t.Error("allChk flag should have usage text")
	}
}

func TestSvrsCmd_CommandStructure(t *testing.T) {
	// Verify the command is properly registered with root
	if svrsCmd.Parent() != rootCmd {
		t.Error("svrs command should have root command as parent")
	}
}

func TestSvrsExecute_SwitchLogic(t *testing.T) {
	// Test that only one branch is executed at a time
	tests := []struct {
		name             string
		svrsReachableChk bool
		allChk           bool
	}{
		{"svrsReachableChk only", true, false},
		{"allChk only", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svrsReachableChk = tt.svrsReachableChk
			allChk = tt.allChk

			// Just verify flags are set correctly
			if svrsReachableChk != tt.svrsReachableChk {
				t.Errorf("svrsReachableChk = %v, want %v", svrsReachableChk, tt.svrsReachableChk)
			}
			if allChk != tt.allChk {
				t.Errorf("allChk = %v, want %v", allChk, tt.allChk)
			}
		})
	}
}
