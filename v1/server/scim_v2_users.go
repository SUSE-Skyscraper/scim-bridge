package server

import (
	"github.com/suse-skyscraper/openfga-scim-bridge/bridge"
	database2 "github.com/suse-skyscraper/openfga-scim-bridge/database"
	payloads2 "github.com/suse-skyscraper/openfga-scim-bridge/payloads"
	"github.com/suse-skyscraper/openfga-scim-bridge/util"
	"github.com/suse-skyscraper/openfga-scim-bridge/v1/middleware"
	"github.com/suse-skyscraper/openfga-scim-bridge/v1/responses"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

func V1ListUsers(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filterString := r.URL.Query().Get("filter")
		filterList, err := util.ParseFilter(filterString)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadFilter(err))
			return
		}

		page := util.Paginate(r)

		totalCount, users, err := bridge.DB.GetUsers(r.Context(), database2.GetUsersParams{
			Filters: filterList,
			Offset:  page.Offset,
			Limit:   page.Limit,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimUserListResponse(
			bridge,
			users,
			responses.ScimUserListResponseInput{
				StartIndex:   int(page.Offset)/int(page.Limit) + 1,
				TotalResults: int(totalCount),
				ItemsPerPage: int(page.Limit),
			}))
	}
}

func V1GetUser(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(database2.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimUserResponse(bridge, user))
	}
}

func V1CreateUser(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := payloads2.CreateUserPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}

		if payload.Username == "" {
			_ = render.Render(w, r, responses.ErrBadValue(errors.New("Attribute 'userName' is required")))
			return
		}

		user, err := bridge.DB.CreateUser(r.Context(), database2.UserParams{
			Username:    payload.Username,
			Name:        payload.Name,
			Active:      payload.Active,
			Emails:      payload.Emails,
			Locale:      payload.Locale,
			ExternalID:  payload.ExternalID,
			DisplayName: payload.DisplayName,
		})

		if err != nil && errors.Is(err, database2.ErrConflict) {
			w.WriteHeader(http.StatusConflict)
			return
		} else if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusCreated, responses.NewScimUserResponse(bridge, user))
	}
}

func V1DeleteUser(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(database2.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		err := bridge.DB.DeleteUser(r.Context(), user.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func V1UpdateUser(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(database2.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		payload, err := payloads2.CreateUserPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}

		user, err = bridge.DB.UpdateUser(r.Context(), user.ID, database2.UserParams{
			Username:    payload.Username,
			Name:        payload.Name,
			Active:      payload.Active,
			Emails:      payload.Emails,
			Locale:      payload.Locale,
			DisplayName: payload.DisplayName,
			ExternalID:  payload.ExternalID,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimUserResponse(bridge, user))
	}
}

func V1PatchUser(bridge *bridge.Bridge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(database2.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		payload, err := payloads2.UserPatchPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}

		for _, op := range payload.Operations {
			switch op.Op {
			case "replace":
				err = bridge.DB.SetUserActive(r.Context(), user.ID, op.Value.Active)
				if err != nil {
					_ = render.Render(w, r, responses.ErrInternalServerError)
					return
				}
			default:
				_ = render.Render(w, r, responses.ErrBadValue(errors.New("Unsupported operation")))
			}
		}

		user, err = bridge.DB.FindUser(r.Context(), user.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimUserResponse(bridge, user))
	}
}
