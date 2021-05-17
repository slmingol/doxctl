// Package cmd - dns implements DNS diagnostic checks
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Run diagnostics related to DNS servers' (resolvers') configurations",
	Long: `
doxctl's 'dns' subcommand can help triage DNS resovler configuration issues, 
general access to DNS resolvers and name resolution against DNS resolvers.`,
	Run: func(cmd *cobra.Command, args []string) {
		dnsDiag()
	},
}

var resolverChk, pingChk, digChk, allChk bool

func init() {
	rootCmd.AddCommand(dnsCmd)

	dnsCmd.Flags().BoolVarP(&resolverChk, "resolverChk", "r", false, "Check if VPN designated DNS resolvers are configured")
	dnsCmd.Flags().BoolVarP(&pingChk, "pingChk", "p", false, "Check if VPN defined resolvers are pingable & reachable")
	dnsCmd.Flags().BoolVarP(&digChk, "digChk", "d", false, "Check if VPN defined resolvers respond with well-known servers in DCs")
	dnsCmd.Flags().BoolVarP(&allChk, "allChk", "a", false, "Run all the checks in this subcommand module")
}

func dnsDiag() {
	cmd := exec.Command("")

	var verboseCmd string

	if verboseChk {
		verboseCmd = "1"
	} else {
		verboseCmd = "0"
	}

	switch {
	case resolverChk:
		cmd = exec.Command("bash", "-c", ". model_cmds/02_dns.sh; dnsResolverChk"+" "+verboseCmd)
	case pingChk:
		cmd = exec.Command("bash", "-c", ". model_cmds/02_dns.sh; dnsResolverPingChk"+" "+verboseCmd)
	case digChk:
		cmd = exec.Command("bash", "-c", ". model_cmds/02_dns.sh; dnsResolverDigChk"+" "+verboseCmd)
	case allChk:
		cmd = exec.Command("bash", "-c", ". model_cmds/02_dns.sh; dnsResolverChk"+" "+verboseCmd)
		cmd = exec.Command("bash", "-c", ". model_cmds/02_dns.sh; dnsResolverPingChk"+" "+verboseCmd)
		cmd = exec.Command("bash", "-c", ". model_cmds/02_dns.sh; dnsResolverDigChk"+" "+verboseCmd)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
