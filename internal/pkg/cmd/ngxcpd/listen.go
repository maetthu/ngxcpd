package ngxcpd

import (
	"errors"
	"github.com/maetthu/ngxcpd/pkg/zone"
	"github.com/spf13/cobra"
	"log"
	"runtime"
	"time"
)

func init() {
	rootCmd.AddCommand(listenCmd)
}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "TODO...",
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

		log.Println(z.Path)

		log.Println("Ready to watch...")

		t, err := z.Watch(4 * 4096)

		if err != nil {
			log.Fatal(err)
		}

		time.AfterFunc(60*time.Second, func() {
			t.Kill(errors.New("Cancel"))
		})

		if err := t.Wait(); err != nil {
			log.Println(err)
		}

	},
}
