package server

import (
	"github.com/suse-skyscraper/openfga-scim-bridge/example/internal/application"
	"github.com/suse-skyscraper/openfga-scim-bridge/example/internal/scimbridgedb"
	"github.com/suse-skyscraper/openfga-scim-bridge/example/internal/server"
	"github.com/suse-skyscraper/openfga-scim-bridge/example/internal/server/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/bridge"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/router"
)

func NewCmd(app *application.App) *cobra.Command {
	baseURL := "http://localhost:8080"

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := chi.NewRouter()

			// common middleware
			r.Use(chimiddleware.Logger)

			authMiddleware := middleware.BearerAuthorizationHandler(app)

			r.Get("/healthz", server.Health)

			db := scimbridgedb.New(app)
			b := bridge.New(&db, baseURL)
			router.Hook(r, &b, authMiddleware)

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
