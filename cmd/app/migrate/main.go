package migrate

import (
	"embed"
	"log"

	// we need to import the postgres drivers.
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/application"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewCmd(app *application.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "migrates the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := goose.OpenDBWithDriver("pgx", app.Config.DB.GetDSN())
			if err != nil {
				return err
			}

			defer func() {
				if err := db.Close(); err != nil {
					log.Fatalf("goose: failed to close DB: %v\n", err)
				}
			}()

			goose.SetBaseFS(embedMigrations)

			command := args[0]
			var arguments []string
			if len(args) > 1 {
				arguments = append(arguments, args[2:]...)
			}

			if err := goose.Run(command, db, "migrations", arguments...); err != nil {
				log.Fatalf("goose %v: %v", command, err)
			}

			return nil
		},
	}

	return cmd
}
