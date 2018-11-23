package cmd

import (
	"fmt"

	"github.com/maetthu/ngxcpd/internal/lib/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ngxcpd",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ngxcpd %s -- %s\n", version.Version, version.Commit)
	},
}
