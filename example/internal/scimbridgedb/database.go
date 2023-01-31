package scimbridgedb

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/suse-skyscraper/openfga-scim-bridge/example/internal/application"
	"github.com/suse-skyscraper/openfga-scim-bridge/example/internal/db"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/database"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/payloads"
)

var _ database.Bridge = (*DB)(nil)

type DB struct {
	app *application.App
}

func New(app *application.App) DB {
	return DB{
		app: app,
	}
}

func (d *DB) PatchGroup(ctx context.Context, groupID uuid.UUID, operations []payloads.GroupPatchOperation) error {
	tx, err := d.app.Repository.Begin(ctx)
	if err != nil {
		return errors.New("failed to begin transaction")
	}

	defer func(tx db.RepositoryQueries, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	for _, op := range operations {
		switch op.Op {
		case "add":
			err = d.patchAdd(ctx, tx, groupID, op)
			if err != nil {
				return err
			}
		case "remove":
			err = d.patchRemove(ctx, tx, groupID, op)
			if err != nil {
				return err
			}
		case "replace":
			err = d.patchReplace(ctx, tx, groupID, op)
			if err != nil {
				return err
			}
		default:
			return errors.New("unknown operation")
		}
	}

	return tx.Commit(ctx)
}

func (d *DB) patchAdd(ctx context.Context, tx db.RepositoryQueries, groupID uuid.UUID, op payloads.GroupPatchOperation) error {
	newMembers, err := op.GetAddMembersPatch()
	if err != nil {
		return errors.New("failed to get add members patch")
	}

	err = d.app.FGAClient.AddUsersToGroup(ctx, newMembers, groupID)
	if err != nil {
		return errors.New("failed to add members to FGA group")
	}

	err = tx.AddUsersToGroup(ctx, groupID, newMembers)
	if err != nil {
		return errors.New("failed to add members")
	}

	return nil
}

func (d *DB) patchRemove(ctx context.Context, tx db.RepositoryQueries, groupID uuid.UUID, op payloads.GroupPatchOperation) error {
	id, err := op.ParseIDFromPath()
	if err != nil {
		return errors.New("failed to parse id from path")
	}

	err = d.app.FGAClient.RemoveUserFromGroup(ctx, id, groupID)
	if err != nil {
		return err
	}

	return tx.RemoveUserFromGroup(ctx, id, groupID)
}

func (d *DB) patchReplace(ctx context.Context, tx db.RepositoryQueries, groupID uuid.UUID, op payloads.GroupPatchOperation) error {
	switch op.Path {
	case "members":
		newMembers, err := op.GetAddMembersPatch()
		if err != nil {
			return errors.New("failed to get add members patch")
		}

		err = d.app.FGAClient.ReplaceUsersInGroup(ctx, newMembers, groupID)
		if err != nil {
			return errors.New("failed to replace members in FGA")
		}

		err = tx.ReplaceUsersInGroup(ctx, groupID, newMembers)
		if err != nil {
			return errors.New("failed to replace members")
		}
	default:
		patch, err := op.GetPatch()
		if err != nil {
			return errors.New("failed to get patch")
		}

		_, err = tx.UpdateGroup(ctx, db.PatchGroupDisplayNameParams{
			ID:          groupID,
			DisplayName: patch.DisplayName,
		})
		if err != nil {
			return errors.New("failed to patch display name")
		}
	}

	return nil
}

func (d *DB) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	err := d.app.FGAClient.RemoveUsersInGroup(ctx, groupID)
	if err != nil {
		return err
	}

	err = d.app.Repository.DeleteGroup(ctx, groupID.String())
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) CreateGroup(ctx context.Context, displayName string) (database.Group, error) {
	group, err := d.app.Repository.CreateGroup(ctx, displayName)
	if err != nil {
		return database.Group{}, err
	}

	scimGroup := toScimGroup(group)
	return scimGroup, nil
}

func (d *DB) GetGroupMembership(ctx context.Context, groupID uuid.UUID) ([]database.GroupMembership, error) {
	members, err := d.app.Repository.GetGroupMembership(ctx, groupID.String())
	if err != nil {
		return nil, err
	}

	var groupMembers []database.GroupMembership
	for _, member := range members {
		groupMembers = append(groupMembers, database.GroupMembership{
			GroupID:  member.GroupID,
			Username: member.Username,
			UserID:   member.UserID,
		})
	}

	return groupMembers, nil
}

func (d *DB) FindGroup(ctx context.Context, userID uuid.UUID) (database.Group, error) {
	group, err := d.app.Repository.FindGroup(ctx, userID.String())
	if err != nil {
		return database.Group{}, err
	}

	scimGroup := toScimGroup(group)
	return scimGroup, nil
}

func (d *DB) GetGroups(ctx context.Context, limit int32, offset int32) (int64, []database.Group, error) {

	totalCount, groups, err := d.app.Repository.GetGroups(ctx, db.GetGroupsParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return 0, nil, err
	}

	var scimGroups []database.Group
	for _, group := range groups {
		scimGroups = append(scimGroups, toScimGroup(group))
	}

	return totalCount, scimGroups, nil
}

func (d *DB) FindUser(ctx context.Context, userID uuid.UUID) (database.User, error) {
	user, err := d.app.Repository.FindUser(ctx, userID.String())
	if err != nil {
		return database.User{}, err
	}

	scimUser, err := toScimUser(user)
	if err != nil {
		return database.User{}, err
	}
	return scimUser, nil
}

func (d *DB) SetUserActive(ctx context.Context, userID uuid.UUID, active bool) error {
	err := d.app.Repository.ScimPatchUser(ctx, db.PatchUserParams{
		ID:     userID,
		Active: active,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) UpdateUser(ctx context.Context, userID uuid.UUID, arg database.UserParams) (database.User, error) {
	name, err := parseJSONB(arg.Name)
	if err != nil {
		return database.User{}, err
	}

	emails, err := parseJSONB(arg.Emails)
	if err != nil {
		return database.User{}, err
	}

	user, err := d.app.Repository.UpdateUser(ctx, userID, db.UpdateUserParams{
		ID:       userID,
		Username: arg.Username,
		Name:     name,
		Active:   arg.Active,
		Emails:   emails,
		Locale: sql.NullString{
			String: arg.Locale,
			Valid:  arg.Locale != "",
		},
		DisplayName: sql.NullString{
			String: arg.DisplayName,
			Valid:  arg.DisplayName != "",
		},
		ExternalID: sql.NullString{
			String: arg.ExternalID,
			Valid:  arg.ExternalID != "",
		},
	})
	if err != nil {
		return database.User{}, err
	}

	scimUser, err := toScimUser(user)
	if err != nil {
		return database.User{}, err
	}

	return scimUser, nil
}

func (d *DB) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	tx, err := d.app.Repository.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(tx db.RepositoryQueries, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	err = tx.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	err = d.app.FGAClient.RemoveUser(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) CreateUser(ctx context.Context, arg database.UserParams) (database.User, error) {
	tx, err := d.app.Repository.Begin(ctx)
	if err != nil {
		return database.User{}, err
	}

	defer func(tx db.RepositoryQueries, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	name, err := parseJSONB(arg.Name)
	if err != nil {
		return database.User{}, err
	}

	emails, err := parseJSONB(arg.Emails)
	if err != nil {
		return database.User{}, err
	}

	user, err := tx.CreateUser(ctx, db.CreateUserParams{
		Username: arg.Username,
		Name:     name,
		Active:   arg.Active,
		Emails:   emails,
		Locale: sql.NullString{
			String: arg.Locale,
			Valid:  arg.Locale != "",
		},
		ExternalID: sql.NullString{
			String: arg.ExternalID,
			Valid:  arg.ExternalID != "",
		},
		DisplayName: sql.NullString{
			String: arg.DisplayName,
			Valid:  arg.DisplayName != "",
		},
	})

	if err != nil {
		return database.User{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return database.User{}, err
	}

	scimUser, err := toScimUser(user)
	if err != nil {
		return database.User{}, err
	}

	return scimUser, nil
}

func (d *DB) GetUsers(ctx context.Context, input database.GetUsersParams) (int64, []database.User, error) {
	count, users, err := d.app.Repository.GetScimUsers(ctx, db.GetScimUsersInput{
		Filters: input.Filters,
		Offset:  input.Offset,
		Limit:   input.Limit,
	})
	if err != nil {
		return 0, nil, err
	}

	var scimUsers []database.User
	for _, user := range users {
		scimUser, err := toScimUser(user)
		if err != nil {
			return 0, nil, err
		}

		scimUsers = append(scimUsers, scimUser)
	}

	return count, scimUsers, nil
}

func toScimGroup(group db.Group) database.Group {
	scimGroup := database.Group{
		ID:          group.ID,
		DisplayName: group.DisplayName,
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
	}

	return scimGroup
}

func toScimUser(user db.User) (database.User, error) {
	var name map[string]string
	if user.Name.Bytes != nil {
		err := json.Unmarshal(user.Name.Bytes, &name)
		if err != nil {
			return database.User{}, err
		}
	}

	var emails []payloads.UserEmail
	if user.Emails.Bytes != nil {
		err := json.Unmarshal(user.Emails.Bytes, &emails)
		if err != nil {
			return database.User{}, err
		}
	}

	scimUser := database.User{
		ID:          user.ID,
		Username:    user.Username,
		ExternalID:  user.ExternalID,
		Name:        name,
		DisplayName: user.DisplayName,
		Locale:      user.Locale,
		Active:      user.Active,
		Emails:      emails,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return scimUser, nil
}

func parseJSONB(arg interface{}) (pgtype.JSONB, error) {
	if arg == nil {
		return pgtype.JSONB{Bytes: nil, Status: pgtype.Null}, nil
	}

	name := pgtype.JSONB{}
	err := name.Set(arg)
	if err != nil {
		return pgtype.JSONB{}, err
	}

	return name, nil
}
