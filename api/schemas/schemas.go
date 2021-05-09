package schemas

import (
	gql "github.com/mattdamon108/gqlmerge/lib"
)

func NewSchema() *string {
	schema := gql.Merge("", "./schemas")

	return schema
}
