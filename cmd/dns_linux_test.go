//go:build linux
// +build linux

/*
Package cmd - Linux-specific DNS tests

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

func TestGetResolverIPsWithDeps_Linux_Success(t *testing.T) {
	setupDNSTestConfig()

	mockFile := &mockFileReader{
		files: map[string][]byte{
			"/etc/resolv.conf": []byte(`# DNS Configuration
nameserver 10.20.1.1
nameserver 10.20.1.2
nameserver 192.168.1.1
search example.com
domain example.com
`),
		},
	}

	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{},
	}

	ips := getResolverIPsWithDeps(mockExec, mockFile)

	// Should only return IPs matching the configured pattern
	if len(ips) != 2 {
		t.Errorf("Expected 2 resolver IPs, got %d", len(ips))
	}

	if ips[0] != "10.20.1.1" {
		t.Errorf("Expected first IP to be 10.20.1.1, got %s", ips[0])
	}

	if ips[1] != "10.20.1.2" {
		t.Errorf("Expected second IP to be 10.20.1.2, got %s", ips[1])
	}
}

func TestGetResolverIPsWithDeps_Linux_FileError(t *testing.T) {
	setupDNSTestConfig()

	mockFile := &mockFileReader{
		err: errors.New("file not found"),
	}

	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{},
	}

	ips := getResolverIPsWithDeps(mockExec, mockFile)

	if len(ips) != 0 {
		t.Errorf("Expected empty resolver IPs, got %d", len(ips))
	}
}

func TestGetResolverIPsWithDeps_Linux_NoNameservers(t *testing.T) {
	setupDNSTestConfig()

	mockFile := &mockFileReader{
		files: map[string][]byte{
			"/etc/resolv.conf": []byte(`# DNS Configuration
search example.com
domain example.com
`),
		},
	}

	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{},
	}

	ips := getResolverIPsWithDeps(mockExec, mockFile)

	if len(ips) != 0 {
		t.Errorf("Expected 0 resolver IPs, got %d", len(ips))
	}
}

func TestGetVPNInterfaceWithDeps_Linux_Success(t *testing.T) {
	mockFile := &mockFileReader{
		files: map[string][]byte{
			"/proc/net/route": []byte(`Iface	Destination	Gateway 	Flags	RefCnt	Use	Metric	Mask		MTU	Window	IRTT
tun0	00000000	00000000	0003	0	0	0	00000000	0	0	0
eth0	00000000	0A14010A	0003	0	0	100	00000000	0	0	0
`),
		},
	}

	iface := getVPNInterfaceWithDeps(mockFile)

	if iface != "tun0" {
		t.Errorf("Expected VPN interface to be tun0, got %s", iface)
	}
}

func TestGetVPNInterfaceWithDeps_Linux_NotFound(t *testing.T) {
	mockFile := &mockFileReader{
		files: map[string][]byte{
			"/proc/net/route": []byte(`Iface	Destination	Gateway 	Flags	RefCnt	Use	Metric	Mask		MTU	Window	IRTT
eth0	00000000	0A14010A	0003	0	0	100	00000000	0	0	0
`),
		},
	}

	iface := getVPNInterfaceWithDeps(mockFile)

	if iface != "N/A" {
		t.Errorf("Expected VPN interface to be N/A, got %s", iface)
	}
}

func TestGetVPNInterfaceWithDeps_Linux_FileError(t *testing.T) {
	mockFile := &mockFileReader{
		err: errors.New("file not found"),
	}

	iface := getVPNInterfaceWithDeps(mockFile)

	if iface != "N/A" {
		t.Errorf("Expected VPN interface to be N/A, got %s", iface)
	}
}

func TestGetDNSConfigWithDeps_Linux_AllSet(t *testing.T) {
	setupDNSTestConfig()

	mockFile := &mockFileReader{
		files: map[string][]byte{
			"/etc/resolv.conf": []byte(`nameserver 10.20.1.1
nameserver 10.20.1.2
search example.com sub.example.com
domain example.com
`),
		},
	}

	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{},
	}

	domainName, searchDomains, serverAddresses := getDNSConfigWithDeps(mockExec, mockFile)

	if domainName != "set" {
		t.Errorf("Expected domainName to be 'set', got %s", domainName)
	}

	if searchDomains != "set" {
		t.Errorf("Expected searchDomains to be 'set', got %s", searchDomains)
	}

	if serverAddresses != "set" {
		t.Errorf("Expected serverAddresses to be 'set', got %s", serverAddresses)
	}
}

func TestGetDNSConfigWithDeps_Linux_NoneSet(t *testing.T) {
	setupDNSTestConfig()

	mockFile := &mockFileReader{
		files: map[string][]byte{
			"/etc/resolv.conf": []byte(`# Empty configuration
`),
		},
	}

	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{},
	}

	domainName, searchDomains, serverAddresses := getDNSConfigWithDeps(mockExec, mockFile)

	if domainName != "unset" {
		t.Errorf("Expected domainName to be 'unset', got %s", domainName)
	}

	if searchDomains != "unset" {
		t.Errorf("Expected searchDomains to be 'unset', got %s", searchDomains)
	}

	if serverAddresses != "unset" {
		t.Errorf("Expected serverAddresses to be 'unset', got %s", serverAddresses)
	}
}

func TestGetDNSConfigWithDeps_Linux_FileError(t *testing.T) {
	setupDNSTestConfig()

	mockFile := &mockFileReader{
		err: errors.New("file not found"),
	}

	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{},
	}

	domainName, searchDomains, serverAddresses := getDNSConfigWithDeps(mockExec, mockFile)

	if domainName != "unset" {
		t.Errorf("Expected domainName to be 'unset', got %s", domainName)
	}

	if searchDomains != "unset" {
		t.Errorf("Expected searchDomains to be 'unset', got %s", searchDomains)
	}

	if serverAddresses != "unset" {
		t.Errorf("Expected serverAddresses to be 'unset', got %s", serverAddresses)
	}
}
