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
	"context"
	"doxctl/internal/output"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// svrsCmd represents the svrs command
var svrsCmd = &cobra.Command{
	Use:   "svrs",
	Short: "Run diagnostics verifying connectivity to well known servers thru a VPN connection",
	Long: `
doxctl's 'svrs' subcommand can help triage & test connectivity to 'well known servers'
thru a VPN connection to servers which have been defined in your '.doxctl.yaml' 
configuration file. 
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Process config, environment variables, and flags
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		viper.AutomaticEnv()

		// Ensure config is loaded
		if conf == nil {
			// Try to read config if not already loaded
			if err := viper.ReadInConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "\nError: Configuration file not found\n")
				fmt.Fprintf(os.Stderr, "Details: %v\n\n", err)
				fmt.Fprintf(os.Stderr, "Run 'doxctl config init' to create a sample configuration file.\n\n")
				os.Exit(1)
			}

			conf = &config{}
			err := viper.Unmarshal(&conf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\nError: Failed to parse configuration file\n")
				fmt.Fprintf(os.Stderr, "Details: %v\n\n", err)
				fmt.Fprintf(os.Stderr, "Please check your configuration file for syntax errors.\n")
				fmt.Fprintf(os.Stderr, "See .doxctl.yaml.example for a sample configuration.\n\n")
				os.Exit(1)
			}

			// Set defaults and validate
			conf.setDefaults()
			if err := conf.Validate(); err != nil {
				fmt.Fprintf(os.Stderr, "\nError: Invalid configuration\n")
				fmt.Fprintf(os.Stderr, "Details: %v\n\n", err)
				fmt.Fprintf(os.Stderr, "Please fix the configuration errors above.\n")
				fmt.Fprintf(os.Stderr, "See .doxctl.yaml.example for a sample configuration.\n\n")
				os.Exit(1)
			}
		}
	},
	Run: svrsExecute,
}

var svrsReachableChk bool

func init() {
	rootCmd.AddCommand(svrsCmd)

	svrsCmd.Flags().BoolVarP(&svrsReachableChk, "svrsReachableChk", "s", false, "Check if well known servers are reachable")
	svrsCmd.Flags().BoolVarP(&allChk, "allChk", "a", false, "Run all the checks in this subcommand module")
}

func svrsExecute(cmd *cobra.Command, args []string) {
	switch {
	case svrsReachableChk:
		svrsReachChk()
	case allChk:
		svrsReachChk()
	default:
		_ = cmd.Usage()
		fmt.Printf("\n")
		os.Exit(1)
	}
}

// Check if well known servers are pingable & reachable
func svrsReachChk() {
	svrsReachChkWithDeps(NewDNSResolver(), NewBraceExpander())
}

// svrsReachChkWithDeps allows dependency injection for testing
func svrsReachChkWithDeps(resolver DNSResolver, expander BraceExpander) {
	/* Walk through list of hosts, attempt to ping 'em.
	 *
	 * 1 - Loop through the list of svcs in .doxctl.yaml file
	 * 2 - Expand brace definitions of hosts determining all the `permutations`
	 * 3 - Go through perms. attempt to ping each and confirm that it was reached
	 * 4 - Confirm response packet was received (PacketLoss & PacketRecv)
	 * 5 - If more than FailThreshold occurs for either Packet* stop trying, call the rest failed
	 *
	 */
	pingFailures := 0
	reachFailures := 0
	var serverResults []output.ServerCheckResult

	// Build list of all targets to test
	type targetInfo struct {
		host    string
		service string
	}
	var targets []targetInfo
	for _, i := range conf.Svcs {
		for _, j := range i.Svrs {
			permutations := expander.Expand(j)
			for _, permutation := range permutations {
				targets = append(targets, targetInfo{host: permutation, service: i.Svc})
			}
		}
	}

	// Test each server with progressive spinner
	err := RunWithSpinnerProgress("Checking server reachability", len(targets), func(index int) error {
		target := targets[index]

		// Attempt to resolve hostname prior to ping
		ctx, cancel := context.WithTimeout(context.Background(), (conf.DNSLookupTimeout * time.Millisecond))
		defer cancel() // important to avoid a resource leak
		ip, err := resolver.LookupHost(ctx, target.host)

		if err != nil || len(ip) == 0 {
			serverResults = append(serverResults, output.ServerCheckResult{
				Host:        target.host,
				Service:     target.service,
				Reachable:   false,
				Performance: "N/A",
			})
			return nil
		}

		// Attempt to ping each host, any that fail keep a tally
		pinger, err := NewPinger(target.host)
		if err != nil {
			serverResults = append(serverResults, output.ServerCheckResult{
				Host:        target.host,
				Service:     target.service,
				Reachable:   false,
				Performance: "N/A",
			})
			pingFailures++
			return nil
		}

		pinger.SetTimeout(conf.PingTimeout * time.Millisecond)
		_ = pinger.Run()
		stats := pinger.Statistics()
		pingPerf := fmt.Sprintf("rnd-trp avg = %v", stats.AvgRtt)

		// Tally fails due to failed/missing responses
		packetAck := (stats.PacketLoss == 0 && stats.PacketsRecv > 0)
		if !packetAck {
			reachFailures++
		}

		serverResults = append(serverResults, output.ServerCheckResult{
			Host:        target.host,
			Service:     target.service,
			Reachable:   packetAck,
			Performance: pingPerf,
		})
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during server checks: %v\n", err)
		os.Exit(1)
	}

	// For JSON/YAML output
	if outputFormat != "table" {
		result := output.ServerReachabilityCheckResult{
			Timestamp:     time.Now(),
			Servers:       serverResults,
			PingFailures:  pingFailures,
			ReachFailures: reachFailures,
		}
		output.Print(outputFormat, result)
		return
	}

	// Table output
	fmt.Printf("\n\n   ...one sec, preparing `ping` results...\n\n")

	if pingFailures > conf.FailThreshold || reachFailures > conf.FailThreshold {
		fmt.Println("")
		color.Warn.Tips("More than %d hosts appear to be unreachable....\n\n", conf.FailThreshold)
	}

	time.Sleep(4 * time.Second)

	// Build rows with datacenter-based grouping
	headers := []string{"Host", "Service", "Reachable?", "Ping Performance"}
	var rows [][]string
	var separators []TableSeparator

	// Group results by service name, then by datacenter
	type datacenterGroup struct {
		datacenter string
		results    []output.ServerCheckResult
	}
	type serviceGroup struct {
		name        string
		datacenters map[string]*datacenterGroup
		dcOrder     []string
	}

	groups := make(map[string]*serviceGroup)
	groupOrder := []string{}

	for _, r := range serverResults {
		if groups[r.Service] == nil {
			groups[r.Service] = &serviceGroup{
				name:        r.Service,
				datacenters: make(map[string]*datacenterGroup),
				dcOrder:     []string{},
			}
			groupOrder = append(groupOrder, r.Service)
		}

		// Extract datacenter from hostname
		datacenter := extractDatacenterFromHostname(r.Host)

		if groups[r.Service].datacenters[datacenter] == nil {
			groups[r.Service].datacenters[datacenter] = &datacenterGroup{
				datacenter: datacenter,
				results:    []output.ServerCheckResult{},
			}
			groups[r.Service].dcOrder = append(groups[r.Service].dcOrder, datacenter)
		}

		groups[r.Service].datacenters[datacenter].results = append(groups[r.Service].datacenters[datacenter].results, r)
	}

	// Sort services alphabetically
	sort.Strings(groupOrder)

	// Build rows with separators between service and datacenter groups
	rowCount := 0
	for svcIdx, serviceName := range groupOrder {
		group := groups[serviceName]

		// Sort datacenters alphabetically
		sort.Strings(group.dcOrder)

		// Count total items in this service
		totalItemsInService := 0
		for _, dc := range group.dcOrder {
			totalItemsInService += len(group.datacenters[dc].results)
		}

		// If service has <= 5 items, no separators within service
		// If > 5 items, group into chunks of 4-5 at datacenter boundaries
		itemsInCurrentChunk := 0

		for dcIdx, dc := range group.dcOrder {
			dcGroup := group.datacenters[dc]

			// Sort results within each datacenter alphabetically by host
			sort.Slice(dcGroup.results, func(i, j int) bool {
				return dcGroup.results[i].Host < dcGroup.results[j].Host
			})

			for _, r := range dcGroup.results {
				row := []string{
					r.Host,
					r.Service,
					fmt.Sprintf("%t", r.Reachable),
					r.Performance,
				}

				rows = append(rows, row)
				rowCount++
				itemsInCurrentChunk++
			}

			// Add light separator if:
			// 1. This service has > 5 items total
			// 2. We've accumulated 4-5 items in current chunk
			// 3. This isn't the last datacenter in the service
			if totalItemsInService > 5 && dcIdx < len(group.dcOrder)-1 {
				if itemsInCurrentChunk >= 4 {
					separators = append(separators, TableSeparator{
						RowIndex: rowCount - 1,
						Type:     LightSeparator,
					})
					itemsInCurrentChunk = 0
				}
			}
		}

		// Add heavy separator after each service group (except last)
		if svcIdx < len(groupOrder)-1 {
			separators = append(separators, TableSeparator{
				RowIndex: rowCount - 1,
				Type:     HeavySeparator,
			})
		}
	}

	fmt.Print(createStyledTableWithTypedSeparators(headers, rows, "Well known Servers Reachable Checks", separators))

	if pingFailures > 0 || reachFailures > 0 {
		fmt.Println("")
		color.Warn.Tips(`

   Your VPN client does not appear to be functioning properly, it's likely one or more of the following:

      - Well known servers are unreachable via ping   --- try running 'doxctl vpn -h'
      - Servers are unresovlable in DNS               --- try running 'doxctl dns -h'
      - VPN client is otherwise misconfigured!
	`)
	}
}

// extractDatacenterFromHostname extracts the datacenter identifier from a hostname
func extractDatacenterFromHostname(host string) string {
	// Remove port if present
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// Check if it's an IP address
	if strings.Count(host, ".") == 3 {
		parts := strings.Split(host, ".")
		if len(parts) == 4 {
			isIP := true
			for _, p := range parts {
				if len(p) == 0 || len(p) > 3 {
					isIP = false
					break
				}
			}
			if isIP {
				return "ip"
			}
		}
	}

	parts := strings.Split(host, ".")
	if len(parts) >= 3 {
		// For patterns like:
		// - api.app1.lab1.ocp.bandwidth.com -> lab1
		// - es-master-01d.lab1.bwnet.us -> lab1
		// - idm-01a.lab1.bandwidthclec.local -> lab1
		// - idm-01a.bru1.bwnet.us -> bru1

		if strings.HasPrefix(host, "api.app1.") {
			// api.app1.lab1.ocp... -> lab1
			return parts[2]
		} else if strings.HasPrefix(host, "es-master-") {
			// es-master-01d.lab1.bwnet.us -> lab1
			return parts[1]
		} else if strings.HasPrefix(host, "idm-") {
			// idm-01a.lab1.bandwidthclec.local -> lab1
			return parts[1]
		}
	}

	return "unknown"
}
