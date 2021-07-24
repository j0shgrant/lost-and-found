package aws

import "github.com/j0shgrant/lost-and-found/internal/tags"

type TaggableService interface {
	Get([]tags.Tag) ([]TaggableResource, error)
}

type TaggableResource interface {
	Tags() []tags.Tag
}
