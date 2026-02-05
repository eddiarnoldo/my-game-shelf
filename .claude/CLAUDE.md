# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Backend (Go)
- **Run backend:** `go run ./cmd/` (from repo root)
- **Run tests:** `go test ./src/api/handlers/...` (add `-v` for verbose)
- **Create a new DB migration:** `migrate create -ext sql -dir db/migrations -seq <name>`

### Frontend (React)
- **Run dev server:** `cd web && npm run dev`
- **Build:** `cd web && npm run build`
- **Lint:** `cd web && npm run lint`

### Database & Docker
- **Start DB only (dev):** `docker-compose -f docker-compose.dev.yml up -d`
- **Stop DB:** `docker-compose -f docker-compose.dev.yml down`
- **Connect to DB:** `psql -h localhost -p 5432 -U mygameshelf my_game_shelf`
- **Full production stack:** `docker-compose up -d`

## Architecture

This is a full-stack board game shelf app: a **Go/Gin backend** with a **React/TypeScript/Vite frontend**, backed by **PostgreSQL**.

### Backend layout (`src/`)
- `cmd/main.go` — Entry point. Loads `.env`, runs migrations, connects to DB, wires repositories into handlers, starts the Gin server on port 8080.
- `api/` — HTTP layer: `api.go` (server init), `handlers/` (request handlers), `middleware/cors.go`, `router/routes.go` (route registration).
- `internal/` — Domain logic: `models/` (structs), `repository/` (SQL queries via pgx, interface-based), `helpers/` (thumbnail generation).
- `db/migrations/` — SQL migrations run automatically at startup via `migrate.go`.
- `config/config.go` — Reads environment variables.

Handlers depend on repository **interfaces** (`BoardGameRepo`, `BoardGameImageRepo`), not concrete types. This is the seam for future mocking/testing.

### Frontend layout (`web/src/`)
- `App.tsx` — React Router setup: `/` → HomePage, `/boardgame/:id` → DetailPage, `/add` → AddGamePage, `/*` → NotFoundPage.
- `components/Layout.tsx` — Sidebar navigation wrapper used by all pages.
- `components/BoardGameCard.tsx` — Reusable card shown on the home grid; loads cover image from the API with a fallback placeholder.
- `pages/` — Each page uses local `useState` for state and the native Fetch API for data; no global state library.
- Vite proxies `/api` requests to `localhost:8080` during development (see `vite.config.ts`).

### Image handling
Images (cover + gameplay) are stored as raw bytes (`BYTEA`) in the `board_game_images` table. On upload the backend generates a thumbnail (max 300px width) via `internal/helpers`. The cover endpoint (`GET /api/boardgame/:id/images/cover`) returns the **thumbnail** with a 24-hour cache header — this is what the frontend uses for cards.

### Notes.md — target API shape
`Notes.md` documents the desired response shape for board game objects including `coverImageUrl` and an `images` array with gameplay images. This is the direction the API is being evolved toward; current endpoints return images as separate requests.

## Environment
Copy `.env` to set DB credentials and `APP_PORT`. `ALLOWED_ORIGINS` controls CORS (defaults to `*`).
