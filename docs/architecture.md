# Architecture

The application uses manual dependency injection and keeps startup, wiring,
HTTP transport, and feature logic separate.

## Request Flow

```text
cmd/api/main.go
-> internal/app/app.go
-> internal/http/router.go
-> internal/<feature>
```

## Package Ownership

- `cmd/api` loads configuration, creates the logger and app, and owns the HTTP server lifecycle.
- `internal/app` opens shared dependencies and wires repositories, services, handlers, and the router.
- `internal/http` owns the Gin router, global middleware, health check, and not-found response.
- `internal/auth` owns authentication routes, handlers, token logic, and services.
- `internal/example` is the reference feature for adding endpoints.
- `internal/user` and `internal/role` own their persistence models; `internal/user` also owns its repository.
- `internal/database` provides shared database infrastructure.
- `internal/http/response` and `internal/http/validation` are used only at the HTTP boundary.

Keep feature business logic out of `cmd/api`, `internal/app`, and the root HTTP
router. A feature owns its handlers, service, route registration, and optional
repository.

Services return domain or standard errors. Handlers map those errors to HTTP
responses so feature logic does not depend on Gin or HTTP status codes.
