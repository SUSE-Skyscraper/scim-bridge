package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/bridge"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/internal/middleware"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/internal/server"
)

type AuthorizationMiddleware func(next http.Handler) http.Handler

func Hook(r *chi.Mux, bridge *bridge.Bridge, authHandler AuthorizationMiddleware) {
	r.Route("/scim/v2", func(r chi.Router) {
		r.Use(authHandler)

		r.Get("/Users", server.V2ListUsers(bridge))
		r.Post("/Users", server.V2CreateUser(bridge))
		r.Route("/Users/{id}", func(r chi.Router) {
			scimUserCtx := middleware.UserCtx(bridge)

			r.Use(scimUserCtx)

			r.Get("/", server.V2GetUser(bridge))
			r.Put("/", server.V2UpdateUser(bridge))
			r.Patch("/", server.V2PatchUser(bridge))
			r.Delete("/", server.V2DeleteUser(bridge))
		})

		r.Get("/Groups", server.V2ListGroups(bridge))
		r.Post("/Groups", server.V2CreateGroup(bridge))
		r.Route("/Groups/{id}", func(r chi.Router) {
			scimGroupCtx := middleware.GroupCtx(bridge)

			r.Use(scimGroupCtx)

			r.Get("/", server.V2GetGroup(bridge))
			r.Patch("/", server.V2PatchGroup(bridge))
			r.Delete("/", server.V2DeleteGroup(bridge))
		})
	})
}
