//go:build darwin
// +build darwin

/*
Package cmd - macOS-specific DNS implementations

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
	"doxctl/internal/cmdhelp"
	"os/exec"
	"strings"
)

// getResolverIPs returns DNS resolver IPs using macOS scutil
func getResolverIPs() []string {
	cmdBase := `printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil`
	cmdExe1 := exec.Command("bash", "-c", cmdBase)
	
	// Use default pattern if conf is nil or DomAddrChk is empty
	domAddrPattern := "."
	if conf != nil && conf.DomAddrChk != "" {
		domAddrPattern = conf.DomAddrChk
	}
	
	cmdGrep1 := `grep -A3 'ServerAddresses' | grep -E '` + domAddrPattern + `' | cut -d':' -f2`
	exeGrep1 := exec.Command("bash", "-c", cmdGrep1)
	output1, _, _ := cmdhelp.Pipeline(cmdExe1, exeGrep1)

	resolverIPs := strings.Split(strings.TrimRight(string(output1), "\n"), "\n")

	for i := 0; i < len(resolverIPs); i++ {
		resolverIPs[i] = strings.TrimSpace(resolverIPs[i])
	}

	return resolverIPs
}

// getVPNInterface returns VPN interface name using macOS scutil
func getVPNInterface() string {
	cmdBase := `printf "get State:/Network/Service/com.cisco.anyconnect/IPv4\nd.show\n" | scutil`
	cmdExe1 := exec.Command("bash", "-c", cmdBase)
	cmdGrep1 := `grep 'InterfaceName' | awk '{print $3}'`
	exeGrep1 := exec.Command("bash", "-c", cmdGrep1)
	output1, _, _ := cmdhelp.Pipeline(cmdExe1, exeGrep1)

	vpnInterface := strings.TrimRight(string(output1), "\n")

	if len(vpnInterface) == 0 {
		vpnInterface = "N/A"
	}

	return vpnInterface
}

// getDNSConfig returns DNS configuration status using macOS scutil
func getDNSConfig() (domainName, searchDomains, serverAddresses string) {
	cmdBase := `printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil`

	// Use default patterns if conf is nil or fields are empty
	domNamePattern := "."
	domSearchPattern := "."
	domAddrPattern := "."
	
	if conf != nil {
		if conf.DomNameChk != "" {
			domNamePattern = conf.DomNameChk
		}
		if conf.DomSearchChk != "" {
			domSearchPattern = conf.DomSearchChk
		}
		if conf.DomAddrChk != "" {
			domAddrPattern = conf.DomAddrChk
		}
	}

	cmdExe1 := exec.Command("bash", "-c", cmdBase)
	cmdGrep1 := `grep -q 'DomainName.*` + domNamePattern + `' && echo "DomainName set" || echo "DomainName unset"`
	exeGrep1 := exec.Command("bash", "-c", cmdGrep1)
	output1, _, _ := cmdhelp.Pipeline(cmdExe1, exeGrep1)

	cmdExe2 := exec.Command("bash", "-c", cmdBase)
	cmdGrep2 := `grep -A1 'SearchDomains' | grep -qE '` + domSearchPattern + `' && echo "SearchDomains set" || echo "SearchDomains unset"`
	exeGrep2 := exec.Command("bash", "-c", cmdGrep2)
	output2, _, _ := cmdhelp.Pipeline(cmdExe2, exeGrep2)

	cmdExe3 := exec.Command("bash", "-c", cmdBase)
	cmdGrep3 := `grep -A3 'ServerAddresses' | grep -qE '` + domAddrPattern + `' && echo "ServerAddresses set" || echo "ServerAddresses unset"`
	exeGrep3 := exec.Command("bash", "-c", cmdGrep3)
	output3, _, _ := cmdhelp.Pipeline(cmdExe3, exeGrep3)

	domainName = strings.Fields(string(output1))[1]
	searchDomains = strings.Fields(string(output2))[1]
	serverAddresses = strings.Fields(string(output3))[1]

	return domainName, searchDomains, serverAddresses
}
