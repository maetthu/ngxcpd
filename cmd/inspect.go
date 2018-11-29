package cmd

import (
	"fmt"
	"github.com/maetthu/ngxcpd/internal/lib/proxycache"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func init() {
	inspectCmd.Flags().BoolP(
		"response",
		"r",
		false,
		"Also display response headers")

	inspectCmd.Flags().BoolP(
		"testfixture",
		"t",
		false,
		"Display data formatted for inclusion in test data")

	rootCmd.AddCommand(inspectCmd)
}

var inspectCmd = &cobra.Command{
	Use:   "inspect <filename> [<filename>...]",
	Short: "Inspect displays meta data of a cache file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, f := range args {
			e, err := proxycache.FromFile(f)

			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}

			if t, err := cmd.Flags().GetBool("testfixture"); err == nil && t {
				fmt.Printf(`	{
		Filename:     %s,
		Version:      5,
		Expire:       time.Unix(%d, 0),
		LastModified: time.Unix(%d, 0),
		Date:         time.Unix(%d, 0),
		Etag:         %s,
		Key:          %s,
		HeaderStart:  %d,
		BodyStart:    %d,
		RawHeader:    %s,
	},
`,
					strconv.Quote(e.Filename),
					e.Expire.Unix(),
					e.LastModified.Unix(),
					e.Date.Unix(),
					strconv.Quote(e.Etag),
					strconv.Quote(e.Key),
					e.HeaderStart,
					e.BodyStart,
					strconv.Quote(e.RawHeader),
				)
			} else {
				fmt.Printf(`%s
  Hash: %s
  Key: %s
  Etag: %s
  Expires: %v
  Last Modified: %v
  Date: %v
`,
					e.Filename,
					func() string { h, _ := e.Hash(); return h }(),
					e.Key,
					e.Etag,
					e.Expire,
					e.LastModified,
					e.Date,
				)

				if r, err := cmd.Flags().GetBool("response"); err == nil && r {
					fmt.Printf("  Response:\n")

					res, err := e.Response()

					if err != nil {
						fmt.Printf("    Error parsing headers: %s\n", err)
					}

					fmt.Printf("    %s %s\n", res.Proto, res.Status)
					fmt.Printf("    Headers:\n")

					for h, v := range res.Header {
						fmt.Printf("      %s: %s\n", h, strings.Join(v, ";"))
					}
				}
			}
		}
	},
}
