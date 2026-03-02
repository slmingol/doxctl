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

// TestSvrsCmd tests the svrs command initialization
func TestSvrsCmd(t *testing.T) {
	if svrsCmd == nil {
		t.Fatal("svrsCmd should not be nil")
	}

	if svrsCmd.Use != "svrs" {
		t.Errorf("Expected Use to be 'svrs', got '%s'", svrsCmd.Use)
	}

	if svrsCmd.Short == "" {
		t.Error("Expected Short description to be non-empty")
	}

	if svrsCmd.Long == "" {
		t.Error("Expected Long description to be non-empty")
	}
}

// TestSvrsCmdFlags tests that svrs command has expected flags
func TestSvrsCmdFlags(t *testing.T) {
	expectedFlags := []string{"svrsReachableChk", "allChk"}

	for _, flagName := range expectedFlags {
		flag := svrsCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' to be defined", flagName)
		}
	}
}
