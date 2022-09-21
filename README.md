# OpenFGA SCIM Bridge

This is a PoC (work in progress) for an OpenFGA SCIM Bridge taken from the skyscraper project.

## Development Environment

### Prerequisites

**Local Environment:**

* PostgresSQL
* make
* golang 1.18+
* sqlc
   * `go install github.com/kyleconroy/sqlc/cmd/sqlc@latest`

### Database Migrations

The database migration files are at [cmd/app/migrate/migrations](cmd/app/migrate/migrations). They're embedded into the binary, and we read them in the `migrate` command.

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

1. Ensure that you have a PostgresSQL server that you can connect to locally.
2. Copy `config.yaml.example` to `config.yaml` and fill in the values.
3. Build the server:
   ```bash
   go build ./cmd/main.go
   ```
4. Run database migrations:
   ```bash
   go run ./cmd/main.go migrate up
   ```
5. Run the server:
   ```bash
   go run ./cmd/main.go server
   ```
