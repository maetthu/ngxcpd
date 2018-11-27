package cmd

import (
	"github.com/maetthu/ngxcpd/internal/lib/zone"
	"github.com/spf13/cobra"
	"log"
	"runtime"
)

func init() {
	rootCmd.AddCommand(listenCmd)
}

var listenCmd = &cobra.Command{
	Use: "listen",
	Run: func(cmd *cobra.Command, args []string) {

		z, err := zone.NewZone(args[0])

		if err != nil {
			log.Fatal(err)
		}

		log.Println("Start initial index")

		go func() {
			if err := z.Warmup(runtime.NumCPU()); err != nil {
				log.Fatal(err)
			} else {
				log.Println("Finished initial directory scan")
			}

			log.Println(z.Cache.ItemCount())
		}()

		log.Println("Ready to watch...")

		if err := z.Watch(); err != nil {
			log.Fatal(err)
		}

	},
}
