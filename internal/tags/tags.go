package tags

import (
	"fmt"
	"strings"
)

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

func HasExcludedTags(excludedTags []Tag, actualTags []Tag) bool {
	// return false if no excluded tags are passed
	if len(excludedTags) == 0 {
		return false
	}

	// build map of actual tags
	actualTagMap := make(map[string]string)
	for _, actualTag := range actualTags {
		actualTagMap[actualTag.Key] = actualTag.Value
	}

	// compare tags
	for _, excludedTag := range excludedTags {
		if actualTagValue, exists := actualTagMap[excludedTag.Key]; exists {
			if !excludedTag.Wildcard {

				// return true if there is a complete tag match
				if excludedTag.Value == actualTagValue {
					return true
				}

				continue
			}

			// return true if a wildcard tag key exists
			return true
		}
	}

	// return false if no matches are found
	return false
}
