package aws

import (
	"fmt"
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
