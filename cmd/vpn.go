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
	"doxctl/internal/output"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

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
		_ = cmd.Usage()
		fmt.Printf("\n")
		os.Exit(1)
	}
}

// Test if interface is reported as reachable via 'scutil'
func ifReachChk() {
	ifReachChkWithDeps(NewCommandExecutor())
}

// ifReachChkWithDeps allows dependency injection for testing
func ifReachChkWithDeps(executor CommandExecutor) {
	// Check if running in container with host VPN data
	hostIfCount := os.Getenv("HOST_VPN_IF_COUNT")
	hostNetIfs := os.Getenv("HOST_VPN_NET_IFS")
	hostHasTun := os.Getenv("HOST_VPN_HAS_TUN")
	hostTunIfs := os.Getenv("HOST_VPN_TUN_IFS")
	hostAllReachable := os.Getenv("HOST_VPN_ALL_IFS_REACHABLE")

	if hostIfCount != "" {
		// Use host-provided data
		ifCount, _ := strconv.Atoi(hostIfCount)
		netIfs := strings.Fields(hostNetIfs)
		foundOneTunIf := (hostHasTun == "true")
		tunIfs := strings.Fields(hostTunIfs)
		allInfsReachable := (hostAllReachable == "true")

		// For JSON/YAML output
		if outputFormat != "table" {
			result := output.VPNInterfaceCheckResult{
				Timestamp:              time.Now(),
				InterfaceCount:         ifCount,
				Interfaces:             netIfs,
				HasTunInterface:        foundOneTunIf,
				TunInterfaces:          tunIfs,
				AllInterfacesReachable: allInfsReachable,
			}
			output.Print(outputFormat, result)
			return
		}

		// Table output
		t := table.NewWriter()
		t.SetTitle("Interfaces Reachable Checks")
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		t.AppendHeader(table.Row{"Property Description", "Value", "Notes"})
		t.AppendRow([]interface{}{"How many network interfaces found?", ifCount, netIfs})
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
		return
	}

	// Original logic for non-container execution
	cmdBase := `scutil --nwi`

	output1, err := executor.Execute("bash", "-c", cmdBase+" | grep 'Network interfaces:' | cut -d\" \" -f 3-")
	if err != nil {
		output1 = []byte("")
	}

	netIfs := strings.Split(strings.TrimRight(string(output1), "\n"), " ")

	var tunIfs []string
	for i := 0; i < len(netIfs); i++ {
		netIfs[i] = strings.TrimSpace(netIfs[i])
		if strings.Contains(netIfs[i], "tun") {
			tunIfs = append(tunIfs, netIfs[i])
		}
	}

	var foundOneTunIf = false
	if len(tunIfs) > 0 {
		foundOneTunIf = true
	}

	output2, err := executor.Execute("bash", "-c", cmdBase+" | grep address -B1 -A1 | grep -E \"flags|reach\" | paste - - | column -t | grep -v Reachable | wc -l | tr -d ' '")
	if err != nil {
		output2 = []byte("0")
	}

	reachableIfs := strings.TrimRight(string(output2), "\n")

	var allInfsReachable = false
	if reachableIfs == "0" {
		allInfsReachable = true
	}

	// For JSON/YAML output
	if outputFormat != "table" {
		result := output.VPNInterfaceCheckResult{
			Timestamp:              time.Now(),
			InterfaceCount:         len(netIfs),
			Interfaces:             netIfs,
			HasTunInterface:        foundOneTunIf,
			TunInterfaces:          tunIfs,
			AllInterfacesReachable: allInfsReachable,
		}
		output.Print(outputFormat, result)
		return
	}

	// Table output
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
}

// Test if VPNs interface defines at least MinVpnRoutes routes
func vpnRteChk() {
	vpnRteChkWithDeps(NewCommandExecutor())
}

// vpnRteChkWithDeps allows dependency injection for testing
func vpnRteChkWithDeps(executor CommandExecutor) {
	// Check if running in container with host VPN data
	hostVpnIf := os.Getenv("HOST_VPN_INTERFACE")
	hostRouteCount := os.Getenv("HOST_VPN_ROUTE_COUNT")

	if hostVpnIf != "" && hostRouteCount != "" {
		// Use host-provided data
		vpnIf := hostVpnIf
		vpnRouteCnt, _ := strconv.Atoi(hostRouteCount)

		// For JSON/YAML output
		if outputFormat != "table" {
			result := output.VPNRoutesCheckResult{
				Timestamp:           time.Now(),
				VPNInterface:        vpnIf,
				RouteCount:          vpnRouteCnt,
				MinRoutesRequired:   conf.MinVpnRoutes,
				HasSufficientRoutes: vpnRouteCnt >= conf.MinVpnRoutes,
			}
			output.Print(outputFormat, result)
			return
		}

		// Table output
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
		return
	}

	// Original logic for non-container execution
	output1, err := executor.Execute("bash", "-c", "scutil --nwi | grep 'Network interfaces:' | grep -o utun[0-9] || echo \"NIL\"")
	if err != nil {
		output1 = []byte("NIL")
	}

	vpnIf := strings.Split(strings.TrimRight(string(output1), "\n"), " ")[0]

	output2, err := executor.Execute("bash", "-c", "netstat -r -f inet | grep -c "+vpnIf) // #nosec G204 - vpnIf is from system scutil output
	if err != nil {
		output2 = []byte("0")
	}

	vpnRouteCnt, _ := strconv.Atoi(strings.Split(strings.TrimRight(string(output2), "\n"), " ")[0])

	// For JSON/YAML output
	if outputFormat != "table" {
		result := output.VPNRoutesCheckResult{
			Timestamp:           time.Now(),
			VPNInterface:        vpnIf,
			RouteCount:          vpnRouteCnt,
			MinRoutesRequired:   conf.MinVpnRoutes,
			HasSufficientRoutes: vpnRouteCnt >= conf.MinVpnRoutes,
		}
		output.Print(outputFormat, result)
		return
	}

	// Table output
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
}

// Test if VPN connection status reports as 'connected'
func vpnConnChk() {
	vpnConnChkWithDeps(NewCommandExecutor())
}

// vpnConnChkWithDeps allows dependency injection for testing
func vpnConnChkWithDeps(executor CommandExecutor) {
	// Check if running in container with host VPN data
	hostVpnConnected := os.Getenv("HOST_VPN_CONNECTED")
	hostVpnClient := os.Getenv("HOST_VPN_CLIENT")

	if hostVpnConnected != "" {
		// Use host-provided data
		vpnConnStatus := 0
		if hostVpnConnected == "true" {
			vpnConnStatus = 1
		}

		// For JSON/YAML output
		if outputFormat != "table" {
			result := output.VPNConnectionStatusResult{
				Timestamp:   time.Now(),
				IsConnected: vpnConnStatus > 0,
			}
			output.Print(outputFormat, result)
			return
		}

		// Table output
		t := table.NewWriter()
		t.SetTitle("VPN Connection Status Checks")
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		t.AppendHeader(table.Row{"Property Description", "Value", "Notes"})

		// Customize description based on client type
		var description string
		if hostVpnClient == "anyconnect" {
			description = "VPN Client (AnyConnect) reports connection status as 'Connected'?"
		} else if hostVpnClient == "generic" {
			description = "VPN Connection detected (via TUN interface + routes)?"
		} else {
			description = "VPN Client reports connection status as 'Connected'?"
		}

		t.AppendRow([]interface{}{description, vpnConnStatus > 0})
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
		return
	}

	// Original logic for non-container execution
	cmdBase := `/opt/cisco/anyconnect/bin/vpn`
	if runtime.GOOS == "linux" {
		cmdBase = `/opt/cisco/anyconnect/bin/vpnui`
	}

	output1, err := executor.Execute("bash", "-c", cmdBase+" state | grep -c 'state: Connected'")
	if err != nil {
		output1 = []byte("0")
	}

	vpnConnStatus, _ := strconv.Atoi(strings.Split(strings.TrimRight(string(output1), "\n"), " ")[0])

	// For JSON/YAML output
	if outputFormat != "table" {
		result := output.VPNConnectionStatusResult{
			Timestamp:   time.Now(),
			IsConnected: vpnConnStatus > 0,
		}
		output.Print(outputFormat, result)
		return
	}

	// Table output
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
}
