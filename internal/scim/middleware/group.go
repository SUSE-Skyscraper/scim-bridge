package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/application"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/scim/responses"
)

func GroupCtx(app *application.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idString := chi.URLParam(r, "id")

			group, err := app.Repository.FindGroup(r.Context(), idString)
			if errors.Is(err, pgx.ErrNoRows) {
				_ = render.Render(w, r, responses.ErrNotFound(idString))
				return
			} else if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), Group, group)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
