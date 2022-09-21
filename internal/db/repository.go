package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewRepository(pool *pgxpool.Pool, db Querier) *Repository {
	return &Repository{
		db:           db,
		postgresPool: pool,
	}
}

type RepositoryQueries interface {
	Begin(ctx context.Context) (RepositoryQueries, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	FindGroup(ctx context.Context, id string) (Group, error)
	CreateGroup(ctx context.Context, displayName string) (Group, error)
	DeleteGroup(ctx context.Context, id string) error
	UpdateGroup(ctx context.Context, input PatchGroupDisplayNameParams) (Group, error)
	RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error
	AddUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error
	ReplaceUsersInGroup(ctx context.Context, groupID uuid.UUID, members []uuid.UUID) error
	AddUsersToGroup(ctx context.Context, groupID uuid.UUID, members []uuid.UUID) error
	GetGroupMembership(ctx context.Context, idString string) ([]GetGroupMembershipRow, error)
	GetGroups(ctx context.Context, params GetGroupsParams) (int64, []Group, error)

	FindUser(ctx context.Context, id string) (User, error)
	FindUserByUsername(ctx context.Context, username string) (User, error)
	GetUsers(ctx context.Context, input GetUsersParams) ([]User, error)
	GetScimUsers(ctx context.Context, input GetScimUsersInput) (int64, []User, error)
	CreateUser(ctx context.Context, input CreateUserParams) (User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateUser(ctx context.Context, id uuid.UUID, input UpdateUserParams) (User, error)
	ScimPatchUser(ctx context.Context, input PatchUserParams) error

	InsertScimAPIKey(ctx context.Context, encodedHash string) (ApiKey, error)
	DeleteScimAPIKey(ctx context.Context) error
	FindAPIKey(ctx context.Context, id uuid.UUID) (ApiKey, error)
	FindScimAPIKey(ctx context.Context) (ApiKey, error)
	GetAPIKeys(ctx context.Context) ([]ApiKey, error)
	CreateAPIKey(ctx context.Context, input InsertAPIKeyParams) (ApiKey, error)
}
