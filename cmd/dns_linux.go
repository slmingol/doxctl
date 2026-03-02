//go:build linux
// +build linux

/*
Package cmd - Linux-specific DNS implementations

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
	"regexp"
	"strings"
)

// getResolverIPs returns DNS resolver IPs by parsing /etc/resolv.conf
func getResolverIPs() []string {
	return getResolverIPsWithDeps(NewCommandExecutor(), NewFileReader())
}

// getResolverIPsWithDeps allows dependency injection for testing
func getResolverIPsWithDeps(executor CommandExecutor, fileReader FileReader) []string {
	var resolverIPs []string

	fileContent, err := fileReader.ReadFile("/etc/resolv.conf")
	if err != nil {
		return resolverIPs
	}

	lines := strings.Split(string(fileContent), "\n")
	nameserverRegex := regexp.MustCompile(`^nameserver\s+(.+)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if matches := nameserverRegex.FindStringSubmatch(line); matches != nil {
			ip := strings.TrimSpace(matches[1])
			// Filter by configured domain address check if available
			if conf != nil && conf.DomAddrChk != "" {
				matched, _ := regexp.MatchString(conf.DomAddrChk, ip)
				if matched {
					resolverIPs = append(resolverIPs, ip)
				}
			} else {
				resolverIPs = append(resolverIPs, ip)
			}
		}
	}

	return resolverIPs
}

// getVPNInterface returns VPN interface name on Linux
// This attempts to detect the VPN interface by common naming patterns
func getVPNInterface() string {
	return getVPNInterfaceWithDeps(NewFileReader())
}

// getVPNInterfaceWithDeps allows dependency injection for testing
func getVPNInterfaceWithDeps(fileReader FileReader) string {
	// Common VPN interface names: tun0, ppp0, vpn0, etc.
	vpnPatterns := []string{"tun", "ppp", "vpn", "wg"}

	fileContent, err := fileReader.ReadFile("/proc/net/route")
	if err != nil {
		return "N/A"
	}

	lines := strings.Split(string(fileContent), "\n")
	// Skip header line
	for i, line := range lines {
		if i == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			iface := fields[0]
			for _, pattern := range vpnPatterns {
				if strings.HasPrefix(iface, pattern) {
					return iface
				}
			}
		}
	}

	return "N/A"
}

// getDNSConfig returns DNS configuration status by parsing /etc/resolv.conf
func getDNSConfig() (domainName, searchDomains, serverAddresses string) {
	return getDNSConfigWithDeps(NewCommandExecutor(), NewFileReader())
}

// getDNSConfigWithDeps allows dependency injection for testing
func getDNSConfigWithDeps(executor CommandExecutor, fileReader FileReader) (domainName, searchDomains, serverAddresses string) {
	fileContent, err := fileReader.ReadFile("/etc/resolv.conf")
	if err != nil {
		return "unset", "unset", "unset"
	}

	var foundDomain, foundSearch, foundNameserver bool

	lines := strings.Split(string(fileContent), "\n")
	domainRegex := regexp.MustCompile(`^domain\s+(.+)$`)
	searchRegex := regexp.MustCompile(`^search\s+(.+)$`)
	nameserverRegex := regexp.MustCompile(`^nameserver\s+(.+)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for domain
		if matches := domainRegex.FindStringSubmatch(line); matches != nil {
			domain := strings.TrimSpace(matches[1])
			if conf != nil && conf.DomNameChk != "" {
				matched, _ := regexp.MatchString(conf.DomNameChk, domain)
				if matched {
					foundDomain = true
				}
			} else {
				foundDomain = true
			}
		}

		// Check for search domains
		if matches := searchRegex.FindStringSubmatch(line); matches != nil {
			searchList := strings.TrimSpace(matches[1])
			if conf != nil && conf.DomSearchChk != "" {
				matched, _ := regexp.MatchString(conf.DomSearchChk, searchList)
				if matched {
					foundSearch = true
				}
			} else {
				foundSearch = true
			}
		}

		// Check for nameservers
		if matches := nameserverRegex.FindStringSubmatch(line); matches != nil {
			ip := strings.TrimSpace(matches[1])
			if conf != nil && conf.DomAddrChk != "" {
				matched, _ := regexp.MatchString(conf.DomAddrChk, ip)
				if matched {
					foundNameserver = true
				}
			} else {
				foundNameserver = true
			}
		}
	}

	domainName = "unset"
	if foundDomain {
		domainName = "set"
	}

	searchDomains = "unset"
	if foundSearch {
		searchDomains = "set"
	}

	serverAddresses = "unset"
	if foundNameserver {
		serverAddresses = "set"
	}

	return domainName, searchDomains, serverAddresses
}
