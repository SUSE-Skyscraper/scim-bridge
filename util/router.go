package util

import (
	"github.com/suse-skyscraper/openfga-scim-bridge/bridge"
	v1middleware "github.com/suse-skyscraper/openfga-scim-bridge/v1/middleware"
	v2middleware "github.com/suse-skyscraper/openfga-scim-bridge/v2/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	v1 "github.com/suse-skyscraper/openfga-scim-bridge/v1/server"
	v2 "github.com/suse-skyscraper/openfga-scim-bridge/v2/server"
)

type AuthorizationMiddleware func(next http.Handler) http.Handler

func Hook(r *chi.Mux, bridge *bridge.Bridge, authHandler AuthorizationMiddleware) {
	r.Route("/scim/v1", func(r chi.Router) {
		r.Use(authHandler)

		r.Get("/Users", v1.V1ListUsers(bridge))
		r.Post("/Users", v1.V1CreateUser(bridge))
		r.Route("/Users/{id}", func(r chi.Router) {
			scimUserCtx := v1middleware.UserCtx(bridge)

			r.Use(scimUserCtx)

			r.Get("/", v1.V1GetUser(bridge))
			r.Put("/", v1.V1UpdateUser(bridge))
			r.Patch("/", v1.V1PatchUser(bridge))
			r.Delete("/", v1.V1DeleteUser(bridge))
		})

		r.Get("/Groups", v1.V1ListGroups(bridge))
		r.Post("/Groups", v1.V1CreateGroup(bridge))
		r.Route("/Groups/{id}", func(r chi.Router) {
			scimGroupCtx := v1middleware.GroupCtx(bridge)

			r.Use(scimGroupCtx)

			r.Get("/", v1.V1GetGroup(bridge))
			r.Patch("/", v1.V1PatchGroup(bridge))
			r.Delete("/", v1.V1DeleteGroup(bridge))
		})
	})

	r.Route("/scim/v2", func(r chi.Router) {
		r.Use(authHandler)

		r.Get("/Users", v2.V2ListUsers(bridge))
		r.Post("/Users", v2.V2CreateUser(bridge))
		r.Route("/Users/{id}", func(r chi.Router) {
			scimUserCtx := v2middleware.UserCtx(bridge)

			r.Use(scimUserCtx)

			r.Get("/", v2.V2GetUser(bridge))
			r.Put("/", v2.V2UpdateUser(bridge))
			r.Patch("/", v2.V2PatchUser(bridge))
			r.Delete("/", v2.V2DeleteUser(bridge))
		})

		r.Get("/Groups", v2.V2ListGroups(bridge))
		r.Post("/Groups", v2.V2CreateGroup(bridge))
		r.Route("/Groups/{id}", func(r chi.Router) {
			scimGroupCtx := v2middleware.GroupCtx(bridge)

			r.Use(scimGroupCtx)

			r.Get("/", v2.V2GetGroup(bridge))
			r.Patch("/", v2.V2PatchGroup(bridge))
			r.Delete("/", v2.V2DeleteGroup(bridge))
		})
	})
}
