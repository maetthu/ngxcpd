package cmd

import (
	"fmt"
	"github.com/maetthu/ngxcpd/internal/lib/proxycache"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync/atomic"
)

func init() {
	localpurgeCmd.Flags().BoolP(
		"dryrun",
		"n",
		false,
		"Just tell what to delete, but don't do it")

	rootCmd.AddCommand(localpurgeCmd)
}

var localpurgeCmd = &cobra.Command{
	Use:   "localpurge <dir> <pattern>",
	Short: "Scan directory and remove cache files which key matches regular expression pattern",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]

		if f, err := os.Stat(dir); err != nil {
			log.Fatal(err)
		} else if !f.IsDir() {
			log.Fatalf("%s is not a directory", dir)
		}

		re, err := regexp.Compile(args[1])

		if err != nil {
			log.Fatal(err)
		}

		dryrun, _ := cmd.Flags().GetBool("dryrun")

		if dryrun {
			fmt.Print("Dry run: not deleting any files, but following entries match:\n\n")
			fmt.Printf("Filename%s\tKey\n", strings.Repeat(" ", 24))
		}

		var count uint32

		callback := func(entry *proxycache.Entry) {
			base := filepath.Base(entry.Filename)

			log.Debugf("[%s] Start processing", base)
			log.Debugf("[%s] Key: %s", base, entry.Key)

			if re.MatchString(entry.Key) {
				if dryrun {
					log.Infof("[%s] Key \"%s\" is matching", base, entry.Key)
					fmt.Printf("%s\t%s\n", base, entry.Key)
					atomic.AddUint32(&count, 1)
				} else {
					log.Infof("[%s] Key \"%s\" is matching, deleting", base, entry.Key)

					if err := os.Remove(entry.Filename); err != nil {
						log.Error(err)
					} else {
						atomic.AddUint32(&count, 1)
					}
				}
			} else {
				log.Debugf("[%s] No match", base)
			}

			log.Debugf("[%s] End processing", base)
		}

		log.Infof("Start scanning %s", dir)

		if err := proxycache.ScanDir(dir, callback, runtime.NumCPU()); err != nil {
			log.Fatal(err)
		}

		log.Infof("Finished scanning %s", dir)

		if !dryrun {
			if count == 0 {
				fmt.Println("No cache entries matched, nothing removed")
			} else {
				fmt.Printf("Removed %d cache entries\n", count)
			}
		} else {
			if count == 0 {
				fmt.Println("No cache entries matched")
			} else {
				fmt.Printf("A real run would remove %d cache entries\n", count)
			}
		}

	},
}
