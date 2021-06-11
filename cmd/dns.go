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
	"doxctl/internal/cmdhelp"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-ping/ping"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	gobrex "github.com/kujtimiihoxha/go-brace-expansion"
	"github.com/lixiangzhong/dnsutil"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	resolverChk, pingChk, digChk bool
)

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

		// Populate as much info as we can from viper
		err := viper.Unmarshal(&conf)
		if err != nil {
			fmt.Printf("could not retrieve supplied project settings: %s\n", err)
			os.Exit(1)
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
		cmd.Usage()
		fmt.Printf("\n\n\n")
		os.Exit(1)
	}
}

// Check if VPN configured DNS is setup
func dnsResolverChk() {
	type dnsChks struct {
		domainName, searchDomains, serverAddresses string
	}

	cmdBase := `printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil`

	cmdExe1 := exec.Command("bash", "-c", cmdBase)
	cmdGrep1 := `grep -q 'DomainName.*` + conf.DomNameChk + `' && echo "DomainName set" || echo "DomainName unset"`
	exeGrep1 := exec.Command("bash", "-c", cmdGrep1)
	output1, _, _ := cmdhelp.Pipeline(cmdExe1, exeGrep1)

	/*if verboseChk {
		log.Info("log blah", "dnsresolvercmd", cmdBase+"|"+cmdGrep1)
	}*/

	cmdExe2 := exec.Command("bash", "-c", cmdBase)
	cmdGrep2 := `grep -A1 'SearchDomains' | grep -qE '` + conf.DomSearchChk + `' && echo "SearchDomains set" || echo "SearchDomains unset"`
	exeGrep2 := exec.Command("bash", "-c", cmdGrep2)
	output2, _, _ := cmdhelp.Pipeline(cmdExe2, exeGrep2)

	cmdExe3 := exec.Command("bash", "-c", cmdBase)
	cmdGrep3 := `grep -A3 'ServerAddresses' | grep -qE '` + conf.DomAddrChk + `' && echo "ServerAddresses set" || echo "ServerAddresses unset"`
	exeGrep3 := exec.Command("bash", "-c", cmdGrep3)
	output3, _, _ := cmdhelp.Pipeline(cmdExe3, exeGrep3)

	var dns dnsChks

	dns.domainName = strings.Fields(string(output1))[1]
	dns.searchDomains = strings.Fields(string(output2))[1]
	dns.serverAddresses = strings.Fields(string(output3))[1]

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
	type resolverChk struct {
		resolverIP, netInterface                  string
		pingReachable, tcpReachable, udpReachable bool
	}

	var resChk resolverChk
	resChks := list.New()
	resolverIPs := scutilResolverIPs()

	for _, ip := range resolverIPs {
		var pingReachable, tcpReachable, udpReachable bool
		var netInterface string

		cmdPingExe := exec.Command("ping", "-c1", ip, "-W", "200", "-t", "30", "-q")

		if _, err := cmdPingExe.CombinedOutput(); err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				pingReachable = false
			}
		} else {
			pingReachable = true
		}

		pinger, _ := ping.NewPinger(ip)
		pinger.Count = 1
		pinger.Timeout = 30 * time.Second
		pinger.Run()

		resChk = resolverChk{resolverIP: ip, pingReachable: pingReachable}

		if pingReachable {
			cmdExeIPRouteGet := exec.Command("ip", "route", "get", ip)

			if out, err := cmdExeIPRouteGet.CombinedOutput(); err != nil {
				if _, ok := err.(*exec.ExitError); ok {
					netInterface = "N/A"
				}
			} else {
				netInterface = strings.Split(string(out), " ")[4]
			}

			cmdExeNcTCP := exec.Command("nc", "-z", "-v", "-w5", ip, "53")

			if _, err := cmdExeNcTCP.CombinedOutput(); err != nil {
				if _, ok := err.(*exec.ExitError); ok {
					tcpReachable = false
				}
			} else {
				tcpReachable = true
			}

			cmdExeNcUDP := exec.Command("nc", "-z", "-u", "-v", "-w5", ip, "53")

			if _, err := cmdExeNcUDP.CombinedOutput(); err != nil {
				if _, ok := err.(*exec.ExitError); ok {
					udpReachable = false
				}
			} else {
				udpReachable = true
			}

			resChk.netInterface = netInterface
			resChk.tcpReachable = tcpReachable
			resChk.udpReachable = udpReachable
			resChks.PushBack(resChk)
		}
	}

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
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	t := table.NewWriter()
	t.SetTitle("Dig Check against VPN defined DNS Resolvers")
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Hostname to 'dig'", "Resolver IP", "Is resolvable?"}, rowConfigAutoMerge)

	resolverIPs := scutilResolverIPs()

	var dig dnsutil.Dig
	resolverCnt := make(map[string]int)

	for _, i := range conf.Svcs {
		if i.Svc != "idm" {
			continue
		}

		for _, j := range i.Svrs {
			permutations := gobrex.Expand(j)

			for _, permutation := range permutations {
				for _, ip := range resolverIPs {
					dig.SetDNS(ip)
					msg, err := dig.GetMsg(dns.TypeA, permutation)

					isResolvable := false
					if err == nil && msg.Answer != nil {
						isResolvable = true
					}

					t.AppendRow([]interface{}{permutation, ip, isResolvable}, rowConfigAutoMerge)

					if !isResolvable {
						continue
					}

					resolverCnt[ip]++
				}

				t.AppendSeparator()
			}
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

// scutil
func scutilResolverIPs() []string {
	cmdBase := `printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil`
	cmdExe1 := exec.Command("bash", "-c", cmdBase)
	cmdGrep1 := `grep -A3 'ServerAddresses' | grep -E '` + conf.DomAddrChk + `' | cut -d':' -f2`
	exeGrep1 := exec.Command("bash", "-c", cmdGrep1)
	output1, _, _ := cmdhelp.Pipeline(cmdExe1, exeGrep1)

	resolverIPs := strings.Split(strings.TrimRight(string(output1), "\n"), "\n")

	for i := 0; i < len(resolverIPs); i++ {
		resolverIPs[i] = strings.TrimSpace(resolverIPs[i])
	}

	return resolverIPs
}
