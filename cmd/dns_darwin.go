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
	"strings"
)

// getResolverIPs returns DNS resolver IPs using macOS scutil
func getResolverIPs() []string {
	return getResolverIPsWithDeps(NewCommandExecutor(), NewFileReader())
}

// getResolverIPsWithDeps allows dependency injection for testing
func getResolverIPsWithDeps(executor CommandExecutor, fileReader FileReader) []string {
	cmdBase := `printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil`
	out1, err := executor.Execute("bash", "-c", cmdBase)
	if err != nil {
		return []string{}
	}

	// Use default pattern if conf is nil or DomAddrChk is empty
	domAddrPattern := "."
	if conf != nil && conf.DomAddrChk != "" {
		domAddrPattern = conf.DomAddrChk
	}

	cmdGrep1 := `grep -A3 'ServerAddresses' | grep -E '` + domAddrPattern + `' | cut -d':' -f2`
	out2, err := executor.Execute("bash", "-c", "echo \""+string(out1)+"\" | "+cmdGrep1)
	if err != nil {
		return []string{}
	}

	resolverIPs := strings.Split(strings.TrimRight(string(out2), "\n"), "\n")

	for i := 0; i < len(resolverIPs); i++ {
		resolverIPs[i] = strings.TrimSpace(resolverIPs[i])
	}

	return resolverIPs
}

// getVPNInterface returns VPN interface name using macOS scutil
func getVPNInterface() string {
	return getVPNInterfaceWithDeps(NewCommandExecutor())
}

// getVPNInterfaceWithDeps allows dependency injection for testing
func getVPNInterfaceWithDeps(executor CommandExecutor) string {
	cmdBase := `printf "get State:/Network/Service/com.cisco.anyconnect/IPv4\nd.show\n" | scutil`
	out1, err := executor.Execute("bash", "-c", cmdBase)
	if err != nil {
		return "N/A"
	}

	cmdGrep1 := `grep 'InterfaceName' | awk '{print $3}'`
	out2, err := executor.Execute("bash", "-c", "echo \""+string(out1)+"\" | "+cmdGrep1)
	if err != nil {
		return "N/A"
	}

	vpnInterface := strings.TrimRight(string(out2), "\n")

	if len(vpnInterface) == 0 {
		vpnInterface = "N/A"
	}

	return vpnInterface
}

// getDNSConfig returns DNS configuration status using macOS scutil
func getDNSConfig() (domainName, searchDomains, serverAddresses string) {
	return getDNSConfigWithDeps(NewCommandExecutor(), NewFileReader())
}

// getDNSConfigWithDeps allows dependency injection for testing
func getDNSConfigWithDeps(executor CommandExecutor, fileReader FileReader) (domainName, searchDomains, serverAddresses string) {
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

	cmdBase := `printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil`

	out, err := executor.Execute("bash", "-c", cmdBase)
	if err != nil {
		return "unset", "unset", "unset"
	}
	scutilOutput := string(out)

	cmdGrep1 := `grep -q 'DomainName.*` + domNamePattern + `' && echo "DomainName set" || echo "DomainName unset"`
	out1, _ := executor.Execute("bash", "-c", "echo \""+scutilOutput+"\" | "+cmdGrep1)

	cmdGrep2 := `grep -A1 'SearchDomains' | grep -qE '` + domSearchPattern + `' && echo "SearchDomains set" || echo "SearchDomains unset"`
	out2, _ := executor.Execute("bash", "-c", "echo \""+scutilOutput+"\" | "+cmdGrep2)

	cmdGrep3 := `grep -A3 'ServerAddresses' | grep -qE '` + domAddrPattern + `' && echo "ServerAddresses set" || echo "ServerAddresses unset"`
	out3, _ := executor.Execute("bash", "-c", "echo \""+scutilOutput+"\" | "+cmdGrep3)

	if fields := strings.Fields(string(out1)); len(fields) > 1 {
		domainName = fields[1]
	}
	if fields := strings.Fields(string(out2)); len(fields) > 1 {
		searchDomains = fields[1]
	}
	if fields := strings.Fields(string(out3)); len(fields) > 1 {
		serverAddresses = fields[1]
	}

	return domainName, searchDomains, serverAddresses
}
