package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/openfga-scim-bridge/cmd/app/migrate"
	"github.com/suse-skyscraper/openfga-scim-bridge/cmd/app/scim"
	"github.com/suse-skyscraper/openfga-scim-bridge/cmd/app/server"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/application"
)

func main() {
	ctx := context.Background()
	app, err := application.NewApp(application.DefaultConfigDir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create application: %v\n", err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use: "openfga-scim-bridge",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := app.Start(ctx)
			if err != nil {
				return err
			}

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			app.Shutdown(ctx)
		},
	}

	rootCmd.AddCommand(server.NewCmd(app))
	rootCmd.AddCommand(migrate.NewCmd(app))
	rootCmd.AddCommand(scim.NewCmd(app))

	err = rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
