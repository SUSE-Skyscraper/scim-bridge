package db

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/scim/filters"
)

var ErrConflict = errors.New("duplicate key value violates unique constraint")

var _ RepositoryQueries = (*Repository)(nil)

type Repository struct {
	postgresPool *pgxpool.Pool
	db           Querier
	tx           pgx.Tx
}

func (r *Repository) CreateAPIKey(ctx context.Context, input InsertAPIKeyParams) (ApiKey, error) {
	return r.db.InsertAPIKey(ctx, input)
}

func (r *Repository) GetAPIKeys(ctx context.Context) ([]ApiKey, error) {
	return r.db.GetAPIKeys(ctx)
}

func (r *Repository) GetUsers(ctx context.Context, input GetUsersParams) ([]User, error) {
	return r.db.GetUsers(ctx, input)
}

func (r *Repository) FindScimAPIKey(ctx context.Context) (ApiKey, error) {
	return r.db.FindScimAPIKey(ctx)
}

func (r *Repository) FindAPIKey(ctx context.Context, id uuid.UUID) (ApiKey, error) {
	return r.db.FindAPIKey(ctx, id)
}

func (r *Repository) DeleteScimAPIKey(ctx context.Context) error {
	apiKey, err := r.FindScimAPIKey(ctx)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if apiKey.ID != uuid.Nil {
		err = r.db.DeleteAPIKey(ctx, apiKey.ID)
		if err != nil {
			return err
		}
	}

	return r.db.DeleteScimAPIKey(ctx)
}

func (r *Repository) InsertScimAPIKey(ctx context.Context, encodedHash string) (ApiKey, error) {
	apiKey, err := r.db.InsertAPIKey(ctx, InsertAPIKeyParams{
		Encodedhash: encodedHash,
		System:      true,
		Owner:       "SCIM",
		Description: sql.NullString{String: "SCIM API key", Valid: true},
	})
	if err != nil {
		return ApiKey{}, err
	}

	_, err = r.db.InsertScimAPIKey(ctx, apiKey.ID)
	if err != nil {
		return ApiKey{}, err
	}

	return apiKey, nil
}

func (r *Repository) ScimPatchUser(ctx context.Context, input PatchUserParams) error {
	return r.db.PatchUser(ctx, input)
}

func (r *Repository) UpdateUser(ctx context.Context, id uuid.UUID, input UpdateUserParams) (User, error) {
	err := r.db.UpdateUser(ctx, input)
	if err != nil {
		return User{}, err
	}

	user, err := r.db.GetUser(ctx, id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.db.DeleteUser(ctx, id)
}

func (r *Repository) CreateUser(ctx context.Context, input CreateUserParams) (User, error) {
	user, err := r.db.CreateUser(ctx, input)
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return User{}, ErrConflict
	} else if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) GetScimUsers(ctx context.Context, input GetScimUsersInput) (int64, []User, error) {
	if len(input.Filters) == 0 {
		totalCount, err := r.db.GetUserCount(ctx)
		if err != nil {
			return 0, nil, err
		}

		users, err := r.db.GetUsers(ctx, GetUsersParams{
			Offset: input.Offset,
			Limit:  input.Limit,
		})
		if err != nil {
			return 0, nil, err
		}

		return totalCount, users, nil
	}

	// we only support the userName filter for now
	// Okta uses this to see if a userName already exists
	filter := input.Filters[0]
	if filter.FilterField == filters.Username && filter.FilterOperator == filters.Eq {
		user, err := r.db.FindByUsername(ctx, filter.FilterValue)
		switch err {
		case nil:
			return 1, []User{user}, nil
		case pgx.ErrNoRows:
			return 0, []User{}, nil
		default:
			return 0, nil, err
		}
	} else {
		return 0, nil, errors.New("unsupported filter")
	}
}

func (r *Repository) CreateGroup(ctx context.Context, displayName string) (Group, error) {
	group, err := r.db.CreateGroup(ctx, displayName)
	if err != nil {
		return Group{}, err
	}

	return group, nil
}

func (r *Repository) GetGroupMembership(ctx context.Context, idString string) ([]GetGroupMembershipRow, error) {
	id, err := uuid.Parse(idString)
	if err != nil {
		return nil, err
	}

	return r.db.GetGroupMembership(ctx, id)
}

func (r *Repository) GetGroups(ctx context.Context, params GetGroupsParams) (int64, []Group, error) {
	totalCount, err := r.db.GetGroupCount(ctx)
	if err != nil {
		return 0, nil, err
	}

	groups, err := r.db.GetGroups(ctx, params)
	if err != nil {
		return 0, nil, err
	}

	return totalCount, groups, nil
}

func (r *Repository) DeleteGroup(ctx context.Context, idString string) error {
	id, err := uuid.Parse(idString)
	if err != nil {
		return err
	}

	err = r.db.DeleteGroup(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateGroup(ctx context.Context, input PatchGroupDisplayNameParams) (Group, error) {
	err := r.db.PatchGroupDisplayName(ctx, input)
	if err != nil {
		return Group{}, err
	}

	return r.FindGroup(ctx, input.ID.String())
}

func (r *Repository) Rollback(ctx context.Context) error {
	if r.tx == nil {
		return errors.New("no transaction in progress")
	}

	return r.tx.Rollback(ctx)
}

func (r *Repository) Begin(ctx context.Context) (RepositoryQueries, error) {
	tx, err := r.postgresPool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	db := &Queries{db: tx}

	return &Repository{
		postgresPool: r.postgresPool,
		db:           db,
		tx:           tx,
	}, nil
}

func (r *Repository) Commit(ctx context.Context) error {
	if r.tx == nil {
		return errors.New("no transaction in progress")
	}

	return r.tx.Commit(ctx)
}

func (r *Repository) FindUser(ctx context.Context, id string) (User, error) {
	idParsed, err := uuid.Parse(id)
	if err != nil {
		return User{}, err
	}

	return r.db.GetUser(ctx, idParsed)
}

func (r *Repository) FindGroup(ctx context.Context, id string) (Group, error) {
	idParsed, err := uuid.Parse(id)
	if err != nil {
		return Group{}, err
	}

	return r.db.GetGroup(ctx, idParsed)
}

func (r *Repository) FindUserByUsername(ctx context.Context, username string) (User, error) {
	return r.db.FindByUsername(ctx, username)
}

func (r *Repository) AddUsersToGroup(ctx context.Context, groupID uuid.UUID, members []uuid.UUID) error {
	for _, member := range members {
		err := r.AddUserToGroup(ctx, member, groupID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) ReplaceUsersInGroup(ctx context.Context, groupID uuid.UUID, members []uuid.UUID) error {
	err := r.db.DropMembershipForGroup(ctx, groupID)
	if err != nil {
		return err
	}

	for _, member := range members {
		err = r.AddUserToGroup(ctx, member, groupID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) AddUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	err := r.db.CreateMembershipForUserAndGroup(ctx, CreateMembershipForUserAndGroupParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	err := r.db.DropMembershipForUserAndGroup(ctx, DropMembershipForUserAndGroupParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil {
		return err
	}

	return nil
}
