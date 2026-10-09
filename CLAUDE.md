# CLAUDE.md

growth-portal: dashboards for the growth team about the members of KTHAIS's Luma "Community" tier
(graduation year, areas of interest, …). Proof of concept; sign-in is simulated until identity-service is
live (PCR 0006). One repo, one image: a Go server (Gin, GORM, Postgres) serves the API under `/api` and an
embedded Vite + React frontend everywhere else. Luma data comes from luma-service (`../luma-service`),
never from Luma directly.

## Commands

- Local Postgres: `docker compose up -d`
- Config: `cp .env.example .env`, then fill in `LUMA_SERVICE_SECRET`
- luma-service: run it from `../luma-service` (see its README); it listens on 8001
- Backend: `go run ./cmd/server` (on `127.0.0.1:8000`)
- Frontend dev: `cd frontend && npm install && npm run dev` (Vite on 5173, proxies `/api` to 8000)
- All checks: `make test` (needs `docker compose up -d`, and `npm install` in `frontend/` once)
- Image: `docker build -t growth-portal .`

## Layout

- `cmd/server`: wiring and startup
- `internal/config`: env, validated at boot
- `internal/auth`: who is signed in; the dev user until identity-service (PCR 0006)
- `internal/handlers`: `/api` routes, and the frontend for every other path
- `internal/demographics`: pure aggregation of members into counts
- `internal/lumaservice`: client for luma-service
- `internal/database`: Postgres connection and migrations (`migrations/*.sql`, applied at startup)
- `frontend/`: the frontend (Vite, React, TypeScript, Tailwind); `frontend/embed.go` embeds `frontend/dist`

## Rules

- Read `CONTEXT.md` for vocabulary, `docs/adr/` for decisions and `docs/pcr/` for conventions before
  changing a part they cover. Record a new decision as a new ADR or PCR, never by editing an accepted one.
- Sign-in is simulated until identity-service is live (PCR 0006): handlers get the user from
  `auth.Current`, and the API returns counts only, never individual members.
- Schema changes are new SQL files in `internal/database/migrations/`, never `AutoMigrate` (PCR 0007).
- No `VITE_*` or other build-time config: the frontend calls `/api` on its own origin (PCR 0001).
- Use npm in `frontend/`; don't add another lockfile.
- PR titles and commits are Conventional Commits (`feat:`, `fix:`, `docs:`, …): they drive releases.
