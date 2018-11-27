package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(blubberCmd)
}

var blubberCmd = &cobra.Command{
	Use: "blubber",
	Run: func(cmd *cobra.Command, args []string) {
		// whatever
	},
}
