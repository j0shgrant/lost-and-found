package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/j0shgrant/lost-and-found/internal/aws"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Lists EC2 instances that match a passed set of filters",
	Run: func(cmd *cobra.Command, _ []string) {
		// parse flags
		regionFlag := cmd.Flag("region").Value.String()
		requiredTagsFlag := cmd.Flag("required-tags").Value.String()

		// build list of regions
		var regions []string
		for _, region := range strings.Split(regionFlag, ",") {
			if region != "" {
				regions = append(regions, region)
			}
		}

		// build list of tags
		requiredTags, err := aws.ParseTags(requiredTagsFlag)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		// build EC2 service
		ec2, err := aws.NewEC2Service(regions)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "encountered an error initialising EC2 client: %s\n", err.Error())
			os.Exit(1)
		}

		instances, err := ec2.ListInstances()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "encountered an error listing EC2 instances: %s\n", err.Error())
			os.Exit(1)
		}

		// Filter instances
		var filteredInstances []types.Instance
		for _, instance := range instances {
			if !aws.EC2InstanceHasRequiredTags(requiredTags, instance) {
				continue
			}

			filteredInstances = append(filteredInstances, instance)
		}

		for _, instance := range filteredInstances {
			fmt.Println(*instance.InstanceId)
		}
	},
}

func init() {
	rootCmd.AddCommand(ec2Cmd)
}
