package bridge

import (
	"github.com/suse-skyscraper/openfga-scim-bridge/database"
)

type Bridge struct {
	BaseURL string
	DB      database.Bridge
}

func New(db database.Bridge, baseURL string) Bridge {
	return Bridge{
		BaseURL: baseURL,
		DB:      db,
	}
}
