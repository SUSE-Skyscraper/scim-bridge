# SCIM Bridge

This is a PoC (work in progress) for a SCIM Bridge taken from the skyscraper project.

The implementation is still in a rough state and is missing essential tests, but covers SCIM V2 for Okta.

Expect major changes as this project matures.

## Project Goals

* To have a SCIM V2 bridge.
* To initially cover Okta's implementation of SCIM, but eventually fulfil the entire specification.
* To be flexible enough to allow for individual applications to implement their own business logic.

## Development Environment

Please use the example application located in the [example](./example) directory for development. It serves as a sample application using SCIM bridge library and OpenFGA.

### Prerequisites

**Local Environment:**

* PostgresSQL
* golang 1.19+
* sqlc
   * `go install github.com/kyleconroy/sqlc/cmd/sqlc@latest`

### Database Migrations

The database migration files are at [cmd/app/migrate/migrations](example/cmd/app/migrate/migrations). They're embedded into the binary, and we read them in the `migrate` command.

**Migrate Up:**
```bash
go run ./cmd/main.go migrate up
```

**Migrate Down:**
```bash
go run ./cmd/main.go migrate down
```

### Generate Database files

**Notes:**

* The queries are located in the `queries.sql` file.
* The config file for `sqlc` is located at `sqlc.yaml`.
* The database files are generated at `internal/db`.

Run the following command to generate the database files:

```bash
sqlc generate
```

### Deploy Locally

The example application is located in the [example](example) directory.

1. Copy `config.yaml.example` to `config.yaml` and fill in the values.
2. Build the server:
   ```bash
   go build ./cmd/main.go
   ```
3. Run database migrations:
   ```bash
   go run ./cmd/main.go migrate up
   ```
4. Run the server:
   ```bash
   go run ./cmd/main.go server
   ```
