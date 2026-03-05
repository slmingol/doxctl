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
	"container/list"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"doxctl/internal/output"

	"github.com/lixiangzhong/dnsutil"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	resolverChk, pingChk, digChk bool
)

const dnsPort = 53

type tableData struct {
	headers []string
	rows    [][]string
}

// createStyledTable creates a table with Ocean theme colors
func createStyledTable(headers []string, rows [][]string, title string) string {
	return createStyledTableWithSeparators(headers, rows, title, nil)
}

// SeparatorType defines the style of separator
type SeparatorType int

const (
	LightSeparator SeparatorType = iota // Light separator for datacenter boundaries
	HeavySeparator                      // Heavy separator for service boundaries
)

// TableSeparator defines a separator with its position and type
type TableSeparator struct {
	RowIndex int
	Type     SeparatorType
}

// createStyledTableWithSeparators creates a table with optional row separators
func createStyledTableWithSeparators(headers []string, rows [][]string, title string, separatorAfter []int) string {
	// Convert old-style separator indices to new format (all heavy)
	var seps []TableSeparator
	for _, idx := range separatorAfter {
		seps = append(seps, TableSeparator{RowIndex: idx, Type: HeavySeparator})
	}
	return createStyledTableWithTypedSeparators(headers, rows, title, seps)
}

// createStyledTableWithTypedSeparators creates a table with typed separators
func createStyledTableWithTypedSeparators(headers []string, rows [][]string, title string, separators []TableSeparator) string {
	var output strings.Builder

	// Title bar with Ocean theme - Sky blue text on deep blue background
	output.WriteString("\n\033[38;2;135;206;250;48;2;0;0;139;1m " + title + " \033[0m\n")

	// Calculate column widths with extra padding
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = utf8.RuneCountInString(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && utf8.RuneCountInString(cell) > colWidths[i] {
				colWidths[i] = utf8.RuneCountInString(cell)
			}
		}
	}
	// Add modest padding to make tables readable, but cap at max width
	for i := range colWidths {
		padding := colWidths[i] / 10 // 10% padding instead of 20%
		if padding < 2 {
			padding = 2
		}
		colWidths[i] += padding

		// Cap column width at 45 characters to prevent excessive wrapping
		if colWidths[i] > 45 {
			colWidths[i] = 45
		}
	}

	// Ocean theme ANSI RGB colors
	borderColor := "\033[38;2;0;191;255;1m" // Bright deep sky blue borders (bold)
	headerColor := "\033[38;2;0;255;255;1m" // Bright cyan bold headers
	cellColor := "\033[38;2;211;211;211m"   // Light gray text
	reset := "\033[0m"

	// Top border
	output.WriteString(borderColor + "╭")
	for i, w := range colWidths {
		output.WriteString(strings.Repeat("─", w+2))
		if i < len(colWidths)-1 {
			output.WriteString("┬")
		}
	}
	output.WriteString("╮" + reset + "\n")

	// Headers
	output.WriteString(borderColor + "│" + reset)
	for i, h := range headers {
		paddingNeeded := colWidths[i] - utf8.RuneCountInString(h)
		output.WriteString(headerColor + " " + h + strings.Repeat(" ", paddingNeeded) + " " + reset + borderColor + "│" + reset)
	}
	output.WriteString("\n")

	// Header separator
	output.WriteString(borderColor + "├")
	for i, w := range colWidths {
		output.WriteString(strings.Repeat("─", w+2))
		if i < len(colWidths)-1 {
			output.WriteString("┼")
		}
	}
	output.WriteString("┤" + reset + "\n")

	// Data rows
	for rowIdx, row := range rows {
		output.WriteString(borderColor + "│" + reset)
		for i, cell := range row {
			if i < len(colWidths) {
				// Truncate cell if it exceeds column width
				displayCell := cell
				cellWidth := utf8.RuneCountInString(cell)
				if cellWidth > colWidths[i] {
					// Truncate by runes, not bytes
					runes := []rune(cell)
					displayCell = string(runes[:colWidths[i]-3]) + "..."
				}
				paddingNeeded := colWidths[i] - utf8.RuneCountInString(displayCell)
				output.WriteString(cellColor + " " + displayCell + strings.Repeat(" ", paddingNeeded) + " " + reset + borderColor + "│" + reset)
			}
		}
		output.WriteString("\n")

		// Add separator row if requested
		if separators != nil {
			for _, sep := range separators {
				if rowIdx == sep.RowIndex && rowIdx < len(rows)-1 {
					if sep.Type == HeavySeparator {
						// Heavy separator (service boundaries)
						output.WriteString(borderColor + "├")
						for i, w := range colWidths {
							output.WriteString(strings.Repeat("─", w+2))
							if i < len(colWidths)-1 {
								output.WriteString("┼")
							}
						}
						output.WriteString("┤" + reset + "\n")
					} else {
						// Light separator (datacenter boundaries) - gray dashed line, blue connectors
						dimColor := "\033[38;2;100;100;100m" // Dim gray for lines
						output.WriteString(borderColor + "├" + reset)
						for i, w := range colWidths {
							output.WriteString(dimColor + strings.Repeat("╌", w+2) + reset) // Dashed line
							if i < len(colWidths)-1 {
								output.WriteString(borderColor + "┼" + reset)
							}
						}
						output.WriteString(borderColor + "┤" + reset + "\n")
					}
					break
				}
			}
		}
	}

	// Bottom border
	output.WriteString(borderColor + "╰")
	for i, w := range colWidths {
		output.WriteString(strings.Repeat("─", w+2))
		if i < len(colWidths)-1 {
			output.WriteString("┴")
		}
	}
	output.WriteString("╯" + reset + "\n")

	return output.String()
}

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Run diagnostics related to DNS servers (aka. resolvers) configurations",
	Long: `
doxctl's 'dns' subcommand can help triage DNS resovler configuration issues, 
general access to DNS resolvers and name resolution against DNS resolvers.`,
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
	Run: dnsExecute,
}

func init() {
	rootCmd.AddCommand(dnsCmd)

	dnsCmd.Flags().BoolVarP(&resolverChk, "resolverChk", "r", false, "Check if VPN designated DNS resolvers are configured")
	dnsCmd.Flags().BoolVarP(&pingChk, "pingChk", "p", false, "Check if VPN defined resolvers are pingable & reachable")
	dnsCmd.Flags().BoolVarP(&digChk, "digChk", "d", false, "Check if VPN defined resolvers respond with well-known servers in DCs")
	dnsCmd.Flags().BoolVarP(&allChk, "allChk", "a", false, "Run all the checks in this subcommand module")
}

func dnsExecute(cmd *cobra.Command, args []string) {
	switch {
	case resolverChk:
		dnsResolverChk()
	case pingChk:
		dnsResolverPingChk()
	case digChk:
		dnsResolverDigChk()
	case allChk:
		// Display prominent note for full test suite
		fmt.Println()
		fmt.Printf("\033[38;2;0;128;128m%s\033[0m\n", strings.Repeat("─", 80))
		fmt.Printf("\033[38;2;255;215;0;48;2;0;0;139;1m ⚠ NOTE ⚠ \033[0m \033[1;93mFull test suite runs connectivity and DNS resolution checks (may take 30-60s)\033[0m\n")
		fmt.Printf("\033[38;2;0;128;128m%s\033[0m\n", strings.Repeat("─", 80))
		fmt.Println()
		dnsResolverChk()
		dnsResolverPingChk()
		dnsResolverDigChk()
	default:
		_ = cmd.Usage()
		fmt.Printf("\n")
		os.Exit(1)
	}
}

// Check if VPN configured DNS is setup
func dnsResolverChk() {
	dnsResolverChkWithDeps(NewCommandExecutor(), NewFileReader())
}

// dnsResolverChkWithDeps allows dependency injection for testing
func dnsResolverChkWithDeps(executor CommandExecutor, fileReader FileReader) {
	type dnsChks struct {
		domainName, searchDomains, serverAddresses string
	}

	var dns dnsChks
	dns.domainName, dns.searchDomains, dns.serverAddresses = getDNSConfigWithDeps(executor, fileReader)

	// For JSON/YAML output
	if outputFormat != "table" {
		result := output.DNSResolverCheckResult{
			Timestamp:          time.Now(),
			DomainNameSet:      dns.domainName == "set",
			SearchDomainsSet:   dns.searchDomains == "set",
			ServerAddressesSet: dns.serverAddresses == "set",
		}
		output.Print(outputFormat, result)
		return
	}

	// Table output
	headers := []string{"Property Description", "Value"}

	rows := [][]string{
		{"DomainName defined?", dns.domainName},
		{"SearchDomains defined?", dns.searchDomains},
		{"ServerAddresses defined?", dns.serverAddresses},
	}

	fmt.Print(createStyledTableWithTypedSeparators(headers, rows, "VPN defined DNS Resolver Checks", nil))

	fmt.Printf("\n")
	fmt.Printf("\033[36;1mINFO:\033[0m %s\n", "Any values of unset indicate that the VPN client is not defining DNS resolver(s) properly!")
}

// Check if DNS resolvers are pingable & reachable via TCP/UDP
func dnsResolverPingChk() {
	dnsResolverPingChkWithDeps(NewCommandExecutor(), NewFileReader(), NewPinger)
}

// dnsResolverPingChkWithDeps allows dependency injection for testing
func dnsResolverPingChkWithDeps(executor CommandExecutor, fileReader FileReader, pingerFactory func(string) (Pinger, error)) {
	type resolverChk struct {
		resolverIP, netInterface                  string
		pingReachable, tcpReachable, udpReachable bool
	}

	var resChk resolverChk
	resChks := list.New()
	resolverIPs := getResolverIPsWithDeps(executor, fileReader)

	// Wrap the connectivity checks in a spinner
	err := RunWithSpinner("Checking DNS resolver connectivity (ping, TCP, UDP)", func() error {
		for _, ip := range resolverIPs {
			var pingReachable, tcpReachable, udpReachable bool
			var netInterface string

			pinger, err := pingerFactory(ip)
			if err != nil {
				pingReachable = false
			} else {
				pinger.SetTimeout(30 * time.Second)
				err = pinger.Run()
				if err != nil {
					pingReachable = false
				} else {
					pingReachable = true
				}
			}

			resChk = resolverChk{resolverIP: ip, pingReachable: pingReachable}

			if pingReachable {
				switch runtime.GOOS {
				case "linux":
					out, err := executor.Execute("ip", "route", "get", ip) // #nosec G204 - ip is from DNS resolver list
					if err != nil {
						netInterface = "N/A"
					} else {
						netInterface = strings.Split(string(out), " ")[4]
					}
				case "darwin":
					netInterface = getVPNInterface()
				}

				target := fmt.Sprintf("%s:%d", ip, dnsPort)

				// TCP check
				_, errTCP := net.DialTimeout("tcp", target, 5*time.Second)

				tcpReachable = false
				if errTCP == nil {
					tcpReachable = true
				}

				// UDP check
				_, errUDP := net.DialTimeout("udp", target, 5*time.Second)

				udpReachable = false
				if errUDP == nil {
					udpReachable = true
				}

				// Collect results
				resChk.netInterface = netInterface
				resChk.tcpReachable = tcpReachable
				resChk.udpReachable = udpReachable
				resChks.PushBack(resChk)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during connectivity checks: %v\n", err)
	}

	// For JSON/YAML output
	if outputFormat != "table" {
		var resolvers []output.ResolverConnectivityResult
		for e := resChks.Front(); e != nil; e = e.Next() {
			itemResChk := resolverChk(e.Value.(resolverChk))
			resolvers = append(resolvers, output.ResolverConnectivityResult{
				ResolverIP:    itemResChk.resolverIP,
				NetInterface:  itemResChk.netInterface,
				PingReachable: itemResChk.pingReachable,
				TCPReachable:  itemResChk.tcpReachable,
				UDPReachable:  itemResChk.udpReachable,
			})
		}
		result := output.DNSResolverPingCheckResult{
			Timestamp: time.Now(),
			Resolvers: resolvers,
		}
		output.Print(outputFormat, result)
		return
	}

	// Table output
	headers := []string{"Property Description", "IP", "Net i/f", "Value"}

	var rows [][]string
	var separators []TableSeparator
	rowIdx := 0
	for e := resChks.Front(); e != nil; e = e.Next() {
		itemResChk := resolverChk(e.Value.(resolverChk))
		rows = append(rows, []string{
			"Resovler is pingable?",
			itemResChk.resolverIP,
			itemResChk.netInterface,
			fmt.Sprintf("%v", itemResChk.pingReachable),
		})
		rows = append(rows, []string{
			"Reachable via TCP?",
			itemResChk.resolverIP,
			itemResChk.netInterface,
			fmt.Sprintf("%v", itemResChk.tcpReachable),
		})
		rows = append(rows, []string{
			"Reachable via UDP?",
			itemResChk.resolverIP,
			itemResChk.netInterface,
			fmt.Sprintf("%v", itemResChk.udpReachable),
		})
		// Add separator after each resolver's 3 rows (except after the last one)
		if e.Next() != nil {
			separators = append(separators, TableSeparator{
				RowIndex: rowIdx + 2,
				Type:     HeavySeparator,
			})
		}
		rowIdx += 3
	}

	fmt.Print(createStyledTableWithTypedSeparators(headers, rows, "VPN defined DNS Resolver Connectivity Checks", separators))

	if len(resolverIPs) <= 1 {
		fmt.Println("")
		fmt.Printf("\033[1;33mWARNING:\033[0m %s\n%s\n",
			"Your VPN client does not appear to be defining any DNS resolver(s) properly,",
			"you're either not connected via VPN or it's misconfigured!")
	}
}

// Check if DNS resolvers return well known server records
func dnsResolverDigChk() {
	dnsResolverDigChkWithDeps(NewCommandExecutor(), NewFileReader(), NewBraceExpander())
}

// dnsResolverDigChkWithDeps allows dependency injection for testing
func dnsResolverDigChkWithDeps(executor CommandExecutor, fileReader FileReader, expander BraceExpander) {
	resolverIPs := getResolverIPsWithDeps(executor, fileReader)

	var dig dnsutil.Dig
	resolverCnt := make(map[string]int)
	var digResults []output.DigCheckResult

	// Wrap DNS resolution in a spinner
	err := RunWithSpinner("Testing DNS resolution for configured hosts", func() error {
		for _, i := range conf.Svcs {
			if i.Svc != "idm" {
				continue
			}

			for _, j := range i.Svrs {
				permutations := expander.Expand(j)

				for _, permutation := range permutations {
					for _, ip := range resolverIPs {
						_ = dig.SetDNS(ip)
						msg, err := dig.GetMsg(dns.TypeA, permutation)

						isResolvable := false
						if err == nil && msg.Answer != nil {
							isResolvable = true
						}

						// Collect for JSON/YAML output
						digResults = append(digResults, output.DigCheckResult{
							Hostname:     permutation,
							ResolverIP:   ip,
							IsResolvable: isResolvable,
						})

						if !isResolvable {
							continue
						}

						resolverCnt[ip]++
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during DNS resolution: %v\n", err)
	}

	// For JSON/YAML output
	if outputFormat != "table" {
		result := output.DNSResolverDigCheckResult{
			Timestamp: time.Now(),
			Results:   digResults,
			Summary:   resolverCnt,
		}
		output.Print(outputFormat, result)
		return
	}

	// Table output
	headers := []string{"Hostname to 'dig'", "Resolver IP", "Is resolvable?"}

	var rows [][]string
	var separators []TableSeparator
	numResolvers := len(resolverIPs)

	for rowIdx, result := range digResults {
		rows = append(rows, []string{
			result.Hostname,
			result.ResolverIP,
			fmt.Sprintf("%v", result.IsResolvable),
		})

		// Add separator after each hostname's complete set of resolver checks
		// (every numResolvers rows), but not after the last hostname
		if numResolvers > 0 && (rowIdx+1)%numResolvers == 0 && rowIdx < len(digResults)-1 {
			separators = append(separators, TableSeparator{
				RowIndex: rowIdx,
				Type:     HeavySeparator,
			})
		}
	}

	// Add summary row
	var summary string
	idx := 1
	for i, j := range resolverCnt {
		if len(resolverCnt) == 1 {
			i = "N/A"
		}
		if idx > 1 {
			summary += " | "
		}
		summary += fmt.Sprintf("(%s): %d", i, j)
		idx++
	}

	if summary != "" {
		// Add separator before summary row if we have data rows
		if len(rows) > 0 {
			separators = append(separators, TableSeparator{
				RowIndex: len(rows) - 1,
				Type:     HeavySeparator,
			})
		}
		rows = append(rows, []string{"SUCCESSFUL QUERIES", summary, ""})
	}

	fmt.Print(createStyledTableWithTypedSeparators(headers, rows, "Dig Check against VPN defined DNS Resolvers", separators))

	if len(resolverIPs) <= 1 {
		fmt.Println("")
		fmt.Printf("\033[1;33mWARNING:\033[0m %s\n%s\n",
			"Your VPN client does not appear to be defining any DNS resolver(s) properly,",
			"you're either not connected via VPN or it's misconfigured!")
	}

	fmt.Printf("\n")
}
