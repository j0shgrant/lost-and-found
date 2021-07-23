package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	version    string
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print current laf version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(version)
			os.Exit(0)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
