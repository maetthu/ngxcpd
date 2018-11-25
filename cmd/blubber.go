package cmd

import (
	"fmt"
	"github.com/maetthu/ngxcpd/internal/lib/proxycache"
	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(blubberCmd)
}

var blubberCmd = &cobra.Command{
	Use: "blubber",
	Run: func(cmd *cobra.Command, args []string) {

		proxycache.ScanDir(args[0], func(entry *proxycache.Entry) {
			s := litter.Sdump(entry)
			fmt.Println(s)
		})

	},
}
