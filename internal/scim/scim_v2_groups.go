package scim

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/application"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/db"
	pagination2 "github.com/suse-skyscraper/openfga-scim-bridge/internal/pagination"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/scim/middleware"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/scim/patcher"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/scim/payloads"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/scim/responses"
)

func V2ListGroups(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := pagination2.Paginate(r)

		totalCount, groups, err := app.Repository.GetGroups(r.Context(), db.GetGroupsParams{
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimGroupListResponse(
			app.Config,
			groups,
			responses.ScimGroupListResponseInput{
				StartIndex:   int(pagination.Offset)/int(pagination.Limit) + 1,
				TotalResults: int(totalCount),
				ItemsPerPage: int(pagination.Limit),
			}))
	}
}

func V2GetGroup(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		group, ok := r.Context().Value(middleware.Group).(db.Group)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		members, err := app.Repository.GetGroupMembership(r.Context(), group.ID.String())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimGroupResponse(app.Config, group, members))
	}
}

func V2CreateGroup(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := payloads.GroupPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}

		group, err := app.Repository.CreateGroup(r.Context(), payload.DisplayName)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		var members []db.GetGroupMembershipRow
		RenderScimJSON(w, r, http.StatusCreated, responses.NewScimGroupResponse(app.Config, group, members))
	}
}

func V2PatchGroup(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		group, ok := r.Context().Value(middleware.Group).(db.Group)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		payload, err := payloads.GroupPatchPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}

		groupPatcher := patcher.NewGroupPatcher(r.Context(), app)
		err = groupPatcher.Patch(group, payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		group, err = app.Repository.FindGroup(r.Context(), group.ID.String())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		var members []db.GetGroupMembershipRow
		RenderScimJSON(w, r, http.StatusOK, responses.NewScimGroupResponse(app.Config, group, members))
	}
}

func V2DeleteGroup(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		group := r.Context().Value(middleware.Group).(db.Group)

		err := app.FGAClient.RemoveUsersInGroup(r.Context(), group.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		err = app.Repository.DeleteGroup(r.Context(), group.ID.String())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
