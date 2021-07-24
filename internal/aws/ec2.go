package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"sort"
)

type EC2Service struct {
	clients map[string]*ec2.Client
}

func NewEC2Service(regions []string) (*EC2Service, error) {
	configs, err := newConfigsFromRegions(regions)
	if err != nil {
		return nil, err
	}

	clients := newEC2ClientsFromConfigs(configs)

	service := &EC2Service{
		clients: clients,
	}

	return service, err
}

func (s *EC2Service) ListInstances(filters []types.Filter) ([]types.Instance, error) {
	var regions []string
	for region := range s.clients {
		regions = append(regions, region)
	}
	sort.Strings(regions)

	for _, filter := range filters {
		fmt.Println(*filter.Name)
	}

	var instances []types.Instance
	for _, region := range regions {
		// Pull all reservations
		var next *string
		var reservations []types.Reservation
		for {
			output, err := s.clients[region].DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
				Filters: filters,
				NextToken: next,
			})
			if err != nil {
				return nil, err
			}

			reservations = append(reservations, output.Reservations...)

			next = output.NextToken

			if next == nil {
				break
			}
		}

		// Flatten instances from reservations
		var regionInstances []types.Instance
		for _, reservation := range reservations {
			regionInstances = append(regionInstances, reservation.Instances...)
		}

		instances = append(instances, regionInstances...)
	}

	return instances, nil
}

func EC2InstanceHasExcludedTags(excludedTags []Tag, instance types.Instance) bool {
	// return true on empty excluded tag array
	if len(excludedTags) == 0 {
		return false
	}

	// build map of instance tags
	instanceTags := make(map[string]string)
	for _, tag := range instance.Tags {
		instanceTags[*tag.Key] = *tag.Value
	}

	// compare instance tags
	for _, excludedTag := range excludedTags {
		if instanceTagValue, exists := instanceTags[excludedTag.Key]; exists {
			if !excludedTag.Wildcard {
				if excludedTag.Value == instanceTagValue {
					return true
				}

				continue
			}

			return true
		}
	}
	return false
}

func newEC2ClientsFromConfigs(configs map[string]aws.Config) map[string]*ec2.Client {
	clients := make(map[string]*ec2.Client)
	for region, cfg := range configs {
		clients[region] = ec2.NewFromConfig(cfg)
	}

	return clients
}
