package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func newConfigsFromRegions(regions []string) (map[string]aws.Config, error) {
	// load in base AWS config
	baseCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// return default region if regions list is empty
	if len(regions) < 1 {
		configs := map[string]aws.Config{
			baseCfg.Region: baseCfg,
		}

		return configs, nil
	}

	// generate list of existing AWS regions
	output, err := ec2.NewFromConfig(baseCfg).DescribeRegions(context.TODO(), &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}

	var actualRegions []string
	for _, region := range output.Regions {
		actualRegions = append(actualRegions, *region.RegionName)
	}

	// validate passed regions all exist
	for _, region := range regions {
		regionExists := false
		for _, actualRegion := range actualRegions {
			if actualRegion == region {
				regionExists = true
				break
			}
		}

		if !regionExists {
			return nil, fmt.Errorf("passed AWS region %s does not exist", region)
		}
	}

	// build map of region => aws.Config
	configs := make(map[string]aws.Config)
	for _, region := range regions {
		cfg := baseCfg.Copy()
		cfg.Region = region
		configs[region] = cfg
	}

	return configs, nil
}
