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
	"strings"
	"time"

	"github.com/go-ping/ping"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	gobrex "github.com/kujtimiihoxha/go-brace-expansion"
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

		// Populate as much info as we can from viper
		err := viper.Unmarshal(&conf)
		if err != nil {
			fmt.Printf("could not retrieve supplied project settings: %s\n", err)
			os.Exit(1)
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
		cmd.Usage()
		fmt.Printf("\n\n\n")
		os.Exit(1)
	}
}

// Check if well known servers are pingable & reachable
func svrsReachChk() {
	color.Info.Tips("Attempting to ping all well known servers, this may take a few...\n")

	// Table head
	t := table.NewWriter()
	t.SetTitle("Well known Servers Reachable Checks")
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMin: 40},
		{Number: 2, WidthMin: 20},
	})
	t.AppendHeader(table.Row{"Host", "Service", "Reachable?", "Ping Performance"})

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

	for _, i := range conf.Svcs {
		fmt.Printf("   --- Working through svc: %s\n", i.Svc)

		for _, j := range i.Svrs {
			permutations := gobrex.Expand(j)

			for _, permutation := range permutations {

				// if FailThreshold is exceeded, stop trying pings, call the rest failed
				if pingFailures > conf.FailThreshold || reachFailures > conf.FailThreshold {
					t.AppendRow([]interface{}{permutation, i.Svc, false, "N/A"})
					continue
				}

				// Attempt to ping each host, any that fail keep a tally of how many
				pinger, err := ping.NewPinger(permutation)
				if err != nil {
					t.AppendRow([]interface{}{permutation, i.Svc, false, "N/A"})
					pingFailures++
					continue
				}

				pinger.Timeout = conf.PingTimeout * time.Millisecond
				pinger.Run()
				stats := pinger.Statistics()
				pingPerf := fmt.Sprintf("rnd-trp avg = %v", stats.AvgRtt)

				// Tally fails due to failed/missing responses
				packetAck := (stats.PacketLoss == 0 && stats.PacketsRecv > 0)
				if !packetAck {
					reachFailures++
				}

				t.AppendRow([]interface{}{permutation, i.Svc, packetAck, pingPerf})
			}
		}
		t.AppendSeparator()
	}
	fmt.Printf("\n\n   ...one sec, preparing `ping` results...\n\n")

	if pingFailures > conf.FailThreshold || reachFailures > conf.FailThreshold {
		fmt.Println("")
		color.Warn.Tips("More than %d hosts appear to be unreachable, aborting remainder....\n\n", conf.FailThreshold)
	}

	time.Sleep(6 * time.Second)

	t.AppendSeparator()
	t.Render()

	if pingFailures > 0 || reachFailures > 0 {
		fmt.Println("")
		color.Warn.Tips(`

   Your VPN client does not appear to be functioning properly, it's likely one or more of the following:

      - Well known servers are unreachable via ping   --- try running 'doxctl vpn -h'
      - Servers are unresovlable in DNS               --- try running 'doxctl dns -h'
      - VPN client is otherwise misconfigured!
	`)
	}

	fmt.Printf("\n\n\n")
}
