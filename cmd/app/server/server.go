package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/application"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/scim"
	scimmiddleware "github.com/suse-skyscraper/openfga-scim-bridge/internal/scim/middleware"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/server"
)

func NewCmd(app *application.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := chi.NewRouter()

			// common middleware
			r.Use(chimiddleware.Logger)

			r.Get("/healthz", server.Health)

			r.Route("/scim/v2", func(r chi.Router) {
				scimAuthorizer := scimmiddleware.BearerAuthorizationHandler(app)

				r.Use(scimAuthorizer)

				r.Get("/Users", scim.V2ListUsers(app))
				r.Post("/Users", scim.V2CreateUser(app))
				r.Route("/Users/{id}", func(r chi.Router) {
					scimUserCtx := scimmiddleware.UserCtx(app)

					r.Use(scimUserCtx)

					r.Get("/", scim.V2GetUser(app))
					r.Put("/", scim.V2UpdateUser(app))
					r.Patch("/", scim.V2PatchUser(app))
					r.Delete("/", scim.V2DeleteUser(app))
				})

				r.Get("/Groups", scim.V2ListGroups(app))
				r.Post("/Groups", scim.V2CreateGroup(app))
				r.Route("/Groups/{id}", func(r chi.Router) {
					scimGroupCtx := scimmiddleware.GroupCtx(app)

					r.Use(scimGroupCtx)

					r.Get("/", scim.V2GetGroup(app))
					r.Patch("/", scim.V2PatchGroup(app))
					r.Delete("/", scim.V2DeleteGroup(app))
				})
			})

			s := &http.Server{
				Addr:         ":8080",
				Handler:      r,
				ReadTimeout:  2 * time.Second,
				WriteTimeout: 2 * time.Second,
			}
			err := s.ListenAndServe()
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
