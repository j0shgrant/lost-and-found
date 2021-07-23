package aws

import (
	"context"
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

func (s *EC2Service) ListInstances() ([]types.Instance, error) {
	var regions []string
	for region := range s.clients {
		regions = append(regions, region)
	}
	sort.Strings(regions)

	var instances []types.Instance
	for _, region := range regions {
		// Pull all reservations
		var next *string
		var reservations []types.Reservation
		for {
			output, err := s.clients[region].DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
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

func newEC2ClientsFromConfigs(configs map[string]aws.Config) map[string]*ec2.Client {
	clients := make(map[string]*ec2.Client)
	for region, cfg := range configs {
		clients[region] = ec2.NewFromConfig(cfg)
	}

	return clients
}
