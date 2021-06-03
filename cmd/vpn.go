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
	"doxctl/internal/cmdhelp"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// vpnCmd represents the vpn command
var vpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "Run diagnostics related to VPN connections, network interfaces & configurations",
	Long: `
doxctl's 'vpn' subcommand can help triage VPN related configuration issues,
& routes related to a VPN connection.`,
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
	Run: vpnExecute,
}

var ifReachableChk, vpnRoutesChk, vpnStatusChk bool

func init() {
	rootCmd.AddCommand(vpnCmd)

	vpnCmd.Flags().BoolVarP(&ifReachableChk, "ifReachableChk", "i", false, "Check if network interfaces are reachable")
	vpnCmd.Flags().BoolVarP(&vpnRoutesChk, "vpnRoutesChk", "r", false, "Check if >5 VPN routes are defined")
	vpnCmd.Flags().BoolVarP(&vpnStatusChk, "vpnStatusChk", "s", false, "Check if VPN client's status reports as 'Connected'")
	vpnCmd.Flags().BoolVarP(&allChk, "allChk", "a", false, "Run all the checks in this subcommand module")
}

func vpnExecute(cmd *cobra.Command, args []string) {
	switch {
	case ifReachableChk:
		ifReachChk()
	case vpnRoutesChk:
		vpnRteChk()
	case vpnStatusChk:
		vpnConnChk()
	case allChk:
		ifReachChk()
		vpnRteChk()
		vpnConnChk()
	default:
		cmd.Usage()
		fmt.Printf("\n\n\n")
		os.Exit(1)
	}
}

// Test if interface is reported as reachable via 'scutil'
func ifReachChk() {
	cmdBase := `scutil --nwi`

	cmdExe1 := exec.Command("bash", "-c", cmdBase)
	cmdGrep1 := `grep 'Network interfaces:' | cut -d" " -f 3-`
	exeGrep1 := exec.Command("bash", "-c", cmdGrep1)
	output1, _, _ := cmdhelp.Pipeline(cmdExe1, exeGrep1)

	netIfs := strings.Split(strings.TrimRight(string(output1), "\n"), " ")

	var tunIfs []string
	for i := 0; i < len(netIfs); i++ {
		netIfs[i] = strings.TrimSpace(netIfs[i])
		if strings.Contains(netIfs[i], "tun") {
			tunIfs = append(tunIfs, netIfs[i])
		}
	}

	var foundOneTunIf bool = false
	if len(tunIfs) > 0 {
		foundOneTunIf = true
	}

	cmdExe2 := exec.Command("bash", "-c", cmdBase)
	cmdGrep2 := `grep address -B1 -A1 | grep -E "flags|reach" | paste - - | column -t | grep -v Reachable | wc -l | tr -d ' '`
	exeGrep2 := exec.Command("bash", "-c", cmdGrep2)
	output2, _, _ := cmdhelp.Pipeline(cmdExe2, exeGrep2)

	reachableIfs := strings.TrimRight(string(output2), "\n")

	var allInfsReachable bool = false
	if reachableIfs == "0" {
		allInfsReachable = true
	}

	t := table.NewWriter()
	t.SetTitle("Interfaces Reachable Checks")
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Property Description", "Value", "Notes"})
	t.AppendRow([]interface{}{"How many network interfaces found?", len(netIfs), netIfs})
	t.AppendRow([]interface{}{"At least 1 interface's a utun device?", foundOneTunIf, tunIfs})
	t.AppendRow([]interface{}{"All active interfaces are reporting as reachable?", allInfsReachable})
	t.AppendSeparator()
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMin: 50},
		{Number: 2, WidthMin: 30},
	})
	t.Render()

	if len(tunIfs) < 1 {
		fmt.Println("")
		color.Warn.Tips(`

   Your VPN client does not appear to be defining a TUN interface properly,
   your VPN is either not connected or it's misconfigured!`)
	}

	fmt.Printf("\n\n\n")
}

// Test if VPNs interface defines at least MinVpnRoutes routes
func vpnRteChk() {
	cmdExe1 := exec.Command("bash", "-c", "scutil --nwi")
	cmdGrep1 := `grep 'Network interfaces:' | grep -o utun[0-9] || echo "NIL"`
	exeGrep1 := exec.Command("bash", "-c", cmdGrep1)
	output1, _, _ := cmdhelp.Pipeline(cmdExe1, exeGrep1)

	vpnIf := strings.Split(strings.TrimRight(string(output1), "\n"), " ")[0]

	cmdExe2 := exec.Command("bash", "-c", `netstat -r -f inet`)
	cmdGrep2 := `grep -c ` + vpnIf
	exeGrep2 := exec.Command("bash", "-c", cmdGrep2)
	output2, _, _ := cmdhelp.Pipeline(cmdExe2, exeGrep2)

	vpnRouteCnt, _ := strconv.Atoi(strings.Split(strings.TrimRight(string(output2), "\n"), " ")[0])

	t := table.NewWriter()
	t.SetTitle("VPN Interface Route Checks")
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Property Description", "Value", "Notes"})
	t.AppendRow([]interface{}{fmt.Sprintf("At least [%d] routes using interface [%s]?", conf.MinVpnRoutes, vpnIf), vpnRouteCnt >= conf.MinVpnRoutes, vpnRouteCnt})
	t.AppendSeparator()
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMin: 50},
		{Number: 2, WidthMin: 30},
	})
	t.Render()

	if vpnRouteCnt < conf.MinVpnRoutes {
		fmt.Println("")
		color.Warn.Tips(`

   Your VPN client does not appear to be defining a TUN interface properly,
   it's either not connected or it's misconfigured!`)
	}

	fmt.Printf("\n\n\n")
}

// Test if VPN connection status reports as 'connected'
func vpnConnChk() {
	cmdBase := `/opt/cisco/anyconnect/bin/vpn`
	if runtime.GOOS == "linux" {
		cmdBase = `/opt/cisco/anyconnect/bin/vpnui`
	}

	cmdExe1 := exec.Command("bash", "-c", cmdBase+" state")
	cmdGrep1 := `grep -c 'state: Connected'`
	exeGrep1 := exec.Command("bash", "-c", cmdGrep1)
	output1, _, _ := cmdhelp.Pipeline(cmdExe1, exeGrep1)

	vpnConnStatus, _ := strconv.Atoi(strings.Split(strings.TrimRight(string(output1), "\n"), " ")[0])

	t := table.NewWriter()
	t.SetTitle("VPN Connection Status Checks")
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Property Description", "Value", "Notes"})
	t.AppendRow([]interface{}{"VPN Client reports connection status as 'Connected'?", vpnConnStatus > 0})
	t.AppendSeparator()
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMin: 50},
		{Number: 2, WidthMin: 30},
	})
	t.Render()

	if vpnConnStatus == 0 {
		fmt.Println("")
		color.Warn.Tips(`

   Your VPN client's does not appear to be a state of 'connected',
   it's either down or misconfigured!`)
	}

	fmt.Printf("\n\n\n")
}
