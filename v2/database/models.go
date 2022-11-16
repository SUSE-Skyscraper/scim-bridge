package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/filters"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/payloads"
)

type User struct {
	ID          uuid.UUID
	Username    string
	ExternalID  sql.NullString
	Name        map[string]string
	DisplayName sql.NullString
	Locale      sql.NullString
	Active      bool
	Emails      []payloads.UserEmail
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Group struct {
	ID          uuid.UUID
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GroupMembership struct {
	GroupID  uuid.UUID
	UserID   uuid.UUID
	Username sql.NullString
}

type UserParams struct {
	Username    string
	Name        map[string]string
	DisplayName string
	Emails      []payloads.UserEmail
	Active      bool
	Locale      string
	ExternalID  string
}

type GetUsersParams struct {
	Filters []filters.Filter
	Offset  int32
	Limit   int32
}
