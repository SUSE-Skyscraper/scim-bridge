package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/payloads"
)

var ErrConflict = errors.New("duplicate key value violates unique constraint")

type Bridge interface {
	FindUser(ctx context.Context, userID uuid.UUID) (User, error)
	CreateUser(ctx context.Context, arg UserParams) (User, error)
	GetUsers(ctx context.Context, arg GetUsersParams) (int64, []User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	UpdateUser(ctx context.Context, userID uuid.UUID, arg UserParams) (User, error)
	SetUserActive(ctx context.Context, userID uuid.UUID, active bool) error

	FindGroup(ctx context.Context, groupID uuid.UUID) (Group, error)
	CreateGroup(ctx context.Context, displayName string) (Group, error)
	GetGroups(ctx context.Context, limit int32, offset int32) (int64, []Group, error)
	GetGroupMembership(ctx context.Context, groupID uuid.UUID) ([]GroupMembership, error)
	DeleteGroup(ctx context.Context, groupID uuid.UUID) error
	PatchGroup(ctx context.Context, groupID uuid.UUID, operations []payloads.GroupPatchOperation) error
}
