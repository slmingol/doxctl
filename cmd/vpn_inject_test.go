/*
Package cmd - Comprehensive tests for VPN functions

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
	"errors"
	"testing"
)

// ========== Test Helpers ==========

func setupVPNTestConfig() {
	conf = &config{
		MinVpnRoutes: 5,
	}
}

// ========== Interface Reachability Tests ==========

func TestIfReachChkWithDeps_AllReachableWithTun(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "json"
	
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			"bash -c scutil --nwi | grep 'Network interfaces:' | cut -d\" \" -f 3-":                                                                           []byte("en0 utun1 en1\n"),
			"bash -c scutil --nwi | grep address -B1 -A1 | grep -E \"flags|reach\" | paste - - | column -t | grep -v Reachable | wc -l | tr -d ' '": []byte("0\n"),
		},
	}
	
	// Should not panic
	ifReachChkWithDeps(mockExec)
}

func TestIfReachChkWithDeps_NoTunInterface(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "json"
	
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			"bash -c scutil --nwi | grep 'Network interfaces:' | cut -d\" \" -f 3-":                                                                           []byte("en0 en1\n"),
			"bash -c scutil --nwi | grep address -B1 -A1 | grep -E \"flags|reach\" | paste - - | column -t | grep -v Reachable | wc -l | tr -d ' '": []byte("0\n"),
		},
	}
	
	// Should not panic even without tun interface
	ifReachChkWithDeps(mockExec)
}

func TestIfReachChkWithDeps_SomeUnreachable(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "json"
	
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			"bash -c scutil --nwi | grep 'Network interfaces:' | cut -d\" \" -f 3-":                                                                           []byte("en0 utun1 en1\n"),
			"bash -c scutil --nwi | grep address -B1 -A1 | grep -E \"flags|reach\" | paste - - | column -t | grep -v Reachable | wc -l | tr -d ' '": []byte("2\n"),
		},
	}
	
	// Should not panic with unreachable interfaces
	ifReachChkWithDeps(mockExec)
}

func TestIfReachChkWithDeps_CommandError(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "json"
	
	mockExec := &mockCommandExecutor{
		err: errors.New("scutil not available"),
	}
	
	// Should not panic even with command errors
	ifReachChkWithDeps(mockExec)
}

func TestIfReachChkWithDeps_TableOutput(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "table"
	
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			"bash -c scutil --nwi | grep 'Network interfaces:' | cut -d\" \" -f 3-":                                                                           []byte("en0 utun1\n"),
			"bash -c scutil --nwi | grep address -B1 -A1 | grep -E \"flags|reach\" | paste - - | column -t | grep -v Reachable | wc -l | tr -d ' '": []byte("0\n"),
		},
	}
	
	// Should render table without panic
	ifReachChkWithDeps(mockExec)
}

// ========== VPN Routes Tests ==========

func TestVpnRteChkWithDeps_SufficientRoutes(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "json"
	
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			"bash -c scutil --nwi | grep 'Network interfaces:' | grep -o utun[0-9] || echo \"NIL\"": []byte("utun1\n"),
			"bash -c netstat -r -f inet | grep -c utun1":                                            []byte("10\n"),
		},
	}
	
	// Should not panic with sufficient routes
	vpnRteChkWithDeps(mockExec)
}

func TestVpnRteChkWithDeps_InsufficientRoutes(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "json"
	
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			"bash -c scutil --nwi | grep 'Network interfaces:' | grep -o utun[0-9] || echo \"NIL\"": []byte("utun1\n"),
			"bash -c netstat -r -f inet | grep -c utun1":                                            []byte("2\n"),
		},
	}
	
	// Should not panic with insufficient routes
	vpnRteChkWithDeps(mockExec)
}

func TestVpnRteChkWithDeps_NoVpnInterface(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "json"
	
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			"bash -c scutil --nwi | grep 'Network interfaces:' | grep -o utun[0-9] || echo \"NIL\"": []byte("NIL\n"),
			"bash -c netstat -r -f inet | grep -c NIL":                                              []byte("0\n"),
		},
	}
	
	// Should not panic when no VPN interface found
	vpnRteChkWithDeps(mockExec)
}

func TestVpnRteChkWithDeps_CommandError(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "json"
	
	mockExec := &mockCommandExecutor{
		err: errors.New("scutil failed"),
	}
	
	// Should not panic with command errors
	vpnRteChkWithDeps(mockExec)
}

func TestVpnRteChkWithDeps_TableOutput(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "table"
	
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			"bash -c scutil --nwi | grep 'Network interfaces:' | grep -o utun[0-9] || echo \"NIL\"": []byte("utun1\n"),
			"bash -c netstat -r -f inet | grep -c utun1":                                            []byte("8\n"),
		},
	}
	
	// Should render table without panic
	vpnRteChkWithDeps(mockExec)
}

// ========== VPN Connection Status Tests ==========

func TestVpnConnChkWithDeps_Connected(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "json"
	
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			"bash -c /opt/cisco/anyconnect/bin/vpn state | grep -c 'state: Connected'": []byte("1\n"),
		},
	}
	
	// Should not panic when connected
	vpnConnChkWithDeps(mockExec)
}

func TestVpnConnChkWithDeps_NotConnected(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "json"
	
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			"bash -c /opt/cisco/anyconnect/bin/vpn state | grep -c 'state: Connected'": []byte("0\n"),
		},
	}
	
	// Should not panic when not connected
	vpnConnChkWithDeps(mockExec)
}

func TestVpnConnChkWithDeps_CommandError(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "json"
	
	mockExec := &mockCommandExecutor{
		err: errors.New("vpn command not found"),
	}
	
	// Should not panic with command errors
	vpnConnChkWithDeps(mockExec)
}

func TestVpnConnChkWithDeps_TableOutput(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "table"
	
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			"bash -c /opt/cisco/anyconnect/bin/vpn state | grep -c 'state: Connected'": []byte("1\n"),
		},
	}
	
	// Should render table without panic
	vpnConnChkWithDeps(mockExec)
}

// ========== Edge Cases ==========

func TestIfReachChkWithDeps_EmptyInterfaceList(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "json"
	
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			"bash -c scutil --nwi | grep 'Network interfaces:' | cut -d\" \" -f 3-":                                                                           []byte("\n"),
			"bash -c scutil --nwi | grep address -B1 -A1 | grep -E \"flags|reach\" | paste - - | column -t | grep -v Reachable | wc -l | tr -d ' '": []byte("0\n"),
		},
	}
	
	// Should handle empty interface list
	ifReachChkWithDeps(mockExec)
}

func TestVpnRteChkWithDeps_ZeroRoutes(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "json"
	
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			"bash -c scutil --nwi | grep 'Network interfaces:' | grep -o utun[0-9] || echo \"NIL\"": []byte("utun0\n"),
			"bash -c netstat -r -f inet | grep -c utun0":                                            []byte("0\n"),
		},
	}
	
	// Should handle zero routes gracefully
	vpnRteChkWithDeps(mockExec)
}

func TestIfReachChkWithDeps_MultipleTunInterfaces(t *testing.T) {
	setupVPNTestConfig()
	outputFormat = "json"
	
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			"bash -c scutil --nwi | grep 'Network interfaces:' | cut -d\" \" -f 3-":                                                                           []byte("en0 utun0 utun1 utun2\n"),
			"bash -c scutil --nwi | grep address -B1 -A1 | grep -E \"flags|reach\" | paste - - | column -t | grep -v Reachable | wc -l | tr -d ' '": []byte("0\n"),
		},
	}
	
	// Should handle multiple tun interfaces
	ifReachChkWithDeps(mockExec)
}
