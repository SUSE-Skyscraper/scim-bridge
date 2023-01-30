package server

import (
	"github.com/suse-skyscraper/openfga-scim-bridge/bridge"
	"github.com/suse-skyscraper/openfga-scim-bridge/database"
	payloads2 "github.com/suse-skyscraper/openfga-scim-bridge/payloads"
	pagination2 "github.com/suse-skyscraper/openfga-scim-bridge/util"
	"github.com/suse-skyscraper/openfga-scim-bridge/v1/middleware"
	"github.com/suse-skyscraper/openfga-scim-bridge/v1/responses"
	"net/http"

	"github.com/go-chi/render"
)

func V1ListGroups(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := pagination2.Paginate(r)

		totalCount, groups, err := bridge.DB.GetGroups(r.Context(), pagination.Limit, pagination.Offset)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimGroupListResponse(
			bridge,
			groups,
			responses.ScimGroupListResponseInput{
				StartIndex:   int(pagination.Offset)/int(pagination.Limit) + 1,
				TotalResults: int(totalCount),
				ItemsPerPage: int(pagination.Limit),
			}))
	}
}

func V1GetGroup(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		group, ok := r.Context().Value(middleware.Group).(database.Group)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		members, err := bridge.DB.GetGroupMembership(r.Context(), group.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimGroupResponse(bridge, group, members))
	}
}

func V1CreateGroup(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := payloads2.GroupPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}

		group, err := bridge.DB.CreateGroup(r.Context(), payload.DisplayName)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// A new group has no members, just make an empty list
		var members []database.GroupMembership

		RenderScimJSON(w, r, http.StatusCreated, responses.NewScimGroupResponse(bridge, group, members))
	}
}

func V1PatchGroup(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		group, ok := r.Context().Value(middleware.Group).(database.Group)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		payload, err := payloads2.GroupPatchPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}

		var operations []payloads2.GroupPatchOperation
		for _, op := range payload.Operations {
			operations = append(operations, payloads2.GroupPatchOperation{
				Op:    op.Op,
				Path:  op.Path,
				Value: op.Value,
			})
		}

		err = bridge.DB.PatchGroup(r.Context(), group.ID, operations)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		group, err = bridge.DB.FindGroup(r.Context(), group.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// displaying the membership is optional after a patch
		var members []database.GroupMembership

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimGroupResponse(bridge, group, members))
	}
}

func V1DeleteGroup(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		group := r.Context().Value(middleware.Group).(database.Group)

		err := bridge.DB.DeleteGroup(r.Context(), group.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
