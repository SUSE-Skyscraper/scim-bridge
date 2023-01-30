package database

import (
	"database/sql"
	"github.com/suse-skyscraper/openfga-scim-bridge/payloads"
	"github.com/suse-skyscraper/openfga-scim-bridge/util"
	"time"

	"github.com/google/uuid"
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
	Filters []util.Filter
	Offset  int32
	Limit   int32
}
