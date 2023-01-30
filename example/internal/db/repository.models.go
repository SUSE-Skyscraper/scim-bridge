package db

import (
	"github.com/suse-skyscraper/openfga-scim-bridge/util"
)

type GetScimUsersInput struct {
	Filters []util.Filter
	Offset  int32
	Limit   int32
}
