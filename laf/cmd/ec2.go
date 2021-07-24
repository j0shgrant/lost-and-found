package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/j0shgrant/lost-and-found/internal/aws"
	"github.com/j0shgrant/lost-and-found/internal/aws/ec2/instance"
	"github.com/j0shgrant/lost-and-found/internal/tags"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"sync"
)

var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Lists EC2 instances that match a passed set of filters",
	Run: func(cmd *cobra.Command, _ []string) {
		// parse flags
		regionFlag := cmd.Flag("region").Value.String()
		requiredTagsFlag := cmd.Flag("required-tags").Value.String()
		excludedTagsFlag := cmd.Flag("excluded-tags").Value.String()

		// build list of regions
		var regions []string
		for _, region := range strings.Split(regionFlag, ",") {
			if region != "" {
				regions = append(regions, region)
			}
		}

		// build list of tags
		requiredTags, err := tags.ParseTags(requiredTagsFlag)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		excludedTags, err := tags.ParseTags(excludedTagsFlag)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		// build configs for regions
		configs, err := aws.NewConfigsFromRegions(regions)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		// create instance services for regions
		var services []instance.Service
		for _, config := range configs {
			services = append(services, instance.NewServiceForRegion(config))
		}

		// list filtered instances by region
		var wg sync.WaitGroup
		var instances []instance.Instance
		for _, service := range services {
			service := service
			wg.Add(1)
			go func() {
				is, err := service.Get(requiredTags)
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, err.Error())
					os.Exit(1)
				}

				instances = append(instances, is...)
				wg.Done()
			}()
		}
		wg.Wait()

		var filteredInstances []types.Instance
		for _, instance := range instances {
			if !tags.HasExcludedTags(excludedTags, instance.Tags()) {
				filteredInstances = append(filteredInstances, instance.Instance)
			}
		}

		for _, instance := range filteredInstances {
			fmt.Println(*instance.InstanceId)
		}
	},
}

func init() {
	rootCmd.AddCommand(ec2Cmd)
}
