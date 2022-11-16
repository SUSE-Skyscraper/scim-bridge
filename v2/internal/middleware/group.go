package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	openfga_scim_bridge "github.com/suse-skyscraper/openfga-scim-bridge/v2/bridge"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/internal/responses"
)

func GroupCtx(bridge *openfga_scim_bridge.Bridge) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idString := chi.URLParam(r, "id")

			id, err := uuid.Parse(idString)
			if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}

			group, err := bridge.DB.FindGroup(r.Context(), id)
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
