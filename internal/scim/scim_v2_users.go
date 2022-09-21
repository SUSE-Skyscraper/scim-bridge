package scim

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/application"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/db"
	pagination2 "github.com/suse-skyscraper/openfga-scim-bridge/internal/pagination"
	filters2 "github.com/suse-skyscraper/openfga-scim-bridge/internal/scim/filters"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/scim/middleware"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/scim/payloads"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/scim/responses"
)

func V2ListUsers(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filterString := r.URL.Query().Get("filter")
		filters, err := filters2.ParseFilter(filterString)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadFilter(err))
			return
		}

		pagination := pagination2.Paginate(r)
		totalCount, users, err := app.Repository.GetScimUsers(r.Context(), db.GetScimUsersInput{
			Filters: filters,
			Offset:  pagination.Offset,
			Limit:   pagination.Limit,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimUserListResponse(
			app.Config,
			users,
			responses.ScimUserListResponseInput{
				StartIndex:   int(pagination.Offset)/int(pagination.Limit) + 1,
				TotalResults: int(totalCount),
				ItemsPerPage: int(pagination.Limit),
			}))
	}
}

func V2GetUser(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(db.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimUserResponse(app.Config, user))
	}
}

func V2CreateUser(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := payloads.UserPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}

		if payload.Username == "" {
			_ = render.Render(w, r, responses.ErrBadValue(errors.New("Attribute 'userName' is required")))
			return
		}

		tx, err := app.Repository.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		defer func(tx db.RepositoryQueries, ctx context.Context) {
			_ = tx.Rollback(ctx)
		}(tx, r.Context())

		user, err := tx.CreateUser(r.Context(), db.CreateUserParams{
			Username: payload.Username,
			Name:     payload.GetJSONName(),
			Active:   payload.Active,
			Emails:   payload.GetJSONEmails(),
			Locale: sql.NullString{
				String: payload.Locale,
				Valid:  payload.Locale != "",
			},
			ExternalID: sql.NullString{
				String: payload.ExternalID,
				Valid:  payload.ExternalID != "",
			},
			DisplayName: sql.NullString{
				String: payload.DisplayName,
				Valid:  payload.DisplayName != "",
			},
		})
		if err != nil && errors.Is(err, db.ErrConflict) {
			w.WriteHeader(http.StatusConflict)
			return
		} else if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		err = tx.Commit(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusCreated, responses.NewScimUserResponse(app.Config, user))
	}
}

func V2DeleteUser(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(db.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		tx, err := app.Repository.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		defer func(tx db.RepositoryQueries, ctx context.Context) {
			_ = tx.Rollback(ctx)
		}(tx, r.Context())

		err = tx.DeleteUser(r.Context(), user.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		err = app.FGAClient.RemoveUser(r.Context(), user.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		err = tx.Commit(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func V2UpdateUser(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(db.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		payload, err := payloads.UserPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}
		user, err = app.Repository.UpdateUser(r.Context(), user.ID, db.UpdateUserParams{
			ID:       user.ID,
			Username: payload.Username,
			Name:     payload.GetJSONName(),
			Active:   payload.Active,
			Emails:   payload.GetJSONEmails(),
			Locale: sql.NullString{
				String: payload.Locale,
				Valid:  payload.Locale != "",
			},
			DisplayName: sql.NullString{
				String: payload.DisplayName,
				Valid:  payload.DisplayName != "",
			},
			ExternalID: sql.NullString{
				String: payload.ExternalID,
				Valid:  payload.ExternalID != "",
			},
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimUserResponse(app.Config, user))
	}
}

func V2PatchUser(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(db.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		payload, err := payloads.UserPatchPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}

		for _, op := range payload.Operations {
			switch op.Op {
			case "replace":
				err = app.Repository.ScimPatchUser(r.Context(), db.PatchUserParams{
					ID:     user.ID,
					Active: op.Value.Active,
				})
				if err != nil {
					_ = render.Render(w, r, responses.ErrInternalServerError)
					return
				}
			default:
				_ = render.Render(w, r, responses.ErrBadValue(errors.New("Unsupported operation")))
			}
		}

		user, err = app.Repository.FindUser(r.Context(), user.ID.String())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimUserResponse(app.Config, user))
	}
}
