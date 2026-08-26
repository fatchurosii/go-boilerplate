# Database, Migrations, and Seed

The application uses GORM with Postgres. `internal/app` owns the API database
connection, and feature repositories receive `*gorm.DB`.

Set `DATABASE_URL` before running local database commands. Docker Compose maps
`DOCKER_DATABASE_URL` to `DATABASE_URL` inside tool containers.

## Migrations

Migration files live in `migrations/` as matching `up.sql` and `down.sql` pairs.
Run commands from the project root.

Create a migration pair:

```bash
go run ./cmd/migrate create create_users_table
docker compose --profile tools run --rm migrate create create_users_table
```

The generated `up.sql` and `down.sql` files contain placeholders. Fill both
files with the forward and rollback SQL before applying the migration.

Apply all pending migrations or one step:

```bash
go run ./cmd/migrate up
go run ./cmd/migrate up 1
docker compose --profile tools run --rm migrate up
```

Roll back one step and inspect the current version:

```bash
go run ./cmd/migrate down 1
go run ./cmd/migrate version
docker compose --profile tools run --rm migrate down 1
docker compose --profile tools run --rm migrate version
```

Use `force` only after manually correcting a failed migration:

```bash
go run ./cmd/migrate force 1
```

## Seed

Seed the default admin role and user after applying migrations:

```bash
go run ./cmd/seed
docker compose --profile tools run --rm seed
```

The command is idempotent and leaves existing admin records unchanged.

```text
username: admin
password: admin
```

> **Warning:** These credentials are for local development. Do not use the
> default password in a non-local environment.

Repositories should parameterize user input in queries rather than concatenate
it into SQL strings.
