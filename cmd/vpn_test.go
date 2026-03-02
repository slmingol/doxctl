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

// TestVpnCmd tests the vpn command initialization
func TestVpnCmd(t *testing.T) {
	if vpnCmd == nil {
		t.Fatal("vpnCmd should not be nil")
	}

	if vpnCmd.Use != "vpn" {
		t.Errorf("Expected Use to be 'vpn', got '%s'", vpnCmd.Use)
	}

	if vpnCmd.Short == "" {
		t.Error("Expected Short description to be non-empty")
	}

	if vpnCmd.Long == "" {
		t.Error("Expected Long description to be non-empty")
	}
}

// TestVpnCmdFlags tests that vpn command has expected flags
func TestVpnCmdFlags(t *testing.T) {
	expectedFlags := []string{"ifReachableChk", "vpnRoutesChk", "vpnStatusChk", "allChk"}
	
	for _, flagName := range expectedFlags {
		flag := vpnCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' to be defined", flagName)
		}
	}
}
