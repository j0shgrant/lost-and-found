package instance

import (
	"context"
	config "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	filters "github.com/j0shgrant/lost-and-found/internal/aws/ec2"
	"github.com/j0shgrant/lost-and-found/internal/tags"
)

type Service struct {
	Client *ec2.Client
}

func NewServiceForRegion(cfg config.Config) Service {
	return Service{
		Client: ec2.NewFromConfig(cfg),
	}
}

func (s Service) Get(tags []tags.Tag) ([]Instance, error) {
	// parse EC2 filters
	instanceFilters := filters.FiltersFromTags(tags)

	// query for all reservations matching filters
	var next *string
	var reservations []types.Reservation
	for {
		output, err := s.Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
			Filters:   instanceFilters,
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

	// flatten reservations into instances
	var instances []types.Instance
	for _, reservation := range reservations {
		instances = append(instances, reservation.Instances...)
	}

	// wrap instances as TaggableResources
	var taggableInstances []Instance
	for _, instance := range instances {
		taggableInstances = append(taggableInstances, Instance{Instance: instance})
	}

	return taggableInstances, nil
}
