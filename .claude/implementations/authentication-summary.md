# Authentication Implementation Summary

## What was built

Full JWT-based authentication with invite-only registration, refresh token rotation, and session-bound expiry.

---

## Backend

### New DB Migrations (`src/db/migrations/`)
| File | Description |
|------|-------------|
| `000003_create_users_table` | `users` table — id, email, username, password_hash, role, timestamps |
| `000004_create_sessions_table` | `sessions` table — UUID pk, user_id FK, expires_at (3 days) |
| `000005_create_refresh_tokens_table` | `refresh_tokens` — hashed token, used flag, session FK (1 day expiry) |
| `000006_create_invites_table` | `invites` — code, email, created_by, used, used_by, expires_at (7 days) |
| `000007_seed_admin_user` | Default admin: username `admin`, password `mygameshelf12345##@@` (bcrypt cost 12) |

### New Models (`src/internal/models/auth.go`)
`User`, `Session`, `RefreshToken`, `Invite` — sensitive fields (`PasswordHash`, `TokenHash`) tagged `json:"-"`.

### New Repositories (`src/internal/repository/`)
| File | Interface | Key methods |
|------|-----------|-------------|
| `user_repository.go` | `UserRepo` | `Create`, `GetByUsername`, `GetByID` |
| `session_repository.go` | `SessionRepo` | `Create`, `GetByID`, `Delete` |
| `refresh_token_repository.go` | `RefreshTokenRepo` | `Create`, `GetByTokenHash`, `MarkUsed`, `DeleteBySessionID` |
| `invite_repository.go` | `InviteRepo` | `Create`, `GetByCode`, `MarkUsed`, `DeletePendingByEmail` |

Added auth sentinel errors to `errors.go`: `ErrUserNotFound`, `ErrDuplicateUsername`, `ErrDuplicateEmail`, `ErrSessionNotFound`, `ErrTokenNotFound`, `ErrInviteNotFound`.

### Email Service (`src/internal/services/email.go`)
- Interface: `EmailService` with `SendInvite(toEmail, inviteCode string) error`
- Implementation: `SMTPEmailService` via `net/smtp.SendMail`
- Skips SMTP auth when `SMTP_USER` is empty (Mailpit local dev)
- Invite link format: `{APP_BASE_URL}/join?code={code}`

### JWT Middleware (`src/api/middleware/auth.go`)
- `AuthRequired(jwtSecret)` — validates `Authorization: Bearer <token>`, injects `user_id`/`username`/`role` into Gin context
- `AdminRequired()` — checks `role == "admin"`, returns 403 otherwise
- Algorithm: HS256 with `JWT_SECRET` env var

### Auth Handler (`src/api/handlers/auth.go`)
| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/register` | POST | none | Validate invite code → bcrypt password → create user → mark invite used |
| `/api/login` | POST | none | bcrypt verify → create session → issue access token (15min) + refresh token (1 day) |
| `/api/refresh` | POST | none | Validate refresh token hash → token theft detection → rotate tokens |
| `/api/logout` | POST | none | Hash token → delete session (cascades all refresh tokens) → 204 |
| `/api/invites` | POST | admin JWT | Delete pending invite for email → new UUID code → save → send email |

**Token theft detection:** if a refresh token is presented but already marked `used`, the entire session is deleted and 401 returned — forcing re-login.

**Security notes:**
- Login returns 401 (not 404) on unknown username — avoids user enumeration
- Refresh tokens stored as SHA-256 hex hashes — raw UUID only transmitted over the wire once
- Access token carries `user_id`, `username`, `role` — no DB lookup per request

### New Go Dependency
`github.com/golang-jwt/jwt/v5`

### Modified Files
- `src/api/middleware/cors.go` — added `Authorization` to allowed headers, handles OPTIONS preflight with 204
- `src/api/router/routes.go` — public vs protected route groups; `RegisterAuthRoutes` function
- `src/api/api.go` — accepts new repos + email service, reads `JWT_SECRET`
- `src/cmd/main.go` — wires all new repos and email service

---

## Frontend

### AuthContext (`web/src/context/AuthContext.tsx`)
- Stores `access_token` + `refresh_token` in `sessionStorage`
- Restores session on mount by decoding JWT payload (no library — pure `atob`)
- Auto-refreshes token 60s before expiry via `setTimeout`
- On token theft (refresh fails): clears all auth state
- Exports `useAuth()` hook: `{ user, accessToken, login, logout }`

### ProtectedRoute (`web/src/components/ProtectedRoute.tsx`)
- `<ProtectedRoute>` — redirects to `/login` (with `state.from`) if not logged in
- `<ProtectedRoute requireAdmin>` — additionally redirects to `/` if not admin

### New Pages
| Route | Page | Description |
|-------|------|-------------|
| `/login` | `LoginPage.tsx` | Username + password, redirects to `state.from` on success |
| `/join` | `RegisterPage.tsx` | Reads `?code` query param, invite code + username + password + confirm |
| `/invite` | `InvitePage.tsx` | Admin-only, sends invite email via API |

### Updated Pages
- `App.tsx` — new routes, `AuthProvider` wrapper, `ProtectedRoute` on `/add`, `/boardgame/:id/edit`, `/invite`
- `Layout.tsx` — conditional nav (Add Game + Invite = admin only), user avatar + logout at sidebar bottom, Sign In link when logged out
- `AddGamePage.tsx` — `Authorization: Bearer` header on POST /api/boardgame and POST image
- `EditGamePage.tsx` — `Authorization: Bearer` header on PUT
- `BoardGameDetailPage.tsx` — `Authorization: Bearer` on DELETE; Edit + Delete buttons wrapped in `{user && ...}`

---

## Local Email (Mailpit)

Mailpit added to both Docker Compose files for catching outgoing emails without a real SMTP server.

**Dev:** `docker-compose.dev.yml`
- SMTP on `localhost:1025` (no auth)
- Web UI at `http://localhost:8025`

**Prod:** `docker-compose.yml`
- Mailpit on internal network only (no exposed ports)
- `app` service configured with `SMTP_HOST=mailpit`, `SMTP_PORT=1025`

`.env` / `.env.example` default to Mailpit settings for local dev.

---

## Environment Variables Added

| Variable | Required | Description |
|----------|----------|-------------|
| `JWT_SECRET` | Yes | 32+ char random string. Generate: `openssl rand -base64 32` |
| `SMTP_HOST` | No | SMTP server host. Dev default: `localhost` (Mailpit) |
| `SMTP_PORT` | No | SMTP port. Dev default: `1025` |
| `SMTP_USER` | No | SMTP username. Empty = no auth (Mailpit) |
| `SMTP_PASSWORD` | No | SMTP password |
| `SMTP_FROM` | No | From address. Dev default: `noreply@mygameshelf.local` |
| `APP_BASE_URL` | No | Base URL for invite links. Dev default: `http://localhost:5173` |

---

## Quick Test

```bash
# Start services
docker compose -f docker-compose.dev.yml up -d

# Run backend
cd src && go run ./cmd/

# Login as default admin
curl -s -X POST localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"mygameshelf12345##@@"}' | jq .

# Send invite (use access_token from above)
curl -s -X POST localhost:8080/api/invites \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{"email":"test@example.com"}' | jq .

# View invite email at http://localhost:8025
# Or get code directly: psql -h localhost -U mygameshelf my_game_shelf -c "SELECT code FROM invites;"
```
