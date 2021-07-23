package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	region, tag string

	rootCmd = &cobra.Command{
		Use:  "laf",
		Long: "laf - lost-and-found is a CLI utility for finding AWS resources that are untagged, non-compliant, or in some way misbehaving.",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "A comma-separated list of AWS regions to include: \"eu-west-1,eu-west-2\"")
	rootCmd.PersistentFlags().StringVar(&tag, "required-tags", "", "A comma-separated list of AWS tags to match on: eg. only return resources with tags \"foo=bar,bar,qux=quux\"")
}
