package patcher

import (
	"context"
	"errors"

	"github.com/suse-skyscraper/openfga-scim-bridge/internal/application"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/db"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/scim/payloads"
)

type GroupPatcher struct {
	ctx context.Context
	app *application.App
}

func NewGroupPatcher(ctx context.Context, app *application.App) *GroupPatcher {
	return &GroupPatcher{
		ctx: ctx,
		app: app,
	}
}

func (p *GroupPatcher) Patch(group db.Group, payload *payloads.GroupPatchPayload) error {
	repo, err := p.app.Repository.Begin(p.ctx)
	if err != nil {
		return errors.New("failed to begin transaction")
	}
	defer func(repo db.RepositoryQueries, ctx context.Context) {
		_ = repo.Rollback(ctx)
	}(repo, p.ctx)

	for _, op := range payload.Operations {
		switch op.Op {
		case "add":
			err = p.patchAdd(repo, group, op)
			if err != nil {
				return err
			}
		case "remove":
			err = p.patchRemove(repo, group, op)
			if err != nil {
				return err
			}
		case "replace":
			err = p.patchReplace(repo, group, op)
			if err != nil {
				return err
			}
		default:
			return errors.New("unknown operation")
		}
	}

	return repo.Commit(p.ctx)
}

func (p *GroupPatcher) patchAdd(repo db.RepositoryQueries, group db.Group, op *payloads.GroupPatchOperation) error {
	newMembers, err := op.GetAddMembersPatch()
	if err != nil {
		return errors.New("failed to get add members patch")
	}

	err = p.app.FGAClient.AddUsersToGroup(p.ctx, newMembers, group.ID)
	if err != nil {
		return errors.New("failed to add members to FGA group")
	}

	err = repo.AddUsersToGroup(p.ctx, group.ID, newMembers)
	if err != nil {
		return errors.New("failed to add members")
	}

	return nil
}

func (p *GroupPatcher) patchRemove(repo db.RepositoryQueries, group db.Group, op *payloads.GroupPatchOperation) error {
	id, err := op.ParseIDFromPath()
	if err != nil {
		return errors.New("failed to parse id from path")
	}

	err = p.app.FGAClient.RemoveUserFromGroup(p.ctx, id, group.ID)
	if err != nil {
		return err
	}

	return repo.RemoveUserFromGroup(p.ctx, id, group.ID)
}

func (p *GroupPatcher) patchReplace(repo db.RepositoryQueries, group db.Group, op *payloads.GroupPatchOperation) error {
	switch op.Path {
	case "members":
		newMembers, err := op.GetAddMembersPatch()
		if err != nil {
			return errors.New("failed to get add members patch")
		}

		err = p.app.FGAClient.ReplaceUsersInGroup(p.ctx, newMembers, group.ID)
		if err != nil {
			return errors.New("failed to replace members in FGA")
		}

		err = repo.ReplaceUsersInGroup(p.ctx, group.ID, newMembers)
		if err != nil {
			return errors.New("failed to replace members")
		}
	default:
		patch, err := op.GetPatch()
		if err != nil {
			return errors.New("failed to get patch")
		}

		_, err = repo.UpdateGroup(p.ctx, db.PatchGroupDisplayNameParams{
			ID:          group.ID,
			DisplayName: patch.DisplayName,
		})
		if err != nil {
			return errors.New("failed to patch display name")
		}
	}

	return nil
}
