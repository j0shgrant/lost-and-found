package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"strings"
)

type Tag struct {
	Key      string
	Value    string
	Wildcard bool
}

func ParseTags(flagValue string) ([]Tag, error) {
	// return empty array for empty flag string
	if flagValue == "" {
		return make([]Tag, 0), nil
	}

	// parse tags
	var tags []Tag
	for _, tagString := range strings.Split(flagValue, ",") {
		tag := Tag{}
		tokens := strings.Split(tagString, "=")
		tag.Key = tokens[0]

		if len(tokens) == 1 {
			tag.Wildcard = true
		}
		if len(tokens) > 1 {
			tag.Wildcard = false
			tag.Value = tokens[1]
		}

		if len(tokens) > 2 {
			return nil, fmt.Errorf("tag %s is invalid, must be one of the following: \"KEY=VALUE\" or \"KEY\"", tagString)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func EC2FiltersFromTags(tags []Tag) []types.Filter {
	var filters []types.Filter
	for _, tag := range tags {
		var filter types.Filter
		// create tag-key filter if tag is wildcard
		if tag.Wildcard {
			filter = types.Filter{
				Name:   aws.String("tag-key"),
				Values: []string{tag.Key},
			}
		} else {
			filter = types.Filter{
				Name:   aws.String(fmt.Sprintf("tag:%s", tag.Key)),
				Values: []string{tag.Value},
			}
		}

		filters = append(filters, filter)
	}

	return filters
}
