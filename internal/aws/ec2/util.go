package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/j0shgrant/lost-and-found/internal/tags"
)

func FiltersFromTags(tags []tags.Tag) []types.Filter {
	var filters []types.Filter
	for _, t := range tags {
		var filter types.Filter
		// create tag-key filter if tag is wildcard
		if t.Wildcard {
			filter = types.Filter{
				Name:   aws.String("tag-key"),
				Values: []string{t.Key},
			}
		} else {
			filter = types.Filter{
				Name:   aws.String(fmt.Sprintf("tag:%s", t.Key)),
				Values: []string{t.Value},
			}
		}

		filters = append(filters, filter)
	}

	return filters
}
