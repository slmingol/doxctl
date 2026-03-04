/*
Package cmd - Network performance testing and SLO validation

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
	"sort"
	"strings"
	"time"

	"doxctl/internal/output"

	"github.com/spf13/cobra"
)

// Network performance result
type netPerfResult struct {
	Timestamp    time.Time `json:"timestamp" yaml:"timestamp"`
	Target       string    `json:"target" yaml:"target"`
	AvgLatencyMs float64   `json:"avgLatencyMs" yaml:"avgLatencyMs"`
	MinLatencyMs float64   `json:"minLatencyMs" yaml:"minLatencyMs"`
	MaxLatencyMs float64   `json:"maxLatencyMs" yaml:"maxLatencyMs"`
	JitterMs     float64   `json:"jitterMs" yaml:"jitterMs"`
	PacketLoss   float64   `json:"packetLoss" yaml:"packetLoss"`
	MeetsSLO     bool      `json:"meetsSLO" yaml:"meetsSLO"`
	SLOThreshold float64   `json:"sloThreshold" yaml:"sloThreshold"`
}

type netPerfOutput struct {
	Timestamp time.Time       `json:"timestamp" yaml:"timestamp"`
	Results   []netPerfResult `json:"results" yaml:"results"`
	Summary   struct {
		TotalTargets int `json:"totalTargets" yaml:"totalTargets"`
		Passing      int `json:"passing" yaml:"passing"`
		Failing      int `json:"failing" yaml:"failing"`
	} `json:"summary" yaml:"summary"`
}

var (
	netPerfChk     bool
	netSLOMs       float64
	netPacketCount int
)

// netCmd represents the net command
var netCmd = &cobra.Command{
	Use:   "net",
	Short: "Network performance testing and SLO validation",
	Long: `Test network connectivity and performance against defined SLO thresholds.

This command measures:
  - Average, minimum, and maximum latency
  - Jitter (latency variance)
  - Packet loss percentage
  - SLO compliance (latency threshold)

Examples:
  # Test network performance to configured targets
  doxctl net --perf

  # Set custom SLO threshold (default: 50ms)
  doxctl net --perf --slo 100

  # Specify number of packets to send (default: 10)
  doxctl net --perf --packets 20`,
	Run: func(cmd *cobra.Command, args []string) {
		if allChk || netPerfChk {
			netPerformanceCheck()
		} else {
			_ = cmd.Usage()
			fmt.Printf("\n")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(netCmd)

	netCmd.Flags().BoolVarP(&netPerfChk, "perf", "p", false, "Run network performance tests")
	netCmd.Flags().Float64VarP(&netSLOMs, "slo", "s", 50.0, "SLO threshold in milliseconds")
	netCmd.Flags().IntVarP(&netPacketCount, "packets", "n", 10, "Number of packets to send")
}

// netPerformanceCheck tests network performance against SLO thresholds
func netPerformanceCheck() {
	netPerformanceCheckWithDeps(conf, netSLOMs, netPacketCount, NewPinger)
}

// netPerformanceCheckWithDeps allows dependency injection for testing
func netPerformanceCheckWithDeps(config *config, sloMs float64, packetCount int, pingerFactory func(string) (Pinger, error)) {
	result := netPerfOutput{
		Timestamp: time.Now(),
		Results:   []netPerfResult{},
	}

	// Get targets from config - use actual service hosts instead of site names
	var targets []string
	expander := NewBraceExpander()

	// Collect unique hosts from all services
	hostMap := make(map[string]bool)
	for _, service := range config.Svcs {
		for _, server := range service.Svrs {
			// Expand brace patterns
			expanded := expander.Expand(server)
			for _, h := range expanded {
				if !hostMap[h] {
					hostMap[h] = true
					targets = append(targets, h)
				}
			}
		}
	}

	if len(targets) == 0 {
		fmt.Println("")
		fmt.Printf("\033[1;33mWARNING:\033[0m No network targets configured in services\n")
		fmt.Printf("Please add services to your configuration file to run network performance tests.\n")
		return
	}

	var pingerErrors int
	var pingerErrorMsgs []string

	// Test each target with progressive spinner
	err := RunWithSpinnerProgress("Testing network performance", len(targets), func(index int) error {
		pingTarget := targets[index]

		pinger, err := pingerFactory(pingTarget)
		if err != nil {
			// Track pinger creation errors
			pingerErrors++
			pingerErrorMsgs = append(pingerErrorMsgs, fmt.Sprintf("%s: %v", pingTarget, err))
			return nil
		}

		pinger.SetCount(packetCount)
		pinger.SetTimeout(10 * time.Second)

		err = pinger.Run()

		stats := pinger.Statistics()

		perfResult := netPerfResult{
			Timestamp:    time.Now(),
			Target:       pingTarget,
			PacketLoss:   stats.PacketLoss,
			SLOThreshold: sloMs,
		}

		if err == nil && stats.PacketsRecv > 0 {
			perfResult.AvgLatencyMs = float64(stats.AvgRtt.Microseconds()) / 1000.0
			perfResult.MinLatencyMs = float64(stats.MinRtt.Microseconds()) / 1000.0
			perfResult.MaxLatencyMs = float64(stats.MaxRtt.Microseconds()) / 1000.0
			perfResult.JitterMs = float64(stats.StdDevRtt.Microseconds()) / 1000.0
			perfResult.MeetsSLO = perfResult.AvgLatencyMs <= sloMs && perfResult.PacketLoss < 5.0
		} else {
			perfResult.MeetsSLO = false
		}

		result.Results = append(result.Results, perfResult)

		if perfResult.MeetsSLO {
			result.Summary.Passing++
		} else {
			result.Summary.Failing++
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during network performance tests: %v\n", err)
		os.Exit(1)
	}

	result.Summary.TotalTargets = len(result.Results)

	// Check if no results were generated
	if len(result.Results) == 0 {
		fmt.Println("")
		fmt.Printf("\033[1;31mERROR:\033[0m Unable to run network performance tests\n\n")
		if pingerErrors > 0 {
			fmt.Printf("Failed to create ping instances for all %d target(s).\n\n", pingerErrors)
			if len(pingerErrorMsgs) > 0 {
				fmt.Printf("\033[1;31mDetailed errors:\033[0m\n")
				for _, errMsg := range pingerErrorMsgs {
					fmt.Printf("  • %s\n", errMsg)
				}
				fmt.Println()
			}
			fmt.Printf("\033[1;33mCommon causes:\033[0m\n")
			fmt.Printf("  • Running in container without CAP_NET_RAW capability\n")
			fmt.Printf("  • Insufficient permissions to create raw sockets\n")
			fmt.Printf("  • Network isolation preventing ICMP packets\n\n")
			fmt.Printf("\033[1;36mSolutions:\033[0m\n")
			fmt.Printf("  • Run container with: --cap-add=CAP_NET_RAW\n")
			fmt.Printf("  • Or run with: --privileged (less secure)\n")
			fmt.Printf("  • Or run doxctl directly on the host (not in container)\n")
		} else {
			fmt.Printf("No targets responded to ping requests.\n")
		}
		fmt.Println()
		return
	}

	// Output results
	switch outputFormat {
	case "json":
		output.PrintJSON(result)
	case "yaml":
		output.PrintYAML(result)
	default:
		printNetPerfTable(result)
	}
}

func printNetPerfTable(result netPerfOutput) {
	headers := []string{"Target", "Avg (ms)", "Min (ms)", "Max (ms)", "Jitter (ms)", "Loss %", "SLO", "Status"}
	var rows [][]string
	var separators []TableSeparator

	// Group results by service type, then by datacenter
	type datacenterGroup struct {
		datacenter string
		results    []netPerfResult
	}
	type serviceGroup struct {
		name        string
		datacenters map[string]*datacenterGroup
		dcOrder     []string
	}

	groups := make(map[string]*serviceGroup)
	groupOrder := []string{}

	// Categorize each result by service and datacenter
	for _, r := range result.Results {
		serviceName := detectServiceType(r.Target)
		datacenter := extractDatacenter(r.Target)

		if groups[serviceName] == nil {
			groups[serviceName] = &serviceGroup{
				name:        serviceName,
				datacenters: make(map[string]*datacenterGroup),
				dcOrder:     []string{},
			}
			groupOrder = append(groupOrder, serviceName)
		}

		if groups[serviceName].datacenters[datacenter] == nil {
			groups[serviceName].datacenters[datacenter] = &datacenterGroup{
				datacenter: datacenter,
				results:    []netPerfResult{},
			}
			groups[serviceName].dcOrder = append(groups[serviceName].dcOrder, datacenter)
		}

		groups[serviceName].datacenters[datacenter].results = append(groups[serviceName].datacenters[datacenter].results, r)
	}

	// Build rows with separators between service and datacenter groups
	rowCount := 0
	for svcIdx, serviceName := range groupOrder {
		group := groups[serviceName]

		// Sort datacenters alphabetically
		sort.Strings(group.dcOrder)

		for dcIdx, dc := range group.dcOrder {
			dcGroup := group.datacenters[dc]

			// Sort results within each datacenter alphabetically by target
			sort.Slice(dcGroup.results, func(i, j int) bool {
				return dcGroup.results[i].Target < dcGroup.results[j].Target
			})

			for _, r := range dcGroup.results {
				status := "✓ PASS"
				if !r.MeetsSLO {
					status = "✗ FAIL"
				}

				rows = append(rows, []string{
					r.Target,
					fmt.Sprintf("%.2f", r.AvgLatencyMs),
					fmt.Sprintf("%.2f", r.MinLatencyMs),
					fmt.Sprintf("%.2f", r.MaxLatencyMs),
					fmt.Sprintf("%.2f", r.JitterMs),
					fmt.Sprintf("%.1f", r.PacketLoss),
					fmt.Sprintf("%.0f ms", r.SLOThreshold),
					status,
				})
				rowCount++
			}

			// Add light separator after each datacenter (except last in service)
			if dcIdx < len(group.dcOrder)-1 {
				separators = append(separators, TableSeparator{
					RowIndex: rowCount - 1,
					Type:     LightSeparator,
				})
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

	fmt.Print(createStyledTableWithTypedSeparators(headers, rows, "Network Performance & SLO Validation", separators))
	// Print summary
	fmt.Printf("\nSummary: %d/%d targets meeting SLO (%.1f%% success rate)\n",
		result.Summary.Passing,
		result.Summary.TotalTargets, float64(result.Summary.Passing)/float64(result.Summary.TotalTargets)*100)
	fmt.Println()
}

// detectServiceType identifies the service based on hostname pattern
func detectServiceType(hostname string) string {
	if strings.HasPrefix(hostname, "api.app1.") {
		return "openshift"
	} else if strings.HasPrefix(hostname, "es-master-") {
		return "elastic"
	} else if strings.HasPrefix(hostname, "idm-") {
		return "idm"
	}
	return "other"
}

// extractDatacenter extracts the datacenter identifier from a hostname
func extractDatacenter(hostname string) string {
	// For api.app1.lab1.ocp.bandwidth.com -> lab1
	// For es-master-01d.lab1.bwnet.us -> lab1
	// For idm-01a.lab1.bandwidthclec.local -> lab1
	// For idm-01a.bru1.bwnet.us -> bru1
	// For IP addresses -> "ip"

	// Check if it's an IP address
	if strings.Count(hostname, ".") == 3 {
		if parts := strings.Split(hostname, "."); len(parts) == 4 {
			// Simple check if all parts are numeric-ish
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

	parts := strings.Split(hostname, ".")
	if len(parts) >= 3 {
		// For patterns like api.app1.lab1... or idm-01a.lab1...
		// The datacenter is typically the second or third component
		if strings.HasPrefix(hostname, "api.app1.") {
			// api.app1.lab1.ocp... -> lab1
			return parts[2]
		} else if strings.HasPrefix(hostname, "es-master-") {
			// es-master-01d.lab1.bwnet.us -> lab1
			return parts[1]
		} else if strings.HasPrefix(hostname, "idm-") {
			// idm-01a.lab1.bandwidthclec.local -> lab1
			return parts[1]
		}
	}

	return "unknown"
}
