package db

import (
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/filters"
)

type GetScimUsersInput struct {
	Filters []filters.Filter
	Offset  int32
	Limit   int32
}
