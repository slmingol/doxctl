// Package cmd - vpn implements VPN diagnstic checks
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
	Run: func(cmd *cobra.Command, args []string) {
		vpnDiag()
	},
}

var ifReachableChk, vpnRoutesChk bool

func init() {
	rootCmd.AddCommand(vpnCmd)

	vpnCmd.Flags().BoolVarP(&ifReachableChk, "ifReachableChk", "i", false, "Check if network interfaces are reachable")
	vpnCmd.Flags().BoolVarP(&vpnRoutesChk, "vpnRoutesChk", "r", false, "Check if >5 VPN routes are defined")
	vpnCmd.Flags().BoolVarP(&allChk, "allChk", "a", false, "Run all the checks in this subcommand module")
}

func vpnDiag() {
	cmd := exec.Command("")

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
	}

	cmd = exec.Command("bash", "-c", cmdString)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
