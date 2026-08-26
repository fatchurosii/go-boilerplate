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

## Documentation

- [Architecture](docs/architecture.md)
- [Configuration and Docker](docs/configuration.md)
- [Authentication](docs/authentication.md)
- [Database, migrations, and seed](docs/database.md)
- [HTTP, responses, and validation](docs/http.md)
- [Adding a feature](docs/adding-a-feature.md)
