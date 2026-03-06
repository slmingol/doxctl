package cmd
package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"



























}	rootCmd.AddCommand(versionCmd)func init() {}	},		fmt.Printf("  Platform:  %s/%s\n", runtime.GOOS, runtime.GOARCH)		fmt.Printf("  Go:        %s\n", runtime.Version())		fmt.Printf("  Built by:  %s\n", builtBy)		fmt.Printf("  Built:     %s\n", date)		fmt.Printf("  Commit:    %s\n", commit)		fmt.Printf("doxctl %s\n", version)	Run: func(cmd *cobra.Command, args []string) {	Long:  `Print the version number, build date, and commit information for doxctl.`,	Short: "Print the version number of doxctl",	Use:   "version",var versionCmd = &cobra.Command{// versionCmd represents the version command)	builtBy = "unknown"	date    = "unknown"	commit  = "none"	version = "dev"var ()