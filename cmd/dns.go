// Package dnsCmd implements DNS diagnostic checks
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// dnsCmd represents the dns command
var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		diag()
	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dnsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dnsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func diag() {
	fmt.Println("dns diag")

	// get `model_cmds/02_dns.sh` executable
	cmdExecutable, _ := exec.LookPath("model_cmds/02_dns.sh")

	// `model_cmds/02_dns.sh` command
	cmdDNS := &exec.Cmd{
		Path:   cmdExecutable,
		Args:   []string{cmdExecutable, ""},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	fmt.Println(cmdDNS.String())

	// see command represented by `cmdDns`
	if err := cmdDNS.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
