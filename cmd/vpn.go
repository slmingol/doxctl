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
	"os/exec"

	"github.com/spf13/cobra"
)

// vpnCmd represents the vpn command
var vpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "Run diagnostics related to VPN connections, net i/fs & configurations",
	Long: `
doxctl's 'vpn' subcommand can help triage VPN related configuration issues,
& routes related to a VPN connection.`,
	Run: vpnExecute,
}

var ifReachableChk, vpnRoutesChk bool

func init() {
	rootCmd.AddCommand(vpnCmd)

	vpnCmd.Flags().BoolVarP(&ifReachableChk, "ifReachableChk", "i", false, "Check if network interfaces are reachable")
	vpnCmd.Flags().BoolVarP(&vpnRoutesChk, "vpnRoutesChk", "r", false, "Check if >5 VPN routes are defined")
	vpnCmd.Flags().BoolVarP(&allChk, "allChk", "a", false, "Run all the checks in this subcommand module")
}

func vpnExecute(cmd *cobra.Command, args []string) {
	exeCmd := exec.Command("")

	var verboseCmd string

	if verboseChk {
		verboseCmd = "1"
	} else {
		verboseCmd = "0"
	}

	var cmdString string

	switch {
	case ifReachableChk:
		cmdString = ". model_cmds/01_vpn.sh; netInterfacesReachableChk" + " " + verboseCmd
	case vpnRoutesChk:
		cmdString = ". model_cmds/01_vpn.sh; vpnInterfaceRoutesChk" + " " + verboseCmd
	case allChk:
		cmdString = ". model_cmds/01_vpn.sh" +
			"; netInterfacesReachableChk" + " " + verboseCmd +
			"; vpnInterfaceRoutesChk" + " " + verboseCmd
	default:
		cmd.Usage()
		os.Exit(1)
	}

	exeCmd = exec.Command("bash", "-c", cmdString)
	exeCmd.Stdout = os.Stdout
	exeCmd.Stderr = os.Stdout

	if err := exeCmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
