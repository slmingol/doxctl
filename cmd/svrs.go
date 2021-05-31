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

	gobrex "github.com/kujtimiihoxha/go-brace-expansion"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// svrsCmd represents the svrs command
var svrsCmd = &cobra.Command{
	Use:   "svrs",
	Short: "Run diagnostics related to server accessiblity through a VPN connection",
	Long: `
doxctl's 'svrs' subcommand can help triage & test connectivity to 'well-known servers'
through a VPN connection to servers which have been defined in your '.doxctl.yaml' 
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

	svrsCmd.Flags().BoolVarP(&svrsReachableChk, "svrsReachableChk", "s", false, "Check if well-known servers are reachable")
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
		os.Exit(1)
	}
}

func svrsReachChk() {
	for _, i := range conf.Svcs {
		fmt.Printf("svc: %s\n", i.Svc)
		fmt.Printf("svrs: %s\n", i.Svrs)
		for _, j := range i.Svrs {
			permutations := gobrex.Expand(j)
			for _, permutation := range permutations {
				fmt.Println(permutation)
			}
		}
	}
}
