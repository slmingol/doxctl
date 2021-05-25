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
	"bytes"
	"container/list"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lixiangzhong/dnsutil"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Run diagnostics related to DNS servers' (resolvers') configurations",
	Long: `
doxctl's 'dns' subcommand can help triage DNS resovler configuration issues, 
general access to DNS resolvers and name resolution against DNS resolvers.`,
	Run: dnsExecute,
}

var resolverChk, pingChk, digChk bool

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
		os.Exit(1)
	}
}

func dnsResolverChk() {
	type dnsChks struct {
		domainName, searchDomains, serverAddresses string
	}

	cmdBase := `printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil`

	cmdExe1 := exec.Command("bash", "-c", cmdBase)
	cmdGrep1 := `grep -q 'DomainName.*bandwidth.local' && echo "DomainName set" || echo "DomainName unset"`
	exeGrep1 := exec.Command("bash", "-c", cmdGrep1)
	output1, _, _ := pipeline(cmdExe1, exeGrep1)

	cmdExe2 := exec.Command("bash", "-c", cmdBase)
	cmdGrep2 := `grep -A1 'SearchDomains' | grep -qE '[0-1].*bandwidth' && echo "SearchDomains set" || echo "SearchDomains unset"`
	exeGrep2 := exec.Command("bash", "-c", cmdGrep2)
	output2, _, _ := pipeline(cmdExe2, exeGrep2)

	cmdExe3 := exec.Command("bash", "-c", cmdBase)
	cmdGrep3 := `grep -A3 'ServerAddresses' | grep -qE '[0-1].*10.5' && echo "ServerAddresses set" || echo "ServerAddresses unset"`
	exeGrep3 := exec.Command("bash", "-c", cmdGrep3)
	output3, _, _ := pipeline(cmdExe3, exeGrep3)

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

	fmt.Println("\n** NOTE:** Any values of unset indicate that the VPN client is not defining DNS resolver(s) properly!\n\n")
}

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
		fmt.Println("** WARN:** Your VPN client does not appear to be defining any DNS resolver(s) properly,")
		fmt.Println("           you're either not connected via VPN or it's misconfigured!")
	}

	fmt.Println("\n\n")
}

func dnsResolverDigChk() {
	t := table.NewWriter()
	t.SetTitle("Dig Check against VPN defined DNS Resolvers")
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Hostname to 'dig'", "Resolver IP", "Is resolvable?"})

	resolverIPs := scutilResolverIPs()
	sites := []string{"lab1", "rdu1", "atl1", "dfw1", "lax2", "jfk1"}

	var dig dnsutil.Dig
	cntA := 0
	cntB := 0

	for _, site := range sites {
		serverA := "idm-01a." + site + ".bandwidthclec.local"
		serverB := "idm-01b." + site + ".bandwidthclec.local"

		for _, ip := range resolverIPs {
			dig.SetDNS(ip)
			msgA, errA := dig.GetMsg(dns.TypeA, serverA)
			msgB, errB := dig.GetMsg(dns.TypeA, serverB)

			isResolvable := false
			if errA == nil && errB == nil && msgA.Answer != nil && msgB.Answer != nil {
				isResolvable = true
			}

			t.AppendRow([]interface{}{serverA, ip, isResolvable})
			t.AppendRow([]interface{}{serverB, ip, isResolvable})

			if !isResolvable {
				continue
			}

			cntA++
			cntB++
		}

		t.AppendSeparator()
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMin: 40},
		{Number: 2, WidthMin: 15},
		{Number: 3, WidthMin: 15},
	})
	summary1 := fmt.Sprintf("resolver #1: %d", cntA)
	summary2 := fmt.Sprintf("resolver #2: %d", cntB)
	t.AppendFooter([]interface{}{"successesful queries", summary1 + "\n" + summary2})
	t.Render()

	if len(resolverIPs) <= 1 {
		fmt.Println("")
		fmt.Println("** WARN:** Your VPN client does not appear to be defining any DNS resolver(s) properly,")
		fmt.Println("           you're either not connected via VPN or it's misconfigured!")
	}

	fmt.Println("\n\n")
}

func pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	// Require at least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	// Collect the output from the command(s)
	var output bytes.Buffer
	var stderr bytes.Buffer

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		var err error
		// Connect each command's stdin to the previous command's stdout
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return nil, nil, err
		}
		// Connect each command's stderr to a buffer
		cmd.Stderr = &stderr
	}

	// Connect the output and error for the last command
	cmds[last].Stdout, cmds[last].Stderr = &output, &stderr

	// Start each command
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Wait for each command to complete
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Return the pipeline output and the collected standard error
	return output.Bytes(), stderr.Bytes(), nil
}

func scutilResolverIPs() []string {
	cmdBase := `printf "get State:/Network/Service/com.cisco.anyconnect/DNS\nd.show\n" | scutil`
	cmdExe1 := exec.Command("bash", "-c", cmdBase)
	cmdGrep1 := `grep -A3 'ServerAddresses' | grep -E '[0-1].*10.5' | cut -d':' -f2`
	exeGrep1 := exec.Command("bash", "-c", cmdGrep1)
	output1, _, _ := pipeline(cmdExe1, exeGrep1)

	resolverIPs := strings.Split(strings.TrimRight(string(output1), "\n"), "\n")

	for i := 0; i < len(resolverIPs); i++ {
		resolverIPs[i] = strings.TrimSpace(resolverIPs[i])
	}

	return resolverIPs
}
