package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ngxcpd",
	Short: "nginx cache purge daemon.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	_ = rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP(
		"verbose",
		"v",
		false,
		"Be more verbose")

	rootCmd.PersistentFlags().BoolP(
		"debug",
		"d",
		false,
		"Be very verbose")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// init logger
		log.SetLevel(log.WarnLevel)

		if v, err := cmd.Flags().GetBool("verbose"); err == nil && v {
			log.SetLevel(log.InfoLevel)
		}

		if d, err := cmd.Flags().GetBool("debug"); err == nil && d {
			log.SetLevel(log.DebugLevel)
		}
	}

}
