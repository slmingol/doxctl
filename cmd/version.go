package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of doxctl",
	Long:  `Print the version number, build date, and commit information for doxctl.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("doxctl %s\n", version)
		fmt.Printf("  Commit:    %s\n", commit)
		fmt.Printf("  Built:     %s\n", date)
		fmt.Printf("  Built by:  %s\n", builtBy)
		fmt.Printf("  Go:        %s\n", runtime.Version())
		fmt.Printf("  Platform:  %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
