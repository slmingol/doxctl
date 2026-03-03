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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	homedir "github.com/mitchellh/go-homedir"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage doxctl configuration",
	Long: `
The config command helps you manage your doxctl configuration file (.doxctl.yaml).

Available subcommands:
  - validate: Validate your configuration file
  - show:     Display current configuration
  - init:     Create an example configuration file
	`,
}

// configValidateCmd validates the configuration file
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration file",
	Long:  `Validates the .doxctl.yaml configuration file for syntax and required fields.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Try to read config
		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "\nError: Could not read configuration file\n")
			fmt.Fprintf(os.Stderr, "Details: %v\n\n", err)
			fmt.Fprintf(os.Stderr, "Run 'doxctl config init' to create a sample configuration file.\n\n")
			os.Exit(1)
		}

		// Unmarshal config
		c := &config{}
		if err := viper.Unmarshal(c); err != nil {
			fmt.Fprintf(os.Stderr, "\nError: Failed to parse configuration file '%s'\n", viper.ConfigFileUsed())
			fmt.Fprintf(os.Stderr, "Details: %v\n\n", err)
			fmt.Fprintf(os.Stderr, "Please check your configuration file for syntax errors.\n\n")
			os.Exit(1)
		}

		// Set defaults
		c.setDefaults()

		// Validate
		if err := c.Validate(); err != nil {
			fmt.Fprintf(os.Stderr, "\nError: Invalid configuration in '%s'\n", viper.ConfigFileUsed())
			fmt.Fprintf(os.Stderr, "Details: %v\n\n", err)
			os.Exit(1)
		}

		fmt.Println("")
		color.Success.Tips("Configuration file '%s' is valid!\n", viper.ConfigFileUsed())
	},
}

// configShowCmd shows the current configuration
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	Long:  `Displays the current configuration loaded from .doxctl.yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Try to read config
		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "\nError: Could not read configuration file\n")
			fmt.Fprintf(os.Stderr, "Details: %v\n\n", err)
			fmt.Fprintf(os.Stderr, "Run 'doxctl config init' to create a sample configuration file.\n\n")
			os.Exit(1)
		}

		// Unmarshal config
		c := &config{}
		if err := viper.Unmarshal(c); err != nil {
			fmt.Fprintf(os.Stderr, "\nError: Failed to parse configuration file '%s'\n", viper.ConfigFileUsed())
			fmt.Fprintf(os.Stderr, "Details: %v\n\n", err)
			os.Exit(1)
		}

		// Set defaults
		c.setDefaults()

		// Validate (but don't exit on error, just warn)
		if err := c.Validate(); err != nil {
			fmt.Println("")
			color.Warn.Tips("Warning: Configuration has validation errors: %v\n", err)
		}

		// Display configuration
		fmt.Println("")
		color.Info.Tips("Configuration loaded from: %s\n", viper.ConfigFileUsed())
		fmt.Println("")
		fmt.Printf("VPN Configuration:\n")
		fmt.Printf("  Min VPN Routes:      %d\n", c.MinVpnRoutes)
		fmt.Printf("\n")
		fmt.Printf("DNS Configuration:\n")
		fmt.Printf("  Domain Name Check:   %s\n", c.DomNameChk)
		fmt.Printf("  Search Check:        %s\n", c.DomSearchChk)
		fmt.Printf("  Address Check:       %s\n", c.DomAddrChk)
		fmt.Printf("  Domain Name:         %s\n", c.DomainName)
		fmt.Printf("  Probe Server A:      %s\n", c.ServerA)
		fmt.Printf("  Probe Server B:      %s\n", c.ServerB)
		fmt.Printf("  DNS Lookup Timeout:  %v\n", c.DNSLookupTimeout*time.Millisecond)
		fmt.Printf("\n")
		fmt.Printf("Sites:\n")
		for _, site := range c.Sites {
			fmt.Printf("  - %s\n", site)
		}
		fmt.Printf("\n")
		fmt.Printf("Well-Known Services (%d):\n", len(c.Svcs))
		for _, svc := range c.Svcs {
			fmt.Printf("  - %s (%d servers)\n", svc.Svc, len(svc.Svrs))
		}
		fmt.Printf("\n")
		fmt.Printf("Ping Configuration:\n")
		fmt.Printf("  Ping Timeout:        %v\n", c.PingTimeout*time.Millisecond)
		fmt.Printf("  Fail Threshold:      %d\n", c.FailThreshold)
		fmt.Printf("\n")
	},
}

// configInitCmd creates a sample configuration file
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an example configuration file",
	Long:  `Creates a sample .doxctl.yaml configuration file in your home directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Find home directory
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError: Could not determine home directory\n")
			fmt.Fprintf(os.Stderr, "Details: %v\n\n", err)
			os.Exit(1)
		}

		targetPath := filepath.Join(home, ".doxctl.yaml")

		// Check if file already exists
		if _, err := os.Stat(targetPath); err == nil {
			fmt.Fprintf(os.Stderr, "\nError: Configuration file already exists at '%s'\n", targetPath)
			fmt.Fprintf(os.Stderr, "Please remove or rename the existing file first.\n\n")
			os.Exit(1)
		}

		// Read the example config from the repository
		exampleContent := getExampleConfig()

		// Write the example config to the target path
		if err := os.WriteFile(targetPath, []byte(exampleContent), 0600); err != nil {
			fmt.Fprintf(os.Stderr, "\nError: Could not create configuration file\n")
			fmt.Fprintf(os.Stderr, "Details: %v\n\n", err)
			os.Exit(1)
		}

		fmt.Println("")
		color.Success.Tips("Created example configuration file at '%s'\n", targetPath)
		fmt.Println("")
		color.Info.Prompt("Please edit this file to match your environment before using doxctl.\n\n")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configInitCmd)
}

// getExampleConfig returns the example configuration content
func getExampleConfig() string {
	return strings.TrimSpace(`
#----------------
# VPN
#----------------
minVpnRoutes: 5

#----------------
# DNS
#----------------
# dom chks
domNameChk: "bandwidth.local"
domSearchChk: "[0-1].*bandwidth"
domAddrChk: "[0-1].*10.5"

# dig chks
domainName: "bandwidthclec.local"
digProbeServerA: "idm-01a"
digProbeServerB: "idm-01b"

#----------------
# SITES
#----------------
# Sites are used by:
#   - dns: For dig checks across datacenters
#   - net: For network performance testing (SLO validation)
sites:
  - lab1
  - rdu1
  - atl1
  - dfw1
  - lax2
  - jfk1
  - lhr1
  - fra1

#----------------
# SERVICES
#----------------
# Well-known services used by:
#   - svrs: Server reachability checks (ping-based)
#   - svcs: Service health checks (HTTP/HTTPS endpoint checks)
# 
# For svcs command, health endpoints are checked at:
#   https://<server>:<port><path>
# 
# Port defaults to 6443 if not specified (OpenShift API pattern)
# Path defaults to /healthz if not specified
# 
# Use brace expansion for multiple servers:
#   {a,b,c} expands to: a, b, c
#   {lab1,rdu1} expands to: lab1, rdu1
wellKnownSvcs:
  - 
    svc: openshift
    port: 6443
    path: /healthz
    svrs:
      - ocp-master-01{a,b,c}.{lab1,rdu1,dfw1,lax2,jfk1}.bandwidthclec.local
      - ocp-master-01{a,b,c}.{lhr1,fra1}.bwnet.us
  - 
    svc: elastic
    port: 9200
    path: /_cluster/health
    svrs:
      - es-master-01{a,b,c}.{lab1,rdu1}.bandwidthclec.local
  - 
    svc: idm
    port: 443
    path: /ipa/ui/
    svrs:
      - idm-01{a,b}.{lab1,rdu1,dfw1,lax2,jfk1}.bandwidthclec.local
      - idm-01{a,b}.{lhr1,fra1}.bwnet.us

#----------------
# TIMEOUTS & THRESHOLDS
#----------------
# Ping timeout in milliseconds (used by: svrs, net)
pingTimeout: 250

# Ping or reach failure threshold (used by: svrs)
failThreshold: 5

# DNS lookup timeout in milliseconds (used by: dns, svrs)
dnsLookupTimeout: 100
`) + "\n"
}
