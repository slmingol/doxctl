//go:build darwin
// +build darwin

/*
Package cmd - macOS-specific DNS tests

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
	"strings"
	"testing"
)

func TestGetResolverIPsWithDeps_Darwin_Success(t *testing.T) {
	setupDNSTestConfig()

	scutilOutput := `<dictionary> {
  DomainName : example.com
  SearchDomains : <array> {
    0 : example.com
  }
  ServerAddresses : <array> {
    0 : 10.20.1.1
    1 : 10.20.1.2
  }
}`

	mockExec := &mockCommandExecutor{
		execFunc: func(name string, args ...string) ([]byte, error) {
			if name == "bash" && len(args) >= 2 && args[0] == "-c" {
				cmd := args[1]
				// First call - scutil
				if strings.Contains(cmd, "scutil") && !strings.Contains(cmd, "echo") {
					return []byte(scutilOutput), nil
				}
				// Second call - grep and cut
				if strings.Contains(cmd, "grep") && strings.Contains(cmd, "ServerAddresses") {
					return []byte("       10.20.1.1\n       10.20.1.2\n"), nil
				}
			}
			return []byte(""), nil
		},
	}

	mockFile := &mockFileReader{
		files: map[string][]byte{},
	}

	ips := getResolverIPsWithDeps(mockExec, mockFile)

	if len(ips) != 2 {
		t.Errorf("Expected 2 resolver IPs, got %d: %v", len(ips), ips)
	}

	if len(ips) > 0 && ips[0] != "10.20.1.1" {
		t.Errorf("Expected first IP to be 10.20.1.1, got %s", ips[0])
	}

	if len(ips) > 1 && ips[1] != "10.20.1.2" {
		t.Errorf("Expected second IP to be 10.20.1.2, got %s", ips[1])
	}
}

func TestGetResolverIPsWithDeps_Darwin_CommandError(t *testing.T) {
	setupDNSTestConfig()

	mockExec := &mockCommandExecutor{
		err: errors.New("command failed"),
	}

	mockFile := &mockFileReader{
		files: map[string][]byte{},
	}

	ips := getResolverIPsWithDeps(mockExec, mockFile)

	if len(ips) != 0 {
		t.Errorf("Expected empty resolver IPs, got %d", len(ips))
	}
}

func TestGetVPNInterfaceWithDeps_Darwin_Success(t *testing.T) {
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			`bash -c printf "get State:/Network/Service/com.cisco.anyconnect/IPv4\nd.show\n" | scutil`: []byte(`<dictionary> {
  InterfaceName : utun1
}`),
			`bash -c echo "<dictionary> {
  InterfaceName : utun1
}" | grep 'InterfaceName' | awk '{print $3}'`: []byte("utun1\n"),
		},
	}

	iface := getVPNInterfaceWithDeps(mockExec)

	if iface != "utun1" {
		t.Errorf("Expected VPN interface to be utun1, got %s", iface)
	}
}

func TestGetVPNInterfaceWithDeps_Darwin_NotFound(t *testing.T) {
	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			`bash -c printf "get State:/Network/Service/com.cisco.anyconnect/IPv4\nd.show\n" | scutil`: []byte(""),
			`bash -c echo "" | grep 'InterfaceName' | awk '{print $3}'`:                                []byte(""),
		},
	}

	iface := getVPNInterfaceWithDeps(mockExec)

	if iface != "N/A" {
		t.Errorf("Expected VPN interface to be N/A, got %s", iface)
	}
}

func TestGetVPNInterfaceWithDeps_Darwin_CommandError(t *testing.T) {
	mockExec := &mockCommandExecutor{
		err: errors.New("command failed"),
	}

	iface := getVPNInterfaceWithDeps(mockExec)

	if iface != "N/A" {
		t.Errorf("Expected VPN interface to be N/A, got %s", iface)
	}
}

func TestGetDNSConfigWithDeps_Darwin_AllSet(t *testing.T) {
	setupDNSTestConfig()

	scutilOutput := `<dictionary> {
  DomainName : example.com
  SearchDomains : <array> {
    0 : example
  }
  ServerAddresses : <array> {
    0 : 10.20.1.1
  }
}`

	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			`bash -c printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil`:                                                              []byte(scutilOutput),
			`bash -c echo "` + scutilOutput + `" | grep -q 'DomainName.*example.com' && echo "DomainName set" || echo "DomainName unset"`:                          []byte("DomainName set\n"),
			`bash -c echo "` + scutilOutput + `" | grep -A1 'SearchDomains' | grep -qE 'example' && echo "SearchDomains set" || echo "SearchDomains unset"`:        []byte("SearchDomains set\n"),
			`bash -c echo "` + scutilOutput + `" | grep -A3 'ServerAddresses' | grep -qE '10\.20\.' && echo "ServerAddresses set" || echo "ServerAddresses unset"`: []byte("ServerAddresses set\n"),
		},
	}

	mockFile := &mockFileReader{
		files: map[string][]byte{},
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

func TestGetDNSConfigWithDeps_Darwin_NoneSet(t *testing.T) {
	setupDNSTestConfig()

	scutilOutput := `<dictionary> {
}`

	mockExec := &mockCommandExecutor{
		commands: map[string][]byte{
			`bash -c printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil`:                                                              []byte(scutilOutput),
			`bash -c echo "` + scutilOutput + `" | grep -q 'DomainName.*example.com' && echo "DomainName set" || echo "DomainName unset"`:                          []byte("DomainName unset\n"),
			`bash -c echo "` + scutilOutput + `" | grep -A1 'SearchDomains' | grep -qE 'example' && echo "SearchDomains set" || echo "SearchDomains unset"`:        []byte("SearchDomains unset\n"),
			`bash -c echo "` + scutilOutput + `" | grep -A3 'ServerAddresses' | grep -qE '10\.20\.' && echo "ServerAddresses set" || echo "ServerAddresses unset"`: []byte("ServerAddresses unset\n"),
		},
	}

	mockFile := &mockFileReader{
		files: map[string][]byte{},
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

func TestGetDNSConfigWithDeps_Darwin_CommandError(t *testing.T) {
	setupDNSTestConfig()

	mockExec := &mockCommandExecutor{
		err: errors.New("command failed"),
	}

	mockFile := &mockFileReader{
		files: map[string][]byte{},
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
