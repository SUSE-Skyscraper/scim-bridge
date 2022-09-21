package application

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	openfga "github.com/openfga/go-sdk"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/db"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/fga"
)

type App struct {
	Config       Config
	Repository   db.RepositoryQueries
	postgresPool *pgxpool.Pool
	FGAClient    fga.Authorizer
}

func NewApp(configDir string) (*App, error) {
	configurator := NewConfigurator(configDir)
	config, err := configurator.Parse()
	if err != nil {
		return &App{}, err
	}
	return &App{
		Config: config,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	database, pool, err := setupDatabase(ctx, a.Config)
	if err != nil {
		return err
	}
	a.Repository = db.NewRepository(pool, database)
	a.postgresPool = pool

	apiClient, err := setupFGA(ctx, a.Config)
	if err != nil {
		return err
	}
	a.FGAClient = fga.NewClient(apiClient)

	return nil
}

func (a *App) Shutdown(_ context.Context) {
	if a.postgresPool != nil {
		a.postgresPool.Close()
	}
}

func setupDatabase(ctx context.Context, config Config) (*db.Queries, *pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(config.DB.GetDSN())
	if err != nil {
		return nil, nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, nil, err
	}
	database := db.New(pool)

	return database, pool, nil
}

func setupFGA(_ context.Context, config Config) (*openfga.APIClient, error) {
	configuration, err := openfga.NewConfiguration(openfga.Configuration{
		ApiScheme: config.FGAConfig.APIScheme,
		ApiHost:   config.FGAConfig.APIHost,
		StoreId:   config.FGAConfig.StoreID,
	})
	if err != nil {
		return nil, err
	}

	apiClient := openfga.NewAPIClient(configuration)

	return apiClient, nil
}
