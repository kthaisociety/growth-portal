# CLAUDE.md

growth-portal: dashboards for the growth team about the members of KTHAIS's Luma "Community" tier
(graduation year, areas of interest, …). Proof of concept, no sign-in. One repo, one image: a Go server
(Gin, GORM, Postgres) serves the API under `/api` and an embedded Vite + React frontend everywhere else.
Luma data comes from luma-service (`../luma-service`), never from Luma directly.

## Commands

- Local Postgres: `docker compose up -d`
- Config: `cp .env.example .env`, then fill in `LUMA_SERVICE_SECRET`
- luma-service: run it from `../luma-service` (see its README); it listens on 8001
- Backend: `go run ./cmd/server` (on `127.0.0.1:8000`)
- Frontend dev: `cd ui && npm install && npm run dev` (Vite on 5173, proxies `/api` to 8000)
- Go tests: `go vet ./... && go test -race ./...`
- Frontend checks: `cd ui && npm run lint && npm test && npm run build`
- Image: `docker build -t growth-portal .`

## Layout

- `cmd/server`: wiring and startup
- `internal/config`: env, validated at boot
- `internal/handlers`: `/api` routes, and the frontend for every other path
- `internal/demographics`: pure aggregation of members into counts
- `internal/lumaservice`: client for luma-service
- `internal/database`: Postgres connection
- `ui/`: the frontend (Vite, React, TypeScript, Tailwind); `ui/embed.go` embeds `ui/dist`

## Rules

- Read `CONTEXT.md` for vocabulary, `docs/adr/` for decisions and `docs/pcr/` for conventions before
  changing a part they cover. Record a new decision as a new ADR or PCR, never by editing an accepted one.
- The API returns aggregates only, never individual members (PCR 0005): there's no sign-in.
- No `VITE_*` or other build-time config: the frontend calls `/api` on its own origin (PCR 0002).
- Use npm in `ui/`; don't add another lockfile.
- PR titles and commits are Conventional Commits (`feat:`, `fix:`, `docs:`, …): they drive releases.
