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

	"doxctl/internal/output"

	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lixiangzhong/dnsutil"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	resolverChk, pingChk, digChk bool
)

const dnsPort = 53

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
		dnsResolverChk()
		dnsResolverPingChk()
		dnsResolverDigChk()
	default:
		_ = cmd.Usage()
		fmt.Printf("\n\n\n")
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
	fmt.Println("")

	t := table.NewWriter()
	t.SetTitle("VPN defined DNS Resolver Checks")
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Property Description", "Value"})
	t.AppendRow([]interface{}{"DomainName defined?", dns.domainName})
	t.AppendRow([]interface{}{"SearchDomains defined?", dns.searchDomains})
	t.AppendRow([]interface{}{"ServerAddresses defined?", dns.serverAddresses})
	t.AppendSeparator()
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMin: 40},
		{Number: 2, WidthMin: 30},
	})
	t.Render()

	fmt.Printf("\n")
	color.Info.Prompt("Any values of unset indicate that the VPN client is not defining DNS resolver(s) properly!\n\n")
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
	t := table.NewWriter()
	t.SetTitle("VPN defined DNS Resolver Connectivity Checks")
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Property Description", "IP", "Net i/f", "Value"})
	for e := resChks.Front(); e != nil; e = e.Next() {
		itemResChk := resolverChk(e.Value.(resolverChk))
		t.AppendRow([]interface{}{"Resovler is pingable?", itemResChk.resolverIP, itemResChk.netInterface, itemResChk.pingReachable})
		t.AppendRow([]interface{}{"Reachable via TCP?", itemResChk.resolverIP, itemResChk.netInterface, itemResChk.tcpReachable})
		t.AppendRow([]interface{}{"Reachable via UDP?", itemResChk.resolverIP, itemResChk.netInterface, itemResChk.udpReachable})
		t.AppendSeparator()
	}
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMin: 40},
		{Number: 2, WidthMin: 13},
		{Number: 3, WidthMin: 13},
	})
	t.Render()

	if len(resolverIPs) <= 1 {
		fmt.Println("")
		color.Warn.Tips(`

   Your VPN client does not appear to be defining any DNS resolver(s) properly,
   you're either not connected via VPN or it's misconfigured!`)
	}

	fmt.Printf("\n\n\n")
}

// Check if DNS resolvers return well known server records
func dnsResolverDigChk() {
	dnsResolverDigChkWithDeps(NewCommandExecutor(), NewFileReader(), NewBraceExpander())
}

// dnsResolverDigChkWithDeps allows dependency injection for testing
func dnsResolverDigChkWithDeps(executor CommandExecutor, fileReader FileReader, expander BraceExpander) {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	resolverIPs := getResolverIPsWithDeps(executor, fileReader)

	var dig dnsutil.Dig
	resolverCnt := make(map[string]int)
	var digResults []output.DigCheckResult

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
	t := table.NewWriter()
	t.SetTitle("Dig Check against VPN defined DNS Resolvers")
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Hostname to 'dig'", "Resolver IP", "Is resolvable?"}, rowConfigAutoMerge)

	for _, result := range digResults {
		t.AppendRow([]interface{}{result.Hostname, result.ResolverIP, result.IsResolvable}, rowConfigAutoMerge)
		// Add separator after each hostname's results
		if result.ResolverIP == resolverIPs[len(resolverIPs)-1] {
			t.AppendSeparator()
		}
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMin: 40, AutoMerge: true},
		{Number: 2, WidthMin: 20},
		{Number: 3, WidthMin: 15},
	})

	var summary string
	idx := 1

	for i, j := range resolverCnt {
		// only 1 resolver?
		if len(resolverCnt) == 1 {
			i = "N/A"
		}

		summary = fmt.Sprintf("%s(%s): %d", summary, i, j)

		// only 1 resolver or the last one?
		if len(resolverCnt) == 1 || len(resolverCnt) == idx {
			break
		}

		idx++
		summary = summary + "\n"
	}

	t.AppendFooter([]interface{}{"successesful queries", summary})
	t.Render()

	if len(resolverIPs) <= 1 {
		fmt.Println("")
		color.Warn.Tips(`

	   Your VPN client does not appear to be defining any DNS resolver(s) properly,
	   you're either not connected via VPN or it's misconfigured!`)
	}

	fmt.Printf("\n\n\n")
}
