# Authentication

Authentication uses username/password login and signed JWT access tokens.

## Routes

| Method | Path | Authentication | Description |
| --- | --- | --- | --- |
| POST | `/api/auth/login` | No | Issue an access token |
| POST | `/api/auth/logout` | Bearer token or cookie | Clear the browser cookie |
| GET | `/api/auth/me` | Bearer token or cookie | Return the current user |

Login with:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"your-username","password":"your-password"}'
```

Login returns the access token in the JSON response and also sets it as an
`HttpOnly`, `SameSite=Lax` cookie. Browser clients should include credentials:

```js
fetch("http://localhost:8080/api/auth/me", { credentials: "include" })
```

Mobile and other non-browser clients can send the returned token on protected
requests:

```http
Authorization: Bearer <access-token>
```

The auth middleware prefers the bearer token and falls back to the cookie, then
stores its subject as the current user ID. Logout clears the cookie; bearer
tokens remain stateless and must be deleted by the client. The cookie is marked
`Secure` when `APP_ENV=production`. Add server-side revocation only when the
application requires it.

The optional default admin account is documented in [Database](database.md).
