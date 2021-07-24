package instance

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/j0shgrant/lost-and-found/internal/tags"
)

type Instance struct {
	Instance types.Instance
}

func (i Instance) Tags() []tags.Tag {
	var ts []tags.Tag
	for _, instanceTag := range i.Instance.Tags {
		t := tags.Tag{
			Key:      *instanceTag.Key,
			Value:    *instanceTag.Value,
			Wildcard: false,
		}

		ts = append(ts, t)
	}

	return ts
}
