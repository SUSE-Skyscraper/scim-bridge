package server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/bridge"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/database"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/filters"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/internal/middleware"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/internal/pagination"
	responses2 "github.com/suse-skyscraper/openfga-scim-bridge/v2/internal/responses"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/payloads"
)

func V2ListUsers(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filterString := r.URL.Query().Get("filter")
		filterList, err := filters.ParseFilter(filterString)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrBadFilter(err))
			return
		}

		page := pagination.Paginate(r)

		totalCount, users, err := bridge.DB.GetUsers(r.Context(), database.GetUsersParams{
			Filters: filterList,
			Offset:  page.Offset,
			Limit:   page.Limit,
		})
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses2.NewScimUserListResponse(
			bridge,
			users,
			responses2.ScimUserListResponseInput{
				StartIndex:   int(page.Offset)/int(page.Limit) + 1,
				TotalResults: int(totalCount),
				ItemsPerPage: int(page.Limit),
			}))
	}
}

func V2GetUser(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(database.User)
		if !ok {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses2.NewScimUserResponse(bridge, user))
	}
}

func V2CreateUser(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := payloads.Parse(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrBadValue(err))
			return
		}

		if payload.Username == "" {
			_ = render.Render(w, r, responses2.ErrBadValue(errors.New("Attribute 'userName' is required")))
			return
		}

		user, err := bridge.DB.CreateUser(r.Context(), database.UserParams{
			Username:    payload.Username,
			Name:        payload.Name,
			Active:      payload.Active,
			Emails:      payload.Emails,
			Locale:      payload.Locale,
			ExternalID:  payload.ExternalID,
			DisplayName: payload.DisplayName,
		})

		if err != nil && errors.Is(err, database.ErrConflict) {
			w.WriteHeader(http.StatusConflict)
			return
		} else if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusCreated, responses2.NewScimUserResponse(bridge, user))
	}
}

func V2DeleteUser(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(database.User)
		if !ok {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		err := bridge.DB.DeleteUser(r.Context(), user.ID)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func V2UpdateUser(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(database.User)
		if !ok {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		payload, err := payloads.Parse(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrBadValue(err))
			return
		}

		user, err = bridge.DB.UpdateUser(r.Context(), user.ID, database.UserParams{
			Username:    payload.Username,
			Name:        payload.Name,
			Active:      payload.Active,
			Emails:      payload.Emails,
			Locale:      payload.Locale,
			DisplayName: payload.DisplayName,
			ExternalID:  payload.ExternalID,
		})
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses2.NewScimUserResponse(bridge, user))
	}
}

func V2PatchUser(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(database.User)
		if !ok {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		payload, err := payloads.UserPatchPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrBadValue(err))
			return
		}

		for _, op := range payload.Operations {
			switch op.Op {
			case "replace":
				err = bridge.DB.SetUserActive(r.Context(), user.ID, op.Value.Active)
				if err != nil {
					_ = render.Render(w, r, responses2.ErrInternalServerError)
					return
				}
			default:
				_ = render.Render(w, r, responses2.ErrBadValue(errors.New("Unsupported operation")))
			}
		}

		user, err = bridge.DB.FindUser(r.Context(), user.ID)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses2.NewScimUserResponse(bridge, user))
	}
}
