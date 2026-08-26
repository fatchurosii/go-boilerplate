# Go Boilerplate

Minimal Gin API boilerplate with manual dependency injection, GORM/Postgres,
JWT authentication, standard responses, validation, request IDs, and `slog`
logging.

## Quick Start

Copy `.env.example` to `.env`, update `DATABASE_URL` and `JWT_SECRET`, then run
the following commands from the project root:

```bash
go run ./cmd/migrate up
go run ./cmd/api
```

The API listens on `http://localhost:8080` by default. Check it with:

```bash
curl http://localhost:8080/api/health
```

## Use as a Template

After creating a repository from this template, rename the Go module from the
project root:

```bash
bash scripts/rename.sh github.com/your-user/your-project
```

The script reads the current module from `go.mod`, updates `go.mod` and all Go
import paths, then runs `go fmt`, `go mod tidy`, tests, and a full build. It is
safe to keep in the repository for future template users.

Run it with Bash on Linux or macOS, or with Git Bash/WSL on Windows. Commit or
back up local changes first because the rename updates files in place.

Runtime identifiers such as `JWT_ISSUER`, `PROXY_ALIAS`, database names, and
repository documentation are not renamed automatically; update them manually
when the new project needs different values.

## Documentation

- [Architecture](docs/architecture.md)
- [Configuration and Docker](docs/configuration.md)
- [Authentication](docs/authentication.md)
- [Database, migrations, and seed](docs/database.md)
- [HTTP, responses, and validation](docs/http.md)
- [Adding a feature](docs/adding-a-feature.md)
