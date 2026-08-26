# HTTP, Responses, and Validation

`internal/http/router.go` is the single router entry point. It creates Gin,
registers global middleware, groups feature routes under `/api`, and provides
the health and not-found handlers.

## Middleware

Global middleware provides:

- Request IDs from `X-Request-ID` or a generated value.
- Structured request logging with Go's `log/slog`.
- Panic recovery with the standard error response.
- CORS configured from environment variables.

`X-Request-ID` is the single source of truth. It is returned in the response
header and included in logs, but not duplicated in the response body. Auth
middleware is applied only to protected route groups.

```http
X-Request-ID: request-123
```

Clients should use this response header when correlating API errors with
server logs.

## Responses

Use `internal/http/response` from handlers.

```go
response.Success(c, http.StatusOK, "success", result)
response.Error(c, response.BadRequest("invalid request"))
```

Successful responses use this shape:

```json
{
  "message": "success",
  "success": true,
  "data": {}
}
```

Errors use this shape:

```json
{
  "message": "validation failed",
  "success": false,
  "errors": [
    {
      "field": "email",
      "message": "email is a required field"
    }
  ]
}
```

Paginated responses use `response.SuccessPaginate` with `response.Meta`.
Available error helpers include `BadRequest`, `Unauthorized`, `Forbidden`,
`NotFound`, `InternalServerError`, and `ValidationError`.

## Validation

Add validator tags to request DTOs and validate after binding:

```go
type CreateRequest struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

if errs := validation.Validate(req); len(errs) > 0 {
	response.Error(c, response.ValidationError(errs))
	return
}
```

Validation errors use field names from JSON tags and English translations.

The API runs through `net/http.Server` with read, write, and idle timeouts. It
gracefully drains active requests when receiving an interrupt or termination
signal.
