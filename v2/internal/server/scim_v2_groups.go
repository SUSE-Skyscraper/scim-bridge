package server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/bridge"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/database"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/internal/middleware"
	pagination2 "github.com/suse-skyscraper/openfga-scim-bridge/v2/internal/pagination"
	responses2 "github.com/suse-skyscraper/openfga-scim-bridge/v2/internal/responses"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/payloads"
)

func V2ListGroups(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := pagination2.Paginate(r)

		totalCount, groups, err := bridge.DB.GetGroups(r.Context(), pagination.Limit, pagination.Offset)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses2.NewScimGroupListResponse(
			bridge,
			groups,
			responses2.ScimGroupListResponseInput{
				StartIndex:   int(pagination.Offset)/int(pagination.Limit) + 1,
				TotalResults: int(totalCount),
				ItemsPerPage: int(pagination.Limit),
			}))
	}
}

func V2GetGroup(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		group, ok := r.Context().Value(middleware.Group).(database.Group)
		if !ok {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		members, err := bridge.DB.GetGroupMembership(r.Context(), group.ID)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses2.NewScimGroupResponse(bridge, group, members))
	}
}

func V2CreateGroup(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := payloads.GroupPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrBadValue(err))
			return
		}

		group, err := bridge.DB.CreateGroup(r.Context(), payload.DisplayName)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		// A new group has no members, just make an empty list
		var members []database.GroupMembership

		RenderScimJSON(w, r, http.StatusCreated, responses2.NewScimGroupResponse(bridge, group, members))
	}
}

func V2PatchGroup(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		group, ok := r.Context().Value(middleware.Group).(database.Group)
		if !ok {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		payload, err := payloads.GroupPatchPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrBadValue(err))
			return
		}

		var operations []payloads.GroupPatchOperation
		for _, op := range payload.Operations {
			operations = append(operations, payloads.GroupPatchOperation{
				Op:    op.Op,
				Path:  op.Path,
				Value: op.Value,
			})
		}

		err = bridge.DB.PatchGroup(r.Context(), group.ID, operations)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		group, err = bridge.DB.FindGroup(r.Context(), group.ID)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		// displaying the membership is optional after a patch
		var members []database.GroupMembership

		RenderScimJSON(w, r, http.StatusOK, responses2.NewScimGroupResponse(bridge, group, members))
	}
}

func V2DeleteGroup(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		group := r.Context().Value(middleware.Group).(database.Group)

		err := bridge.DB.DeleteGroup(r.Context(), group.ID)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
