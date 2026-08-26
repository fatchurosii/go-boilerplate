# Configuration and Docker

`.env.example` is the canonical list of supported environment variables and
example values. Copy it to `.env` for local overrides.

## Application Configuration

- `APP_ENV` controls application mode. `production` enables Gin release mode.
- `PORT` is the HTTP server port.
- `DATABASE_URL` is the Postgres connection used by local commands and the API.
- `CORS_ALLOWED_ORIGINS` is a comma-separated list of allowed browser origins.
- `CORS_ALLOW_CREDENTIALS` enables browser credentials for explicit origins.
- `JWT_SECRET`, `JWT_ISSUER`, and `JWT_ACCESS_TOKEN_TTL` configure access tokens.

The API requires `DATABASE_URL` and a non-empty `JWT_SECRET` to start.
Invalid ports, booleans, durations, and non-positive JWT token lifetimes fail
startup instead of silently using defaults.

## Docker Configuration

- `DOCKER_DATABASE_URL` becomes `DATABASE_URL` inside app containers.
- `PROXY_NETWORK` connects the API to shared Nginx.
- `PROXY_ALIAS` is the stable API hostname used by shared Nginx.
- `DB_NETWORK` connects the API and tools to shared Postgres.

Use the Postgres container hostname in `DOCKER_DATABASE_URL`, not `localhost`.

Shared Nginx, Postgres, RustFS, LGTM, and external network setup are documented
in the [shared infra README](../../infra/README.md).

## Docker Commands

Build and start the API:

```bash
docker compose up -d --build
```

View API logs:

```bash
docker compose logs -f api
```

Stop or remove the app containers:

```bash
docker compose stop
docker compose down
```

Configure `GO_BOILERPLATE_DOMAIN` in `../infra/.env` for shared Nginx. Database
tool commands are documented in [Database](database.md).
