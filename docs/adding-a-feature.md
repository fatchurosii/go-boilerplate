# Adding a Feature

`internal/example` is the reference for adding an HTTP feature.

## Example Routes

| Method | Path | Description |
| --- | --- | --- |
| GET | `/api/examples` | List examples |
| GET | `/api/examples/page?page=1&limit=10` | List paginated examples |
| POST | `/api/examples` | Create an example |

Create an example with:

```bash
curl -X POST http://localhost:8080/api/examples \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

## Feature Flow

1. Create a package under `internal/<feature>`.
2. Add a handler for request binding, validation, service calls, and responses.
3. Add a service for business logic and domain errors.
4. Add a route file with the feature's `RegisterRoutes` method.
5. Add a repository only when the feature needs database access.
6. Wire dependencies in `internal/app/app.go`.
7. Pass the handler to `internal/http/router.go` and register its routes.

Keep endpoint details in the feature route file rather than the root router.
Map service errors to API responses in the handler.
