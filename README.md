# SCIM Bridge

This is a PoC (work in progress) for a SCIM Bridge taken from the skyscraper project.

The implementation is still in a rough state and is missing essential tests, but covers SCIM V2 for Okta.

Expect major changes as this project matures.

## SCIM References

* V1:
  * http://www.simplecloud.info/specs/draft-scim-core-schema-01.html
  * http://www.simplecloud.info/specs/draft-scim-api-01.html

## Project Goals

* To have a SCIM V1 and V2 bridge.
* To initially cover Okta's implementation of SCIM, but eventually fulfil the entire specification.
* To be flexible enough to allow for individual applications to implement their own business logic.

## Development Environment

Please use the example application located in the [example](./example) directory for development. It serves as a sample application using SCIM bridge library and OpenFGA.

### Prerequisites

**Local Environment:**

* PostgresSQL
* golang 1.18+
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
