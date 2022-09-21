package db

import (
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/scim/filters"
)

type GetScimUsersInput struct {
	Filters []filters.Filter
	Offset  int32
	Limit   int32
}
