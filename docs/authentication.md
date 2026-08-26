# Authentication

Authentication uses username/password login and signed JWT access tokens.

## Routes

| Method | Path | Authentication | Description |
| --- | --- | --- | --- |
| POST | `/api/auth/login` | No | Issue an access token |
| POST | `/api/auth/logout` | Bearer token | Confirm client-side logout |
| GET | `/api/auth/me` | Bearer token | Return the current user |

Login with:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"your-username","password":"your-password"}'
```

Send the returned token on protected requests:

```http
Authorization: Bearer <access-token>
```

The auth middleware validates the token and stores its subject as the current
user ID. Logout is stateless: the client deletes its token. Add server-side
revocation only when the application requires it.

The optional default admin account is documented in [Database](database.md).
